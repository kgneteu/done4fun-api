package main

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"net/http"
)

func (app *application) userLoginEndpoint(c echo.Context) (err error) {
	type Credentials struct {
		Email    string `form:"email" json:"email" xml:"email"`
		Password string `form:"password" json:"password" xml:"password"`
	}

	var json Credentials
	if err = c.Bind(&json); err != nil {
		_ = BadRequest(c, err.Error())
		return
	}

	user, err := app.models.GetUserByEmail(json.Email)
	if err != nil {
		_ = Unauthorized(c)
		return
	}

	if !PasswordVerify(json.Password, user.Password) {
		_ = Unauthorized(c)
		return errors.New("bad password")
	}

	var token string
	token, err = createToken(user.ID)
	if err != nil {
		_ = InternalError(c)
		return
	}
	//
	//cookie := new(http.Cookie)
	//cookie.Name = accessTokenCookieName
	//cookie.Value = token
	//cookie.HttpOnly = true
	//cookie.SameSite = 4
	////cookie.Secure = false
	////cookie.MaxAge = 1000000000
	////cookie.Domain = ""
	//cookie.Path = "/"
	//cookie.Expires = time.Now().Add(24 * time.Hour)
	//
	//c.SetCookie(cookie)

	filteredUser := map[string]interface{}{
		"id":         user.ID,
		"role":       user.Role,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"parent_id":  user.ParentId,
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
		"user":  filteredUser,
	})
}

func (app *application) userRegisterEndpoint(c echo.Context) (err error) {
	type Register struct {
		Password  string `form:"password" json:"password" xml:"password"`
		Email     string `form:"email" json:"email" xml:"email"`
		FirstName string `form:"firstname" json:"firstname" xml:"firstname"`
		LastName  string `form:"lastname" json:"lastname" xml:"lastname"`
	}

	var json Register
	if err = c.Bind(&json); err != nil {
		_ = BadRequest(c, err.Error())
		return
	}

	_, err = app.models.CreateUser(json.FirstName, json.LastName, json.Email, json.Password)
	if err != nil {
		if pqErr := err.(*pq.Error); pqErr.Code == "23505" {
			_ = Forbidden(c, "duplicated")
		} else {
			_ = Forbidden(c, err.Error())
		}
		return
	}
	return OK(c, "account created")
}

// ADMIN ZONE
func (app *application) userProfileCreateEndpoint(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "created"})
}

func (app *application) userProfileUpdateEndpoint(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "updated"})
}

func (app *application) userDeleteEndpoint(c echo.Context) (err error) {
	userId := c.Param("id")
	var id uint
	if id, err = toUint(userId); err != nil {
		_ = BadRequest(c, err.Error())
		return
	}
	if err = app.models.DeleteUser(id); err != nil {
		_ = InternalError(c, err.Error())
		return
	}
	return OK(c, "deleted")
}

func (app *application) getUserListEndpoint(c echo.Context) (err error) {
	type PageInfo struct {
		Page  int    `form:"page" json:"page" xml:"page"`
		Limit int    `form:"limit" json:"limit" xml:"limit"`
		Order string `form:"order" json:"order" xml:"order"`
	}
	var json PageInfo
	if err = c.Bind(&json); err != nil {
		_ = BadRequest(c, err.Error())
		return
	}

	userList, err := app.models.GetUserList(json.Page, json.Limit, json.Order)
	if err != nil {
		_ = BadRequest(c, err.Error())
		return
	}
	return c.JSON(http.StatusOK, echo.Map{"users": userList.Users, "total": userList.Total})
}
