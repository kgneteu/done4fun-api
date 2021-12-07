package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
)

func main() {
	db = connectToDB()
	router := echo.New()

	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}", "method=${method}", "uri=${uri}", "status"="${status}"`,
	}))

	users := router.Group("api/user")
	{
		users.POST("/login", userLoginEndpoint)
		users.POST("/register", userRegisterEndpoint)
		users.GET("/:id", userReadEndpoint)
	}
	//
	//// Simple group: v2
	//v2 := router.Group("api/messages")
	//{
	//	v2.GET("/list", userLoginEndpoint)
	//}

	defer db.Close()
	s := http.Server{
		Addr:    ":9000",
		Handler: router,
		//ReadTimeout: 30 * time.Second, // customize http.Server timeouts
	}
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}


}
