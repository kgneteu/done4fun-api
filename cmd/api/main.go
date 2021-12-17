package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"server/models"
	"strconv"
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

	defEnv := os.Getenv("ENV_TYPE")
	var err error
	if defEnv == "" {
		defEnv = "development"
		err = godotenv.Load(".env.development")
	} else {
		err = godotenv.Load(".env")
	}

	if err != nil {
		println(err.Error())
	}
	defPort, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))

	defDSN := os.Getenv("POSTGRES_DSN")

	flag.IntVar(&cfg.port, "port", defPort, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", defEnv, "Application environment (development|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", defDSN, "Postgres connection string")
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
