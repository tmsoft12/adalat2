package news

import (
	"fmt"
	"os"
	"strconv" // Import this package to convert string to int
	config "tm/db"
	model "tm/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
	var news model.NewsSchema
	id := c.Params("id")      // URL'den gelen haber ID'si
	user := c.Cookies("test") // Çerezden gelen kullanıcı bilgisi

	userID, err := strconv.Atoi(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	NewsID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid news ID format",
		})
	}

	// Haber verisini almak için sorgu
	result := config.DB.Where("id = ?", NewsID).First(&news)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("No news found with ID %d", NewsID),
			})
		}
		// Diğer hata durumları
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching news: " + result.Error.Error(),
		})
	}

	var existingView model.Vi
	viewCheck := config.DB.Where("news_id = ? AND user_id = ?", NewsID, userID).First(&existingView)
	if viewCheck.Error == nil {
		return c.Status(200).JSON(fiber.Map{
			"news":    news,
			"message": "View already recorded",
		})
	} else if viewCheck.Error != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking view record: " + viewCheck.Error.Error(),
		})
	}

	// Yeni görüntüleme kaydı oluştur
	view := model.Vi{
		UserID: userID,
		NewsID: NewsID,
	}
	if err := config.DB.Create(&view).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create view: " + err.Error(),
		})
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	news.Image = fmt.Sprintf("http://%s%s/%s", ip, port, news.Image)

	return c.Status(200).JSON(fiber.Map{
		"news": news,
		"user": userID,
	})
}
