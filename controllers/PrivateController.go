package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/restapi_fiber/config"
	"github.com/restapi_fiber/helpers"
	"github.com/restapi_fiber/models"
)

func AuthData(c *fiber.Ctx ) error {
	UID := helpers.DecodeTokenIssuerUserId(c)
	var user models.Users 
	Result :=config.DB.Joins("Roles").Where("user_id = ?",UID).Find(&user)
	if Result.RowsAffected > 0 {
		var GetToken  models.AuthUserTokens
		Token := config.DB.Where("id_user = ?",user.UserId).First(&GetToken)
		if Token.RowsAffected > 0 {
			c.Status(fiber.StatusOK)
			return c.JSON(fiber.Map{
				"Data":user,
			})
		}
		c.ClearCookie("GfAtID", "GfSID")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message":"Could not login !",
		})
	}else{
		c.ClearCookie("GfAtID", "GfSID")
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message":"Could not login !",
		})
	}
}