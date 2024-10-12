package media_controller

import (
	"fmt"
	"strconv"
	"time"
	media_service "tm/controllers/admin/media/controller/services"
	media_model "tm/controllers/admin/media/models"
	"tm/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllMedia(c *fiber.Ctx) error {
	page, pageSize := parsePaginationParams(c)

	media, total, err := media_service.GetAllMedia(page, pageSize)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server internal error"})
	}
	ipAndPort := utils.GetHostAndPort()

	for i := range media {
		media[i].Video = fmt.Sprintf("http://%s/%s", ipAndPort, media[i].Video)
		media[i].Cover = fmt.Sprintf("http://%s/%s", ipAndPort, media[i].Cover)
	}

	return c.Status(200).JSON(fiber.Map{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"media":    media,
	})
}

func CreateMedia(c *fiber.Ctx) error {
	filePath, err := utils.SaveFile(c, "video", "./uploads/media/video")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	photoPath, err := utils.SaveFile(c, "cover", "./uploads/media/cover")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	media, err := parseMediaRequest(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	media.Video = filePath
	media.Cover = photoPath
	media.Date = time.Now().Format("2006-01-02")

	if err := media_service.SaveMedia(media); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(media)
}

func GetById(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	media, err := media_service.FindMediaById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	ipAndPort := utils.GetHostAndPort()
	media.Video = fmt.Sprintf("http://%s/%s", ipAndPort, media.Video)
	media.Cover = fmt.Sprintf("http://%s/%s", ipAndPort, media.Cover)
	return c.Status(200).JSON(media)
}

func UpdateMedia(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	media, err := media_service.FindMediaById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media not found"})
	}

	// Yeni dosya yüklendiyse eski dosyayı sil ve yenisini yükle
	if newFilePath, err := utils.SaveFile(c, "video", "./uploads/media/video"); err == nil {
		utils.DeleteFile(media.Video)
		media.Video = newFilePath
	}

	if err := c.BodyParser(media); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request"})
	}

	if err := media_service.UpdateMedia(media); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update media"})
	}

	ipAndPort := utils.GetHostAndPort()
	media.Video = fmt.Sprintf("http://%s/%s", ipAndPort, media.Video)
	media.Cover = fmt.Sprintf("http://%s/%s", ipAndPort, media.Cover)
	return c.Status(fiber.StatusOK).JSON(media)
}

func DeleteMedia(c *fiber.Ctx) error {
	id, err := parseIDParam(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	media, err := media_service.FindMediaById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Media not found"})
	}

	utils.DeleteFile(media.Video)

	if err := media_service.DeleteMediaFromDB(media); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete media from database"})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Media deleted successfully",
	})
}

func parseIDParam(c *fiber.Ctx) (int, error) {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format")
	}
	return id, nil
}

func parsePaginationParams(c *fiber.Ctx) (int, int) {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	return page, pageSize
}
func parseMediaRequest(c *fiber.Ctx) (*media_model.MediaSchema, error) {
	media := new(media_model.MediaSchema)
	if err := c.BodyParser(media); err != nil {
		return nil, fmt.Errorf("cannot parse request: %v", err)
	}
	return media, nil
}
