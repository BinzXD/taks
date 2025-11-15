package controller

import (
	"task/db"
	models "task/src/model"
	"task/src/validator"

	"fmt"
	"strconv"
	"task/src/helper"

	"github.com/gofiber/fiber/v2"
)

func ListUsers(c *fiber.Ctx) error {
	q := c.Query("q")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * perPage

	filter := ""
	var replacements []interface{}

	if q != "" {
		filter += " AND (u.fullname LIKE ? OR u.username LIKE ?)"
		replacements = append(replacements, "%"+q+"%"+"%", "%"+q+"%")
	}

	var total int64
	if err := db.DB.Raw(fmt.Sprintf("SELECT COUNT(*) FROM users u WHERE u.deleted_at IS NULL %s", filter), replacements...).Scan(&total).Error; err != nil {
		return helper.Result(c, fiber.StatusInternalServerError, nil, err)
	}

	var user []map[string]interface{}
	if err := db.DB.Raw(fmt.Sprintf(`
		SELECT u.id, u.fullname, u.username, u.email, r.name as role
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.deleted_at IS NULL
		%s
		ORDER BY u.created_at DESC
		LIMIT ? OFFSET ?`, filter),
		append(replacements, perPage, offset)...,
	).Scan(&user).Error; err != nil {
		return helper.Result(c, fiber.StatusInternalServerError, nil, err)
	}

	res := map[string]interface{}{
		"count": total,
		"rows":  user,
	}

	return helper.Result(c, fiber.StatusOK, res, nil)
}

func CreateUsers(c *fiber.Ctx) error {
	var body validator.UserCreateRequest
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
	users := models.Users{
		Fullname: body.Fullname,
		Username: body.Username,
		Email:    body.Email,
	}
	if err := db.DB.Create(&users).Error; err != nil {
		return helper.Result(c, 500, nil, err)
	}
	return helper.Result(c, 200, users, nil)
}

func DeleteUsers(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.Users
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return helper.Result(c, 404, nil, err)
	}
	if err := db.DB.Delete(&user).Error; err != nil {
		return helper.Result(c, 500, nil, err)
	}
	return helper.Result(c, 200, nil, nil)
}

func UpdateUsers(c *fiber.Ctx) error {
	id := c.Params("id")
	var body validator.UserUpdateRequest
	if err := c.BodyParser(&body); err != nil {
		return helper.Result(c, 400, nil, err)
	}
	if err := validator.Validate.Struct(body); err != nil {
		return helper.Result(c, 422, nil, err)
	}
	var user models.Users
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return helper.Result(c, 400, nil, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}
	user.Fullname = body.Fullname
	user.RoleID = body.RoleID

	if err := db.DB.Save(&user).Error; err != nil {
		return helper.Result(c, 500, nil, err)
	}
	return helper.Result(c, 200, user, nil)
}

func ShowUsers(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.Users
	if err := db.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return helper.Result(c, 400, nil, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}
	return helper.Result(c, 200, user, nil)
}
