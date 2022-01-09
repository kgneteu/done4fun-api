package models

import (
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
)

func (model *DBModel) GetUserById(Id uint) (user *User, err error) {
	sqlStatement := `SELECT * FROM users WHERE ID=$1;`
	user = &User{}
	err = model.Db.Get(user, sqlStatement, Id)
	return user, err
}

func (model *DBModel) GetUserByEmail(email string) (user *User, err error) {
	sqlStatement := `SELECT * FROM users WHERE email=$1;`
	user = &User{}
	err = model.Db.Get(user, sqlStatement, email)
	return user, err
}

type UserList struct {
	Total uint
	Users *[]User
}

func (model *DBModel) GetUserList(page int, limit int, order string) (userList UserList, err error) {
	const fields = "id, first_name, last_name, email, role, picture, verified"
	sqlStatement := `SELECT ` + fields + ` FROM users ORDER BY $3 LIMIT $1 OFFSET $2 ;`
	if limit == 0 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}
	if order == "" {
		order = "first_name, last_name"
	}
	offset := (page - 1) * limit

	if err = model.Db.Get(&userList.Total, `SELECT COUNT(*) as total FROM users`); err != nil {
		return
	}

	userList.Users = &[]User{}
	err = model.Db.Select(userList.Users, sqlStatement, limit, offset, order)
	return
}

func (model *DBModel) GetSubUserList(page int, limit int, order string, parentId uint) (userList UserList, err error) {
	const fields = "id, first_name, last_name, email, role, picture, verified"
	sqlStatement := `SELECT ` + fields + ` FROM users WHERE parent_id=$4 ORDER BY $3 DESC LIMIT $1 OFFSET $2 ;`
	if limit == 0 {
		limit = 10
	}

	if page < 1 {
		page = 1
	}
	if order == "" {
		order = "first_name, last_name"
	}
	offset := (page - 1) * limit

	if err = model.Db.Get(&userList.Total, `SELECT COUNT(*) as total FROM users WHERE parent_id=$1`, parentId); err != nil {
		return
	}

	userList.Users = &[]User{}
	err = model.Db.Select(userList.Users, sqlStatement, limit, offset, order, parentId)
	return
}

func (model *DBModel) SearchUser(limit int, order string, filter string, role string, parentId uint) (userList UserList, err error) {
	const fields = "id, first_name, last_name, email, role, picture, verified"
	where := ""
	if limit == 0 {
		limit = 10
	}
	if order == "" {
		order = "'id'"
	}

	params := map[string]interface{}{
		"order": order,
		"limit": limit,
	}

	if filter != "" || role != "" || parentId > 0 {
		params["parent"] = parentId
		where = "WHERE "
		if parentId > 0 {
			where += "parent_id=:parent"
		}

		if role != "" {
			params["role"] = role
			if where != "" {
				where += " AND "
			}
			where += "role=:role"
		}

		if filter != "" {
			params["filter1"] = "%" + filter + "%"
			params["filter2"] = filter + "%"
			if where != "" {
				where += " AND "
			}
			where += "(email LIKE :filter1 OR first_name LIKE :filter2 OR last_name LIKE :filter2)"
		}
	}
	sqlStatement := `SELECT ` + fields + ` FROM users ` + where + ` ORDER BY :order LIMIT :limit;`
	userList.Users = &[]User{}

	//err = model.Db.Select(userList.Users, sqlStatparams)
	var user User
	var users []User
	var rows *sqlx.Rows
	if rows, err = model.Db.NamedQuery(sqlStatement, params); err != nil {
		return
	}

	for rows.Next() {
		if err = rows.StructScan(&user); err != nil {
			return
		}
		users = append(users, user)
	}
	userList.Users = &users
	userList.Total = uint(len(users))
	rows.Close()
	return
}

func (model *DBModel) CreateNewUser(firstName, lastName, email, password string) (int64, error) {
	var err error
	var id int64

	sqlStatement := `
		INSERT INTO users (first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err = model.Db.QueryRow(sqlStatement, firstName, lastName, email, password).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (model *DBModel) DeleteUser(id uint) (err error) {
	sqlStatement := `DELETE FROM users WHERE id=$1;`
	_, err = model.Db.Exec(sqlStatement, &id)
	return err
}

func (model *DBModel) CreateUser(user map[string]interface{}) (id uint, err error) {
	var fields []string
	var values []interface{}
	var placeholders []string
	i := 1
	for k, v := range user {
		fields = append(fields, k)
		values = append(values, v)
		placeholders = append(placeholders, "$"+strconv.Itoa(i))
		i++
	}

	fString := strings.Join(fields, ", ")
	pString := strings.Join(placeholders, ", ")
	sqlStatement := `
		INSERT INTO users (` + fString + `)
		VALUES (` + pString + `)
		RETURNING id`

	err = model.Db.QueryRow(sqlStatement, values...).Scan(&id)
	return
}

func (model *DBModel) UpdateUser(user map[string]interface{}, userId uint) (err error) {
	var fields []string
	var values []interface{}
	values = append(values, userId)
	i := 2
	for k, v := range user {
		fields = append(fields, k+"=$"+strconv.Itoa(i))
		values = append(values, v)
		i++
	}

	fString := strings.Join(fields, ", ")
	sqlStatement := `UPDATE users SET ` + fString + ", updated_at=NOW() WHERE id=$1"
	_, err = model.Db.Exec(sqlStatement, values...)
	return
}
