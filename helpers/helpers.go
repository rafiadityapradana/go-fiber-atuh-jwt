package helpers

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/restapi_fiber/config"
	"github.com/restapi_fiber/models"
)


type ReqLogin struct {
	Username string `validate:"required,min=6"`
	Password string `validate:"required,min=6"`
}
type ErrorResponseLogin struct {
	FailedField string
	Message       string
}

func ValidateStruct(req ReqLogin) []*ErrorResponseLogin {
	var errors []*ErrorResponseLogin
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var data = []string{"This",err.Field(), "must be", err.Tag(), err.Param()}
			var element ErrorResponseLogin
			element.FailedField = err.Field()	
			element.Message = strings.Join(data, " ")
			errors = append(errors, &element)
		}
	}
	return errors
}
type CustomClaims struct {
	UserData models.Users
	jwt.StandardClaims
}
func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
func GenerateAccessToken(user models.Users) string {
	env:= goDotEnvVariable("TOKEN_SCRET")
	TokenScret := []byte(env)
	Claims:= jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{Issuer: user.UserId,
		ExpiresAt:  time.Now().Add(time.Hour * 15).Unix() })
	Token,err := Claims.SignedString(TokenScret)
	if err != nil{
		return ""
	}	
	return Token
}
func GenerateRefreshToken (user models.Users) string{
	Env:= goDotEnvVariable("TOKEN_SCRET_REFT")
	TokenScretReft:= []byte(Env)
	ClaimsReft:= CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    user.UserId,
		},
	}
	CreateTokenReft := jwt.NewWithClaims(jwt.SigningMethodHS256, ClaimsReft)
	TokenReft,err := CreateTokenReft.SignedString([]byte(TokenScretReft))
	if err != nil{
		return ""
	}
	return TokenReft
}
func DecodeTokenIssuerUserId (c *fiber.Ctx)string {
	cookie := c.Cookies("GfAtID")
	if cookie == ""{
		return ""
	}
	env:= goDotEnvVariable("TOKEN_SCRET")
	Token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(env), nil
	})
	if err != nil {
		return ""
	}
	Isclaims:= Token.Claims.(*jwt.StandardClaims)
	var user models.Users 
	result := config.DB.Where("user_id = ?",Isclaims.Issuer).First(&user)
	if result.RowsAffected > 0 {
		return Isclaims.Issuer
	}else{
		return ""
	}
}

