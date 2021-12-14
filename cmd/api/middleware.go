package main

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"net/http"
)

const accessTokenCookieName = "X-Auth"

//Auth Middleware
func checkToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenCookie, _ := c.Cookie(accessTokenCookieName)
		//todo support for auth headers
		//tokenHeader.

		claims := &jwt.StandardClaims{}
		if tokenCookie != nil {
			//token = decodeJWTToken
			_, err := jwt.ParseWithClaims(tokenCookie.Value, claims,
				func(token *jwt.Token) (interface{}, error) {
					return getJWTSigningKey(), nil
				})
			if err != nil {
				return c.JSON(http.StatusForbidden, echo.Map{"message": "invalid token"})
			}

			c.Set("User", map[string]interface{}{
				"Id":    claims.Id,
				"Token": tokenCookie.Value,
			})
		}
		return next(c)
	}
}
