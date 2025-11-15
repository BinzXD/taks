package routes

import (
	"task/src/controller"
	"task/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func BookingRoutes(router fiber.Router) {
	router.Use(middleware.JWTProtected())
	router.Get("/", controller.ListBooking)
	router.Post("/", controller.CreateBooking, middleware.AdminOnly)
	router.Get("/:id", controller.ShowBooking, middleware.AdminOnly)
	router.Post("/payment/:id", controller.PaymentBooking, middleware.AdminOnly)
}
