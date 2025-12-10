// internal/handlers/news.go
package handlers

import (
	"fmt"
	"service/internal/models"
	"service/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type NewsHandler struct {
	service service.INewsService
	log     *logrus.Logger
}

func NewNewsHandler(service service.INewsService, log *logrus.Logger) NewsHandler {
	return NewsHandler{
		service: service,
		log:     log,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *NewsHandler) EditNews(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		h.log.Warnf("Неверный ID в параметрах: %s", idParam)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: fmt.Sprintf("id is invalid: %s", idParam),
		})
	}

	var editForm models.NewsEditForm
	if err = c.BodyParser(&editForm); err != nil {
		h.log.Warnf("Error while parsing request body: %s", err)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Error while parsing request body",
		})
	}

	editForm.Normalize()
	err = editForm.Validate()
	if err != nil {
		h.log.Warn("Не передано ни одного поля для обновления")
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: fmt.Sprintf("validation failed: %s", err),
		})
	}

	if err = h.service.EditNews(id, editForm); err != nil {
		h.log.Errorf("Ошибка редактирования новости ID=%d: %v", id, err)

		statusCode := fiber.StatusInternalServerError
		//if err.Error() == "новость не найдена" {
		//	statusCode = fiber.StatusNotFound
		//} else if isValidationError(err) {
		//	statusCode = fiber.StatusBadRequest
		//}

		return c.Status(statusCode).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"Success": true,
	})
}

func (h *NewsHandler) ListNews(c *fiber.Ctx) error {
	var newsList []models.NewsWithCategories
	//добавить валидацию
	limit, err := strconv.ParseUint(c.Query("limit", "10"), 10, 64)
	offset, err := strconv.ParseUint(c.Query("offset", "0"), 10, 64)

	newsList, err = h.service.ListNews(limit, offset)
	if err != nil {
		h.log.Errorf("Ошибка получения списка новостей: %v", err)
		//почему то null а не пустой массив возвращается
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewsListsResponse{Success: false, News: newsList})
	}

	return c.Status(fiber.StatusOK).JSON(models.NewsListsResponse{Success: true, News: newsList})
}
