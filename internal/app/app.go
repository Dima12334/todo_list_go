package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo_list_go/internal/config"
	"todo_list_go/internal/db"
	"todo_list_go/internal/handlers"
	"todo_list_go/internal/repository"
	"todo_list_go/internal/server"
	"todo_list_go/internal/service"
	"todo_list_go/pkg/logger"
)

func Run(configDir string) {
	if err := logger.Init(); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	cfg, err := config.Init(configDir)
	if err != nil {
		logger.Errorf("failed to init configs: %v", err.Error())
		return
	}

	dbConn, err := db.ConnectDB(cfg.DB)
	if err != nil {
		logger.Errorf("failed to connect to database: %v", err.Error())
		return
	}
	defer func() {
		if err := dbConn.Close(); err != nil {
			logger.Errorf("error occurred on db connection close: %s", err.Error())
		} else {
			logger.Info("db connection closed successfully")
		}
	}()

	repositories := repository.NewRepositories(dbConn)
	services := service.NewServices(
		service.Deps{
			Repos:          repositories,
			AccessTokenTTL: cfg.Auth.JWT.AccessTokenTTL,
		},
	)
	handler := handlers.NewHandler(services)

	srv := server.NewServer(cfg, handler.Init())

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()
	logger.Info("server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("failed to stop server: %v", err.Error())
	} else {
		logger.Info("server stopped successfully")
	}
}
