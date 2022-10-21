package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tuanhnguyen888/Server_Management/app/controllers"
	"github.com/tuanhnguyen888/Server_Management/app/models"
	"github.com/tuanhnguyen888/Server_Management/flatform"
	"github.com/tuanhnguyen888/Server_Management/pkg/routes"
)

func main() {

	db, err := flatform.NewInit()
	if err != nil {
		fmt.Println("can not connect")
		panic(err)
	}

	r := models.Repository{
		DB: db,
	}

	err = models.MigrateServer(db)
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	routes.PublicRoutes(app, &r)
	routes.PrivateRoutes(app, &r)

	go controllers.Cron(r)

	if err := app.Listen(":5000"); err != nil {
		log.Printf(" Server is not running! Reason: %v", err)
		// should use log.Fatal instead to terminate the program when cannot start service.
		// e.g: log.Fatal(app.Listen(":5000"))
		log.Fatal(err)
	}
	//

}
