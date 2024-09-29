package routers

import (
	"tm/controllers/admin"
	"tm/controllers/home"
	"tm/middleware"
	"tm/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitRouters(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/", middleware.FakeUser, home.Home_Page)
	app.Get("/uploads/media/:filename", utils.Play)
	app.Static("/uploads", "./uploads")

	adminR := app.Group("/api/admin")
	adminR.Post("/banner", admin.CreateBanner)
	adminR.Post("/employer", admin.CreateEmployer)
	adminR.Post("/news", admin.CreateNews)
	adminR.Post("/media", admin.CreateMedia)
}
