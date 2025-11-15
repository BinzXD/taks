package db

import (
	"fmt"
	"log"
	"task/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	host := config.Get("DB_HOST")
	user := config.Get("DB_USER")
	password := config.Get("DB_PASS")
	dbname := config.Get("DB_NAME")
	port := config.Get("DB_PORT")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname,
	)

	var db *gorm.DB
	var err error

	// Retry sampai database siap
	for i := 0; i < 10; i++ { // max 10 percobaan
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Database connected")
			DB = db
			return
		}

		log.Println("Waiting for database to be ready... retry in 3s")
		time.Sleep(3 * time.Second)
	}

	log.Fatal("Failed to connect database after 10 attempts: ", err)
}
