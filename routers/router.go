package routers

import (
	"tm/controllers/admin"
	"tm/controllers/admin/login"
	"tm/controllers/home"
	"tm/controllers/media"
	"tm/controllers/news"
	"tm/middleware"
	admin_middleware "tm/middleware/admin"
	"tm/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitRouters(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:5500", // İzin verilen kökeni buraya yazın
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true, // Kimlik bilgilerini dahil etmek için true olarak ayarlayın
	}))

	loginP := app.Group("api/auth")
	loginP.Post("/login", login.Login)
	loginP.Get("/protected", login.Protected, func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uint)
		return c.JSON(fiber.Map{
			"message": "This is a protected route",
			"userID":  userID,
		})
	})
	loginP.Post("/refresh", login.Refresh)

	homeP := app.Group("api/home")
	homeP.Get("/", middleware.FakeUser, home.Home_Page)
	homeP.Get("/uploads/media/:filename", utils.Play)
	homeP.Static("/uploads", "./uploads")

	newsP := app.Group("/api/news")
	newsP.Get("/", news.GetAllNews)
	newsP.Get("/:id", middleware.FakeUser, news.NewsDetail)

	mediaP := app.Group("api/media")
	mediaP.Get("/", media.GetAllMedia)

	adminR := app.Group("/api/admin", admin_middleware.Protected)
	adminR.Post("/banner", admin.CreateBanner)
	adminR.Post("/employer", admin.CreateEmployer)
	adminR.Post("/news", admin.CreateNews)
	adminR.Post("/media", admin.CreateMedia)

}
