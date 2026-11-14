package server

import (
	"fmt"
	"net/http"

	"student/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

// Router 是业务模块向传输层注册路由的最小契约。
// 传输层只依赖该接口，不感知任何具体模块的 Handler。
type Router interface {
	Routes(r chi.Router)
}

// NewHTTPServer 组装 HTTP 服务：全局中间件 + 各模块路由。
func NewHTTPServer(addr string, modules ...Router) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.JWTVerifier())

	RegisterRoutes(r, modules...)

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

// RegisterRoutes 挂载版本化路由与各业务模块。
// 新增模块只需实现 Router 并在装配处传入，本文件无需修改。
func RegisterRoutes(r chi.Router, modules ...Router) {
	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.JWTAuthenticator())

		for _, m := range modules {
			m.Routes(r)
		}

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		})
		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			claims := middleware.GetClaims(r)
			w.Write(fmt.Appendf([]byte{}, "protected area. hi %v\n", claims["user_id"]))
		})
	})
}
