package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v8"
	"github.com/idos-network/idos-extensions/extension"
)

func main() {
	cfg := &extension.Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	ext, err := extension.NewFractalExt(cfg.ChainURLs())
	if err != nil {
		log.Fatalf("failed to create extension: %v", err)
	}

	logger := log.New(os.Stdout, "idos: ", log.LstdFlags)

	svr, err := ext.BuildServer(logger)
	if err != nil {
		log.Fatalf("failed to construct server: %v", err)
	}

	lis, err := net.Listen("tcp", cfg.ListenAddr())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		logger.Printf("listening on %s\n", cfg.ListenAddr())
		if err := svr.Serve(lis); err != nil {
			logger.Printf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("shutting down")

	err = svr.GracefulStop()
	if err != nil {
		logger.Printf("failed to shutdown: %v", err)
		os.Exit(1)
	}

	logger.Println("shutdown complete")
}
