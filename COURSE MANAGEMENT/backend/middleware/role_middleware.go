package middleware

import (
	"course-management/utils"

	"github.com/gofiber/fiber/v2"
)

func RequireRole(roles ...string) fiber.Handler {

	return func(c *fiber.Ctx) error {

		role, ok := c.Locals("role").(string)

		if !ok {
			return utils.Error(
				c,
				fiber.StatusUnauthorized,
				"Unauthorized",
			)
		}

		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}

		return utils.Error(
			c,
			fiber.StatusForbidden,
			"Permission denied",
		)
	}
}
