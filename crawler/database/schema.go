package database

import (
	"database/sql"
)

func InitSchema(db *sql.DB) error {
	err := InitPages(db)
	err = InitResponses(db)
	err = InitHeaders(db)
	err = ResetFetching(db)

	return err
}

func InitPages(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS pages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL UNIQUE,
			host TEXT,
			crawled_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			status TEXT NOT NULL DEFAULT 'queued'
		);
	`

	_, err := db.Exec(query)
	return err
}

func InitResponses(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS responses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			page_id INTEGER NOT NULL,
			status_code INTEGER,
			content_type TEXT,
			content_length INTEGER,
			body TEXT,

			FOREIGN KEY(page_id) REFERENCES pages(id)
		);
	`

	_, err := db.Exec(query)
	return err
}

func InitHeaders(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS headers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			response_id INTEGER NOT NULL,
			headers TEXT,

			FOREIGN KEY(response_id) REFERENCES responses(id)
		);
	`

	_, err := db.Exec(query)
	return err
}

func ResetFetching(db *sql.DB) error {
	query := `
		UPDATE pages
		SET status = 'queued'
		WHERE status = 'fetching';
	`

	_, err := db.Exec(query)
	return err
}
