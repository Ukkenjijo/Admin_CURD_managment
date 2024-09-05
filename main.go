package main

import (
	"log"
	"webapp/config"
	"webapp/models"
	"webapp/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func CacheMiddleware(c *fiber.Ctx) error {
    c.Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
    c.Set("Pragma", "no-cache")
    c.Set("Expires", "0")
    return c.Next()
}

func main() {
    

    engine := html.New("./templates", ".html")

    app := fiber.New(fiber.Config{
        Views: engine,
    })
    app.Static("/static", "./static")

    app.Use(CacheMiddleware)

    

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
