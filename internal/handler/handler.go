package handler

import (
	"github.com/truxcoder/log"
)

type Handler struct {
	Logger *log.Logger
	Render *Render
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		Logger: logger,
		Render: &Render{},
	}
}
