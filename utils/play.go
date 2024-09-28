package utils

import (
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func Play(c *fiber.Ctx) error {
	fileName := c.Params("filename")
	filePath := filepath.Join("./uploads/media", fileName)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File not found"})
	}

	return c.SendFile(filePath, false)
}
