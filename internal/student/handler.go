package student

import (
	"net/http"
	"student/api"
	"student/internal/handler"
	"student/pkg/dict"
	"student/pkg/util"
)

type StudentHandler struct {
	*handler.Handler
	service StudentService
}

func NewStudentHandler(handler *handler.Handler, service StudentService) *StudentHandler {
	return &StudentHandler{
		Handler: handler,
		service: service,
	}
}

func (h *StudentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := h.Render.Param(r, "id")
	student, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.Logger.Error("StudentHandler.Get Error", err.Error())
		if err.Error() == "sql: no rows in result set" {
			dict.HandleNotFound(w)
			return
		}
		dict.HandleInternal(w)
		return
	}
	if err := h.Render.Render(w, r, ToStudentResponse(&student)); err != nil {
		dict.HandleInternal(w)
		return
	}
}

func (h *StudentHandler) List(w http.ResponseWriter, r *http.Request) {
	students, err := h.service.List(r.Context())
	if err != nil {
		h.Logger.Error("StudentHandler.List Error", err.Error())
		if err.Error() == "sql: no rows in result set" {
			dict.HandleNotFound(w)
			return
		}
		dict.HandleInternal(w)
		return
	}

	if err := h.Render.RenderList(w, r, util.ToRenderers(ToStudentResponses(students))); err != nil {
		dict.HandleInternal(w)
		return
	}
}

func (h *StudentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req api.StudentRequest
	if err := h.Render.Bind(r, &req); err != nil {
		h.Logger.Error("Bind Error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.Create(r.Context(), req.ToCreateParams()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Render.JSON(w, r, map[string]string{"message": "student created"})
}

func (h *StudentHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req api.StudentRequest
	if err := h.Render.Bind(r, &req); err != nil {
		h.Logger.Error("Bind Error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.service.Update(r.Context(), req.ToUpdateParams()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.Render.JSON(w, r, map[string]string{"message": "student updated"})
}

func (h *StudentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := h.Render.Param(r, "id")
	_, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.Logger.Error("StudentHandler.Delete Error", err.Error())
		if err.Error() == "sql: no rows in result set" {
			dict.HandleNotFound(w)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		h.Logger.Error("StudentHandler.Delete Error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Render.JSON(w, r, map[string]string{"message": "student deleted"})
}
