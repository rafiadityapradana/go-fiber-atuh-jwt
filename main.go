package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/restapi_fiber/AuthMiddleware"
	"github.com/restapi_fiber/config"
	"github.com/restapi_fiber/routes"
)
func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
func main() {
	config.ConnectionDataBase()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	Env := goDotEnvVariable("API_KEY")
	app.Use(AuthMiddleware.AuthApi(AuthMiddleware.Config{Key:Env}))
	routes.Setup(app)
	PORT:=goDotEnvVariable("APP_PORT")
    app.Listen(PORT)
}	