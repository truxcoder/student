package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"student/internal/platform/database"
	"student/internal/sqlc"
)

// Base 是所有仓储共享的基座：只负责连接与事务上下文，不得出现任何业务概念。
// 字段由构造函数统一注入，仓储实现应通过 Q(ctx) 访问数据库，
// 以便自动加入调用方开启的事务。
type Base struct {
	DB      *sql.DB
	Queries *sqlc.Queries
}

// New 基于共享数据库句柄创建仓储基座。
func New(handle *database.Handle) *Base {
	return &Base{
		DB:      handle.DB,
		Queries: sqlc.New(handle.DB),
	}
}

// Q 返回与当前事务绑定的查询句柄；ctx 中未开启事务时返回默认句柄。
func (b *Base) Q(ctx context.Context) *sqlc.Queries {
	if tx := database.FromContext(ctx); tx != nil {
		return b.Queries.WithTx(tx)
	}
	return b.Queries
}

// NewDB 打开数据库并返回共享连接句柄。
func NewDB(dsn string) (*database.Handle, error) {
	dir, _ := os.Getwd()
	file := filepath.Join(dir, dsn)
	fmt.Printf("Connecting to database at: %s\n", file)
	conn, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return database.NewHandle(conn), nil
}
