package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/models"
	"time"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *log.Logger
	models models.DBModel
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 9000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", `postgres://ufo:!!!UfO:-)1234!!!@localhost:/done4fun?sslmode=disable`, "Postgres connection string")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := OpenDB(cfg.db.dsn, logger)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.DBModel{Db: db},
	}

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("HTTP server is listening on port", cfg.port)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal(err)
	}

}
