package controller

import (
	"task/db"
	"task/src/helper"
	models "task/src/model"
	"task/src/validator"

	jwt "task/src/middleware"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	var body validator.LoginRequest
	if err := c.BodyParser(&body); err != nil {
		return helper.Result(c, 400, nil, err)
	}

	if err := validator.Validate.Struct(body); err != nil {
		return helper.Result(c, 422, nil, err)
	}

	var user models.Users
	if err := db.DB.Where("email = ?", body.Email).First(&user).Error; err != nil {
		return helper.Result(c, 404, nil, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return helper.Result(c, 400, nil, fiber.NewError(fiber.StatusBadRequest, "Invalid password"))
	}

	token, err := jwt.GenerateToken(user.ID, user.RoleID, user.Email)
	if err != nil {
		return helper.Result(c, 500, nil, err)
	}

	return helper.Result(c, 200, token, nil)
}

func Register(c *fiber.Ctx) error {
	var body validator.RegisterRequest
	if err := c.BodyParser(&body); err != nil {
		return helper.Result(c, 400, nil, err)
	}

	if err := validator.Validate.Struct(body); err != nil {
		return helper.Result(c, 422, nil, err)
	}

	var user models.Users
	if err := db.DB.Where("email = ?", body.Email).First(&user).Error; err == nil {
		return helper.Result(c, 409, nil, fiber.NewError(fiber.StatusConflict, "Email already exists"))
	}

	if err := db.DB.Where("username = ?", body.Username).First(&user).Error; err == nil {
		return helper.Result(c, 409, nil, fiber.NewError(fiber.StatusConflict, "Username already exists"))
	}

	if body.ConfirmPassword != body.Password {
		return helper.Result(c, 400, nil, fiber.NewError(fiber.StatusBadRequest, "Passwords do not match"))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return helper.Result(c, 500, nil, err)
	}

	var role models.Role
	if err := db.DB.Where("name = ?", "user").First(&role).Error; err != nil {
		return helper.Result(c, 404, nil, err)
	}
	user = models.Users{
		Fullname: body.Fullname,
		Username: body.Username,
		Email:    body.Email,
		Password: string(hashedPassword),
		RoleID:   role.ID,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return helper.Result(c, 500, nil, err)
	}
	return helper.Result(c, 200, nil, nil)
}

func ChangePassword(c *fiber.Ctx) error {
	var body validator.ChangePasswordRequest

	if err := c.BodyParser(&body); err != nil {
		return helper.Result(c, 400, nil, err)
	}

	if err := validator.Validate.Struct(body); err != nil {
		return helper.Result(c, 422, nil, err)
	}

	var user models.Users
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword)); err != nil {
		return helper.Result(c, 400, nil, fiber.NewError(fiber.StatusBadRequest, "Invalid password"))
	}

	if body.ConfirmPassword != body.NewPassword {
		return helper.Result(c, 400, nil, fiber.NewError(fiber.StatusBadRequest, "Passwords do not match"))
	}

	users := models.Users{
		Password: body.NewPassword,
	}

	if err := db.DB.Save(&users).Error; err != nil {
		return helper.Result(c, 500, nil, err)
	}
	return helper.Result(c, 200, nil, nil)
}
