package middleware

import (
	//"fmt"
	//"github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/Harry-027/go-notify/api-server/config"
	"github.com/Harry-027/go-notify/api-server/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"strconv"
	//jwtware "github.com/gofiber/jwt/v2"
)

// middleware to verify the Jwt ...
func Protected() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenString := ctx.Get("Authorization")

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(config.GetConfig(config.JWT_SECRET)), nil
		})

		claims := token.Claims.(jwt.MapClaims)
		if claims.Valid() != nil {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "Error", "message": "Invalid jwt token"})
		}

		userId, err := strconv.Atoi(fmt.Sprintf("%v", claims["user_id"]))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "Error", "message": "Internal server error"})
		}

		foundUser, err := repository.GetUserById(userId)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "Error", "message": "Internal server error"})
		}

		if foundUser.Password == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "Error", "message": "User is not authorized"})
		}

		ctx.Locals("user_id", claims["user_id"])
		ctx.Locals("user_role", claims["user_role"])
		ctx.Locals("uuid", claims["uuid"])

		return ctx.Next()
	}
}
