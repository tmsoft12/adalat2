package admin

import (
	"fmt"
	"time"
	config "tm/db"
	model "tm/models"
	"tm/utils"

	"github.com/gofiber/fiber/v2"
)

func CreateNews(c *fiber.Ctx) error {
	filePath, err := utils.SaveFile(c, "image", "./uploads")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot upload image"})
	}

	news := new(model.NewsSchema)
	if err := c.BodyParser(news); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	news.Image = filePath
	news.ID = fmt.Sprintf("%d", time.Now().Unix())
	news.Date = time.Now().Format("2006-01-02")

	if err := config.DB.Create(news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create news"})
	}

	return c.Status(fiber.StatusCreated).JSON(news)
}

func CreateBanner(c *fiber.Ctx) error {
	filePath, err := utils.SaveFile(c, "bannerimg", "./uploads/banners/")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot upload image"})
	}

	banners := new(model.BannerSchema)
	if err := c.BodyParser(banners); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	banners.Image = filePath
	banners.ID = fmt.Sprintf("%d", time.Now().Unix())
	banners.Title = c.FormValue("title")
	banners.Description = c.FormValue("description")
	banners.Link = c.FormValue("link")

	if err := config.DB.Create(banners).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create banner"})
	}

	return c.Status(fiber.StatusCreated).JSON(banners)
}

func CreateEmployer(c *fiber.Ctx) error {
	filePath, err := utils.SaveFile(c, "image", "./uploads/employers")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot upload image"})
	}

	// Employer verilerini body'den alma
	employer := new(model.EmployerSchema)
	if err := c.BodyParser(employer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	// Employer verilerini ayarlama
	employer.Image = filePath
	employer.ID = fmt.Sprintf("%d", time.Now().Unix())

	// Employer veritabanına kaydetme
	if err := config.DB.Create(employer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create employer"})
	}

	// Başarılı olduğunda yanıt
	return c.Status(fiber.StatusCreated).JSON(employer)
}
func CreateMedia(c *fiber.Ctx) error {
	filePath, err := utils.SaveFile(c, "video", "./uploads/media")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot upload video file"})
	}

	media := new(model.MediaSchema)
	if err := c.BodyParser(media); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	// Title alanını kontrol et
	if media.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Title is required"})
	}
	media.Video = filePath
	media.ID = fmt.Sprintf("%d", time.Now().Unix())
	media.Date = time.Now().Format("2006-01-02")

	if err := config.DB.Create(media).Error; err != nil {
		// Hata loglama
		fmt.Println("Database error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create media"})
	}

	return c.Status(fiber.StatusCreated).JSON(media)
}
