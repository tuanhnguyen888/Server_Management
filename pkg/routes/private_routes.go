package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tuanhnguyen888/Server_Management/app/controllers"
	"github.com/tuanhnguyen888/Server_Management/pkg/middleware"
)

func PrivateRoutes(app *fiber.App) {
	route := app.Group("api/v1")

	route.Post("/server", middleware.JWTProtected(), controllers.CreateServer)

}
