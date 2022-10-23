package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthRequired() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		// tokenString, err := c.Cookie("token")
		authentication := ctx.Get("Authorization")
		if authentication == "" {
			ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
			return nil
		}

		var mySigningKey = []byte("secret")

		token, err := jwt.Parse(authentication[7:], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Error parsing token",
			})
			return nil
		}

		timeNow := time.Now().Unix()
		claims, ok := token.Claims.(jwt.MapClaims)
		if !(ok && token.Valid) {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid token",
			})
			return nil
		}

		expires := int64(claims["exp"].(float64))
		if timeNow > expires {
			// Return status 401 and unauthorized error message.
			ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "unauthorized, jwt is expired",
			})
			return nil
		}
		ctx.Next()
		return nil

	}
}
