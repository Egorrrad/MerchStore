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
	cfg.Service.ServerPort = ":" + cfg.Service.ServerPort
	logger.Init(cfg.Service.LogLevel, "text", "stdout")
	logger.Logger.Info("Starting server...", "port", cfg.Service.ServerPort)

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	store, db := storage.NewDataStorage(dsn)
	defer func() {
		if err := postgres.CloseDB(db); err != nil {
			logger.Logger.Error("Database closure error", "error", err)
		} else {
			logger.Logger.Info("Database connection closed successfully")
		}
	}()

	redisAdr := fmt.Sprintf("%s:%s", cfg.Cache.Host, cfg.Cache.Port)
	redisRepo := storage.NewCacheStorage(redisAdr)
	logger.Logger.Info("Connected to Redis", "address", redisAdr)

	resp := repository.NewRepository(store, redisRepo)

	options := generated.StdHTTPServerOptions{
		Middlewares: []generated.MiddlewareFunc{
			middleware.AuthMiddleware(resp, cfg.Service.SecretKey),
			middleware.LoggingMiddleware,
		},
	}
	server := handlers.NewServer(resp)
	h := generated.HandlerWithOptions(server, options)
	s := &http.Server{
		Addr:    cfg.Service.ServerPort,
		Handler: h,
	}

	go func() {
		logger.Logger.Info("Service started", "port", cfg.Service.ServerPort)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Error("Error while starting server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-quit
	logger.Logger.Warn("Received shutdown signal", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		logger.Logger.Error("Error during server shutdown", "error", err)
	} else {
		logger.Logger.Info("Server gracefully stopped")
	}
}
