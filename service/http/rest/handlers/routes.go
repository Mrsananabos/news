package handlers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"orderService/http/rest/handlers/order"
	"orderService/http/rest/middleware"
	"orderService/internal/service"
)

func Register(gin *gin.Engine, orderService service.IOrderService) {
	orderHandler := order.NewHandler(orderService)

	gin.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	gin.GET("/order/:uid", middleware.RequestIdMiddleware("getOrderById"), middleware.SetCors(), orderHandler.GetOrderById)

}
