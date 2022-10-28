package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tuanhnguyen888/Server_Management/app/controllers"
	"github.com/tuanhnguyen888/Server_Management/app/models"
	"github.com/tuanhnguyen888/Server_Management/flatform"
	"github.com/tuanhnguyen888/Server_Management/pkg/routes"

	// load API Docs files (Swagger)
	_ "github.com/tuanhnguyen888/Server_Management/docs"
)

// @title Server Management
// @version 1.0
// @description Công ty VCS hiện tại có gồm khoảng 10000 server. App xây dựng 1 hệ thống quản lý trạng thái On/Off của danh sách server này.
// @host localhost:5000
// @BasePath

func main() {

	db, err := flatform.NewInit()
	if err != nil {
		controllers.ErrorLogger.Println("can not connect DB")
		log.Fatal(err)
	}

	r := models.Repository{
		DB: db,
	}

	err = models.MigrateServer(db)
	if err != nil {
		log.Fatal(err)
	}

	esclient, err := flatform.GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}

	redisClient, err := flatform.NewInitResdis()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	routes.SwaggerRoute(app)
	routes.PublicRoutes(app, &r)
	routes.PrivateRoutes(app, &r, redisClient)
	routes.NotFoundRoute(app)

	go controllers.Cron(r, redisClient, esclient)

	if err := app.Listen(":5000"); err != nil {
		controllers.ErrorLogger.Printf(" Server is not running! Reason: %v", err)
		// should use log.Fatal instead to terminate the program when cannot start service.
		// e.g: log.Fatal(app.Listen(":5000"))
		log.Fatal(err)
	}
	//

}
