package internal

import (
	"context"
	"service/internal/configs"
	"service/internal/handlers"
	handler "service/internal/handlers/news"
	"service/internal/repository"
	"service/internal/service"
	"service/pkg/db"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

type Server struct {
	config configs.Config
	reform *reform.DB
	ctx    context.Context
}

func RunServer(ctx context.Context) {
	log := logrus.New()
	cnf, err := configs.NewParsedConfig()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	db, err := db.InitReformDB(cnf.Database)
	if err != nil {
		log.Fatalf("failed to init reform db: %v", err)
	}

	repo := repository.NewNewsRepository(db, log, ctx)
	service := service.NewNewsService(repo, log)
	handler := handler.NewNewsHandler(service, log)

	app := fiber.New()

	//// Middleware
	//app.Use(recover.New())
	//app.Use(cors.New())
	//app.Use(handlers.LoggingMiddleware(log))

	// Роуты с авторизацией
	//api := app.Group("/api", handlers.AuthMiddleware(cfg.AuthToken, log))

	handlers.SetupRoutes(app, handler)

	if err = app.Listen(":" + cnf.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
