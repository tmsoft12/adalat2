package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func FakeUser(c *fiber.Ctx) error {
	// Check if the "id" cookie exists
	username := c.Cookies("id")

	// If the cookie doesn't exist, create and set it
	if username == "" {
		cookie := new(fiber.Cookie)
		cookie.Name = "id"
		cookie.Value = "kerim" // Set default value
		cookie.Expires = time.Now().Add(24 * time.Hour)
		cookie.HTTPOnly = true // Optional for security
		c.Cookie(cookie)
	}

	// Continue to the next handler (e.g., Home_Page)
	return c.Next()
}
