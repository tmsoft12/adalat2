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
	// CORS configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Public Routes
	app.Get("/", middleware.FakeUser, home.Home_Page)
	app.Get("/uploads/media/:filename", utils.Play)

	// Admin routes
	adminR := app.Group("/api/admin")
	adminR.Post("/banner", admin.CreateBanner)
	adminR.Post("/employer", admin.CreateEmployer)
	adminR.Post("/news", admin.CreateNews) // Changed route to avoid conflict
	adminR.Post("/media", admin.CreateMedia)
}
