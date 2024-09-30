package news

import (
	"fmt"
	"os"
	"strconv" // Import this package to convert string to int
	config "tm/db"
	model "tm/models"

	"github.com/gofiber/fiber/v2"
)

func GetAllNews(c *fiber.Ctx) error {
	lang := c.Get("Accept-Language", "tm")
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	var news []model.NewsSchema
	if err := config.DB.Find(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch news"})
	}
	for i := range news {
		news[i].Image = fmt.Sprintf("http://%s%s/%s", ip, port, news[i].Image)
		switch lang {
		case "en":
			news[i].Title = news[i].EN_title
			news[i].Description = news[i].EN_description
		case "ru":
			news[i].Title = news[i].RU_title
			news[i].Description = news[i].RU_description
		default:
		}
	}
	return c.JSON(news)
}

func ViewsAll(c *fiber.Ctx) error {
	var news model.NewsSchema // Tek bir haber için struct
	var views model.ViewsNews // Görüntüleme kaydı struct
	id := c.Params("id")      // URL'den gelen haber ID'si
	user := c.Cookies("test") // Çerezden gelen kullanıcı bilgisi

	userID, err := strconv.Atoi(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Haber ID'sini string'ten int'e dönüştürme
	ID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid news ID",
		})
	}

	// Haber var mı kontrol et
	result := config.DB.Where("id = ?", ID).First(&news)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "News not found",
		})
	}

	// Görüntüleme kaydının zaten olup olmadığını kontrol et
	var existingView model.ViewsNews
	viewCheck := config.DB.Where("id = ? AND user_id = ?", ID, userID).First(&existingView)
	if viewCheck.Error == nil { // Eğer kayıt varsa
		return c.Status(200).JSON(fiber.Map{
			"message": "View already recorded",
		})
	}

	views.UserID = userID // Kullanıcı kimliği
	views.ID = ID         // Hangi habere ait olduğunu belirt

	if err := config.DB.Create(&views).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create view",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"news": news,
		"user": userID,
	})
}
