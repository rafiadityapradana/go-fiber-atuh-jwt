package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/restapi_fiber/AuthMiddleware"
	"github.com/restapi_fiber/config"
	"github.com/restapi_fiber/routes"
)

func main() {
	
	config.ConnectionDataBase()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	app.Use(AuthMiddleware.AuthApi(AuthMiddleware.Config{Key:"APP_KEY06e91258-3238-4463-b04e-9a900a17f744"}))
	routes.Setup(app)
    app.Listen(":8070")
}	