package middleware

import (
	//"fmt"
	//"github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/Harry-027/go-notify/api-server/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	//jwtware "github.com/gofiber/jwt/v2"
)

// middleware to verify the Jwt ...
func Protected() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenString := ctx.Get("Authorization")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if err != nil {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "Error", "message": "Jwt is malformed or Invalid"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx.Locals("user_id", claims["user_id"])
			ctx.Locals("user_role", claims["user_role"])
			ctx.Locals("uuid", claims["uuid"])
			uuid := claims["uuid"]
			currentSession := repository.IfUsersCurrentSession(uuid)

			if !currentSession {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "Error", "message": "Invalid Session"})
			}
			return ctx.Next()
		}

		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "Error", "message": "Invalid jwt token"})
	}
}
