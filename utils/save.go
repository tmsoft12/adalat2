package utils

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/rand"
)

func SaveFile(c *fiber.Ctx, fieldName, dir string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(file.Filename)

	newFileName := fmt.Sprintf("%d_%d%s", time.Now().Unix(), rand.Intn(1000), ext)

	filePath := fmt.Sprintf("%s/%s", dir, newFileName)

	if err := c.SaveFile(file, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}
