package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request) error
}

// Binder interface for managing request payloads.
type Binder interface {
	Bind(r *http.Request) error
}

type Render struct {
}

func (rd *Render) Bind(r *http.Request, v Binder) error {
	return render.Bind(r, v)
}

func (rd *Render) Render(w http.ResponseWriter, r *http.Request, v Renderer) error {
	var vv render.Renderer = v
	return render.Render(w, r, vv)
}

func (rd *Render) RenderList(w http.ResponseWriter, r *http.Request, l []Renderer) error {
	list := make([]render.Renderer, len(l))
	for i, v := range l {
		list[i] = v
	}
	return render.RenderList(w, r, list)
}

func (rd *Render) JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	render.JSON(w, r, v)
}

func (rd *Render) Param(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
