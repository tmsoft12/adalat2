package media

import (
	config "tm/db"
	model "tm/models"
	"tm/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllMedia(c *fiber.Ctx) error {
	var media []model.MediaSchema
	if err := config.DB.Find(&media).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch news"})
	}
	utils.UrlCom(media)

	return c.Status(200).JSON(media)
}
