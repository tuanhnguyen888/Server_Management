package controllers

import (
	"errors"
	"log"
	"net"
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
1. (DONE) The DB instance should be initialized when init app/service and be reused for API.
2. Golang tends to deal every thing with errors instead of raising errors. Therefor every logic function should return
an error if existing then we will define an appropriate behavior for each error.
3. Controller needs to handle errors and response status code correctly.
4. Read the comments in GetServers then correct all the mistakes for all APIs.
*/

// GetServers func gets all exists server or HTTP 400.
// @Description Get all exists server.
// @Summary get all exists Server
// @Tags Server Public
// @Accept json
// @Produce json
// @Param page query string false "search by page"
// @Param sort query string false "field names to sort"
// @Param kind query string false "sort type"
// @Failure 400 {string} status "400"
// @Success 200 {object} models.Server
// @Router /api/v1/servers [get]
func GetServers(r *models.Repository) func(c *fiber.Ctx) error {
	// connect db
	// TODO (DONE) : init a db instance in main.go then use it for all APIs instead of creating as many as instances for each user request
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

		// TODO (DONE): refactor exportToExcel to return an error then handle it
		//err := exportToExcel(*servers)
		//if err != nil {
		//	c.Status(http.StatusBadRequest).JSON(
		//		&fiber.Map{
		//			"message": "could not export excel",
		//		})
		//	return err
		//}
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

// SearchServers func Search server by name or HTTP 400.
// @Description Search server by name.
// @Summary Search server by name
// @Tags Server Public
// @Accept json
// @Produce json
// @Param page query string false "search by page"
// @Param name query string true "name to search"
// @Failure 400 {string} status "400"
// @Success 200 {object} models.Server
// @Router /api/v1/search [get]
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

		v := c.Query("value")

		r.DB.Where("name LIKE ?", "%"+v+"%").Order("name").Offset(offset).Limit(perPage).Find(&servers)

		if len(*servers) == 0 {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{
					"message": "Can't find the right servers",
				})
			return nil
		}

		c.Status(http.StatusOK).JSON(
			&fiber.Map{
				"message": "server fetched",
				"data":    servers,
			})
		return nil
	}
}

// GetServerByID func gets Server by given ID or 400 error.
// @Description Get server by given ID.
// @Summary get Server by given ID
// @Tags Server Public
// @Accept json
// @Produce json
// @Param id path string true "Server ID"
// @Success 200 {object} models.Server
// @Failure 400 {string} status "400"
// @Router /api/v1/server/{id} [get]
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

// Login method for login and create a new access token.
// @Description Create a new access token.
// @Summary create a new access token
// @Tags Token
// @Accept json
// @Produce json
// @Success 200 {string} status "ok"
// @Failure 400 {string} status "400"
// @Router /api/v1/login [post]
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
		return err
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
		return err
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

// CreateServer func for creates a new Server.
// @Description Create a new Server.
// @Summary create a new Server
// @Tags Server Private
// @Accept json
// @Produce json
// @Param name body string true "name Server "
// @Param ipv4 body string true "ipv4 Server "
// @Param status body string false "status Server"
// @Success 200 {object} models.Server
// @Failure 400 {string} status "400"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /api/v1/server [post]
func CreateServer(r *models.Repository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		// ---------------
		server := Server{}
		err := c.BodyParser(&server)
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "not parser body data"})
			return err
		}

		err = CheckName(r, server.Name)
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "Server is duplicated"})
			return err
		}

		server.ID, err = uuid.NewV1()
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "not create uuid"})
			return err
		}

		err = checkIPAddress(*server.Ipv4)
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "Address Invalid"})
			return err
		}
		// pinggg net
		// _, err = exec.Command("ping", *server.Ipv4).Output()
		// if err != nil {
		// 	server.Status = false
		// } else {
		// 	server.Status = true
		// }

		server.CreatedAt = time.Now().UnixMilli()
		server.UpdatedAt = time.Now().UnixMilli()
		// -------
		err = r.DB.Create(&server).Error
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
}

// UpdateServer func for updates Server by given ID.
// @Description Update Server.
// @Summary update Server
// @Tags Server Private
// @Accept json
// @Produce json
// @Param id path string true "Server ID"
// @Param name body string true "name Server "
// @Param status body string true "status Server"
// @Param ipv4 body string true "ipv4 Server "
// @Failure 400 {string} status "400"
// @Success 200 {string} status "ok"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /api/v1/server/{id} [post]
func UpdateServer(r *models.Repository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		// ----
		id := c.Params("id")

		u3, err := uuid.FromString(id)
		if err != nil {
			log.Fatalf("failed to parse UUID %q: %v", id, err)
		}

		server := Server{}
		err = c.BodyParser(&server)
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "not request"})
			return err
		}

		err = checkIPAddress(*server.Ipv4)
		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "Address Invalid"})
			return err
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

}

// DeleteServer func for deletes Server by given ID.
// @Description Delete Server by given ID.
// @Summary delete Server by given ID
// @Tags Server Private
// @Accept json
// @Produce json
// @Param id path string true "Server ID"
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Router /api/v1/server/{id} [delete]
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
			return nil
		}

		c.Status(http.StatusOK).JSON(
			&fiber.Map{
				"message": "server has been deleted",
			})
		return nil
	}

}

// validate name
func CheckName(r *models.Repository, name *string) error {
	allServers := []Server{}
	r.DB.Find(&allServers)
	sameName := false
	for _, s := range allServers {
		if *name == *s.Name {
			sameName = true
		}
	}
	if sameName {
		return errors.New("server name is duplicated")
	}
	return nil
}

// validate ip
func checkIPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		ErrorLogger.Printf("IP Address: %s - Invalid\n", ip)
		return errors.New("Address Invalid")
	} else {
		ErrorLogger.Printf("IP Address: %s - Valid\n", ip)
		return nil
	}
}
