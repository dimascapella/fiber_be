package router

import (
	"fiber_be/app/handlers"

	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {

	api := app.Group("/api")

	api.Get("/products", handlers.GetProducts)
	api.Post("/products", handlers.CreateProduct)
	api.Get("/products/:code", handlers.GetProductById)
	api.Put("/products/:code", handlers.UpdateProduct)
	api.Delete("/products/:code", handlers.DeleteProduct)

	api.Post("/post/note", handlers.PostNote)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "404: Not Found",
		})
	})
}
