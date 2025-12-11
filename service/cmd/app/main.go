package main

import (
	"context"
	//_ "orderService/docs"
	"os"
	"os/signal"
	"service/internal"
	"syscall"

	"github.com/sirupsen/logrus"
)

// @title Order Service
// @version 1.0
// @description Order service API Docs

// @host 	localhost:8080
// @BasePath /api
func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	signals := []os.Signal{
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT,
	}

	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()

	go func() {
		internal.RunServer(ctx, log)
	}()

	<-ctx.Done()
	log.Println("Завершение работы...")

	// Выполняем graceful shutdown
	if err := internal.Stop(); err != nil {
		log.Printf("Ошибка при shutdown: %v", err)
	}

	log.Println("Сервер остановлен.")

}
