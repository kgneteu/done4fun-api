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
	router.Use(middleware.CORS())

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
	admin := router.Group("/api/admin")
	{
		//admin.Group()

		admin.GET("/user/list", app.getUserListEndpoint)
		admin.POST("/user/create", userProfileCreateEndpoint)
		admin.DELETE("/user/:id", userProfileDeleteEndpoint)
		admin.PATCH("/user/:id", userProfileUpdateEndpoint)
	}

	return router
}
