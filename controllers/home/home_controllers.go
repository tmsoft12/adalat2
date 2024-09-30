package home

import (
	"fmt"
	"os"
	config "tm/db"
	model "tm/models"

	"github.com/gofiber/fiber/v2"
)

func Home_Page(c *fiber.Ctx) error {
	lang := c.Get("Accept-Language", "tm")

	var news []model.NewsSchema
	if err := config.DB.Find(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch news"})
	}

	var banner []model.BannerSchema
	if err := config.DB.Find(&banner).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch banner"})
	}

	var media []model.MediaSchema
	if err := config.DB.Find(&media).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch media"})
	}

	var employ []model.EmployerSchema
	if err := config.DB.Find(&employ).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch employ"})
	}

	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	// Process media URLs
	for i := range media {
		media[i].Video = fmt.Sprintf("http://%s%s/%s", ip, port, media[i].Video)
	}

	for i := range banner {
		banner[i].Image = fmt.Sprintf("http://%s%s/%s", ip, port, banner[i].Image)
	}

	for i := range news {
		news[i].Image = fmt.Sprintf("http://%s%s/%s", ip, port, news[i].Image)
		// Apply language-specific logic
		switch lang {
		case "en":
			// Only return English title and description
			news[i].Title = news[i].EN_title
			news[i].Description = news[i].EN_description
		case "ru":
			// Only return Russian title and description
			news[i].Title = news[i].RU_title
			news[i].Description = news[i].RU_description
		default:
		}
	}

	for i := range employ {
		employ[i].Image = fmt.Sprintf("http://%s%s/%s", ip, port, employ[i].Image)
	}

	return c.JSON(fiber.Map{
		"news":   news,
		"banner": banner,
		"media":  media,
		"employ": employ,
	})
}
