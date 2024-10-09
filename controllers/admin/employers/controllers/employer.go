package employer_contrllers

import (
	"fmt"
	"os"
	"strconv"
	employer_models "tm/controllers/admin/employers/models"
	employer_utils "tm/controllers/admin/employers/utils"
	config "tm/db"

	"github.com/gofiber/fiber/v2"
)

func GetAllEmployers(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.Query("pageSize", "10"))
	if err != nil || pageSize < 1 {
		page = 10
	}
	var employer []employer_models.EmployerSchema
	var total int64
	if err := config.DB.Model(&employer_models.EmployerSchema{}).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server internal error"})
	}
	if err := config.DB.Offset((page - 1) * page).Limit(pageSize).Find(&employer).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Server internal error"})
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	for i := range employer {
		employer[i].Image = fmt.Sprintf("http://%s%s/%s", ip, port, employer[i].Image)
	}

	return c.Status(200).JSON(fiber.Map{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"employer": employer,
	})
}

func CreateEmployer(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}
	var employer employer_models.EmployerSchema
	if err := config.DB.First(&employer, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Employer not found",
		})
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	employer.Image = fmt.Sprintf("http://%s%s/%s", ip, port, employer.Image)

	return c.Status(200).JSON(employer)
}

func GetByIdEmployer(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}
	var employer employer_models.EmployerSchema
	if err := config.DB.First(&employer, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Employer not found",
		})
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")

	employer.Image = fmt.Sprintf("http://%s%s/%s", ip, port, employer.Image)

	return c.Status(200).JSON(employer)

}

func handleFileUpload(c *fiber.Ctx) (string, error) {
	filePath, err := employer_utils.SaveFile(c, "image", "./uploads/employers")
	if err != nil {
		return "", fmt.Errorf("cannot upload image file: %v", err)
	}
	return filePath, nil
}
func UpdateEmployer(c *fiber.Ctx) error {
	idParm := c.Params("id")
	id, err := strconv.Atoi(idParm)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}
	var employer employer_models.EmployerSchema
	if err := config.DB.First(&employer, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Employer not found",
		})
	}
	if newFilePath, err := handleFileUpload(c); err == nil {
		if employer.Image != "" {
			if _, err := os.Stat(employer.Image); err == nil {
				if err := os.Remove(employer.Image); err != nil {
					fmt.Println("Error deleting old file:", err)

				}
			}
		}
		employer.Image = newFilePath
	}
	if err := config.DB.Save(&employer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update new",
		})
	}
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	employer.Image = fmt.Sprintf("http://%s%s/%s", ip, port, employer.Image)
	return c.Status(fiber.StatusOK).JSON(employer)

}

func DeleteEmployer(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}
	var employer employer_models.EmployerSchema
	if err := config.DB.First(&employer, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "employer not found",
		})
	}
	if employer.Image != "" {
		if _, err := os.Stat(employer.Image); os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File not found"})
		}
		if err := os.Remove(employer.Image); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete the employer file"})
		}
	}
	if err := config.DB.Delete(&employer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete employer from database",
		})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
		"message": "Employer deleted successfully",
	})
}
