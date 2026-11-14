package student

import "github.com/go-chi/chi/v5"

// Routes 注册学生模块的 HTTP 路由。
// 模块自己拥有路由表，传输层只负责把子路由挂载到版本前缀下。
func (h *Handler) Routes(r chi.Router) {
	r.Route("/student", func(r chi.Router) {
		r.Get("/", h.List)
		r.Post("/", h.Create)
		r.Patch("/", h.Update)
		r.Get("/{id}", h.Get)
		r.Delete("/{id}", h.Delete)
	})
}
