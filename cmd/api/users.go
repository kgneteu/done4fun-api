package main

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"net/http"
	"time"
)

func (app *application) userLoginEndpoint(c echo.Context) error {
	type Credentials struct {
		Email    string `form:"email" json:"email" xml:"email"`
		Password string `form:"password" json:"password" xml:"password"`
	}

	var json Credentials
	if err := c.Bind(&json); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return err
	}

	user, err := app.models.GetUserByEmail(json.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthorized"})
		return err
	}

	if !PasswordVerify(json.Password, user.Password){
		c.JSON(http.StatusUnauthorized, echo.Map{"message": "unauthorized"})
		return errors.New("Bad password")
	}

	var token string
	token, err = createToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"message": "unauthorized"})
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = accessTokenCookieName
	cookie.Value = token
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour)

	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
		"user": user,
	})
}

func (app *application) userRegisterEndpoint(c echo.Context) error {
	type Register struct {
		Password  string `form:"password" json:"password" xml:"password"`
		Email     string `form:"email" json:"email" xml:"email"`
		FirstName string `form:"firstname" json:"firstname" xml:"firstname"`
		LastName  string `form:"lastname" json:"lastname" xml:"lastname"`
	}

	var json Register
	if err := c.Bind(&json); err != nil {
		c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
		return err
	}

	_, err := app.models.CreateUser(json.FirstName, json.LastName, json.Email, json.Password)
	if err != nil {
		if pqErr := err.(*pq.Error); pqErr.Code == "23505" {
			c.JSON(http.StatusForbidden, echo.Map{"message": "duplicated"})
		} else {
			c.JSON(http.StatusForbidden, echo.Map{"message": err})
		}
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "account created"})
}

// ADMIN ZONE
func userProfileCreateEndpoint(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "created"})
}

func userProfileUpdateEndpoint(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "updated"})
}

func userProfileDeleteEndpoint(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "deleted"})
}

func getUserListEndpoint(c echo.Context) error {
	type PageLogin struct {
		Page  int    `form:"page" json:"page" xml:"page"`
		Limit int    `form:"limit" json:"limit" xml:"limit"`
		Order string `form:"order" json:"order" xml:"order"`
	}
	c.Param("page")
	return c.JSON(http.StatusOK, echo.Map{"message": "deleted"})
}
