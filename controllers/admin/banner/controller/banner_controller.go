package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"tm/controllers/admin/banner/models"
	config "tm/db"

	"github.com/gofiber/fiber/v2"
)

func GetAllBanner(c *fiber.Ctx) error {
	var banners []models.BannerSchema
	if err := config.DB.Find(&banners).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server internal error"})
	}

	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	for i := range banners {
		banners[i].Image = fmt.Sprintf("http://%s%s/%s", ip, port, banners[i].Image)
	}

	return c.Status(200).JSON(banners)
}

func GetBannerById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var banner models.BannerSchema
	if err := config.DB.First(&banner, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Banner not found",
		})
	}

	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	banner.Image = fmt.Sprintf("http://%s%s/%s", ip, port, banner.Image)

	return c.Status(200).JSON(banner)
}

func CreateBanner(c *fiber.Ctx) error {
	uploadDir := "./uploads/banners/"

	if err := createUploadDir(uploadDir); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Cannot create upload directory"})
	}

	filePath, err := uploadBannerFile(c, uploadDir)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Cannot upload image"})
	}

	banner, err := parseBannerRequest(c, filePath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if err := saveBannerToDB(banner); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create banner"})
	}

	return c.Status(fiber.StatusCreated).JSON(banner)
}

func createUploadDir(uploadDir string) error {
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func uploadBannerFile(c *fiber.Ctx, uploadDir string) (string, error) {
	file, err := c.FormFile("bannerimg")
	if err != nil {
		return "", err
	}

	// Generate a unique file name
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath := filepath.Join(uploadDir, fileName)

	if err := c.SaveFile(file, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

func parseBannerRequest(c *fiber.Ctx, filePath string) (*models.BannerSchema, error) {
	banner := new(models.BannerSchema)
	if err := c.BodyParser(banner); err != nil {
		return nil, err
	}

	banner.Image = filePath
	banner.Link = c.FormValue("link")
	banner.IsActive = true

	return banner, nil
}

// saveBannerToDB saves the banner to the database
func saveBannerToDB(banner *models.BannerSchema) error {
	if err := config.DB.Create(banner).Error; err != nil {
		return err
	}
	return nil
}

func DeleteBanner(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var banner models.BannerSchema
	if err := config.DB.First(&banner, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Banner not found",
		})
	}

	if banner.Image != "" {
		if err := os.Remove(banner.Image); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete the image file"})
		}
	}

	if err := config.DB.Delete(&banner).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete banner from database",
		})
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Banner deleted successfully",
	})
}

func UpdateBanner(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var banner models.BannerSchema
	if err := config.DB.First(&banner, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Banner not found",
		})
	}

	if err := c.BodyParser(&banner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
		})
	}

	if isActiveStr := c.FormValue("is_active"); isActiveStr != "" {
		banner.IsActive = isActiveStr == "true"
	}

	if file, err := c.FormFile("bannerimg"); err == nil {
		uploadDir := "./uploads/banners/"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Cannot create upload directory"})
		}

		if banner.Image != "" {
			oldImagePath := banner.Image // Change to use the stored path directly
			if err := os.Remove(oldImagePath); err != nil {
				fmt.Println("Cannot delete old image:", err)
			}
		}

		// Generate a unique file name
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		filePath := filepath.Join(uploadDir, fileName)
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Cannot upload image"})
		}
		banner.Image = filePath
	}

	if err := config.DB.Save(&banner).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update banner"})
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	banner.Image = fmt.Sprintf("http://%s%s/%s", ip, port, banner.Image)

	return c.Status(fiber.StatusOK).JSON(banner)
}
