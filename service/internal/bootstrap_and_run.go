package internal

import (
	"context"
	"database/sql"
	"errors"
	"service/internal/configs"
	"service/internal/handlers"
	handler "service/internal/handlers/news"
	"service/internal/repository"
	"service/internal/service"
	"service/pkg/db"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var server Server

type Server struct {
	app *fiber.App
	db  *sql.DB
}

func RunServer(ctx context.Context, log *logrus.Logger) {
	cnf, err := configs.NewParsedConfig()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	db, reform, err := db.InitReformDB(cnf.Database)
	if err != nil {
		log.Fatalf("failed to init reform db: %v", err)
	}
	repo := repository.NewNewsRepository(reform, log, ctx)
	service := service.NewNewsService(repo, log)
	handler := handler.NewNewsHandler(service, log)

	app := fiber.New()

	//// Middleware
	//app.Use(recover.New())
	//app.Use(cors.New())
	//app.Use(handlers.LoggingMiddleware(log))

	// Роуты с авторизацией
	//api := app.Group("/api", handlers.AuthMiddleware(cfg.AuthToken, log))

	server = Server{app: app, db: db}
	handlers.SetupRoutes(app, handler)

	if err = app.Listen(":" + cnf.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func Stop() error {
	var errs []error

	if err := server.app.Shutdown(); err != nil {
		errs = append(errs, err)
	}

	if err := server.db.Close(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
