package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (app *application) routes() *echo.Echo {
	router := echo.New()
	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}", "method=${method}", "uri=${uri}", "status"="${status}"`,
	}))
	//router.Use(middleware.Recover())

	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		//AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	router.Use(checkToken)

	users := router.Group("api/user")
	{
		users.POST("/login", app.userLoginEndpoint)
		users.POST("/register", app.userRegisterEndpoint)
	}

	admin := router.Group("/api/admin", app.adminAuth)
	{
		admin.GET("/user/list", app.getUserListEndpoint)
		admin.DELETE("/user/:id", app.userDeleteEndpoint)
	}

	auth := router.Group("/api/auth", app.authenticated)
	{
		auth.PATCH("/user/:id", app.userUpdateEndpoint)
		auth.POST("/user/create", app.userCreateEndpoint)
	}

	return router
}
