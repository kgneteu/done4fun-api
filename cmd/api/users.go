package main

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"net/http"
	"server/models"
	"strconv"
	"strings"
	"time"
)

func (app *application) authenticatedUser(c echo.Context, email string, password string) (err error) {
	var user *models.User
	user, err = app.models.GetUserByEmail(email)
	if err != nil {
		_ = Unauthorized(c)
		return
	}

	if !PasswordVerify(password, user.Password) {
		_ = Unauthorized(c)
		return errors.New("bad password")
	}

	var token string
	token, err = createToken(user.ID)
	if err != nil {
		_ = InternalError(c)
		return
	}

	var refreshToken string
	refreshToken, err = createRefreshToken(user.ID, email, password)
	if err != nil {
		_ = InternalError(c)
		return
	}

	var filteredUser echo.Map
	if filteredUser, err = app.filterUserData(user); err != nil {
		_ = InternalError(c)
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token":         token,
		"refresh_token": refreshToken,
		"user":          filteredUser,
	})
}

func (app *application) userLoginEndpoint(c echo.Context) (err error) {
	type Credentials struct {
		Email    string `form:"email" json:"email" xml:"email"`
		Password string `form:"password" json:"password" xml:"password"`
	}

	var jsonCredentials Credentials
	if err = c.Bind(&jsonCredentials); err != nil {
		_ = BadRequest(c, err.Error())
		return
	}
	return app.authenticatedUser(c, jsonCredentials.Email, jsonCredentials.Password)
}

func (app *application) userRefreshTokenEndpoint(c echo.Context) (err error) {
	type RefreshToken struct {
		Token string `form:"token" json:"token" xml:"token"`
	}

	type RefreshClaims struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		jwt.StandardClaims
	}

	var jsonRefreshToken RefreshToken
	if err = c.Bind(&jsonRefreshToken); err != nil {
		_ = BadRequest(c, err.Error())
		return
	}

	claims := &RefreshClaims{}
	_, err = jwt.ParseWithClaims(jsonRefreshToken.Token, claims,
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

	return app.authenticatedUser(c, claims.Email, claims.Password)
}

//todo database timestamps verified
func (app *application) userRegisterEndpoint(c echo.Context) (err error) {
	type Register struct {
		Password  string `form:"password" json:"password" xml:"password"`
		Email     string `form:"email" json:"email" xml:"email"`
		FirstName string `form:"first_name" json:"first_name" xml:"first_name"`
		LastName  string `form:"last_name" json:"last_name" xml:"last_name"`
	}

	var jsonRegister Register
	if err = c.Bind(&jsonRegister); err != nil {
		_ = BadRequest(c, err.Error())
		return
	}

	if strings.Trim(jsonRegister.FirstName, "") == "" || strings.Trim(jsonRegister.LastName, "") == "" || strings.Trim(jsonRegister.Password, "") == "" || strings.Trim(jsonRegister.Email, "") == "" {
		_ = BadRequest(c)
		return errors.New("missing required fields")
	}

	jsonRegister.Password, err = PasswordHash(jsonRegister.Password)
	if err != nil {
		_ = InternalError(c)
		return
	}

	_, err = app.models.CreateNewUser(jsonRegister.FirstName, jsonRegister.LastName, jsonRegister.Email, jsonRegister.Password)
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

func (app *application) filterUserData(user *models.User) (map[string]interface{}, error) {
	var data = make(map[string]interface{})
	b, err := json.Marshal(user)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return data, err
	}
	delete(data, "created_at")
	delete(data, "deleted_at")
	delete(data, "updated_at")
	delete(data, "password")
	return data, nil
}

//todo filter data
func (app *application) getUserListEndpoint(c echo.Context) (err error) {
	type PageInfo struct {
		Page  int    `form:"page" json:"page" xml:"page"`
		Limit int    `form:"limit" json:"limit" xml:"limit"`
		Order string `form:"order" json:"order" xml:"order"`
	}
	var jsonPageInfo PageInfo
	if err = c.Bind(&jsonPageInfo); err != nil {
		_ = BadRequest(c, err.Error())
		return
	}

	userList, err := app.models.GetUserList(jsonPageInfo.Page, jsonPageInfo.Limit, jsonPageInfo.Order)
	if err != nil {
		_ = BadRequest(c, err.Error())
		return
	}
	return c.JSON(http.StatusOK, echo.Map{"users": userList.Users, "total": userList.Total})
}

