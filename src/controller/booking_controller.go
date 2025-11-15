package controller

import (
	"task/db"
	models "task/src/model"
	"task/src/validator"

	jwt "github.com/golang-jwt/jwt/v4"

	"fmt"
	"strconv"
	"task/src/helper"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateBooking(c *fiber.Ctx) error {
	var body validator.BookingCreateRequest

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	userIDFloat := claims["id"].(float64)
	userID := uint(userIDFloat)

	var user models.Users
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		return helper.Result(c, 404, nil, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}
	if err := c.BodyParser(&body); err != nil {
		return helper.Result(c, 400, nil, err)
	}

	if err := validator.Validate.Struct(body); err != nil {
		return helper.Result(c, 422, nil, err)
	}

	start, try := time.Parse(time.RFC3339, body.StartTime)
	if try != nil {
		return helper.Result(c, 400, nil, try)
	}
	end, try := time.Parse(time.RFC3339, body.EndTime)
	if try != nil {
		return helper.Result(c, 400, nil, try)
	}

	var checkBooking models.Booking
	if err := db.DB.
		Where("field_id = ? AND NOT (end_time <= ? OR start_time >= ?)", body.FieldID, start, end).
		First(&checkBooking).Error; err == nil {
		return helper.Result(c, 409, nil, fiber.NewError(fiber.StatusConflict, "Field is already booked in this time range"))
	}

	prefix := time.Now().Format("012006")
	var lastBooking models.Booking

	err := db.DB.
		Where("DATE_FORMAT(created_at, '%m%Y') = ?", prefix).
		Order("created_at DESC").
		First(&lastBooking).Error

	nextIncrement := 1
	if err == nil && lastBooking.Name != "" {
		lastNumber := lastBooking.Name[len(lastBooking.Name)-3:]
		if n, err := strconv.Atoi(lastNumber); err == nil {
			nextIncrement = n + 1
		}
	}

	newBooking := fmt.Sprintf("BOOK-%s-%03d", prefix, nextIncrement)
	booking := models.Booking{
		Name:      newBooking,
		CreatedBy: user.ID,
		Status:    "pending",
		FieldID:   body.FieldID,
		StartTime: start,
		EndTime:   end,
	}

	if err := db.DB.Create(&booking).Error; err != nil {
		return helper.Result(c, 500, nil, err)
	}

	return helper.Result(c, 200, booking, nil)
}

func ShowBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	var booking models.Booking
	if err := db.DB.Where("id = ?", id).First(&booking).Error; err != nil {
		return helper.Result(c, 400, nil, fiber.NewError(fiber.StatusNotFound, "Booking not found"))
	}
	return helper.Result(c, 200, booking, nil)
}

func ListBooking(c *fiber.Ctx) error {
	q := c.Query("q")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	offset := (page - 1) * perPage

	filter := ""
	var replacements []interface{}

	if q != "" {
		filter += " AND b.name LIKE ?"
		replacements = append(replacements, "%"+q+"%")
	}

	if status != "" {
		filter += " AND b.status = ?"
		replacements = append(replacements, status)
	}

	var total int64
	if err := db.DB.Raw(fmt.Sprintf("SELECT COUNT(*) FROM bookings b WHERE b.deleted_at IS NULL %s", filter), replacements...).Scan(&total).Error; err != nil {
		return helper.Result(c, fiber.StatusInternalServerError, nil, err)
	}

	var booking []map[string]interface{}
	if err := db.DB.Raw(fmt.Sprintf(`
		SELECT b.id, b.name, f.name as field_name, b.start_time, b.end_time, b.status
		FROM bookings b
		JOIN fields f ON b.field_id = f.id
		WHERE b.deleted_at IS NULL
		%s
		ORDER BY b.created_at DESC
		LIMIT ? OFFSET ?`, filter),
		append(replacements, perPage, offset)...,
	).Scan(&booking).Error; err != nil {
		return helper.Result(c, fiber.StatusInternalServerError, nil, err)
	}

	res := map[string]interface{}{
		"count": total,
		"rows":  booking,
	}

	return helper.Result(c, fiber.StatusOK, res, nil)
}

func PaymentBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	var booking models.Booking
	if err := db.DB.Where("id = ?", id).First(&booking).Error; err != nil {
		return helper.Result(c, 400, nil, fiber.NewError(fiber.StatusNotFound, "Booking not found"))
	}
	booking.Status = "done"

	if err := db.DB.Save(&booking).Error; err != nil {
		return helper.Result(c, 500, nil, err)
	}
	return helper.Result(c, 200, booking, nil)
}
