package rest

import (
	"context"
	"fmt"
	"log"
	"service/http/rest/handlers"
	"service/internal/configs"
	"service/internal/repository"
	"service/internal/service"
	"service/pkg/db"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/reform.v1"
)

type Server struct {
	config configs.Config
	reform *reform.DB
	ctx    context.Context
}

func NewServer(ctx context.Context) (*Server, error) {
	cnf, err := configs.NewParsedConfig()
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	reform, err := db.InitReformDB(cnf.Database)
	if err != nil {
		return nil, err
	}

	repo := repository.NewRepository(reform)
	newsService := service.NewService(repo, lruCache)

	// Создание Fiber приложения
	app := fiber.New(fiber.Config{
		ErrorHandler: handlers.Er(log),
	})

	//// Middleware
	//app.Use(recover.New())
	//app.Use(cors.New())
	//app.Use(handlers.LoggingMiddleware(log))

	// Роуты с авторизацией
	api := app.Group("/api", handlers.AuthMiddleware(cfg.AuthToken, log))

	api.Post("/edit/:id", newsHandler.EditNews)
	api.Get("/list", newsHandler.ListNews)

	// Запуск сервера
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Infof("Запуск сервера на %s", addr)

	if err := app.Listen(addr); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func (s *Server) Run() error {
	go s.consumer.Start(s.ctx)
	err := s.gin.Run(s.config.Port)

	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}
