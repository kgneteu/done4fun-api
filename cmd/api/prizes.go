package main

import (
	"github.com/labstack/echo"
	"net/http"
	"server/models"
)

func (app *application) getAvailablePrizesEndpoint(c echo.Context) (err error) {
	targetUser := c.Get("TargetUserInfo").(*models.User)
	var prizes *[]models.Prize
	prizes, err = app.models.GetAvailablePrizes(targetUser.ID)
	if err != nil {
		_ = InternalError(c, err.Error())
		return
	}
	return c.JSON(http.StatusOK, echo.Map{"prizes": prizes})
}
