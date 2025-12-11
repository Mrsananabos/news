package handlers

import (
	handler "service/internal/handlers/news"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes настраивает все роуты приложения
func SetupRoutes(app *fiber.App, newsHandler handler.NewsHandler) {
	// API группа с авторизацией
	//api := app.Group("/api", AuthMiddleware(authToken, log))
	api := app.Group("/")

	// Роуты для работы с новостями
	api.Post("edit/:id", newsHandler.EditNews)
	api.Get("list", newsHandler.ListNews)
	api.Post("create", newsHandler.CreateNews)
}

//import (
//	"github.com/gin-gonic/gin"
//	swaggerFiles "github.com/swaggo/files"
//	ginSwagger "github.com/swaggo/gin-swagger"
//	"orderService/http/rest/handlers/order"
//	"orderService/http/rest/middleware"
//	"orderService/internal/service"
//)
//
//func Register(gin *gin.Engine, orderService service.IOrderService) {
//	orderHandler := order.NewHandler(orderService)
//
//	gin.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
//
//	gin.GET("/order/:uid", middleware.RequestIdMiddleware("getOrderById"), middleware.SetCors(), orderHandler.GetOrderById)
//
//}
