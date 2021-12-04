package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/google/uuid"
	"github.com/restapi_fiber/config"
	"github.com/restapi_fiber/helpers"
	"github.com/restapi_fiber/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	Req := new(helpers.ReqLogin)
	if err := c.BodyParser(Req); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
    }
	errors := helpers.ValidateStruct(*Req)
    if errors != nil {return c.JSON(errors)}
	var user models.Users 
	result := config.DB.Where("username = ?",Req.Username).First(&user)
	if result.RowsAffected > 0 {
		if err := bcrypt.CompareHashAndPassword(user.Password, []byte(Req.Password)); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"message":"Incorrect password !"})
		}else{
			var Token = helpers.GenerateAccessToken(user)
			var TokenReft = helpers.GenerateRefreshToken(user)
			var GetToken  models.AuthUserTokens
			OldToken := config.DB.Where("id_user = ?",user.UserId).First(&GetToken)
			if OldToken.RowsAffected > 0 {
				config.DB.Model(models.AuthUserTokens{}).Where("user_id = ?",  user.UserId).Updates(models.AuthUserTokens{AccessToken: Token,RefeshToken:TokenReft})
			}else{
				NewToken := models.AuthUserTokens{TokenId:uuid.New().String(), IdUser:user.UserId, AccessToken: Token, RefeshToken: TokenReft}
				config.DB.Create(&NewToken)
			}
			store := session.New(session.Config{
				Expiration:     15 * time.Hour,
				Storage:        nil,
				KeyLookup:      "cookie:GfSID",
				CookieDomain:   "",
				CookiePath:     "",
				CookieSecure:   true,
				CookieHTTPOnly: true,
				CookieSameSite: "",
				KeyGenerator:   utils.UUIDv4,
				CookieName:     "",
			})
			cookie := fiber.Cookie{
				Name:     "GfAtID",
				Value:    Token,
				Path:     "",
				Domain:   "",
				MaxAge:   0,
				Expires:  time.Now().Add(time.Hour * 15),
				Secure:   true,
				HTTPOnly: true,
				SameSite: "lax",
			}
			SessionStore, err := store.Get(c)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{"message":" Could not login !"})
			}
			if errSaveSession := SessionStore.Save(); errSaveSession != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{"message":errSaveSession})
			}
			SessionStore.Set("GfSID", Token)
			c.Cookie(&cookie)			
			c.Status(fiber.StatusOK)
			return c.JSON(fiber.Map{
				"AccessToken":Token,
				"RefeshToken":TokenReft,
			})		
		}	
	}else{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{"message":"Usename not found !"})
	}
}		