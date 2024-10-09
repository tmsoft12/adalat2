package routers

import (
	"tm/controllers/admin"
	"tm/controllers/admin/banner/controller"
	employer_contrllers "tm/controllers/admin/employers/controllers"
	media_controller "tm/controllers/admin/media/controller"
	news_controller "tm/controllers/admin/news/controller"
	"tm/controllers/home"
	"tm/controllers/media"
	news_page "tm/controllers/news"
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
	uploads := app.Group("/")
	uploads.Get("/uploads/media/:filename", utils.Play)
	uploads.Static("/uploads", "./uploads")

	homeP := app.Group("api/home")
	homeP.Get("/", middleware.FakeUser, home.Home_Page)

	banner := app.Group("/api/banner")
	banner.Get("/", controller.GetAllBanner)
	banner.Get("/:id", controller.GetBannerById)
	banner.Post("/", controller.CreateBanner)
	banner.Delete("/:id", controller.DeleteBanner)
	banner.Put("/:id", controller.UpdateBanner)

	mediaP := app.Group("/api/media")
	mediaP.Get("/", media_controller.GetAllMedia)
	mediaP.Get("/:id", media_controller.GetById)
	mediaP.Post("/", media_controller.CreateMedia)
	mediaP.Delete("/:id", media_controller.DeleteMedia)
	mediaP.Put("/:id", media_controller.UpdateMedia)

	news := app.Group("api/news")
	news.Get("/", news_controller.GetAllNews)
	news.Get("/:id", news_controller.GetByIdNews)
	news.Post("/", news_controller.CreateNews)
	news.Delete("/:id", news_controller.DeleteNew)
	news.Put("/:id", news_controller.UpdateNews)

	employerP := app.Group("/api/employer")
	employerP.Get("/", employer_contrllers.GetAllEmployers)
	employerP.Post("/", employer_contrllers.CreateEmployer)
	employerP.Delete("/:id", employer_contrllers.DeleteEmployer)
	employerP.Put("/:id", employer_contrllers.UpdateEmployer)
	employerP.Get("/:id", employer_contrllers.GetByIdEmployer)

	newsP := app.Group("/api/home")
	newsP.Get("/", home.Home_Page)
	newsP.Get("/", news_page.GetAllNews)

	adminR := app.Group("/api/admin", admin_middleware.Protected)
	adminR.Post("/banner", admin.CreateBanner)
	adminR.Post("/employer", admin.CreateEmployer)
	adminR.Post("/news", admin.CreateNews)
	adminR.Post("/media", admin.CreateMedia)

	app.Get("/:id", middleware.FakeUser, media.MediaDetail)
	app.Get("/:id", middleware.FakeUser, news_page.NewsDetail)

}
