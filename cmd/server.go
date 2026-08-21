package main

import (
	"context"
	"fmt"
	"student/internal/handler"
	"student/internal/repository"
	"student/internal/server"
	"student/internal/service"
	"student/internal/student"
	"student/pkg/app"
	"student/pkg/config"
	"student/pkg/logger"
	"student/pkg/myhttp"

	_ "modernc.org/sqlite"
)

func main() {
	l := logger.New()
	cfg, err := config.NewFromFile("config/config.json")
	if err != nil {
		l.Fatal("Failed to load config", "err", err)
	}

	name := config.GetOr(cfg, "app.name", "student")
	dns := config.GetOr(cfg, "db.dns", "")
	host := config.GetOr(cfg, "app.host", "localhost")
	port := config.GetOr(cfg, "app.port", 3000)
	addr := fmt.Sprintf("%s:%d", host, port)

	ctx := context.Background()

	db, err := repository.NewDB(dns)
	if err != nil {
		l.Fatal("Failed to connect to database", "err", err)
	}
	repo := repository.NewRepository(db)
	studentRepository := student.NewStudentRepository(repo)
	svc := service.NewService(l)
	studentService := student.NewStudentService(svc, studentRepository)
	hdl := handler.NewHandler(l)
	studentHandler := student.NewStudentHandler(hdl, studentService)
	chiServer := server.NewHTTPServer(addr, studentHandler)
	httpServer := myhttp.NewHTTPServer(chiServer, l, func() error {
		l.Info("cleanup...")
		if err := db.Close(); err != nil {
			l.Error("Failed to close database", "err", err)
			return err
		}
		l.Info("Database closed successfully")
		return nil
	})

	a := app.NewApp(app.WithName(name), app.WithServer(httpServer))
	err = a.Run(ctx)
	if err != nil {
		l.Fatal("Failed to run app", "err", err)
	}
}
