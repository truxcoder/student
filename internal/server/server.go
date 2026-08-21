package server

import "github.com/truxcoder/log"

type Server struct {
	logger *log.Logger
}

func New(logger *log.Logger) *Server {
	return &Server{logger: logger}
}
