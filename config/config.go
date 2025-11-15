package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
}

func Get(key string) string {
	return os.Getenv(key)
}
