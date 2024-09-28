package routers

import (
	"tm/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitRouters(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	newsRoutes := app.Group("/news")

	newsRoutes.Post("/", controllers.CreateNews)
	newsRoutes.Post("/emp", controllers.CreateEmployer)
	newsRoutes.Get("/", controllers.Home_Page)
	app.Get("/uploads/media/:filename", controllers.Play)

	newsRoutes.Post("/banner", controllers.CreateBanner)
	newsRoutes.Post("/media", controllers.CreateMedia)
}
