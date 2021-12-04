package AuthMiddleware

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/restapi_fiber/config"
	"github.com/restapi_fiber/models"
)
const (DefaultHeaderKeyIdentifier string = "Go-Key")

type Config struct {
	Skip func(*fiber.Ctx) bool
	Key string 
	ValidatorFunc func(*fiber.Ctx, Config) bool
}

func DefaultValidatorFunc(c *fiber.Ctx, cfg Config) bool {
	headerKey := c.Get(DefaultHeaderKeyIdentifier)
	if headerKey == "" {
		return false
	}
	if headerKey != "" && headerKey == cfg.Key {
		return true
	}
	return false
}

var defaultConfig = Config{
	ValidatorFunc: DefaultValidatorFunc,
}
func AuthApi(config ...Config) fiber.Handler {
	var cfg Config
	if len(config) == 0 {
		cfg = defaultConfig
	} else {
		cfg = config[0]
		if cfg.ValidatorFunc == nil {
			cfg.ValidatorFunc = DefaultValidatorFunc
		}
	}
	return func(c *fiber.Ctx) error {
		if cfg.Skip != nil && cfg.Skip(c) {
			return c.Next()
		}
		pass := cfg.ValidatorFunc(c, cfg)
		if !pass {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message":"Invalid Go Key",
			})
		}
		return c.Next()
	}
}


const (
	DefaultHeaderAuththentication string = "Authorization"
)
type ConfigAuthorization struct {
	Skip          func(*fiber.Ctx) bool
	Key           string
	ValidatorFunc func(*fiber.Ctx, ConfigAuthorization) bool
}
func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
func DefaultValidatorAuthorizationFunc(c *fiber.Ctx, cfg ConfigAuthorization) bool {
	headerKey := c.Get(DefaultHeaderAuththentication)
	cookie := c.Cookies("GfAtID")
	if headerKey == "" {
		return false
	}
	if cookie == "" { 
		return false
	}
	if headerKey != "" && headerKey == cookie {
		var ENV = goDotEnvVariable("TOKEN_SCRET")
		Token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(ENV), nil
		})
		if err != nil {
			return false
		}
		Isclaims:= Token.Claims.(*jwt.StandardClaims)
		var user models.Users 
		result := config.DB.Where("user_id = ?",Isclaims.Issuer).First(&user)
		if result.RowsAffected > 0 {
			return true
		}else{
			return false
		}
	}
	return false
}
var defaultConfigauthorization = ConfigAuthorization{
	ValidatorFunc: DefaultValidatorAuthorizationFunc,
}
func AuthAuthorization(config ...ConfigAuthorization) fiber.Handler {
	var cfg ConfigAuthorization
	if len(config) == 0 {
		cfg = defaultConfigauthorization
	} else {
		cfg = config[0]
		if cfg.ValidatorFunc == nil {
			cfg.ValidatorFunc = DefaultValidatorAuthorizationFunc
		}
	}
	return func(c *fiber.Ctx) error {
		pass := cfg.ValidatorFunc(c, cfg)
		if !pass {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message":"Unauthorized",
			})
		}
		store := session.New(session.Config{
			Expiration:     15 * time.Hour,
			Storage:        nil,
			KeyLookup:      "cookie:GfSID",
			CookieDomain:   "",
			CookiePath:     "",
			CookieSecure:   false,
			CookieHTTPOnly: true,
			CookieSameSite: "",
			KeyGenerator:   utils.UUIDv4,
			CookieName:     "",
		})
		SessionStore, err := store.Get(c)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message":" Could not login !",
			})
		}
		if SessionStore.ID() != "" && SessionStore.ID() == c.Cookies("GfSID"){
			return c.Next()
		}else{
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message":"Unauthorized",
			})
		}	
	}
}
