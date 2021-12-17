package main

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"strconv"
	"strings"
)

//Auth Middleware
func checkToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//tokenCookie, _ := c.Cookie(accessTokenCookieName)
		//claims := &jwt.StandardClaims{}
		//if tokenCookie != nil {
		//	_, err := jwt.ParseWithClaims(tokenCookie.Value, claims,
		//		func(token *jwt.Token) (interface{}, error) {
		//			return getJWTSigningKey(), nil
		//		})
		//	if err != nil {
		//		return Unauthorized(c, "Invalid token")
		//	}
		//
		//	c.Set("User", map[string]string{
		//		"Id":    claims.Id,
		//		"Token": tokenCookie.Value,
		//	})
		//	return next(c)
		//}
		//check headers

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 {
				_ = BadRequest(c, "invalid header")
				return errors.New("invalid auth header")
			}
			if headerParts[0] != "Bearer" {
				_ = BadRequest(c, "invalid header")
				return errors.New("invalid auth header")
			}

			token := headerParts[1]

			claims := &jwt.StandardClaims{}
			_, err := jwt.ParseWithClaims(token, claims,
				func(token *jwt.Token) (interface{}, error) {
					return getJWTSigningKey(), nil
				})
			if err != nil {
				return Unauthorized(c, "Invalid token")
			}

			c.Set("User", map[string]string{
				"Id":    claims.Id,
				"Token": token,
			})
		}
		return next(c)
	}
}

//Auth Middleware
func (app *application) adminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userI := c.Get("User")
		if userI == nil {
			return Unauthorized(c)
		}
		user := userI.(map[string]string)
		uId, err := strconv.ParseUint(user["Id"], 10, 32)
		if err != nil {
			return Unauthorized(c)
		}
		userData, err := app.models.GetUserById(uint(uId))
		if err != nil {
			return Unauthorized(c)
		}
		if userData.Role != "admin" {
			return Forbidden(c)
		}
		return next(c)
	}
}

//sets models.User in context UserInfo
func (app *application) authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("User").(map[string]string)
		if user == nil {
			return Unauthorized(c)
		}
		uId, err := strconv.ParseUint(user["Id"], 10, 32)
		if err != nil {
			return Unauthorized(c)
		}
		userData, err := app.models.GetUserById(uint(uId))
		if err != nil {
			return Unauthorized(c)
		}
		c.Set("UserInfo", userData)
		return next(c)
	}
}
