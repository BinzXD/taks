package middleware

import (
	"task/db"
	models "task/src/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AdminOnly(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	roleIDFloat, ok := claims["role"].(float64)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Forbidden",
			"errors": map[string]interface{}{
				"code":    403,
				"message": "Role claim missing or invalid",
			},
			"metadata": nil,
			"data":     nil,
		})
	}

	var role models.Role
	if err := db.DB.First(&role, "id = ?", roleIDFloat).Error; err != nil || role.Name != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"message": "Forbidden",
			"errors": map[string]interface{}{
				"code":    403,
				"message": "Only admin can access",
			},
			"metadata": nil,
			"data":     nil,
		})
	}

	return c.Next()
}
