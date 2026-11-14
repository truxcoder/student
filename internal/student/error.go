package student

import "errors"

// 领域哨兵错误：由仓储/服务返回，handler 通过 errors.Is 映射为 HTTP 状态码。
var (
	// ErrNotFound 表示学生记录不存在。
	ErrNotFound = errors.New("student not found")
	// ErrInvalidInput 表示请求内容不合法（如请求体解析失败）。
	ErrInvalidInput = errors.New("invalid student input")
)
