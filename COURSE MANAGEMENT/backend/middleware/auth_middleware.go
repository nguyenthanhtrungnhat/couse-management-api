package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	UserIDKey = "userID"
	RoleKey   = "role"
)

func AuthMiddleware(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "authorization header is required",
		})
	}

	parts := strings.SplitN(authHeader, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid authorization header",
		})
	}

	tokenString := parts[1]

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "JWT_SECRET is not configured",
		})
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil, fiber.ErrUnauthorized
			}

			return []byte(secret), nil
		},
	)

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid or expired token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid token claims",
		})
	}

	// Get user ID from JWT "sub"
	sub, ok := claims["sub"].(string)

	if !ok || sub == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "user id not found in token",
		})
	}

	userID, err := uuid.Parse(sub)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid user id",
		})
	}

	// Get role
	role, _ := claims["role"].(string)

	// Store authentication information
	c.Locals(UserIDKey, userID)
	c.Locals(RoleKey, role)

	return c.Next()
}