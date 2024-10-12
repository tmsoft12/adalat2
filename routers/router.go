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
	homeP.Get("/:id", middleware.FakeUser, news_page.NewsDetail)
	homeP.Get("media/:id", media.MediaDetail)

	bannerAdmin := app.Group("/api/admin/banner")
	bannerAdmin.Get("/", controller.GetAllBanner)       //ok
	bannerAdmin.Get("/:id", controller.GetBannerById)   //ok
	bannerAdmin.Post("/", controller.CreateBanner)      //ok
	bannerAdmin.Delete("/:id", controller.DeleteBanner) //ok
	bannerAdmin.Put("/:id", controller.UpdateBanner)    //ok

	mediaAdmin := app.Group("/api/admin/media")
	mediaAdmin.Get("/", media_controller.GetAllMedia)       //ok
	mediaAdmin.Get("/:id", media_controller.GetById)        //ok
	mediaAdmin.Post("/", media_controller.CreateMedia)      //ok
	mediaAdmin.Delete("/:id", media_controller.DeleteMedia) //ok
	mediaAdmin.Put("/:id", media_controller.UpdateMedia)    //ok

	news := app.Group("api/news")
	news.Get("/", news_controller.GetAllNews)
	news.Get("/:id", news_controller.GetByIdNews)
	news.Post("/", news_controller.CreateNews)
	news.Delete("/:id", news_controller.DeleteNew)
	news.Put("/:id", news_controller.UpdateNews)

	employerAdmin := app.Group("/api/admin/employers")
	employerAdmin.Get("/", employer_contrllers.GetAllEmployers)
	employerAdmin.Post("/", employer_contrllers.CreateEmployer)
	employerAdmin.Delete("/:id", employer_contrllers.DeleteEmployer)
	employerAdmin.Put("/:id", employer_contrllers.UpdateEmployer)
	employerAdmin.Get("/:id", employer_contrllers.GetByIdEmployer)

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
