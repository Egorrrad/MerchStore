///go:generate oapi-codegen --config=../../api/cfg.yaml ../../api/schema.yaml

package main

import (
	"MerchStore/src/cmd"
	"MerchStore/src/internal/generated"
	"MerchStore/src/internal/handlers"
	"MerchStore/src/internal/logger"
	"MerchStore/src/internal/middleware"
	"MerchStore/src/internal/repository"
	"MerchStore/src/internal/storage"
	"MerchStore/src/internal/storage/postgres"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := cmd.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	logger.Init(cfg.Service.LogLevel)
	slog.Info("Starting server...")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)
	store, db := storage.NewDataStorage(dsn)
	defer postgres.CloseDB(db)

	redisAdr := fmt.Sprintf("%s:%s", cfg.Cache.Host, cfg.Cache.Port)
	redisRepo := storage.NewCacheStorage(redisAdr)
	resp := repository.NewRepository(store, redisRepo)
	server := handlers.NewServer(resp)

	// Создаем опции с middleware
	options := generated.StdHTTPServerOptions{
		Middlewares: []generated.MiddlewareFunc{
			middleware.AuthMiddleware(resp, cfg.Service.SecretKey),
			middleware.LoggingMiddleware,
		},
	}

	h := generated.HandlerWithOptions(server, options)

	s := &http.Server{
		Addr:    cfg.Service.ServerPort,
		Handler: h,
	}

	// Start server in a new goroutine
	go func() {
		slog.Info(fmt.Sprintf("Service started on port %s", cfg.Service.ServerPort))
		if err = s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Info(fmt.Sprintf("Error while starting server: %v", err))
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	// Shutdown server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = s.Shutdown(ctx); err != nil {
		slog.Info(fmt.Sprintf("Error during server shutdown: %v", err))
	}
	slog.Info(fmt.Sprintln("Server gracefully stopped"))
}
