package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"tm/controllers/admin/banner/models"
	utilsBanner "tm/controllers/admin/banner/utils"
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
		banners[i].Image = fmt.Sprintf("http://%s%s/api/home/%s", ip, port, banners[i].Image)
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
	banner.Image = fmt.Sprintf("http://%s%s/api/home/%s", ip, port, banner.Image)

	return c.Status(200).JSON(banner)
}

func CreateBanner(c *fiber.Ctx) error {
	uploadDir := "./uploads/banners/"

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Cannot create upload directory"})
	}

	filePath, err := utilsBanner.SaveFile(c, "bannerimg", uploadDir)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Cannot upload image"})
	}

	banner := new(models.BannerSchema)
	if err := c.BodyParser(banner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	banner.Image = filePath
	banner.ID = int(time.Now().Unix())
	banner.Link = c.FormValue("link")

	if err := config.DB.Create(banner).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create banner"})
	}

	return c.Status(fiber.StatusCreated).JSON(banner)
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
		fmt.Println("Invalid ID format:", err) // Terminal çıktısı
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var banner models.BannerSchema
	// Banner'ı veritabanından al
	if err := config.DB.First(&banner, id).Error; err != nil {
		fmt.Println("Banner not found:", err) // Terminal çıktısı
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Banner not found",
		})
	}

	// Yeni değerleri al
	if err := c.BodyParser(&banner); err != nil {
		fmt.Println("Cannot parse request:", err) // Terminal çıktısı
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse request",
		})
	}

	// "is_active" alanını güncelle
	if isActiveStr := c.FormValue("is_active"); isActiveStr != "" {
		banner.IsActive = isActiveStr == "true" // "true" string değerine göre boolean'a çevir
	}

	// Eğer yeni bir resim yükleniyorsa
	if file, err := c.FormFile("bannerimg"); err == nil {
		uploadDir := "./uploads/banners/"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			fmt.Println("Cannot create upload directory:", err) // Terminal çıktısı
			return c.Status(500).JSON(fiber.Map{"error": "Cannot create upload directory"})
		}

		// Eski resmi sil
		if banner.Image != "" {
			oldImagePath := filepath.Join(uploadDir, filepath.Base(banner.Image))
			fmt.Println("Deleting old image at:", oldImagePath) // Terminal çıktısı
			if err := os.Remove(oldImagePath); err != nil {
				fmt.Println("Cannot delete old image:", err) // Terminal çıktısı
			}
		}

		// Yeni resmi yükle ve dosya yolunu güncelle
		filePath := filepath.Join(uploadDir, file.Filename)
		if err := c.SaveFile(file, filePath); err != nil {
			fmt.Println("Cannot upload image:", err)
			return c.Status(500).JSON(fiber.Map{"error": "Cannot upload image"})
		}
		banner.Image = filePath // Veritabanı için yeni dosya yolunu güncelle
		fmt.Println("New image uploaded at:", filePath)
	}

	if err := config.DB.Save(&banner).Error; err != nil {
		fmt.Println("Cannot update banner:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot update banner"})
	}

	fmt.Println("Banner updated successfully:", banner)
	return c.Status(fiber.StatusOK).JSON(banner)
}
