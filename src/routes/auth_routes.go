package routes

import (
	"task/src/controller"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router fiber.Router) {
	router.Post("/login", controller.Login)
	router.Post("/register", controller.Register)
	router.Post("/change-password", controller.ChangePassword)
}
