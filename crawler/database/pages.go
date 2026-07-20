package database

import (
	"database/sql"
)

type Page struct {
	ID            int
	URL           string
	Title         string
	Content       string
	CrawledAt     string
	ContentLength int
}

func InsertPage(db *sql.DB, url string, title string, content string, content_length int64) error {
	query := `
		INSERT INTO pages (url, title, content, content_length)
		VALUES (?, ?, ?, ?);
	`

	_, err := db.Exec(query, url, title, content, content_length)
	return err
}

func GetAllPages(db *sql.DB) ([]Page, error) {
	query := `
		SELECT 
			id,
			url,
			title,
			content,
			crawled_at,
			content_length
		FROM pages;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []Page

	for rows.Next() {
		var page Page
		err := rows.Scan(&page.ID, &page.URL, &page.Title, &page.Content, &page.CrawledAt, &page.ContentLength)

		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	return pages, rows.Err()
}
