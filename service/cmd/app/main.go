package main

import (
	"context"
	//_ "orderService/docs"
	"os"
	"os/signal"
	"service/internal"
	"syscall"
)

// @title Order Service
// @version 1.0
// @description Order service API Docs

// @host 	localhost:8080
// @BasePath /api
func main() {
	signals := []os.Signal{
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT,
	}

	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()

	internal.RunServer(ctx)
}
