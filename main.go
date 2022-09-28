package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
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

	err = models.MigrateServer(db)
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	//
	routes.PrivateRoutes(app)
	routes.PublicRoutes(app)

	//

	//

	if err := app.Listen(":5000"); err != nil {
		log.Printf(" Server is not running! Reason: %v", err)
	}

}
