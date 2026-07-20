package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func New() (*sql.DB, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(
		home,
		".local",
		"share",
		"wtk",
		"db",
		"crawler.db",
	)

	db, err := sql.Open("sqlite3", dbPath)
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
