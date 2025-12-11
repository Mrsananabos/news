package main

import (
	"context"
	"time"

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

	// Контекст для запуска сервера
	ctx := context.Background()

	// Создаем канал для graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)

	// Запускаем сервер
	server, err := internal.NewServer(ctx, log)
	if err != nil {
		log.Fatal(err)
	}

	// Запускаем в горутине
	serverErrors := make(chan error, 1)
	go func() {
		log.Info("Service starting...")
		serverErrors <- server.Start()
	}()

	// Ожидаем сигнал завершения или ошибку сервера
	select {
	case err = <-serverErrors:
		log.Fatal(err)
	case sig := <-shutdown:
		log.Infof("Getting shutdown signal : %v", sig)

		// Создаем контекст с таймаутом для graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Выполняем graceful shutdown
		if err := server.Stop(ctx); err != nil {
			log.Errorf("Ошибка при graceful shutdown: %v", err)
			// Принудительное завершение
			log.Fatal("Принудительное завершение работы")
		}

		log.Info("Сервер успешно остановлен")
	}
}
