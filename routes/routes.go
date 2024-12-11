package routes

import (
    "github.com/gofiber/fiber/v2"
    "webapp/controllers"
    "webapp/middlewares"
)

func SetupRoutes(app *fiber.App) {

    // Public routes
    app.Post("/api/signup", controllers.Signup)
    app.Post("/api/login", controllers.Login)
    app.Get("/api/home",middlewares.AuthRequired(), controllers.Home)

    // Admin routes
    app.Post("/api/admin/login", controllers.AdminLogin)
    app.Get("/api/admin/dashboard", middlewares.AdminAuthRequired(), controllers.AdminPanel)
    app.Post("/api/admin/user", middlewares.AdminAuthRequired(), controllers.CreateUser)
    app.Get("/api/admin/user/:id", middlewares.AdminAuthRequired(), controllers.GetUser)
    app.Get("/api/admin/search", middlewares.AdminAuthRequired(), controllers.SearchUser)
    app.Put("/api/admin/user/:id", middlewares.AdminAuthRequired(), controllers.EditUser)
    app.Delete("/api/admin/user/:id", middlewares.AdminAuthRequired(), controllers.DeleteUser)
  
}
