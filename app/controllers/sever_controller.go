package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"update_at"`
}

func GetServers(c *fiber.Ctx) error {
	// connect db
	db, err := flatform.NewInit()
	if err != nil {
		fmt.Println("can not connect")
		panic(err)
	}

	servers := &[]models.Server{}

	// ------ pagination ----------
	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}
	perPage := 9
	offset := (page - 1) * perPage

	//  ------ sort -------
	sort := c.Query("sort")
	k := c.Query("kind")
	if sort != "" {
		db.Order(sort + " " + k).Offset(offset).Limit(perPage).Find(&servers)
	} else {
		db.Offset(offset).Limit(perPage).Find(&servers)
	}

	exportToExcel(*servers)

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "server fetched",
			"data":    servers,
		})
	return nil
}

func Search(c *fiber.Ctx) error {
	// connect db
	db, err := flatform.NewInit()
	if err != nil {
		fmt.Println("can not connect")
		panic(err)
	}

	// ------ pagination ----------
	page, _ := strconv.Atoi(c.Query("page"))
	if page == 0 {
		page = 1
	}
	perPage := 9
	offset := (page - 1) * perPage

	// ----search

	servers := &[]models.Server{}
	n := c.Query("name")
	s := c.Query("status")

	if s != "" && n != "" {
		db.Where("name LIKE ? and status LIKE ? ", "%"+n+"%", s+"%").Offset(offset).Limit(perPage).Find((&servers))
	} else {
		db.Offset(offset).Limit(perPage).Find((&servers))
	}

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
		"exp":   time.Now().Add(time.Minute * time.Duration(20)).Unix(),
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
	timeNow := time.Now().Unix()

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	expires := int64(claims["exp"].(float64))

	if timeNow > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}
	// ---------------
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
	server.CreatedAt = time.Now().UnixMilli()
	server.UpdatedAt = time.Now().UnixMilli()

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
	timeNow := time.Now().Unix()

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	expires := int64(claims["exp"].(float64))

	if timeNow > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}
	// ----
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

	server.UpdatedAt = time.Now().UnixMilli()

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
	server.ID = u3
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
