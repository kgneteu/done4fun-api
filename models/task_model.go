package models

import (
	"strconv"
	"strings"
)

func (model *DBModel) GetTask(taskId uint) (task *Task, err error) {
	sqlStatement := `SELECT FROM tasks WHERE id=$1;`
	task = &Task{}
	err = model.Db.Get(task, sqlStatement, taskId)
	return task, err
}

func (model *DBModel) DeleteTask(taskId uint) (err error) {
	sqlStatement := `DELETE FROM tasks WHERE id=$1;`
	_, err = model.Db.Exec(sqlStatement, &taskId)
	return err
}

func (model *DBModel) CreateTask(task map[string]string) (id uint, err error) {
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

func (model *DBModel) UpdateTask(task map[string]string, taskId uint) (err error) {
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
