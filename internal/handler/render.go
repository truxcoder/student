package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Render 是传输层统一的请求解析与响应写出工具。
// 只提供与框架交互的最小方法，DTO 不需要实现任何框架接口。
type Render struct{}

// Decode 按请求的 Content-Type 解析请求体。
func (rd *Render) Decode(r *http.Request, v interface{}) error {
	return render.Decode(r, v)
}

// JSON 以 JSON 格式写出响应。
func (rd *Render) JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	render.JSON(w, r, v)
}

// Error 以统一的 JSON 格式写出错误响应，避免与 http.Error 的纯文本格式混用。
func (rd *Render) Error(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, map[string]string{"error": message})
}

// Param 读取 URL 路由参数。
func (rd *Render) Param(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
