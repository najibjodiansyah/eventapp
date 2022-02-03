package middlewares

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)
type User struct {
	Id   int
	Role string
}

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte("R4HASIA"),
		Skipper: func(c echo.Context) bool {
			// if c.Request().Header.Get("Authorization") == "" {
			// 	return true
			// }
			// return false
			return c.Request().Header.Get("Authorization") == ""
		}, SuccessHandler: func(c echo.Context) {
			c.Set("INFO", &User{ExtractToken(c), "admin"})
			fmt.Println(c)
		},
	})
}

func ExtractToken(e echo.Context) int {
	token := e.Get("user").(*jwt.Token)
	if token != nil && token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		id := claims["id"]
		switch id.(type) {
		case float64:
			return int(id.(float64))
		default:
			return id.(int)
		}
	}
	return -1 //invalid
}

func ValidateToken(e echo.Context) bool {
	login := e.Get("user").(*jwt.Token)

	return login.Valid
}

func CreateToken(id int) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("R4HASIA"))
}