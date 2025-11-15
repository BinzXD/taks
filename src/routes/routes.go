package routes

import (
	"task/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is running 🚀")
	})
	api := app.Group("/api")

	auth := api.Group("/auth")
	AuthRoutes(auth)

	field := api.Group("field")
	FieldRoutes(field)

	admin := api.Group("/admin")
	admin.Use(middleware.JWTProtected())
	admin.Use(middleware.AdminOnly)
	UserRoutes(admin)
	RoleRoutes(admin)

	booking := api.Group("/booking")
	BookingRoutes(booking)
}
