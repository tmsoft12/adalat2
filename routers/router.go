package routers

import (
	"tm/controllers/admin"
	"tm/controllers/admin/banner/controller"
	media_controller "tm/controllers/admin/media/controller"
	"tm/controllers/home"
	"tm/controllers/news"
	"tm/middleware"
	admin_middleware "tm/middleware/admin"
	"tm/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitRouters(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	homeP := app.Group("api/home")
	homeP.Get("/", middleware.FakeUser, home.Home_Page)
	homeP.Get("/uploads/media/:filename", utils.Play)
	homeP.Static("/uploads", "./uploads")

	banner := app.Group("/api/banner")
	banner.Get("/all", controller.GetAllBanner)
	banner.Get("/by/:id", controller.GetBannerById)
	banner.Post("/create", controller.CreateBanner)
	banner.Delete("/delete/:id", controller.DeleteBanner)
	banner.Put("/update/:id", controller.UpdateBanner)

	mediaP := app.Group("/api/media")
	mediaP.Get("/all", media_controller.GetAllMedia)
	mediaP.Get("/byId/:id", media_controller.GetById)
	mediaP.Post("/create", media_controller.CreateMedia)
	mediaP.Delete("/delete/:id", media_controller.DeleteMedia)
	mediaP.Put("/:id", media_controller.UpdateMedia)

	newsP := app.Group("/api/news")
	newsP.Get("/", news.GetAllNews)
	newsP.Get("/:id", middleware.FakeUser, news.NewsDetail)

	adminR := app.Group("/api/admin", admin_middleware.Protected)
	adminR.Post("/banner", admin.CreateBanner)
	adminR.Post("/employer", admin.CreateEmployer)
	adminR.Post("/news", admin.CreateNews)
	adminR.Post("/media", admin.CreateMedia)

}
