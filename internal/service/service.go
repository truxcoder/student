package service

import (
	"student/internal/platform/database"

	"github.com/truxcoder/log"
)

// Base 是应用服务的共享基座，只承载横切能力，不得出现任何业务概念。
// 字段由构造函数统一注入，模块内只读，不得覆写。
type Base struct {
	Logger *log.Logger
	Tx     database.TxManager
}

func New(logger *log.Logger, tx database.TxManager) *Base {
	return &Base{
		Logger: logger,
		Tx:     tx,
	}
}
