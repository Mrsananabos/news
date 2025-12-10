package order

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"orderService/internal/service"
)

type Handler struct {
	service service.IOrderService
}

func NewHandler(service service.IOrderService) Handler {
	return Handler{
		service: service,
	}
}

// FindByIdTags 		godoc
// @Summary				Get Order by id
// @Param				id path string true "Get order by id"
// @Description			Return order by id
// @Produce				application/json
// @Tags				order
// @Success				200 {object} models.OrderView
// @Router				/order/{id} [get]
func (h Handler) GetOrderById(c *gin.Context) {
	uidStr := c.Param("uid")
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uid is not UUID format"})
		log.Printf("uid %s is not UUID format", uidStr)
		return
	}

	order, err := h.service.GetById(uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}
