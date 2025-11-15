package controller

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"task/db"
	models "task/src/model"
	roleValidator "task/src/validator"

	"task/src/helper"

	"github.com/gofiber/fiber/v2"
)

func CreateRole(c *fiber.Ctx) error {
	var body roleValidator.RoleCreateRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := roleValidator.Validate.Struct(body); err != nil {
		return helper.Result(c, 422, nil, errors.New(strings.Join(helper.ParseValidationErrors(err), ", ")))
	}

	if err := db.DB.Where("name = ?", body.Name).First(&models.Role{}).Error; err == nil {
		return helper.Result(c, 409, nil, fiber.NewError(fiber.StatusConflict, "Role already exists"))
	}

	role := models.Role{
		Name:   body.Name,
		Status: body.Status,
	}

	if err := db.DB.Create(&role).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return helper.Result(c, fiber.StatusOK, role, nil)
}

func ListRoles(c *fiber.Ctx) error {
	q := c.Query("q")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * perPage

	filter := ""
	var replacements []interface{}

	if q != "" {
		filter += " AND r.name LIKE ?"
		replacements = append(replacements, "%"+q+"%")
	}

	if status != "" {
		filter += " AND r.status = ?"
		replacements = append(replacements, status)
	}

	var total int64
	if err := db.DB.Raw(fmt.Sprintf("SELECT COUNT(*) FROM roles r WHERE r.deleted_at IS NULL %s", filter), replacements...).Scan(&total).Error; err != nil {
		return helper.Result(c, fiber.StatusInternalServerError, nil, err)
	}

	var roles []map[string]interface{}
	if err := db.DB.Raw(fmt.Sprintf(`
		SELECT r.id, r.name, r.status
		FROM roles r
		WHERE r.deleted_at IS NULL
		%s
		ORDER BY r.created_at DESC
		LIMIT ? OFFSET ?`, filter),
		append(replacements, perPage, offset)...,
	).Scan(&roles).Error; err != nil {
		return helper.Result(c, fiber.StatusInternalServerError, nil, err)
	}

	res := map[string]interface{}{
		"count": total,
		"rows":  roles,
	}

	return helper.Result(c, fiber.StatusOK, res, nil)
}
