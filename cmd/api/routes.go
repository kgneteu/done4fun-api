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
		users.POST("/login", app.userLoginEndpoint)          //todo refactor
		users.POST("/register", app.userRegisterEndpoint)    //todo refactor
		users.POST("/refresh", app.userRefreshTokenEndpoint) //todo refactor
	}

	admin := router.Group("/api/admin", app.authenticated, app.adminAuth)
	{
		admin.GET("/user/list", app.getUserListEndpoint) //todo refactor
	}

	auth := router.Group("/api/auth/user", app.authenticated)
	{
		auth.GET("/:id", app.getUserEndpoint, app.ownerAuth)
		auth.DELETE("/:id", app.deleteUserEndpoint, app.ownerAuth)
		auth.PATCH("/:id", app.updateUserEndpoint, app.ownerAuth)
		auth.POST("/create", app.createUserEndpoint)

		auth.GET("/:id/prizes/available", app.getAvailablePrizesEndpoint, app.ownerAuth)
		auth.GET("/:id/tasks/available", app.getAvailableTasksEndpoint, app.ownerAuth)

		auth.GET("/children/list", app.getSubUserListEndpoint) //todo refactor
		auth.GET("/search", app.searchUserEndpoint)            //todo refactor

		auth.POST("/:id/task/create", app.createTaskEndpoint, app.ownerAuth)
		auth.POST("/:id/prize/create", app.createPrizeEndpoint, app.ownerAuth)
	}

	prize := router.Group("/api/auth/prize", app.authenticated)
	{
		prize.GET("/:prizeId", app.getPrizeEndpoint)
		prize.DELETE("/:prizeId", app.deletePrizeEndpoint)
		prize.PATCH("/:prizeId", app.updatePrizeEndpoint)
	}

	task := router.Group("/api/auth/task", app.authenticated)
	{
		task.GET("/:taskId", app.getTaskEndpoint)
		task.DELETE("/:taskId", app.deleteTaskEndpoint)
		task.PATCH("/:taskId", app.updateTaskEndpoint)
	}

	return router
}
