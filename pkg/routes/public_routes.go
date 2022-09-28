package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tuanhnguyen888/Server_Management/app/controllers"
)

func PublicRoutes(app *fiber.App) {
	route := app.Group("api/v1")

	route.Get("/servers", controllers.GetServers)
	route.Get("/server/:id", controllers.GetServerById)
}
