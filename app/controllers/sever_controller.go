package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/tuanhnguyen888/Server_Management/app/models"
	"github.com/tuanhnguyen888/Server_Management/flatform"
)

type Server struct {
	ID        uuid.UUID `json:"id"`
	Name      *string   ` json:"name" `
	Status    *string   `json:"status" `
	Ipv4      *string   `json:"ipvd4" `
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}

func GetServers(c *fiber.Ctx) error {
	// connect db
	db, err := flatform.NewInit()
	if err != nil {
		fmt.Println("can not connect")
		panic(err)
	}

	servers := &[]models.Server{}
	db.Find(&servers)

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "server fetched",
			"data":    servers,
		})
	return nil
}

func GetServerById(c *fiber.Ctx) error {
	id := c.Params("id")
	// connect db
	db, err := flatform.NewInit()
	if err != nil {
		fmt.Println("can not connect")
		panic(err)
	}

	server := &models.Server{}
	db.Where("id = ?", id).Find(&server)

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "server fetched",
			"data":    server,
		})
	return nil
}
func CreateServer(c *fiber.Ctx) error {
	server := Server{}
	err := c.BodyParser(&server)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "not request"})
		return err
	}

	server.ID, err = uuid.NewV1()
	if err != nil {
		panic(err)
	}
	server.CreatedAt = time.Now()
	server.UpdatedAt = time.Now()

	// connect db
	db, err := flatform.NewInit()
	if err != nil {
		fmt.Println("can not connect")
		panic(err)
	}

	err = db.Create(&server).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create Server"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message":      "server has been added",
			"server added": server,
		})
	return nil
}
