package database

import (
	"context"
	"database/sql"
)

// TxManager 定义事务边界。
//
// 应用服务通过 Do 开启事务，仓储通过 FromContext 自动加入当前事务，
// 无需显式传递 *sql.Tx，也不需要在方法签名上暴露 tx 参数。
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

// Handle 是运行期共享的数据库句柄：连接池 + 事务管理器。
type Handle struct {
	DB *sql.DB
	Tx TxManager
}

// NewHandle 基于已建立的连接创建共享句柄。
func NewHandle(db *sql.DB) *Handle {
	return &Handle{
		DB: db,
		Tx: &manager{db: db},
	}
}

// Close 关闭连接池。
func (h *Handle) Close() error {
	return h.DB.Close()
}

type manager struct {
	db *sql.DB
}

func (m *manager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(withTx(ctx, tx)); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

type txKey struct{}

func withTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// FromContext 返回 ctx 中绑定的事务，未处于事务中时返回 nil。
func FromContext(ctx context.Context) *sql.Tx {
	tx, _ := ctx.Value(txKey{}).(*sql.Tx)
	return tx
}
