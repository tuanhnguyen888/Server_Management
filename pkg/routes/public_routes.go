package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tuanhnguyen888/Server_Management/app/controllers"
	"github.com/tuanhnguyen888/Server_Management/app/models"
)

func PublicRoutes(app *fiber.App, r *models.Repository) {
	route := app.Group("api/v1")

	route.Get("/servers", controllers.GetServers(r))
	route.Get("/server/:id", controllers.GetServerById(r))
	route.Get("/search", controllers.Search(r))
}
