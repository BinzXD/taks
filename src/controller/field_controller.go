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

func ListFields(c *fiber.Ctx) error {
	q := c.Query("q")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * perPage

	filter := ""
	var replacements []interface{}

	if q != "" {
		filter += " AND f.name LIKE ?"
		replacements = append(replacements, "%"+q+"%")
	}

	var total int64
	if err := db.DB.Raw(fmt.Sprintf("SELECT COUNT(*) FROM fields f WHERE f.deleted_at IS NULL %s", filter), replacements...).Scan(&total).Error; err != nil {
		return helper.Result(c, fiber.StatusInternalServerError, nil, err)
	}

	var field []map[string]interface{}
	if err := db.DB.Raw(fmt.Sprintf(`
		SELECT f.id, f.name, f.price as price_per_hour, f.location
		FROM fields f
		WHERE f.deleted_at IS NULL
		%s
		ORDER BY f.created_at DESC
		LIMIT ? OFFSET ?`, filter),
		append(replacements, perPage, offset)...,
	).Scan(&field).Error; err != nil {
		return helper.Result(c, fiber.StatusInternalServerError, nil, err)
	}

	res := map[string]interface{}{
		"count": total,
		"rows":  field,
	}

	return helper.Result(c, fiber.StatusOK, res, nil)
}

func CreateField(c *fiber.Ctx) error {
	var body validator.FieldCreateRequest
	if err := c.BodyParser(&body); err != nil {
		return helper.Result(c, 400, nil, err)
	}
	if err := validator.Validate.Struct(body); err != nil {
		return helper.Result(c, 422, nil, err)
	}
	field := models.Field{
		Name:     body.Name,
		Price:    body.Price,
		Location: body.Location,
	}
	if err := db.DB.Create(&field).Error; err != nil {
		return helper.Result(c, 500, nil, err)
	}
	return helper.Result(c, 200, field, nil)
}

func DeleteField(c *fiber.Ctx) error {
	id := c.Params("id")
	var field models.Field
	if err := db.DB.Where("id = ?", id).First(&field).Error; err != nil {
		return helper.Result(c, 404, nil, err)
	}
	if err := db.DB.Delete(&field).Error; err != nil {
		return helper.Result(c, 500, nil, err)
	}
	return helper.Result(c, 200, nil, nil)
}

func UpdateField(c *fiber.Ctx) error {
	id := c.Params("id")
	var body validator.FieldUpdateRequest
	if err := c.BodyParser(&body); err != nil {
		return helper.Result(c, 400, nil, err)
	}
	if err := validator.Validate.Struct(body); err != nil {
		return helper.Result(c, 422, nil, err)
	}
	var field models.Field
	if err := db.DB.Where("id = ?", id).First(&field).Error; err != nil {
		return helper.Result(c, 400, nil, fiber.NewError(fiber.StatusNotFound, "Field not found"))
	}
	field.Name = body.Name
	field.Price = body.Price
	field.Location = body.Location
	if err := db.DB.Save(&field).Error; err != nil {
		return helper.Result(c, 500, nil, err)
	}
	return helper.Result(c, 200, field, nil)
}

func ShowField(c *fiber.Ctx) error {
	id := c.Params("id")
	var field models.Field
	if err := db.DB.Where("id = ?", id).First(&field).Error; err != nil {
		return helper.Result(c, 400, nil, fiber.NewError(fiber.StatusNotFound, "Field not found"))
	}
	return helper.Result(c, 200, field, nil)
}
