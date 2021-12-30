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
		users.POST("/refresh", app.userRefreshTokenEndpoint)
	}

	admin := router.Group("/api/admin", app.authenticated, app.adminAuth)
	{
		admin.GET("/user/list", app.getUserListEndpoint)
		admin.DELETE("/user/:id", app.userDeleteEndpoint)
	}

	auth := router.Group("/api/auth/user", app.authenticated)
	{
		auth.GET("/children/list", app.getSubUserListEndpoint)
		auth.GET("/search", app.searchUserEndpoint)
		auth.GET("/:id", app.userGetEndpoint)
		auth.PATCH("/:id", app.userUpdateEndpoint)
		auth.POST("/create", app.userCreateEndpoint)
		auth.GET("/:id/prizes/available", app.getAvailablePrizesEndpoint, app.ownerAuth)
		//auth.POST("/:id/prize/:prizeId", app.addPrizeEndpoint)
	}
	//
	//prize := router.Group("/api/auth/prize", app.authenticated)
	//{
	//	//prize.DELETE("/:prizeId", app.deletePrizeEndpoint)
	//	//prize.PATCH("/:prizeId", app.updatePrizeEndpoint)
	//}

	return router
}
