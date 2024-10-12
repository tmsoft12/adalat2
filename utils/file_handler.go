package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func SaveFile(c *fiber.Ctx, formFieldName, uploadDir string) (string, error) {
	file, err := c.FormFile(formFieldName)
	if err != nil {
		return "", fmt.Errorf("file not found: %v", err)
	}

	// Upload directory oluştur
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("cannot create upload directory: %v", err)
	}

	filePath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return "", fmt.Errorf("cannot save file: %v", err)
	}

	return filePath, nil
}

func DeleteFile(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("error deleting file: %v", err)
		}
	}
	return nil
}
func GetHostAndPort() string {
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if ip == "" {
		ip = "localhost"
	}
	if port == "" {
		port = "3000"
	}

	return fmt.Sprintf("%s%s", ip, port)
}
