package routes

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/tuanhnguyen888/Server_Management/app/controllers"
	"github.com/tuanhnguyen888/Server_Management/app/models"
	"github.com/tuanhnguyen888/Server_Management/pkg/middleware"
)

func PrivateRoutes(app *fiber.App, r *models.Repository, redisClient *redis.Client) {
	route := app.Group("api/v1")

	route.Post("/login", controllers.Login)

	route.Post("/server", middleware.AuthRequired(), controllers.CreateServer(r))
	route.Post("/importServer", middleware.AuthRequired(), controllers.ImportExcel(r))
	route.Post("/server/:id", middleware.AuthRequired(), controllers.UpdateServer(r))
	route.Delete("/server/:id", middleware.AuthRequired(), controllers.DeleteServer(r))
	route.Post("/sendEmail", middleware.AuthRequired(), controllers.CustomSendEmail(redisClient))

}
