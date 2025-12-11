package handlers

import (
	"errors"
	"service/internal/apperrors"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Success bool
	Error   string `validate:"omitempty"`
}

// CustomErrorHandler - простой и понятный обработчик
func ErrorHandler(log *logrus.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Дефолтные значения
		code := fiber.StatusInternalServerError
		message := "Internal server error"

		// Проверяем тип ошибки
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			// Наша кастомная ошибка
			code = appErr.StatusCode
			message = appErr.Message

			// Логируем в зависимости от типа
			if code >= 500 {
				// Серверные ошибки - ERROR уровень
				log.WithFields(logrus.Fields{
					"method": c.Method(),
					"path":   c.Path(),
					"error":  err.Error(),
				}).Error("Internal server error")
			} else {
				// Клиентские ошибки - WARN уровень
				log.WithFields(logrus.Fields{
					"method": c.Method(),
					"path":   c.Path(),
					"error":  message,
				}).Warn("Client error")
			}
		} else {
			// Fiber ошибка или неожиданная ошибка
			var fiberErr *fiber.Error
			if errors.As(err, &fiberErr) {
				code = fiberErr.Code
				message = fiberErr.Message
			}

			// Логируем неожиданные ошибки
			log.WithFields(logrus.Fields{
				"method": c.Method(),
				"path":   c.Path(),
				"error":  err.Error(),
			}).Error("Unexpected error")
		}

		// Возвращаем JSON в едином формате
		return c.Status(code).JSON(ErrorResponse{
			Success: false,
			Error:   message,
		})
	}
}
