package main

import (
	"log"
	"os"
	"task/config"
	"task/db"
	"task/src/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.Init()
	db.ConnectDB()
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			if err := db.DB.AutoMigrate(db.AutoMigrateModels()...); err != nil {
				log.Fatal("Failed to migrate:", err)
			}
			log.Println("Migration completed")
			return
		case "seed":
			db.Seed()
			log.Println("Seeding completed")
			return
		}
	}
	app := fiber.New()
	port := config.Get("APP_PORT")
	routes.SetupRoutes(app)
	log.Println("Server running on port: " + port)

	app.Listen(":" + port)
}
