package controllers

import (
	"fmt"
	"log"
	"time"
	"webapp/config"
	"webapp/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" // Updated to latest version
	"golang.org/x/crypto/bcrypt"
)


func ToLogin(c *fiber.Ctx) error{
    return c.Redirect("/login")
}
// Render the login page
func ShowLoginPage(c *fiber.Ctx) error {
    tokenString := c.Cookies("token")
    fmt.Println(tokenString)
    if tokenString != "" {
        // Parse the token to check if it's valid
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Validate the signing method and return the secret key
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte("your_secret_key"), nil // Replace with your secret key
        })
        fmt.Println(err," ",token.Valid)
        

        // If the token is valid, redirect to the home page
        if err == nil && token.Valid {
            fmt.Println("Token valid. Redirecting to /home")
            return c.Redirect("/home")
        }
    }
    
    return c.Render("login", nil)
}

// Handle login and issue JWT token
func Login(c *fiber.Ctx) error {
    
    // Extract login details from form data
    type LoginData struct {
        Username string `form:"username"`
        Email    string `form:"email"`
        Password string `form:"password"`
    }
    var data LoginData
    if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
    }

    var user models.User
    config.DB.Where("username = ?", data.Username).First(&user)

    // Check if user exists and if password matches
    if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)) != nil {
        return c.Render("login",fiber.Map{
            "Error":"Invalid Credintials",
        })
    }

    // Create JWT claims
    claims := jwt.MapClaims{
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 72).Unix(),
    }

    // Create the token with claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign the token with a secret key
    tokenString, err := token.SignedString([]byte("your_secret_key")) // Replace with your secret key
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Could not login")
    }

    // Set the token in a cookie
    c.Cookie(&fiber.Cookie{
        Name:     "token",
        Value:    tokenString,
        Expires:  time.Now().Add(time.Hour * 72),
        HTTPOnly: true, // More secure
    })

    return c.Redirect("/home")
}

// Render the signup page
func ShowSignupPage(c *fiber.Ctx) error {
    tokenString := c.Cookies("token")
    fmt.Println(tokenString)
    if tokenString != "" {
        // Parse the token to check if it's valid
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Validate the signing method and return the secret key
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte("your_secret_key"), nil // Replace with your secret key
        })
        fmt.Println(err," ",token.Valid)
        

        // If the token is valid, redirect to the home page
        if err == nil && token.Valid {
            fmt.Println("Token valid. Redirecting to /home")
            return c.Redirect("/home")
        }
    }
    return c.Render("signup", nil)
}

// Handle user signup
func Signup(c *fiber.Ctx) error {
    type SignupData struct {
        Username string `form:"username"`
        Email    string `form:"email"`
        Password string `form:"password"`
    }
    var data SignupData
    if err := c.BodyParser(&data); err != nil {
        log.Println("Error parsing body:", err)
        return c.Status(fiber.StatusBadRequest).SendString("Invalid request")
    }

    // Check if the username already exists
    var existingUser models.User
    config.DB.Where("username = ?", data.Username).First(&existingUser)
    if existingUser.ID != 0 {
        // Username already exists
        return c.Render("signup",fiber.Map{
            "Error":"Username Alredy Taken",
        })
    }

    // Hash the user's password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Error creating user")
    }

    // Create a new user record
    user := models.User{
        Username: data.Username,
        Password: string(hashedPassword),
        Email: data.Email,
    }
    config.DB.Create(&user)

    return c.Redirect("/login")
}

// Render the home page
func Home(c *fiber.Ctx) error {
    return c.Render("home", nil)
}
func Logout(c *fiber.Ctx) error {
    // Clear the JWT token by deleting the cookie
    c.Cookie(&fiber.Cookie{
        Name:     "token",
        Value:    "",
        Expires:  time.Now().Add(-time.Hour), // Set expiration to the past
        HTTPOnly: true,
    })

    // Redirect to the login page
    return c.Redirect("/login")
}

