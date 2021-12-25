package main

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"net/http"
	"server/models"
	"strconv"
	"strings"
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
		FirstName string `form:"first_name" json:"first_name" xml:"first_name"`
		LastName  string `form:"last_name" json:"last_name" xml:"last_name"`
	}

	var json Register
	if err = c.Bind(&json); err != nil {
		_ = BadRequest(c, err.Error())
		return
	}

	if strings.Trim(json.FirstName, "") == "" || strings.Trim(json.LastName, "") == "" || strings.Trim(json.Password, "") == "" || strings.Trim(json.Email, "") == "" {
		_ = BadRequest(c)
		return errors.New("missing required fields")
	}

	json.Password, err = PasswordHash(json.Password)
	if err != nil {
		_ = InternalError(c)
		return
	}

	_, err = app.models.CreateNewUser(json.FirstName, json.LastName, json.Email, json.Password)
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
func (app *application) userCreateEndpoint(c echo.Context) (err error) {
	if c.Get("UserInfo") == nil {
		_ = Unauthorized(c)
		return errors.New("bad user")
	}
	var userInfo *models.User
	userInfo = c.Get("UserInfo").(*models.User)
	if userInfo == nil {
		_ = Forbidden(c)
		return errors.New("need user")
	}
	if !(userInfo.Role == "admin" || userInfo.Role == "parent") {
		_ = Forbidden(c)
		return errors.New("invalid user")
	}

	fields := map[string]string{}
	if err = c.Bind(&fields); err != nil {
		_ = BadRequest(c)
		return err
	}

	requiredFields := []string{"email", "password", "first_name", "last_name"}
	for _, field := range requiredFields {
		f := fields[field]
		if strings.Trim(f, "") == "" {
			_ = BadRequest(c)
			return errors.New("missing required fields")
		}
	}

	fields["password"], err = PasswordHash(fields["password"])
	if err != nil {
		_ = InternalError(c)
		return
	}

	if userInfo.Role == "parent" {
		fields["parent_id"] = strconv.FormatInt(int64(userInfo.ID), 10)
		fields["role"] = "kid"
	}

	_, err = app.models.CreateUser(fields)
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

func (app *application) userUpdateEndpoint(c echo.Context) (err error) {
	var userInfo *models.User
	var targetId uint

	if c.Get("UserInfo") == nil {
		_ = Unauthorized(c)
		return errors.New("bad user")
	}

	userInfo = c.Get("UserInfo").(*models.User)
	if userInfo == nil {
		_ = Forbidden(c)
		return errors.New("need user")
	}

	fields := map[string]string{}
	if err = c.Bind(&fields); err != nil {
		_ = BadRequest(c)
		return err
	}

	userId := c.Param("id")
	if targetId, err = toUint(userId); err != nil {
		_ = BadRequest(c, "invalid id")
		return
	}

	if !(userInfo.Role == "admin" || userInfo.Role == "parent") {
		if targetId != userInfo.ID {
			_ = Forbidden(c)
			return errors.New("invalid target")
		}
	}

	if userInfo.Role == "parent" && targetId != userInfo.ID {
		fields["role"] = "kid"
		var subUser *models.User
		subUser, err = app.models.GetUserById(targetId)
		if err != nil {
			_ = BadRequest(c, "invalid id")
			return
		}
		if subUser.ParentId != userInfo.ID {
			_ = Forbidden(c)
			return errors.New("invalid target")
		}
	}

	if _, ok := fields["password"]; ok {
		if userInfo.ID == targetId {
			oldPassword := fields["old_password"]
			if !PasswordVerify(oldPassword, userInfo.Password) {
				_ = Unauthorized(c)
				return errors.New("bad password")
			}
		}

		fields["password"], err = PasswordHash(fields["password"])
		if err != nil {
			_ = InternalError(c)
			return
		}
		delete(fields, "old_password")
	}

	err = app.models.UpdateUser(fields, targetId)
	if err != nil {
		if pqErr := err.(*pq.Error); pqErr.Code == "23505" {
			_ = Forbidden(c, "duplicated")
		} else {
			_ = Forbidden(c, err.Error())
		}
		return
	}
	return OK(c, "changed")
}

func (app *application) userDeleteEndpoint(c echo.Context) (err error) {
	userId := c.Param("id")
	var id uint
	if id, err = toUint(userId); err != nil {
		_ = BadRequest(c, "invalid id")
		return
	}
	if err = app.models.DeleteUser(id); err != nil {
		_ = InternalError(c, err.Error())
		return
	}
	return OK(c, "deleted")
}

//todo filter data
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

//todo filter data
func (app *application) getSubUserListEndpoint(c echo.Context) (err error) {
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

	if c.Get("UserInfo") == nil {
		_ = Unauthorized(c)
		return errors.New("bad user")
	}

	var userInfo *models.User
	userInfo = c.Get("UserInfo").(*models.User)
	if userInfo == nil {
		_ = Forbidden(c)
		return errors.New("need user")
	}

	userList, err := app.models.GetSubUserList(json.Page, json.Limit, json.Order, userInfo.ID)
	if err != nil {
		_ = BadRequest(c, err.Error())
		return
	}
	return c.JSON(http.StatusOK, echo.Map{"users": userList.Users, "total": userList.Total})
}

func (app *application) searchUserEndpoint(c echo.Context) (err error) {
	type PageInfo struct {
		Limit  int    `form:"limit" json:"limit" xml:"limit"`
		Order  string `form:"order" json:"order" xml:"order"`
		Role   string `form:"role" json:"role" xml:"role"`
		Filter string `form:"filter" json:"filter" xml:"filter"`
	}

	var json PageInfo
	if err = c.Bind(&json); err != nil {
		_ = BadRequest(c, err.Error())
		return
	}

	if c.Get("UserInfo") == nil {
		_ = Unauthorized(c)
		return errors.New("bad user")
	}

	var userInfo *models.User
	userInfo = c.Get("UserInfo").(*models.User)
	if userInfo == nil {
		_ = Forbidden(c)
		return errors.New("need user")
	}

	if userInfo.Role == "kid" {
		_ = Forbidden(c)
		return errors.New("forbidden")
	}

	var parentId uint
	if userInfo.Role == "parent" {
		parentId = userInfo.ID
	}

	userList, err := app.models.SearchUser(json.Limit, json.Order, json.Filter, json.Role, parentId)
	if err != nil {
		_ = BadRequest(c, err.Error())
		return
	}
	return c.JSON(http.StatusOK, echo.Map{"users": userList.Users, "total": userList.Total})
}
