package middlewares

import (
	

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/contrib/jwt"   // For Fiber integration
	"github.com/golang-jwt/jwt/v5" // Updated to latest version
)

// AuthRequired checks for JWT authentication and redirects unauthenticated users
func AuthRequired() fiber.Handler {
    return jwtware.New(jwtware.Config{
        SigningKey:  jwtware.SigningKey{Key: []byte("your_secret_key")}, // Replace with your secret key
        TokenLookup: "header:Authorization",            // Look for the JWT in the "Authorization" header
        AuthScheme:  "Bearer ",                          // Look for the "Bearer " prefix in the "Authorization" header
        ContextKey:  "user",                            // Store the validated token in the context under "user"
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            return c.Status(fiber.StatusUnauthorized).SendString(err.Error()) // Return 401 Unauthorized if authentication fails
        },
    })
}
// AdminAuthRequired checks for JWT authentication and admin role
func AdminAuthRequired() fiber.Handler {
	
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte("secret")},   // Replace with your actual secret key
		TokenLookup:  "header:Authorization",        // Look for the JWT in the "Authorization" header with the "Bearer " prefix
		AuthScheme:   "Bearer ",                      // Look for the "Bearer " prefix in the "Authorization" header
		ContextKey:   "admin_user",                   // Store the validated token in the context under "admin_user"
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error()) // Return 401 Unauthorized if authentication fails
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			// Retrieve the validated token from the context
			userToken := c.Locals("admin_user").(*jwt.Token)
			
			// Extract claims from the token
			claims := userToken.Claims.(jwt.MapClaims)
			
			// Check if the user has the admin role
			isAdmin, ok := claims["is_admin"].(bool)
			if !ok || !isAdmin {
				return c.Status(fiber.StatusForbidden).SendString("Only admins are allowed to access this resource") // Return 403 Forbidden if the user is not an admin
			}

			// Proceed to the next middleware or route if the user is authenticated and has the admin role
			return c.Next()
		},
	})
}
