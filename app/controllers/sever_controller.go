package controllers

import (
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tuanhnguyen888/Server_Management/app/models"
)

type Server struct {
	ID        uuid.UUID `json:"id"`
	Name      *string   ` json:"name" `
	Status    bool      `json:"status" `
	Ipv4      *string   `json:"ipv4" `
	CreatedAt int64     `json:"created_at"`
	UpdatedAt int64     `json:"update_at"`
}

/* Summarize:
1. The DB instance should be initialized when init app/service and be reused for API.
2. Golang tends to deal every thing with errors instead of raising errors. Therefor every logic function should return
an error if existing then we will define an appropriate behavior for each error.
3. Controller needs to handle errors and response status code correctly.
4. Read the comments in GetServers then correct all the mistakes for all APIs.
*/

func GetServers(r *models.Repository) func(c *fiber.Ctx) error {
	// connect db
	// TODO: init a db instance in main.go then use it for all APIs instead of creating as many as instances for each user request
	return func(c *fiber.Ctx) error {
		servers := &[]models.Server{}

		// ------ pagination ----------
		page, _ := strconv.Atoi(c.Query("page"))
		if page == 0 {
			page = 1
		}
		perPage := 10
		offset := (page - 1) * perPage

		//  ------ sort -------
		sort := c.Query("sort")
		k := c.Query("kind")
		if sort != "" {
			r.DB.Order(sort + " " + k).Offset(offset).Limit(perPage).Find(&servers)
		} else {
			r.DB.Offset(offset).Limit(perPage).Find(&servers)
		}

		// TODO: refactor exportToExcel to return an error then handle it
		err := exportToExcel(*servers)
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{
					"message": "could not export excel",
				})
		}
		// TODO: Handle error for JSON method
		c.Status(http.StatusOK).JSON(
			&fiber.Map{
				"message": "server fetched",
				"data":    servers,
			})

		return nil
		// TODO: if an error happens, handle it with correct response status codes.
		// e.g: an API should return 200 for successful, 400 if request params are invalid, 500 if an error like DB connection error occurs
		// For more details, visit https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
	}
}

func Search(r *models.Repository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// ------ pagination ----------
		page, _ := strconv.Atoi(c.Query("page"))
		if page == 0 {
			page = 1
		}
		perPage := 10
		offset := (page - 1) * perPage

		// ----search

		servers := &[]models.Server{}
		n := c.Query("name")

		if n != "" {
			r.DB.Where("name LIKE ?  ", "%"+n+"%").Offset(offset).Limit(perPage).Find((&servers))
		} else {
			r.DB.Offset(offset).Limit(perPage).Find((&servers))
		}

		if len(*servers) == 0 {
			c.Status(http.StatusOK).JSON(
				&fiber.Map{
					"message": "Can't find the right servers",
				})
		}

		c.Status(http.StatusOK).JSON(
			&fiber.Map{
				"message": "server fetched",
				"data":    servers,
			})
		return nil
	}
}

func GetServerById(r *models.Repository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		server := &models.Server{}
		r.DB.Where("id = ?", id).Find(&server)

		c.Status(http.StatusOK).JSON(
			&fiber.Map{
				"message": "server fetched",
				"data":    server,
			})
		return nil
	}
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
	}

	if body.Email != "tuanh@gmail.com" || body.Password != "khong123" {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Bad Credentials",
		})
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"admin": true,
		"exp":   time.Now().Add(time.Minute * time.Duration(50)).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.

	s, err := token.SignedString([]byte("secret"))
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot singed tolen",
		})
	}

	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": s,
		"admin": claims["admin"],
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

func CreateServer(r *models.Repository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		timeNow := time.Now().Unix()

		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		expires := int64(claims["exp"].(float64))

		if timeNow > expires {
			// Return status 401 and unauthorized error message.
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "unauthorized, check expiration time of your token",
			})
		}
		// ---------------
		server := Server{}
		err := c.BodyParser(&server)
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "not request"})
		}

		server.ID, err = uuid.NewV1()
		if err != nil {
			panic(err)
		}
		server.CreatedAt = time.Now().UnixMilli()
		server.UpdatedAt = time.Now().UnixMilli()

		// pinggg net
		_, err = exec.Command("ping", *server.Ipv4).Output()
		if err != nil {
			server.Status = true
		} else {
			server.Status = false
		}

		err = r.DB.Create(&server).Error
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "could not create Server"})
		}

		c.Status(http.StatusOK).JSON(
			&fiber.Map{
				"message":      "server has been added",
				"server added": server,
			})
		return nil
	}
}

func UpdateServer(r *models.Repository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		timeNow := time.Now().Unix()

		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		expires := int64(claims["exp"].(float64))

		if timeNow > expires {
			// Return status 401 and unauthorized error message.
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
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
		}

		_, err = exec.Command("ping", *server.Ipv4).Output()
		if err != nil {
			server.Status = false
		} else {
			server.Status = true
		}
		server.UpdatedAt = time.Now().UnixMilli()

		err = r.DB.Where("id = ? ", u3).Updates(&server).Error
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "could not update Server"})
		}

		server.ID = u3
		c.Status(http.StatusOK).JSON(
			&fiber.Map{
				"message":       "server has been update",
				"server update": server,
			})
		return nil
	}

}

func DeleteServer(r *models.Repository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		u3, err := uuid.FromString(id)
		if err != nil {
			log.Fatalf("failed to parse UUID %q: %v", id, err)
		}

		server := Server{}

		err = r.DB.Delete(&server, u3).Error
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "could not delete Server"})
		}

		c.Status(http.StatusOK).JSON(
			&fiber.Map{
				"message": "server has been deleted",
			})
		return nil
	}

}