//todo filter data
func (app *application) getSubUserListEndpoint(c echo.Context) (err error) {
	var jsonPageInfo PageInfo
	if err = c.Bind(&jsonPageInfo); err != nil {
		_ = BadRequest(c, err.Error())
		return
	}

	if c.Get(UserInfo) == nil {
		_ = Unauthorized(c)
		return errors.New("bad user")
	}

	var userInfo *models.User
	userInfo = c.Get(UserInfo).(*models.User)
	if userInfo == nil {
		_ = Forbidden(c)
		return errors.New("need user")
	}

	userList, err := app.models.GetSubUserList(jsonPageInfo.Page, jsonPageInfo.Limit, jsonPageInfo.Order, userInfo.ID)
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

	var jsonPageInfo PageInfo
	if err = c.Bind(&jsonPageInfo); err != nil {
		_ = BadRequest(c, err.Error())
		return
	}

	if c.Get(UserInfo) == nil {
		_ = Unauthorized(c)
		return errors.New("bad user")
	}

	var userInfo *models.User
	userInfo = c.Get(UserInfo).(*models.User)
	if userInfo == nil {
		_ = Forbidden(c)
		return errors.New("need user")
	}

	if userInfo.Role == KidRole {
		_ = Forbidden(c)
		return errors.New("forbidden")
	}

	var parentId uint

	if userInfo.Role == ParentRole {
		parentId = userInfo.ID
	}

	userList, err := app.models.SearchUser(jsonPageInfo.Limit, jsonPageInfo.Order, jsonPageInfo.Filter, jsonPageInfo.Role, parentId)
	if err != nil {
		_ = BadRequest(c, err.Error())
		return
	}
	return c.JSON(http.StatusOK, echo.Map{"users": userList.Users, "total": userList.Total})
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (app *application) createUserEndpoint(c echo.Context) (err error) {

	var user *models.User
	user = c.Get(UserInfo).(*models.User)
	if !(user.Role == AdminRole || user.Role == ParentRole) {
		_ = Forbidden(c)
		return errors.New("invalid user")
	}

	fields := map[string]interface{}{}
	if err = c.Bind(&fields); err != nil {
		_ = BadRequest(c)
		return err
	}

	requiredFields := []string{"email", "password", "first_name", "last_name"}
	for _, field := range requiredFields {
		if _, ok := fields[field]; !ok {
			_ = BadRequest(c)
			return errors.New("missing required field: " + field)
		}
	}

	fields["password"], err = PasswordHash(fields["password"].(string))
	if err != nil {
		_ = InternalError(c)
		return
	}

	delete(fields, "id")

	if user.Role == ParentRole {
		fields["parent_id"] = strconv.FormatInt(int64(user.ID), 10)
		fields["role"] = KidRole //default is parent
	}

	var id uint
	id, err = app.models.CreateUser(fields)
	if err != nil {
		if pqErr := err.(*pq.Error); pqErr.Code == "23505" {
			_ = Forbidden(c, "duplicated")
		} else {
			_ = Forbidden(c, err.Error())
		}
		return
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "created", "id": id})
}

func (app *application) updateUserEndpoint(c echo.Context) (err error) {
	user := c.Get(UserInfo).(*models.User)
	targetUser := c.Get(TargetUserInfo).(*models.User)

	fields := map[string]interface{}{}
	if err = c.Bind(&fields); err != nil {
		_ = BadRequest(c)
		return err
	}

	//role can't be changed
	if _, ok := fields["role"]; ok {
		_ = BadRequest(c, "invalid id")
		return errors.New("admins only")
	}

	if _, ok := fields["password"]; ok {
		if user.ID == targetUser.ID {
			if _, ok = fields["old_password"]; ok {
				oldPassword := fields["old_password"].(string)
				if !PasswordVerify(oldPassword, targetUser.Password) {
					_ = Unauthorized(c)
					return errors.New("bad password")
				}
			} else {
				_ = Unauthorized(c)
				return errors.New("old password required")
			}
		}

		fields["password"], err = PasswordHash(fields["password"].(string))
		if err != nil {
			_ = InternalError(c)
			return
		}
		delete(fields, "old_password")
	}

	delete(fields, "id")
	delete(fields, "parent_id")
	delete(fields, "created_at")

	err = app.models.UpdateUser(fields, targetUser.ID)
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

func (app *application) getUserEndpoint(c echo.Context) (err error) {
	targetUser := c.Get(TargetUserInfo).(*models.User)
	var filteredUser echo.Map
	if filteredUser, err = app.filterUserData(targetUser); err != nil {
		_ = InternalError(c)
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"user": filteredUser,
	})
}

func (app *application) deleteUserEndpoint(c echo.Context) (err error) {
	targetUser := c.Get(TargetUserInfo).(*models.User)
	if err = app.models.DeleteUser(targetUser.ID); err != nil {
		_ = InternalError(c, err.Error())
		return
	}
	return OK(c, "deleted")
}
