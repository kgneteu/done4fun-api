package models

import (
	"strconv"
	"strings"
)

func (model *DBModel) GetTask(taskId uint) (task *Task, err error) {
	sqlStatement := `SELECT * FROM tasks WHERE id=$1;`
	task = &Task{}
	err = model.Db.Get(task, sqlStatement, taskId)
	return task, err
}

func (model *DBModel) DeleteTask(taskId uint) (err error) {
	sqlStatement := `DELETE FROM tasks WHERE id=$1;`
	_, err = model.Db.Exec(sqlStatement, taskId)
	return err
}

func (model *DBModel) CreateTask(task map[string]interface{}) (id uint, err error) {
	var fields []string
	var values []interface{}
	var placeholders []string
	i := 1
	for k, v := range task {
		fields = append(fields, k)
		values = append(values, v)
		placeholders = append(placeholders, "$"+strconv.Itoa(i))
		i++
	}
	fString := strings.Join(fields, ", ")
	pString := strings.Join(placeholders, ", ")
	sqlStatement := `
		INSERT INTO tasks (` + fString + `)
		VALUES (` + pString + `)
		RETURNING id`

	err = model.Db.QueryRow(sqlStatement, values...).Scan(&id)
	return
}

func (model *DBModel) UpdateTask(task map[string]interface{}, taskId uint) (err error) {
	var fields []string
	var values []interface{}
	values = append(values, taskId)
	i := 2
	for k, v := range task {
		fields = append(fields, k+"=$"+strconv.Itoa(i))
		values = append(values, v)
		i++
	}

	fString := strings.Join(fields, ", ")
	sqlStatement := `UPDATE tasks SET ` + fString + ", updated_at=NOW() WHERE id=$1"
	_, err = model.Db.Exec(sqlStatement, values...)
	return
}

type TaskList struct {
	Total int
	Tasks *[]Task
}

func (model *DBModel) GetAvailableTasks(userId uint, page int, limit int, order string) (taskList TaskList, err error) {
	if page < 1 {
		page = 1
	}
	if order == "" {
		order = "'id'"
	}
	if limit < 0 {
		limit = 0
	}
	offset := (page - 1) * limit
	taskList.Tasks = &[]Task{}

	if limit > 0 {
		sqlStatement := `SELECT * FROM tasks WHERE kid_id=$1 AND published=true ORDER BY $4 OFFSET $2 LIMIT $3`
		err = model.Db.Select(taskList.Tasks, sqlStatement, userId, offset, limit, order)
		if err = model.Db.Get(&taskList.Total, `SELECT COUNT(*) as total FROM tasks WHERE kid_id=$1 and published=true`, userId); err != nil {
			return
		}
	} else {
		sqlStatement := `SELECT * FROM tasks WHERE kid_id=$1 AND published=true ORDER BY $3 OFFSET $2`
		if err = model.Db.Select(taskList.Tasks, sqlStatement, userId, offset, order); err != nil {
			taskList.Total = len(*taskList.Tasks)
		}
	}

	return
}
