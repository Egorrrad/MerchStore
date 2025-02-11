///go:generate oapi-codegen --config=../../api/cfg.yaml ../../api/schema.yaml

package main

import (
	"MerchStore/src/cmd"
	"MerchStore/src/internal/datastorage"
	"MerchStore/src/internal/datastorage/postgres"
	"MerchStore/src/internal/handlers"
	"MerchStore/src/internal/repository"
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

	config, err := cmd.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)

	store, db := datastorage.NewDataStorage(dsn)
	defer postgres.CloseDB(db)

	resp := repository.NewRepository(store)
	server := handlers.NewServer(resp)
	mux := http.NewServeMux()
	h := handlers.HandlerFromMux(server, mux)

	s := &http.Server{
		Addr:    config.ServerPort,
		Handler: h,
	}

	// Start server in a new goroutine
	go func() {
		log.Printf("Service started on port %s", config.ServerPort)
		if err = s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error while starting server: %v", err)
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
		log.Printf("Error during server shutdown: %v", err)
	}
	log.Println("Server gracefully stopped")
}
