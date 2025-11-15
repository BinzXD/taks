package routes

import (
	"task/src/controller"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router fiber.Router) {
	router.Post("/users", controller.CreateUsers)
	router.Get("/users", controller.ListUsers)
	router.Get("/users/:id", controller.ShowUsers)
	router.Put("/users/:id", controller.UpdateUsers)
	router.Delete("/users/:id", controller.DeleteUsers)
}
