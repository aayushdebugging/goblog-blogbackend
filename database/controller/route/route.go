package route

import (
	"github.com/aayushdebugging/blogbackend/database/controller"
	"github.com/aayushdebugging/blogbackend/middleware"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controller.Register)
	app.Post("/api/login", controller.Login)

	app.Use(middleware.IsAuthenticate)
	app.Post("/api/createpost", controller.CreatePost)
	app.Get("/api/allpost", controller.AllPost)
	app.Get("/api/detailpost/:id", controller.DetailPost)
	app.Put("/api/updatepost/:id", controller.UpdatePost)
	app.Get("/api/uniquepost", controller.UniquePost)
	app.Delete("/api/deletepost/:id", controller.DeletePost)
	app.Post("/api/uploads-images",controller.Upload)
	app.Static("/api/uploads" , "./uploads")
}
