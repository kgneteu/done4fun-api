package models

func (db *DBModel) CreateUser(firstName, lastName, email, password string) (int64, error) {
	var err error
	var id int64

	sqlStatement := `
		INSERT INTO users (first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err = db.Db.QueryRow(sqlStatement, firstName, lastName, email, password).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (db *DBModel) GetUserById(Id uint) (user *User, err error) {
	sqlStatement := `SELECT * FROM users WHERE ID=$1;`
	user = &User{}
	err = db.Db.Get(user, sqlStatement, Id)
	return user, err
}

func (db *DBModel) GetUserByEmail(email string) (user *User, err error) {
	sqlStatement := `SELECT * FROM users WHERE email=$1;`
	user = &User{}
	err = db.Db.Get(user, sqlStatement, email)
	return user, err
}

type UserList struct {
	Total uint
	Users *[]User
}

func (db *DBModel) GetUserList(page int, limit int, order string) (userList UserList, err error) {
	sqlStatement := `SELECT * FROM users ORDER BY $3 LIMIT $1 OFFSET $2 ;`
	if limit == 0 {
		limit = 10
	}
	if page < 0 {
		page = 0
	}
	if order == "" {
		order = "'id'"
	}
	page = page * limit

	if err = db.Db.Get(&userList.Total, `SELECT COUNT(*) as total FROM users`); err != nil {
		return
	}

	userList.Users = &[]User{}
	err = db.Db.Select(userList.Users, sqlStatement, limit, page, order)
	return
}
