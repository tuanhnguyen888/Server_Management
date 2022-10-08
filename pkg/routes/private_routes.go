package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tuanhnguyen888/Server_Management/app/controllers"
	"github.com/tuanhnguyen888/Server_Management/pkg/middleware"
)

func PrivateRoutes(app *fiber.App) {
	route := app.Group("api/v1")

	app.Post("/login", controllers.Login)

	route.Post("/server", middleware.AuthRequired(), controllers.CreateServer)
	route.Post("/server/:id", middleware.AuthRequired(), controllers.UpdateServer)
	route.Delete("/server/:id", middleware.AuthRequired(), controllers.DeleteServer)

}
