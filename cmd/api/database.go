package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

func OpenDB(dsn string, logger *log.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	logger.Println("Successfully connected to database")
	return db, nil
}
