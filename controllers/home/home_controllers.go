package home

import (
	"fmt"
	"os"
	config "tm/db"
	model "tm/models"
	"tm/utils"

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

	for i := range banner {
		banner[i].Image = fmt.Sprintf("http://%s%s/%s", ip, port, banner[i].Image)
	}
	utils.UrlCom(media)
	for i := range news {
		news[i].Image = fmt.Sprintf("http://%s%s/%s", ip, port, news[i].Image)
		switch lang {
		case "en":
			news[i].TM_title = news[i].EN_title
			news[i].TM_description = news[i].EN_description
		case "ru":
			news[i].TM_title = news[i].RU_title
			news[i].TM_description = news[i].RU_description
		default:
		}
	}

	for i := range employ {
		employ[i].Image = fmt.Sprintf("http://%s%s/%s", ip, port, employ[i].Image)
	}

	return c.Status(200).JSON(fiber.Map{
		"news":   news,
		"banner": banner,
		"media":  media,
		"employ": employ,
	})
}
