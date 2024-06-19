package main

import (
	"fiber_be/database"
	"fiber_be/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	app := fiber.New()
	router.Init(app)

	app.Listen(":3000")
}
