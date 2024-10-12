package news_page

import (
	"fmt"
	"os"
	"strconv"
	config "tm/db"
	model "tm/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const UserCookieName = "test"

func GetAllNews(c *fiber.Ctx) error {
	lang := c.Get("Accept-Language", "tm")
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	var news []model.NewsSchema

	if err := config.DB.Find(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch news"})
	}

	for i := range news {
		news[i].Image = formatImageURL(ip, port, news[i].Image)
		switch lang {
		case "en":
			news[i].TM_title = news[i].EN_title
			news[i].TM_description = news[i].EN_description
		case "ru":
			news[i].TM_title = news[i].RU_title
			news[i].TM_description = news[i].RU_description
		default:
			// Default to Turkmen, no action needed
		}
	}
	return c.JSON(news)
}

func NewsDetail(c *fiber.Ctx) error {
	var news model.NewsSchema
	id := c.Params("id")
	user := c.Cookies(UserCookieName)

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

	var existingView model.Views
	viewCheck := config.DB.Where("news_id = ? AND user_id = ?", NewsID, userID).First(&existingView)
	if viewCheck.Error == nil {
		return getViewCount(c, NewsID, news)
	} else if viewCheck.Error != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking view record: " + viewCheck.Error.Error(),
		})
	}

	// Create a new view record
	view := model.Views{
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
	countResult := config.DB.Model(&model.Views{}).
		Where("news_id = ?", NewsID).
		Count(&viewCount)

	if countResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error counting views: " + countResult.Error.Error(),
		})
	}

	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	news.Image = formatImageURL(ip, port, news.Image)

	if err := config.DB.Model(&news).Where("id = ?", NewsID).Update("view", int(viewCount)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating view count: " + err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"news": news,
	})
}

// Helper function to format image URLs
func formatImageURL(ip, port, imagePath string) string {
	url := fmt.Sprintf("http://%s", ip)
	if port != "80" {
		url = fmt.Sprintf("%s:%s", url, port)
	}
	return fmt.Sprintf("%s/%s", url, imagePath)
}
