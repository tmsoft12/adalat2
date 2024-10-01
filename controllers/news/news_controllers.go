package news

import (
	"fmt"
	"os"
	"strconv"
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

func NewsDetail(c *fiber.Ctx) error {
	var news model.NewsSchema
	id := c.Params("id")
	user := c.Cookies("test")

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

	result := config.DB.Where("id = ?", NewsID).First(&news)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("No news found with ID %d", NewsID),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching news: " + result.Error.Error(),
		})
	}

	var existingView model.Vi
	viewCheck := config.DB.Where("news_id = ? AND user_id = ?", NewsID, userID).First(&existingView)
	if viewCheck.Error == nil {
		return getViewCount(c, NewsID, news)
	} else if viewCheck.Error != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking view record: " + viewCheck.Error.Error(),
		})
	}

	view := model.Vi{
		UserID: userID,
		NewsID: NewsID,
	}
	if err := config.DB.Create(&view).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create view: " + err.Error(),
		})
	}
	return getViewCount(c, NewsID, news)
}

func getViewCount(c *fiber.Ctx, NewsID int, news model.NewsSchema) error {
	var viewCount int64
	countResult := config.DB.Model(&model.Vi{}).
		Where("news_id = ?", NewsID).
		Count(&viewCount)

	if countResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error counting views: " + countResult.Error.Error(),
		})
	}

	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	news.Image = fmt.Sprintf("http://%s%s/%s", ip, port, news.Image)

	if err := config.DB.Model(&news).Where("id = ?", NewsID).Update("count", int(viewCount)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating view count: " + err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"news": news,
	})
}
