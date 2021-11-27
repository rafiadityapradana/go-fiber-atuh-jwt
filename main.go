package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/restapi_fiber/config"
	"github.com/restapi_fiber/routes"
)
func main() {
	config.ConnectionDataBase()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	routes.Setup(app)
    app.Listen(":8070")
}	