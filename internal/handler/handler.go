package handler

import (
	"github.com/truxcoder/log"
)

// Base 是 HTTP 处理器的共享基座，只承载横切能力，不得出现任何业务概念。
// 字段由构造函数统一注入，模块内只读，不得覆写。
type Base struct {
	Logger *log.Logger
	Render *Render
}

func New(logger *log.Logger) *Base {
	return &Base{
		Logger: logger,
		Render: &Render{},
	}
}
