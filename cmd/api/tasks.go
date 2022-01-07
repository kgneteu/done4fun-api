package main

import (
	"errors"
	"github.com/labstack/echo"
	"net/http"
	"server/models"
	"strconv"
)

func (app *application) getTaskEndpoint(c echo.Context) (err error) {
	var taskId uint
	if taskId, err = toUint(c.Param("taskId")); err != nil {
		_ = BadRequest(c, "invalid task id")
		return
	}

	var task *models.Task
	task, err = app.models.GetTask(taskId)
	if err != nil {
		_ = BadRequest(c, "task not found")
		return
	}

	err = app.CheckOwnerAccess(c.Get(UserInfo).(*models.User), task.KidId)
	if err != nil {
		_ = Forbidden(c, "forbidden")
		return
	}
	return c.JSON(http.StatusOK, echo.Map{
		"task": task,
	})
}

func (app *application) deleteTaskEndpoint(c echo.Context) (err error) {
	var taskId uint
	if taskId, err = toUint(c.Param("taskId")); err != nil {
		_ = BadRequest(c, "invalid task id")
		return
	}

	var task *models.Task
	task, err = app.models.GetTask(taskId)
	if err != nil {
		_ = BadRequest(c, "task not found")
		return
	}

	err = app.CheckOwnerAccess(c.Get("UserInfo").(*models.User), task.KidId)
	if err != nil {
		_ = Forbidden(c, "forbidden")
		return
	}

	if err = app.models.DeleteTask(taskId); err == nil {
		return OK(c, "deleted")
	}
	return
}

func (app *application) updateTaskEndpoint(c echo.Context) (err error) {
	var taskId uint
	if taskId, err = toUint(c.Param("taskId")); err != nil {
		_ = BadRequest(c, "invalid task id")
		return
	}

	var task *models.Task
	task, err = app.models.GetTask(taskId)
	if err != nil {
		_ = BadRequest(c, "task not found")
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

	err = app.CheckOwnerAccess(user, task.KidId)
	if err != nil || user.Role == KidRole {
		_ = Forbidden(c, "forbidden")
		return
	}
	err = app.models.UpdateTask(fields, taskId)
	if err != nil {
		_ = Forbidden(c, err.Error())
		return
	}
	return OK(c, "changed")
}

func (app *application) createTaskEndpoint(c echo.Context) (err error) {
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
	id, err = app.models.CreateTask(fields)
	if err != nil {
		_ = Forbidden(c, err.Error())
		return
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "created", "id": id})

}

//middleware protected
func (app *application) getAvailableTasksEndpoint(c echo.Context) (err error) {
	var pageInfo PageInfo
	if err = c.Bind(&pageInfo); err != nil {
		_ = BadRequest(c, err.Error())
		return
	}

	targetUser := c.Get(TargetUserInfo).(*models.User)
	var taskList models.TaskList
	taskList, err = app.models.GetAvailableTasks(targetUser.ID, pageInfo.Page, pageInfo.Limit, pageInfo.Order)
	if err != nil {
		_ = InternalError(c, err.Error())
		return
	}

	return c.JSON(http.StatusOK, echo.Map{"tasks": taskList.Tasks, "total": taskList.Total})
	//return c.JSON(http.StatusOK, echo.Map{"prizes": prizes})
}

