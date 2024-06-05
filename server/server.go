package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type config struct {
	port int
}

type application struct {
	config    config
	logger    *log.Logger
	errlogger *log.Logger
	router    *httprouter.Router
	db        *sql.DB
}

func Start() error {
	var cfg config
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE,
			password TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	intPort, err := strconv.Atoi(port)
	if err != nil {
		intPort = 8080
	}

	cfg.port = intPort
	logger := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
		db:     db,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  45 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}

	logger.Printf("Server is starting on %s\n", srv.Addr)

	if err = srv.ListenAndServe(); err != nil {
		logger.Printf("Server failed to start: %v\n", err)
		os.Exit(1)
	}

	return err
}
