package routes

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/restapi_fiber/AuthMiddleware"
	"github.com/restapi_fiber/controllers"
)
func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
func Setup(app *fiber.App) {
	Env := goDotEnvVariable("API_KEY")
	api := app.Group("/api", AuthMiddleware.AuthApi(AuthMiddleware.Config{Key:Env}))
	api.Post("/login", controllers.Login)
	api.Get("/auth-data",AuthMiddleware.AuthAuthorization(), controllers.AuthData)
	api.Post("/refresh-token",AuthMiddleware.AuthAuthorizationReft(), controllers.AuthReftToken)
	
}