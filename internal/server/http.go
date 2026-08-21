package server

import (
	"fmt"
	"net/http"

	"student/internal/student"
	"student/pkg/middleware"

	"github.com/go-chi/chi/v5"
)

func NewHTTPServer(addr string, studentHandler *student.StudentHandler) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.JWTVerifier())
	r.Use(middleware.JWTAuthenticator())

	RegisterRoutes(r, studentHandler)
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func RegisterRoutes(r chi.Router, studentHandler *student.StudentHandler) {
	r.Route("/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuthenticator())
			r.Get("/student/{id}", studentHandler.Get)
			r.Get("/student", studentHandler.List)
			r.Post("/student", studentHandler.Create)
			r.Put("/student", studentHandler.Update)
			r.Delete("/student/{id}", studentHandler.Delete)
		})
		r.Group(func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello, World!"))
			})
			r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
				claims := middleware.GetClaims(r)
				w.Write(fmt.Appendf([]byte{}, "protected area. hi %v\n", claims["user_id"]))
			})
		})

	})
}
