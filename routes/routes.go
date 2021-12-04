package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/restapi_fiber/AuthMiddleware"
	"github.com/restapi_fiber/controllers"
)
func Setup(app *fiber.App) {
	app.Post("/api/login", controllers.Login)
	app.Use(AuthMiddleware.AuthAuthorization())
	app.Get("/api/auth-data", controllers.AuthData)
	
}