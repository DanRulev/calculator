package app

import (
	"calculator-go/internal/handler"
	"calculator-go/internal/service"
	"calculator-go/pkg/config"
	"calculator-go/pkg/server"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	timeout = 5 * time.Second
)

func Start(path, file string) {
	log.Println("Starting the application...")
	ctx := context.Background()

	log.Println("Initializing configuration...")
	cfg, err := config.Init(path, file)
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}
	log.Printf("Configuration initialized: %v", cfg)

	log.Println("Initializing service...")
	service := service.NewCalcService()

	log.Println("Initializing handler...")
	hdl := handler.New(service)

	log.Println("Initializing server...")
	srv := server.New(cfg.ServerCfg, hdl.InitRouter())

	log.Println("Starting server...")
	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("Server closed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server shutdown complete.")
}
