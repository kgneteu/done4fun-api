package models

func (db *Database) CreateUser(firstName, lastName, email, password string) (int64, error) {
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

func (db *Database) GetUserById(Id uint) (user *User, err error) {
	sqlStatement := `SELECT * FROM users WHERE ID=$1;`
	user = &User{}
	err = db.Db.Get(user, sqlStatement, Id)
	return user, err
}

func (db *Database) GetUserByEmail(email string) (user *User, err error) {
	sqlStatement := `SELECT * FROM users WHERE email=$1;`
	user = &User{}
	err = db.Db.Get(user, sqlStatement, email)
	return user, err
}

func (db *Database) GetUserList(page int, limit int, order string) (users *[]User, err error) {
	sqlStatement := `SELECT * FROM users LIMIT $2 OFFSET $3 ORDER BY $3;`
	if limit == 0 {
		limit = 10
	}
	if page < 0 {
		page = 0
	}
	if order == "" {
		order = "ID"
	}
	page = page * limit
	users = &[]User{}
	err = db.Db.Select(users, sqlStatement, limit, page, order)
	return users, err
}
