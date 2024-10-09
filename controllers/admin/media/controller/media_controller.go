package media_controller

import (
	"fmt"
	"os"
	"strconv"
	"time"
	media_model "tm/controllers/admin/media/models"
	config "tm/db"
	"tm/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllMedia(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		page = 10
	}
	var media []media_model.MediaSchema
	var total int64
	if err := config.DB.Model(&media_model.MediaSchema{}).Count(&total).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server internal error"})
	}

	if err := config.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&media).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server internal error"})

	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	for i := range media {
		media[i].Video = fmt.Sprintf("http://%s%s/%s", ip, port, media[i].Video)
	}
	for i := range media {
		media[i].Cover = fmt.Sprintf("http://%s%s/%s", ip, port, media[i].Cover)
	}

	return c.Status(200).JSON(fiber.Map{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"media":    media,
	})

}

func CreateMedia(c *fiber.Ctx) error {
	filePath, err := handleFileUpload(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	media, err := parseMediaRequest(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validateMedia(media); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	media.Video = filePath
	media.Date = time.Now().Format("2006-01-02")

	if err := saveMediaToDB(media); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(media)
}

func handleFileUpload(c *fiber.Ctx) (string, error) {
	filePath, err := utils.SaveFile(c, "video", "./uploads/media")
	if err != nil {
		return "", fmt.Errorf("cannot upload video file: %v", err)
	}
	return filePath, nil
}

func parseMediaRequest(c *fiber.Ctx) (*media_model.MediaSchema, error) {
	media := new(media_model.MediaSchema)
	if err := c.BodyParser(media); err != nil {
		return nil, fmt.Errorf("cannot parse request: %v", err)
	}
	return media, nil
}

func validateMedia(media *media_model.MediaSchema) error {
	if media.TM_title == "" {
		return fmt.Errorf("title is required")
	}
	return nil
}

func saveMediaToDB(media *media_model.MediaSchema) error {
	if err := config.DB.Create(media).Error; err != nil {
		fmt.Println("Database error:", err)
		return fmt.Errorf("cannot create media")
	}
	return nil
}

func GetById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}
	var media media_model.MediaSchema
	if err := config.DB.First(&media, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Banner not found",
		})
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	media.Video = fmt.Sprintf("http://%s%s/%s", ip, port, media.Video)

	return c.Status(200).JSON(media)

}
func DeleteMedia(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var media media_model.MediaSchema

	if err := config.DB.First(&media, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Media not found",
		})
	}

	if media.Video != "" {
		// Dosyanın mevcut olup olmadığını kontrol et
		if _, err := os.Stat(media.Video); os.IsNotExist(err) {
			return c.Status(404).JSON(fiber.Map{"error": "File not found"})
		}

		fmt.Println("Attempting to delete file:", media.Video) // Silinmeye çalışılan dosya yolu
		if err := os.Remove(media.Video); err != nil {
			fmt.Println("Error deleting file:", err) // Hata detayını logla
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete the media file"})
		}
	}

	if err := config.DB.Delete(&media).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete media from database",
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Media deleted successfully",
	})
}
func UpdateMedia(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	// Eski medyayı veritabanından getir
	var media media_model.MediaSchema
	if err := config.DB.First(&media, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Media not found",
		})
	}

	// Yeni dosya yüklendiyse eski dosyayı sil ve yeni dosyayı yükle
	if newFilePath, err := handleFileUpload(c); err == nil {
		if media.Video != "" {
			if _, err := os.Stat(media.Video); err == nil {
				fmt.Println("Deleting old file:", media.Video)
				if err := os.Remove(media.Video); err != nil {
					fmt.Println("Error deleting old file:", err)
				}
			}
		}
		// Yeni dosya yolunu güncelle
		media.Video = newFilePath
	}

	// Diğer alanları güncelle
	if err := c.BodyParser(&media); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request",
		})
	}

	if err := config.DB.Save(&media).Error; err != nil {
		fmt.Println("Error while updating media:", err) // Hata mesajı loglanıyor
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update media",
		})
	}

	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	media.Video = fmt.Sprintf("http://%s%s/%s", ip, port, media.Video)

	return c.Status(fiber.StatusOK).JSON(media)

}
