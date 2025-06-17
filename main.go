package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func main() {
	// Load our .env file
	err := godotenv.Load(".evn")
	if err != nil {
		log.Fatalf("Error while loading .evn with: %s", err)
	}

	app := fiber.New()
	app.Listen(":8080")
}
