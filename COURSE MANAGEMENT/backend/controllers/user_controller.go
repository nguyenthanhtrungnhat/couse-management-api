package controllers

import (
	userDTO "course-management/dto/user"
	"course-management/services"
	"course-management/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserController struct {
	service services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		service: services.NewUserService(),
	}
}

func (uc *UserController) Profile(c *fiber.Ctx) error {

	id, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "Invalid user id")
	}

	user, err := uc.service.GetProfile(id)
	if err != nil {
		return utils.Error(c, fiber.StatusNotFound, err.Error())
	}

	return utils.Success(c, "Success", user)
}

func (uc *UserController) UpdateProfile(c *fiber.Ctx) error {

	id, err := uuid.Parse(c.Locals("user_id").(string))
	if err != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "Invalid user id")
	}

	var req userDTO.UpdateProfileRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request")
	}

	user, err := uc.service.UpdateProfile(id, req)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.Success(c, "Profile updated", user)
}