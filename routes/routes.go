package routes

import (
    "github.com/gofiber/fiber/v2"
    "webapp/controllers"
    "webapp/middlewares"
)

func SetupRoutes(app *fiber.App) {
  // Auth routes
  app.Get("/",controllers.ToLogin)
  app.Get("/login", controllers.ShowLoginPage)
  app.Post("/login", controllers.Login)
  app.Get("/signup", controllers.ShowSignupPage)
  app.Post("/signup", controllers.Signup)
  app.Get("/logout", controllers.Logout)

  // User home page
  app.Get("/home", middlewares.AuthRequired(), controllers.Home)

  // Admin routes
  app.Get("/admin/login", controllers.ShowAdminLoginPage)
  app.Post("/admin/login", controllers.AdminLogin)
  app.Get("/admin/logout", middlewares.AdminAuthRequired(), controllers.AdminLogout)
  app.Get("/admin/panel",middlewares.AdminAuthRequired(), controllers.AdminPanel)
  app.Get("/admin/search", middlewares.AdminAuthRequired(), controllers.SearchUser)
  app.Post("/admin/create", middlewares.AdminAuthRequired(), controllers.CreateUser)
  app.Post("/admin/edit/:id", middlewares.AdminAuthRequired(), controllers.EditUser)
  app.Get("/admin/delete/:id", middlewares.AdminAuthRequired(), controllers.DeleteUser)
}
