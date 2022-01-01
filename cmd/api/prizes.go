package main

import (
	"errors"
	"github.com/labstack/echo"
	"net/http"
	"server/models"
	"strconv"
)

//middleware protected
func (app *application) getAvailablePrizesEndpoint(c echo.Context) (err error) {
	targetUser := c.Get(TargetUserInfo).(*models.User)
	var prizes *[]models.Prize
	prizes, err = app.models.GetAvailablePrizes(targetUser.ID)
	if err != nil {
		_ = InternalError(c, err.Error())
		return
	}
	return c.JSON(http.StatusOK, echo.Map{"prizes": prizes})
}

func (app *application) getPrizeEndpoint(c echo.Context) (err error) {
	var prizeId uint
	if prizeId, err = toUint(c.Param("prizeId")); err != nil {
		_ = BadRequest(c, "invalid prize id")
		return
	}

	var prize *models.Prize
	prize, err = app.models.GetPrize(prizeId)
	if err != nil {
		_ = BadRequest(c, "prize not found")
		return
	}

	err = app.CheckOwnerAccess(c.Get(UserInfo).(*models.User), prize.KidId)
	if err != nil {
		_ = Forbidden(c, "forbidden")
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"prize": prize,
	})
}

func (app *application) deletePrizeEndpoint(c echo.Context) (err error) {
	var prizeId uint
	if prizeId, err = toUint(c.Param("prizeId")); err != nil {
		_ = BadRequest(c, "invalid prize id")
		return
	}

	var prize *models.Prize
	prize, err = app.models.GetPrize(prizeId)
	if err != nil {
		_ = BadRequest(c, "prize not found")
		return
	}

	err = app.CheckOwnerAccess(c.Get(UserInfo).(*models.User), prize.KidId)
	if err != nil {
		_ = Forbidden(c, "forbidden")
		return
	}

	if err = app.models.DeletePrize(prizeId); err == nil {
		return OK(c, "deleted")
	}
	return
}

func (app *application) updatePrizeEndpoint(c echo.Context) (err error) {
	var prizeId uint
	if prizeId, err = toUint(c.Param("prizeId")); err != nil {
		_ = BadRequest(c, "invalid prize id")
		return
	}

	fields := map[string]string{}
	if err = c.Bind(&fields); err != nil {
		_ = BadRequest(c)
		return
	}

	user := c.Get(UserInfo).(*models.User)
	delete(fields, "id")
	delete(fields, "created_at")
	if user.Role == ParentRole {
		delete(fields, "kid_id")
	}

	var prize *models.Prize
	prize, err = app.models.GetPrize(prizeId)
	if err != nil {
		_ = BadRequest(c, "prize not found")
		return
	}

	err = app.CheckOwnerAccess(user, prize.KidId)
	if err != nil || user.Role == KidRole {
		_ = Forbidden(c, "forbidden")
		return
	}
	err = app.models.UpdatePrize(fields, prizeId)
	if err != nil {
		_ = Forbidden(c, err.Error())
		return
	}
	return OK(c, "changed")
}

func (app *application) createPrizeEndpoint(c echo.Context) (err error) {
	targetUser := c.Get(TargetUserInfo).(*models.User)

	fields := map[string]string{}
	if err = c.Bind(&fields); err != nil {
		_ = BadRequest(c)
		return
	}

	user := c.Get(UserInfo).(*models.User)

	if user.Role == KidRole {
		_ = Forbidden(c)
		return errors.New("forbidden")
	}

	delete(fields, "id")
	if user.Role == ParentRole {
		delete(fields, "kid_id")
	}
	fields["kid_id"] = strconv.FormatInt(int64(targetUser.ID), 10)

	var id uint
	id, err = app.models.CreatePrize(fields)
	if err != nil {
		_ = Forbidden(c, err.Error())
		return
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "created", "id": id})
}
