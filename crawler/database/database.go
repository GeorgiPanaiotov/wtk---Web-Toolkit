package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func New() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./crawler.db")
	if err != nil {
		fmt.Print("Error while opening the database")
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(10 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		fmt.Print("Error when trying to ping database")
		return nil, err
	}

	return db, nil
}
