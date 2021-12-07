package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "ufo"
	dbPassword = "!!!UfO:-)1234!!!"
	dbDbname   = "done4fun"
)

var db *sql.DB

func connectToDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbDbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}


	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")
	return db
}

func createUser(firstName, lastName, email, password string) (int64, error) {
	var err error
	var id int64

	sqlStatement := `
		INSERT INTO users (first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	err = db.QueryRow(sqlStatement, firstName, lastName, email, password).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func getUserId(email string, password string) (int64, error) {
	var err error
	var id int64
	sqlStatement := `SELECT id FROM users WHERE email=$1 AND password=$2;`
	err = db.QueryRow(sqlStatement, email, password).Scan(&id)
	switch err {
	case nil:
		return id, nil
	default:
		return 0, err
	}
}