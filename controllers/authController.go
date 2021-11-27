package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/restapi_fiber/config"
	"github.com/restapi_fiber/helpers"
	"github.com/restapi_fiber/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	TokenScret := []byte("TokenScret")
 	TokenScretReft:= []byte("TokenScretReft")
	Req := new(helpers.ReqLogin)
	if err := c.BodyParser(Req); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": err.Error(),
        })
    }
	errors := helpers.ValidateStruct(*Req)
    if errors != nil {
       return c.JSON(errors)
    }
	var user models.Users 
	result := config.DB.Where("username = ?",Req.Username).Or("email = ? ",Req.Username).First(&user)
	if result.RowsAffected > 0 {
		if err := bcrypt.CompareHashAndPassword(user.Password, []byte(Req.Password)); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message":"Incorrect password !",
			})
		}else{
			type CustomClaims struct {
				UserData models.Users
				jwt.StandardClaims
			}
			Claims:= jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{Issuer: user.UserId,
				ExpiresAt:  time.Now().Add(time.Hour * 5).Unix() })
			ClaimsReft:= CustomClaims{
				user,
				jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour * 15).Unix(),
					Issuer:    user.UserId,
				},
			}
			CreateTokenReft := jwt.NewWithClaims(jwt.SigningMethodHS256, ClaimsReft)
			Token,errToken := Claims.SignedString(TokenScret) 
			TokenReft, errTokenReft := CreateTokenReft.SignedString([]byte(TokenScretReft))
			if errToken != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(fiber.Map{
					"message":"Could not login !",
				})
			}else {
				if errTokenReft != nil {
					c.Status(fiber.StatusInternalServerError)
					return c.JSON(fiber.Map{
						"message":"Reft Could not login !",
					})
				}else{
					var GetToken  models.AuthUserTokens
					OldToken := config.DB.Where("user_id = ?",user.UserId).First(&GetToken)
					if OldToken.RowsAffected > 0 {
						config.DB.Model(models.AuthUserTokens{}).Where("user_id = ?",  user.UserId).Updates(models.AuthUserTokens{AccessToken: Token,RefeshToken:TokenReft})
					}else{
						NewToken := models.AuthUserTokens{TokenId:uuid.New().String(), UserId:user.UserId, AccessToken: Token, RefeshToken: TokenReft}
						 config.DB.Create(&NewToken)
					}
					cookie := fiber.Cookie{
						Name: "GFATID",
						Value: Token,
						Expires: time.Now().Add(time.Hour * 15),
						HTTPOnly: true,
						SameSite: "lax",
						Secure: true,
					}
					cookieReft := fiber.Cookie{
						Name: "GFRTID",
						Value: TokenReft,
						Expires: time.Now().Add(time.Hour * 24),
						HTTPOnly: true,
						SameSite: "lax",
						Secure: true,
					}
					c.Cookie(&cookie)
					c.Cookie(&cookieReft)
					c.Status(fiber.StatusOK)
					return c.JSON(fiber.Map{
						"AccessToken":Token,
						"RefeshToken":TokenReft,
					})
				}
				
			}
		}
	}else{
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message":"Usename not found !",
		})
	}
}		