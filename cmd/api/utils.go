package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func Forbidden(c echo.Context, msg ...string) error {
	message := "forbidden"
	if len(msg) != 0 {
		message = msg[0]
	}
	return c.JSON(http.StatusForbidden, echo.Map{"message": message})
}

func BadRequest(c echo.Context, msg ...string) error {
	message := "bad request"
	if len(msg) != 0 {
		message = msg[0]
	}
	return c.JSON(http.StatusBadRequest, echo.Map{"message": message})
}

func Unauthorized(c echo.Context, msg ...string) error {
	message := "unauthorized"
	if len(msg) != 0 {
		message = msg[0]
	}
	return c.JSON(http.StatusUnauthorized, echo.Map{"message": message})
}

func InternalError(c echo.Context, msg ...string) error {
	message := "internal error"
	if len(msg) != 0 {
		message = msg[0]
	}
	return c.JSON(http.StatusInternalServerError, echo.Map{"message": message})
}

func OK(c echo.Context, msg ...string) error {
	message := "ok"
	if len(msg) != 0 {
		message = msg[0]
	}
	return c.JSON(http.StatusOK, echo.Map{"message": message})
}

func toUint(num string) (res uint, err error) {
	var r uint64
	r, err = strconv.ParseUint(num, 10, 32)
	if err == nil {
		res = uint(r)
	}
	return
}
