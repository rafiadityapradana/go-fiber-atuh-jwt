package config

import (
	"github.com/restapi_fiber/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
var DB *gorm.DB
func ConnectionDataBase () {
	Connection, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "user=postgres password=root dbname=go_fiber port=5432 sslmode=disable",
		PreferSimpleProtocol: true,
	  }), &gorm.Config{})
	if err != nil {
		panic("Connection Databse Failed")
	}
	DB= Connection
	Connection.AutoMigrate(&models.Users{},&models.AuthUserTokens{})
}
