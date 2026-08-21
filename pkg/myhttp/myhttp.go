package myhttp

import (
	"context"
	"errors"
	"net/http"
	"student/pkg/app"

	"github.com/truxcoder/log"
)

type httpServer struct {
	logger     *log.Logger
	httpServer *http.Server
	cleanup    func() error
}

func NewHTTPServer(server *http.Server, logger *log.Logger, cleanup func() error) app.Server {
	return &httpServer{
		logger:     logger,
		httpServer: server,
		cleanup:    cleanup,
	}
}

func (s *httpServer) Start(ctx context.Context) error {
	s.logger.Info("HTTP Server Starting on " + s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("HTTP Server failed to start", "err", err)
		return err
	}
	return nil
}

func (s *httpServer) Stop(ctx context.Context) error {
	s.logger.Info("HTTP Server Stopping")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	s.logger.Info("HTTP Server Stopped")

	// 清理其它资源，如释放数据库连接等。
	if err := s.cleanup(); err != nil {
		s.logger.Error("Cleanup failed", "err", err)
	}
	return nil
}
