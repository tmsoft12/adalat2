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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch media"})
	}
	utils.UrlCom(media)

	return c.Status(fiber.StatusOK).JSON(media)
}

func MediaDetail(c *fiber.Ctx) error {
	var media model.MediaSchema
	id := c.Params("id")
	userCookie := c.Cookies("test")

	userID, err := strconv.Atoi(userCookie)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID format"})
	}

	mediaID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid media ID format"})
	}

	if err := config.DB.Where("id = ?", mediaID).First(&media).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": fmt.Sprintf("No media found with ID %d", mediaID)})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching media: " + err.Error()})
	}

	// Check if the user has already viewed the media
	// Check if the user has already viewed the media
	var existingView model.ViewsMedia
	viewCheck := config.DB.Where("media_id = ? AND user_id = ?", mediaID, userID).First(&existingView)

	if viewCheck.Error == gorm.ErrRecordNotFound {
		// No existing view, create a new one
		view := model.ViewsMedia{
			UserID:  userID,
			MediaID: mediaID,
		}
		if err := config.DB.Create(&view).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create view: " + err.Error()})
		}
	} else if viewCheck.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking view record: " + viewCheck.Error.Error()})
	}

	// Get the updated view count
	return getViewCount(c, mediaID, media)
}

func getViewCount(c *fiber.Ctx, mediaID int, media model.MediaSchema) error {
	var viewCount int64
	if err := config.DB.Model(&model.ViewsMedia{}).Where("media_id = ?", mediaID).Count(&viewCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error counting views: " + err.Error()})
	}

	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	// Format media video URL
	media.Video = formatURL(ip, port, media.Video)

	// Update view count in the media record
	if err := config.DB.Model(&media).Where("id = ?", mediaID).Update("view", int(viewCount)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating view count: " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"media": media})
}

// Helper function to format media URLs
func formatURL(ip, port, mediaPath string) string {
	url := fmt.Sprintf("http://%s", ip)
	if port != "80" {
		url = fmt.Sprintf("%s:%s", url, port)
	}
	return fmt.Sprintf("%s/%s", url, mediaPath)
}
