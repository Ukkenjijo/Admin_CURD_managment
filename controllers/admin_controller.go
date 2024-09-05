package controllers

import (
	
	"fmt"
	"time"
	"webapp/config"
	"webapp/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Struct for Admin Login
type AdminLoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

// Struct for User Search
type UserSearchForm struct {
	Search string `form:"search"`
}

// Struct for Create/Edit User
type UserForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
	Email    string `form:"email"`
	
}

func ShowAdminLoginPage(c *fiber.Ctx) error {
	tokenString := c.Cookies("admin_token")
    fmt.Println(tokenString)
    if tokenString != "" {
        // Parse the token to check if it's valid
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Validate the signing method and return the secret key
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte("secret"), nil // Replace with your secret key
        })
        fmt.Println(err," ",token.Valid)
        

        // If the token is valid, redirect to the home page
        if err == nil && token.Valid {
            fmt.Println("Token valid. Redirecting to /home")
            return c.Redirect("/admin/panel")
        }
    }
	return c.Render("admin_login", nil)
}

func AdminLogin(c *fiber.Ctx) error {
	var data AdminLoginForm
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
	}

	var adminUser models.User
	config.DB.Where("username = ? AND is_admin = ?", data.Username, true).First(&adminUser)

	if adminUser.ID == 0 || bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(data.Password)) != nil {
		return c.Render("admin_login",fiber.Map{
			"Error":"Admin access doesn't exist in this username",
	})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": adminUser.Username,
		"is_admin": true,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	fmt.Println(token)

	tokenString, _ := token.SignedString([]byte("secret"))
	c.Cookie(&fiber.Cookie{
		Name:  "admin_token",
		Value: tokenString,
	})
    
	return c.Redirect("/admin/panel")
}

func AdminPanel(c *fiber.Ctx) error {
	var users []models.User
	config.DB.Find(&users) // Fetch all users

	return c.Render("admin", fiber.Map{
		"Users": users,
	})
}

func SearchUser(c *fiber.Ctx) error {
    // Fetch the search query from the URL query parameters
    searchQuery := c.Query("search", "")

    // Search for users where the username matches the query
    var users []models.User
    config.DB.Where("username ILIKE ?", "%"+searchQuery+"%").Find(&users)

    // Render the admin template with the search results
    return c.Render("admin", fiber.Map{
        "Users": users, // Return search results in Users
    })
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
		Email: data.Email,// Handle admin flag
	}

	config.DB.Create(&user)

	return c.Redirect("/admin/panel")
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

	// Handle admin flag

	config.DB.Save(&user)

	return c.Redirect("/admin/panel")
}

func DeleteUser(c *fiber.Ctx) error {
    // Admin delete user logic
    id := c.Params("id")
    var user models.User

    // Fetch the user, including soft-deleted ones
    if err := config.DB.Unscoped().First(&user, id).Error; err != nil {
        return c.Status(fiber.StatusNotFound).SendString("User not found")
    }

    // Permanently delete the user from the database
    if err := config.DB.Unscoped().Delete(&user).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete user")
    }

    return c.Redirect("/admin/panel")
}


func AdminLogout(c *fiber.Ctx) error {
	// Clear the admin JWT token by deleting the cookie
	c.Cookie(&fiber.Cookie{
		Name:     "admin_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set expiration to the past
		HTTPOnly: true,
	})

	// Redirect to the admin login page
	return c.Redirect("/admin/login")
}
