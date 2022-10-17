package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tuanhnguyen888/Server_Management/app/controllers"
	"github.com/tuanhnguyen888/Server_Management/app/models"
	"github.com/tuanhnguyen888/Server_Management/pkg/middleware"
)

func PrivateRoutes(app *fiber.App, r *models.Repository) {
	route := app.Group("api/v1")

	app.Post("/login", controllers.Login)

	route.Post("/server", middleware.AuthRequired(), controllers.CreateServer(r))
	route.Post("/importServerFromExcel", middleware.AuthRequired(), controllers.ImportExcel)
	route.Post("/server/:id", middleware.AuthRequired(), controllers.UpdateServer(r))
	route.Delete("/server/:id", middleware.AuthRequired(), controllers.DeleteServer(r))

}
