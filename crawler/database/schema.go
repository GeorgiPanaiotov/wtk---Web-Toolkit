package database

import (
	"database/sql"
)

func InitSchema(db *sql.DB) error {

	query := `
		CREATE TABLE IF NOT EXISTS pages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL UNIQUE,
			title TEXT,
			content TEXT,
			crawled_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			content_length INTEGER
		);
	`

	_, err := db.Exec(query)
	return err
}
