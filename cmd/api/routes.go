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
	//router.Use(middleware.CORS())

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

	//middleware.JWTWithConfig(middleware.JWTConfig{
	//	SigningKey:  getJWTSigningKey(),
	//	TokenLookup: "cookie:" + accessTokenCookieName,
	//})
	admin := router.Group("/api/admin", app.adminAuth)
	{
		//admin.Group()

		admin.GET("/user/list", app.getUserListEndpoint)
		admin.POST("/user/create", app.userProfileCreateEndpoint)
		admin.DELETE("/user/:id", app.userDeleteEndpoint)
		admin.PATCH("/user/:id", app.userProfileUpdateEndpoint)
	}

	return router
}
