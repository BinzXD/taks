package routes

import (
	"task/src/controller"
	"task/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func FieldRoutes(router fiber.Router) {
	router.Use(middleware.JWTProtected())
	router.Get("/", controller.ListFields)
	router.Get("/:id", controller.ShowField)
	router.Post("/", middleware.AdminOnly, controller.CreateField)
	router.Put("/:id", middleware.AdminOnly, controller.UpdateField)
	router.Delete("/:id", middleware.AdminOnly, controller.DeleteField)
}
