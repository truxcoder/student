package student

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "modernc.org/sqlite"

	"student/internal/platform/database"
	"student/internal/repository"
	"student/internal/sqlc"
)

// newTestRepository 为每个测试创建独立的内存数据库，测试之间互不影响。
func newTestRepository(t *testing.T) (Repository, *database.Handle) {
	t.Helper()

	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("打开内存数据库: %v", err)
	}
	// :memory: 库绑定在单条连接上，限制连接数才能保证所有操作看到同一份数据。
	conn.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = conn.Close() })

	schema, err := os.ReadFile(filepath.Join("..", "..", "sqlc", "schema.sql"))
	if err != nil {
		t.Fatalf("读取 schema.sql: %v", err)
	}
	if _, err := conn.Exec(string(schema)); err != nil {
		t.Fatalf("初始化表结构: %v", err)
	}

	handle := database.NewHandle(conn)
	return NewRepository(repository.New(handle)), handle
}

func newCreateParams(id string) sqlc.CreateStudentParams {
	return sqlc.CreateStudentParams{
		ID:        id,
		Name:      "张三",
		Gender:    "男",
		Birthday:  time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
		Classroom: "一班",
	}
}

func TestRepositoryGetNotFound(t *testing.T) {
	repo, _ := newTestRepository(t)

	if _, err := repo.Get(context.Background(), "missing-id"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get 未命中时应返回 ErrNotFound, got %v", err)
	}
}

func TestRepositoryUpdateNotFound(t *testing.T) {
	repo, _ := newTestRepository(t)

	name := "李四"
	params := sqlc.UpdateStudentParams{ID: "missing-id", Name: &name}
	if err := repo.Update(context.Background(), params); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Update 影响 0 行时应返回 ErrNotFound, got %v", err)
	}
}

func TestRepositoryDeleteNotFound(t *testing.T) {
	repo, _ := newTestRepository(t)

	if err := repo.Delete(context.Background(), "missing-id"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Delete 影响 0 行时应返回 ErrNotFound, got %v", err)
	}
}

func TestRepositoryTxCommit(t *testing.T) {
	repo, handle := newTestRepository(t)
	ctx := context.Background()

	if err := handle.Tx.Do(ctx, func(ctx context.Context) error {
		return repo.Create(ctx, newCreateParams("tx-commit"))
	}); err != nil {
		t.Fatalf("事务提交失败: %v", err)
	}

	got, err := repo.Get(ctx, "tx-commit")
	if err != nil {
		t.Fatalf("事务提交后应能读到记录: %v", err)
	}
	if got.ID != "tx-commit" {
		t.Fatalf("ID = %q, want tx-commit", got.ID)
	}
}

func TestRepositoryTxRollback(t *testing.T) {
	repo, handle := newTestRepository(t)
	ctx := context.Background()
	boom := errors.New("boom")

	err := handle.Tx.Do(ctx, func(ctx context.Context) error {
		if err := repo.Create(ctx, newCreateParams("tx-rollback")); err != nil {
			return err
		}
		return boom
	})
	if !errors.Is(err, boom) {
		t.Fatalf("事务应原样返回回调错误, got %v", err)
	}

	if _, err := repo.Get(ctx, "tx-rollback"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("事务回滚后记录不应存在, got %v", err)
	}
}
