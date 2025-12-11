package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// ErrorResponse стандартный формат ответа с ошибкой
type ErrorResponse struct {
	Success bool
	Error   string
}

// ErrorHandler глобальный обработчик ошибок
func ErrorHandler(log *logrus.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		log.WithFields(logrus.Fields{
			"method": c.Method(),
			"path":   c.Path(),
			"error":  err.Error(),
			"status": code,
		}).Error("Ошибка обработки запроса")

		return c.Status(code).JSON(ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
	}
}
