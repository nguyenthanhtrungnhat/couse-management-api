package middleware

import (
	"course-management/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")

	// DEBUG
	println("AUTH HEADER:", authHeader)

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Authorization header is required",
		})
	}

	parts := strings.Fields(authHeader)

	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid authorization format",
		})
	}

	tokenString := parts[1]

	claims, err := utils.ParseToken(tokenString)

	if err != nil {
		println("JWT ERROR:", err.Error())

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	println("USER ID:", claims.UserID)
	println("ROLE:", claims.Role)

	c.Locals("user_id", claims.UserID)
	c.Locals("role", claims.Role)

	return c.Next()
}