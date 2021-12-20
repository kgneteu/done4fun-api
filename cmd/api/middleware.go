package main

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"server/models"
	"strconv"
	"strings"
)

//Auth Middleware
func checkToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 {
				_ = BadRequest(c, "invalid header")
				err = errors.New("invalid auth header")
				return
			}
			if headerParts[0] != "Bearer" {
				_ = BadRequest(c, "invalid header")
				err = errors.New("invalid auth header")
				return
			}

			token := headerParts[1]

			claims := &jwt.StandardClaims{}
			_, err = jwt.ParseWithClaims(token, claims,
				func(token *jwt.Token) (interface{}, error) {
					return getJWTSigningKey(), nil
				})
			if err != nil {
				_ = Unauthorized(c, "Invalid token")
				return
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
	return func(c echo.Context) (err error) {
		if c.Get("User") == nil {
			_ = Unauthorized(c)
			return errors.New("bad user")
		}

		userI := c.Get("User")
		if userI == nil {
			_ = Unauthorized(c)
			return
		}
		user := userI.(map[string]string)

		var uId uint64
		uId, err = strconv.ParseUint(user["Id"], 10, 32)
		if err != nil {
			_ = Unauthorized(c)
			return
		}

		var userData *models.User
		userData, err = app.models.GetUserById(uint(uId))
		if err != nil {
			_ = Unauthorized(c)
			return
		}

		if userData.Role != "admin" {
			_ = Forbidden(c)
			err = errors.New("invalid user")
			return
		}

		c.Set("UserInfo", userData)
		return next(c)
	}
}

//sets models.User in context UserInfo
func (app *application) authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		if c.Get("User") == nil {
			_ = Unauthorized(c)
			return errors.New("bad user")
		}
		user := c.Get("User").(map[string]string)
		if user == nil {
			_ = Unauthorized(c)
			return
		}
		var uId uint64
		uId, err = strconv.ParseUint(user["Id"], 10, 32)
		if err != nil {
			_ = Unauthorized(c)
			return
		}

		var userData *models.User
		userData, err = app.models.GetUserById(uint(uId))
		if err != nil {
			_ = Unauthorized(c)
			return
		}
		c.Set("UserInfo", userData)
		return next(c)
	}
}
