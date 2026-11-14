package student

import (
	"errors"
	"fmt"
	"net/http"

	"student/internal/handler"
	"student/pkg/dict"
)

type Handler struct {
	*handler.Base
	service Service
}

func NewHandler(base *handler.Base, svc Service) *Handler {
	return &Handler{
		Base:    base,
		service: svc,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := h.Render.Param(r, "id")
	student, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	h.Render.JSON(w, r, ToStudentResponse(&student))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	students, err := h.service.List(r.Context())
	if err != nil {
		h.fail(w, r, err)
		return
	}
	h.Render.JSON(w, r, ToStudentResponses(students))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req StudentRequest
	if err := h.Render.Decode(r, &req); err != nil {
		h.fail(w, r, fmt.Errorf("%w: %v", ErrInvalidInput, err))
		return
	}
	if err := h.service.Create(r.Context(), req); err != nil {
		h.fail(w, r, err)
		return
	}
	h.Render.JSON(w, r, map[string]string{"message": "student created"})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var req StudentUpdateRequest
	if err := h.Render.Decode(r, &req); err != nil {
		h.fail(w, r, fmt.Errorf("%w: %v", ErrInvalidInput, err))
		return
	}
	if err := h.service.Update(r.Context(), req); err != nil {
		h.fail(w, r, err)
		return
	}
	h.Render.JSON(w, r, map[string]string{"message": "student updated"})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := h.Render.Param(r, "id")
	if err := h.service.Delete(r.Context(), id); err != nil {
		h.fail(w, r, err)
		return
	}
	h.Render.JSON(w, r, map[string]string{"message": "student deleted"})
}

// fail 把领域错误映射为 HTTP 状态码，并以统一的 JSON 格式写出。
// 具体错误只记录到日志，不返回给客户端。
func (h *Handler) fail(w http.ResponseWriter, r *http.Request, err error) {
	h.Logger.Error("StudentHandler Error", "err", err)
	switch {
	case errors.Is(err, ErrNotFound):
		h.Render.Error(w, r, http.StatusNotFound, dict.ErrNotFound)
	case errors.Is(err, ErrInvalidInput):
		h.Render.Error(w, r, http.StatusBadRequest, dict.ErrBadRequest)
	default:
		h.Render.Error(w, r, http.StatusInternalServerError, dict.ErrInternal)
	}
}
