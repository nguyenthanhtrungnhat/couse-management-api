package middleware

import (
	"course-management/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {

	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return utils.Error(
				c,
				fiber.StatusUnauthorized,
				"Authorization header is required",
			)
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.Error(
				c,
				fiber.StatusUnauthorized,
				"Invalid authorization format",
			)
		}

		claims, err := utils.ParseToken(parts[1])

		if err != nil {
			return utils.Error(
				c,
				fiber.StatusUnauthorized,
				"Invalid or expired token",
			)
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}