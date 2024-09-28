package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	config "tm/db"
	model "tm/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/rand"
)

func saveFile(c *fiber.Ctx, fieldName, dir string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(file.Filename)

	newFileName := fmt.Sprintf("%d_%d%s", time.Now().Unix(), rand.Intn(1000), ext)

	filePath := fmt.Sprintf("%s/%s", dir, newFileName)

	if err := c.SaveFile(file, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

func CreateNews(c *fiber.Ctx) error {
	filePath, err := saveFile(c, "image", "./uploads")
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
	filePath, err := saveFile(c, "bannerimg", "./uploads/banners/")
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
	// Dosya yükleme işlemi
	filePath, err := saveFile(c, "image", "./uploads/employers")
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
	filePath, err := saveFile(c, "video", "./uploads/media")
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
func Play(c *fiber.Ctx) error {
	fileName := c.Params("filename")
	filePath := filepath.Join("./uploads/media", fileName)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File not found"})
	}

	return c.SendFile(filePath, false)
}

func Home_Page(c *fiber.Ctx) error {
	username := c.Cookies("id")

	if username == "" {
		cookie := new(fiber.Cookie)
		cookie.Name = "id"
		cookie.Value = "kerim"
		cookie.Expires = time.Now().Add(24 * time.Hour)

		// Cookie'yi ekle
		c.Cookie(cookie)
		return c.SendString(cookie.Value)

	}

	var news []model.NewsSchema
	if err := config.DB.Find(&news).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch news"})
	}

	var banner []model.BannerSchema
	if err := config.DB.Find(&banner).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch banner"})
	}

	var media []model.MediaSchema
	if err := config.DB.Find(&media).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch media"})
	}

	var employ []model.EmployerSchema
	if err := config.DB.Find(&employ).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch employ"})
	}

	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	for i := range media {
		media[i].Video = fmt.Sprintf("http://%s%s/%s", ip, port, media[i].Video)
	}

	return c.JSON(fiber.Map{
		"news":   news,
		"banner": banner,
		"media":  media,
		"employ": employ,
	})
}
