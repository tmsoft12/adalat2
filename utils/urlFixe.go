package utils

import (
	"fmt"
	"os"
	model "tm/models"
)

func UrlCom(media []model.MediaSchema) {
	ip := os.Getenv("HOST")
	port := os.Getenv("PORT")
	for i := range media {
		media[i].Video = fmt.Sprintf("http://%s%s/api/home/%s", ip, port, media[i].Video)
	}
}
