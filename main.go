package main

import (
	"log"
	"os"
	config "tm/db"
	"tm/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	config.ConnectDatabase()
	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024,
	})

	routers.InitRouters(app)
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("Request: %s %s", c.Method(), c.Path())
		return c.Next()
	})

	log.Fatal(app.Listen(host + port))
}
