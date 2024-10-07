package news_controller

import (
	"fmt"
	"os"
	"strconv"
	"time"
	news_model "tm/controllers/admin/news/models"
	config "tm/db"
	"tm/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllNews(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		page = 10
	}
	var news []news_model.NewsSchema
	var total int64
	if err := config.DB.Model(&news_model.NewsSchema{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server internal error"})
	}
	if err := config.DB.Offset((page - 1) * page).Limit(pageSize).Find(&news).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server internal error"})

	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	for i := range news {
		news[i].Image = fmt.Sprintf("http://%s%s/%s", ip, port, news[i].Image)
	}
	return c.Status(200).JSON(fiber.Map{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"media":    news,
	})
}
func GetByIdNews(c *fiber.Ctx) error {
	idParam := c.Params(("id"))
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}
	var new news_model.NewsSchema
	if err := config.DB.First(&new, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Banner not found",
		})
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	new.Image = fmt.Sprintf("http://%s%s/%s", ip, port, new.Image)

	return c.Status(200).JSON(new)

}
func handleFileUpload(c *fiber.Ctx) (string, error) {
	filePath, err := utils.SaveFile(c, "image", "./uploads/news")
	if err != nil {
		return "", fmt.Errorf("cannot upload image file: %v", err)
	}
	return filePath, nil
}
func parseRequest(c *fiber.Ctx) (*news_model.NewsSchema, error) {
	news := new(news_model.NewsSchema)
	if err := c.BodyParser(news); err != nil {
		return nil, fmt.Errorf("cannot parse request: %v", err)
	}
	return news, nil
}
func validateMedia(news *news_model.NewsSchema) error {
	if news.TM_title == "" {
		return fmt.Errorf("TM title is required")
	}
	if news.TM_description == "" {
		return fmt.Errorf("TM description is required")
	}

	if news.EN_title == "" {
		return fmt.Errorf("EN title is required")
	}
	if news.EN_description == "" {
		return fmt.Errorf("EN description is required")
	}

	if news.RU_title == "" {
		return fmt.Errorf("RU title is required")
	}
	if news.RU_description == "" {
		return fmt.Errorf("RU description is required")
	}

	return nil
}
func UpdateNews(c *fiber.Ctx) error {
	idPram := c.Params("id")
	id, err := strconv.Atoi(idPram)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}
	var new news_model.NewsSchema
	if err := config.DB.First(&new, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Media not found",
		})
	}
	if newFilePath, err := handleFileUpload(c); err == nil {
		if new.Image != "" {
			if _, err := os.Stat(new.Image); err == nil {
				if err := os.Remove(new.Image); err != nil {
					fmt.Println("Error deleting old file:", err)

				}
			}
		}
		new.Image = newFilePath

	}
	if err := config.DB.Save(&new).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update new",
		})
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	new.Image = fmt.Sprintf("http://%s%s/%s", ip, port, new.Image)
	return c.Status(fiber.StatusOK).JSON(new)
}
func CreateNews(c *fiber.Ctx) error {
	filePath, err := handleFileUpload(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	new, err := parseRequest(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := validateMedia(new); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

	}
	new.Image = filePath
	new.Date = time.Now().Format("2006-01-02")
	return c.Status(fiber.StatusCreated).JSON(new)
}

func DeleteNew(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var new news_model.NewsSchema

	if err := config.DB.First(&new, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "new not found",
		})
	}

	if new.Image != "" {
		if _, err := os.Stat(new.Image); os.IsNotExist(err) {
			return c.Status(404).JSON(fiber.Map{"error": "File not found"})
		}

		fmt.Println("Attempting to delete file:", new.Image)
		if err := os.Remove(new.Image); err != nil {
			fmt.Println("Error deleting file:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete the media file"})
		}
	}

	if err := config.DB.Delete(&new).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete new from database",
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "New deleted successfully",
	})
}
