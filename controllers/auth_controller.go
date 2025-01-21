package controllers

import (
	"time"
	"webapp/config"
	"webapp/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" // Updated to latest version
	"golang.org/x/crypto/bcrypt"
)

// Handle login and issue JWT token
func Login(c *fiber.Ctx) error {

	// Extract login details from form data
	type LoginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var data LoginData
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var user models.User
	config.DB.Where("username = ?", data.Username).First(&user)

	// Check if user exists and if password matches
	if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Create JWT claims
	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your_secret_key")) // Replace with your secret key
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}

	// Return the token as response
	return c.JSON(fiber.Map{"token": tokenString})
}

// Handle signup and issue JWT token
func Signup(c *fiber.Ctx) error {

	// Extract signup details from request body
	type SignupData struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	var data SignupData
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Check if username already exists
	var existingUser models.User
	config.DB.Where("username = ?", data.Username).First(&existingUser)
	if existingUser.ID != 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// Create a new user record
	user := models.User{
		Username: data.Username,
		Password: string(hashedPassword),
		Email:    data.Email,
	}
	config.DB.Create(&user)

	//send a statusok response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User created successfully"})
}

func Home(c *fiber.Ctx) error {
    // Get the user from the context
    userToken := c.Locals("user").(*jwt.Token)
    claims := userToken.Claims.(jwt.MapClaims)
    username, ok := claims["username"].(string)
    if !ok {
        return c.Status(fiber.StatusInternalServerError).SendString("Error fetching user")
    }

    // Return the user's name as response
    return c.SendString(username)
}
