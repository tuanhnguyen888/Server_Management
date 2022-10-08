package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tuanhnguyen888/Server_Management/app/models"
	"github.com/tuanhnguyen888/Server_Management/flatform"
)

type Server struct {
	ID        uuid.UUID `json:"id"`
	Name      *string   ` json:"name" `
	Status    *string   `json:"status" `
	Ipv4      *string   `json:"ipv4" `
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

func Login(ctx *fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body request
	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return nil
	}

	if body.Email != "tuanh@gmail.com" || body.Password != "khong123" {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Bad Credentials",
		})
		return nil
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.

	s, err := token.SignedString([]byte("secret"))
	if err != nil {
		ctx.SendStatus(fiber.StatusInternalServerError)
		return nil
	}

	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": s,
		"user": struct {
			Id    int    `json:"id"`
			Email string `json:"email"`
		}{
			Id:    1,
			Email: "tuanh@gmail.com",
		},
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

func UpdateServer(c *fiber.Ctx) error {
	id := c.Params("id")

	u3, err := uuid.FromString(id)
	if err != nil {
		log.Fatalf("failed to parse UUID %q: %v", id, err)
	}

	server := Server{}
	err = c.BodyParser(&server)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "not request"})
		return err
	}

	server.UpdatedAt = time.Now()

	// connect db
	db, err := flatform.NewInit()
	if err != nil {
		fmt.Println("can not connect")
		panic(err)
	}

	err = db.Where("id = ? ", u3).Updates(&server).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not update Server"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message":       "server has been update",
			"server update": server,
		})
	return nil
}

func DeleteServer(c *fiber.Ctx) error {
	id := c.Params("id")

	u3, err := uuid.FromString(id)
	if err != nil {
		log.Fatalf("failed to parse UUID %q: %v", id, err)
	}

	server := Server{}

	// connect db
	db, err := flatform.NewInit()
	if err != nil {
		fmt.Println("can not connect")
		panic(err)
	}

	err = db.Delete(&server, u3).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not delete Server"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "server has been deleted",
		})
	return nil
}
