package media

import (
	"fmt"
	"os"
	"strconv"
	config "tm/db"
	model "tm/models"
	"tm/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllMedia(c *fiber.Ctx) error {
	var media []model.MediaSchema
	if err := config.DB.Find(&media).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch news"})
	}
	utils.UrlCom(media)

	return c.Status(200).JSON(media)
}

func MediaDetail(c *fiber.Ctx) error {
	var media model.MediaSchema
	id := c.Params("id")
	user := c.Cookies("test")

	userID, err := strconv.Atoi(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	MediaID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid media ID format",
		})
	}
	result := config.DB.Where("id = ?", MediaID).First(&media)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("No media found with ID %d", MediaID),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching media: " + result.Error.Error(),
		})
	}
	var existingView model.ViewsMedia
	viewCheck := config.DB.Where("media_id = ? AND user_id", MediaID, userID).First(&existingView)
	if viewCheck.Error == nil {
		return getViewCount(c, MediaID, media)
	} else if viewCheck.Error != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking view record: " + viewCheck.Error.Error(),
		})
	}
	view := model.ViewsMedia{
		UserID:  userID,
		MediaID: MediaID,
	}
	if err := config.DB.Create(&view).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create view: " + err.Error(),
		})
	}
	return getViewCount(c, MediaID, media)

}

func getViewCount(c *fiber.Ctx, MediaID int, media model.MediaSchema) error {
	var viewCount int64
	countResult := config.DB.Model(&model.ViewsMedia{}).
		Where("media_id = ?", MediaID).
		Count(&viewCount)

	if countResult.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error counting views: " + countResult.Error.Error(),
		})
	}

	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	media.Video = fmt.Sprintf("http://%s%s/%s", ip, port, media.Video)

	if err := config.DB.Model(&media).Where("id = ?", MediaID).Update("view", int(viewCount)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating view count: " + err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"Media": media,
	})
}
