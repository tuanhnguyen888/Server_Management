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

	// -------- CRON -----------

	// var ctx *fiber.Ctx

	// c := cron.New()
	// c.AddFunc("@every 0h0m5s", func() { controllers.UpdateServer(ctx) })
	// c.Start()
	// controllers.Cron()
	//

	if err := app.Listen(":5000"); err != nil {
		log.Printf(" Server is not running! Reason: %v", err)
	}

}
