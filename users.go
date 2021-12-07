package main

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

type Login struct {
	Password string `form:"password" json:"password" xml:"password"`
	Email    string `form:"email" json:"email" xml:"email"`
}

func userLoginEndpoint(c echo.Context) error {
	var json Login
	if err := c.Bind(&json); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return nil
	}

	id, err := getUserId(json.Email, json.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return err
	}

	println(id)
	//
	//if json.User != "manu" || json.Password != "123" {
	//	c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
	//	return nil
	//}
	//

	// Set custom claims

	claims := &Claims{
		"Jon Snow",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}


func userReadEndpoint(c echo.Context) error {
	id := c.Param("id")
	c.String(http.StatusOK, "ok1"+id)
	return nil
}

//func userListEndpoint(c *gin.Context) {
//	c.String(http.StatusOK, "list1")
//}

func userRegisterEndpoint(c echo.Context) error {
	type Register struct {
		Password  string `form:"password" json:"password" xml:"password"`
		Email     string `form:"email" json:"email" xml:"email"`
		FirstName string `form:"firstname" json:"firstname" xml:"firstname"`
		LastName  string `form:"lastname" json:"lastname" xml:"lastname"`
	}

	var json Register
	if err := c.Bind(&json); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return nil
	}

	id, err := createUser(json.FirstName, json.LastName, json.Email, json.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return err
	}

	println(id)
	//
	//if json.User != "manu" || json.Password != "123" {
	//	c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
	//	return nil
	//}
	//

	// Set custom claims

	claims := &Claims{
		"Jon Snow",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
