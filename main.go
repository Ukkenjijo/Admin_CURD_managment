package main

import (
	"log"
	"webapp/config"
	"webapp/models"
	"webapp/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func CacheMiddleware(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	c.Set("Pragma", "no-cache")
	c.Set("Expires", "0")
	return c.Next()
}

func main() {

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	// Initialize DB
	config.InitDB()

	err := config.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	log.Fatal(app.Listen(":3000"))
}
