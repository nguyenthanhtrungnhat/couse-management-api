package controllers

import (
	authDTO "course-management/dto/auth"
	"course-management/services"
	"course-management/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	service services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		service: services.NewAuthService(),
	}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {

	var req authDTO.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.Error(
			c,
			fiber.StatusBadRequest,
			"Invalid request body",
		)
	}

	result, err := ac.service.Register(req)

	if err != nil {
		return utils.Error(
			c,
			fiber.StatusBadRequest,
			err.Error(),
		)
	}

	return utils.Success(
		c,
		"Register successfully",
		result,
	)
}

func (ac *AuthController) Login(c *fiber.Ctx) error {

	var req authDTO.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.Error(
			c,
			fiber.StatusBadRequest,
			"Invalid request body",
		)
	}

	result, err := ac.service.Login(req)

	if err != nil {
		return utils.Error(
			c,
			fiber.StatusUnauthorized,
			err.Error(),
		)
	}

	return utils.Success(
		c,
		"Login successfully",
		result,
	)
}

