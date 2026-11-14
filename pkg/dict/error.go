package dict

// 统一错误消息文案，由 handler 在写出错误响应时使用。
// 响应格式统一在 internal/handler.Render.Error 中实现。
const (
	ErrNotFound     = "未找到记录"
	ErrInternal     = "内部错误"
	ErrBadRequest   = "请求错误"
	ErrUnauthorized = "未授权"
	ErrForbidden    = "禁止访问"
)
