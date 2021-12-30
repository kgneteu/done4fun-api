package main

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"server/models"
	"strconv"
	"strings"
	"time"
)

//Auth Middleware - check if token exists and store it with user id
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
				if claims.ExpiresAt < time.Now().Unix() {
					_ = Unauthorized(c, "expired")
				} else {
					_ = Unauthorized(c, "invalid token")
				}
				return
			}

			err = claims.Valid()
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

//requires any user, sets models.User in context UserInfo
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

//Auth Middleware - checks if admin
func (app *application) adminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		if c.Get("UserInfo") == nil {
			_ = Unauthorized(c)
			return errors.New("user required")
		}
		user := c.Get("UserInfo").(*models.User)
		if user.Role != "admin" {
			_ = Forbidden(c)
			return errors.New("admin only")
		}

		return next(c)
	}
}

//Auth Middleware - checks target id - admins always, parents and kid if id belongs to them
func (app *application) ownerAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		if c.Get("UserInfo") == nil {
			_ = Unauthorized(c)
			return errors.New("user required")
		}

		user := c.Get("UserInfo").(*models.User)
		targetUserId := c.Param("id")
		var targetId uint

		if targetId, err = toUint(targetUserId); err != nil {
			_ = BadRequest(c, "invalid id")
			return
		}

		var targetUser *models.User
		if targetId == user.ID {
			targetUser = user
		} else {
			targetUser, err = app.models.GetUserById(targetId)
			if err != nil {
				_ = BadRequest(c, "invalid id")
				return
			}
		}

		c.Set("TargetUserInfo", targetUser)

		if user.Role == "admin" {
			return next(c)
		}

		if user.Role == "kid" {
			if targetId != user.ID {
				_ = Forbidden(c)
				return errors.New("forbidden")
			}
			return next(c)
		}

		if user.Role == "parent" {
			if targetUser.ParentId == nil || *targetUser.ParentId != user.ID {
				_ = Forbidden(c)
				return errors.New("forbidden")
			}
		}

		return next(c)
	}
}
