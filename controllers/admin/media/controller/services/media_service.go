package media_service

import (
	"fmt"
	media_model "tm/controllers/admin/media/models"
	config "tm/db"
)

func GetAllMedia(page int, pageSize int) ([]media_model.MediaSchema, int64, error) {
	var media []media_model.MediaSchema
	var total int64

	if err := config.DB.Model(&media_model.MediaSchema{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := config.DB.Offset((page - 1) * pageSize).Limit(pageSize).Find(&media).Error; err != nil {
		return nil, 0, err
	}

	return media, total, nil
}
func SaveMedia(media *media_model.MediaSchema) error {
	if err := config.DB.Create(media).Error; err != nil {
		return fmt.Errorf("cannot create media: %v", err)
	}
	return nil
}

func UpdateMedia(media *media_model.MediaSchema) error {
	if err := config.DB.Save(media).Error; err != nil {
		return fmt.Errorf("cannot update media: %v", err)
	}
	return nil
}

func FindMediaById(id int) (*media_model.MediaSchema, error) {
	var media media_model.MediaSchema
	if err := config.DB.First(&media, id).Error; err != nil {
		return nil, fmt.Errorf("media not found")
	}
	return &media, nil
}

func DeleteMediaFromDB(media *media_model.MediaSchema) error {
	if err := config.DB.Delete(media).Error; err != nil {
		return fmt.Errorf("cannot delete media from database: %v", err)
	}
	return nil
}
