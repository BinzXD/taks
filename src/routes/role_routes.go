package routes

import (
	"task/src/controller"

	"github.com/gofiber/fiber/v2"
)

func RoleRoutes(router fiber.Router) {
	router.Get("/roles", controller.ListRoles)
	router.Post("/roles", controller.CreateRole)
}
