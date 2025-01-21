package controllers

import (
	"time"
	"webapp/config"
	"webapp/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Struct for Admin Login
type AdminLoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}



// Struct for Create/Edit User
type UserForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func AdminLogin(c *fiber.Ctx) error {
	var data AdminLoginForm
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid request"})
	}

	var adminUser models.User
	config.DB.Where("username = ? AND is_admin = ?", data.Username, true).First(&adminUser)

	if adminUser.ID == 0 || bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(data.Password)) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": adminUser.Username,
		"is_admin": true,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
	})
}

func AdminPanel(c *fiber.Ctx) error {
	var users []models.User
	config.DB.Find(&users) // Fetch all users

	return c.JSON(users)
}

func SearchUser(c *fiber.Ctx) error {
	// Fetch the search query from the URL query parameters
	searchQuery := c.Query("search", "")

	// Search for users where the username matches the query
	var users []models.User
	config.DB.Where("username ILIKE ?", "%"+searchQuery+"%").Find(&users)

	// Return the search results as JSON
	return c.JSON(users)
}
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	return c.JSON(user)

}

func CreateUser(c *fiber.Ctx) error {
	// Admin create user logic
	var data UserForm
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create user")
	}

	user := models.User{
		Username: data.Username,
		Password: string(hashedPassword),
		Email:    data.Email, // Handle admin flag
	}

	config.DB.Create(&user)

	return c.SendStatus(fiber.StatusCreated)
}

func EditUser(c *fiber.Ctx) error {
	// Admin edit user logic
	id := c.Params("id")
	var user models.User

	if err := config.DB.First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	var data UserForm
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	if data.Username != "" {
		user.Username = data.Username
	}

	if data.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	if data.Email != "" {
		user.Email = data.Email
	}

	config.DB.Save(&user)

	return c.JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	// Admin delete user logic
	id := c.Params("id")

	// Fetch the user, including soft-deleted ones
	var user models.User
	if err := config.DB.Unscoped().First(&user, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	// Permanently delete the user from the database
	if err := config.DB.Unscoped().Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete user")
	}

	return c.SendStatus(fiber.StatusNoContent)
}


