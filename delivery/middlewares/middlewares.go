package middlewares

import (
	_config "eventapp/config"

	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	Id   int
	Role string
}

// Custom Middleware to handle JSON Web Token
func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		// Signing method used to check the token's signing algorithm.
		SigningMethod: "HS256",

		// Signing key to validate token.
		SigningKey: []byte(_config.GetConfig().SecretJWT),

		// Skipper defines a function to skip middleware.
		Skipper: func(c echo.Context) bool {
			return c.Request().Header.Get("Authorization") == ""
		},

		// SuccessHandler defines a function which is executed for a valid token.
		SuccessHandler: func(c echo.Context) {
			c.Set("INFO", &User{ExtractToken(c), "admin"})
		},
	})
}

// Extract user id from token
// Return user id if successful, or -1 if token is nil or invalid
func ExtractToken(c echo.Context) int {
	token := c.Get("user").(*jwt.Token)

	if token != nil && token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		id := int(claims["id"].(float64))

		return id
	}

	return -1
}

// Create token by embedding user id
func CreateToken(id int) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(_config.GetConfig().SecretJWT))
}
