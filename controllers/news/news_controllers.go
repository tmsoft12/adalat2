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
	var news []model.NewsSchema // Use a slice to fetch multiple records
	var views model.Views
	id := c.Params("id")
	user := c.Cookies("test")

	// Convert user from string to int
	userID, err := strconv.Atoi(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}
	ID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Query the news based on the user ID (assuming user_id is a foreign key in news_schemas)
	result := config.DB.Where("user_id = ?", userID).Find(&news)

	// If no news is found, return an error
	if result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No news found for the user",
		})
	}

	views.ID = ID
	views.UserID = user
	config.DB.Create(&views)

	// Return the news and the user in the response
	return c.Status(200).JSON(fiber.Map{
		"news": news,
		"user": userID,
	})
}
