package database

import (
	"bytes"
	"database/sql"
	"io"
	"net/http"
)

type Page struct {
	ID            int
	URL           string
	Host          string
	CrawledAt     string
	PageID        int
	StatusCode    int64
	ContentType   string
	ContentLength int64
	Body          string
	ResponseID    int64
	Headers       string
	HeaderID      int64
}

func InsertPageRecord(db *sql.DB, res *http.Response) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	pageResult, err := InsertPage(db, res.Request.URL.String(), res.Request.URL.Host)
	if err != nil {
		return err
	}

	page_id, err := pageResult.LastInsertId()
	if err != nil {
		return err
	}

	responseResult, err := InsertResponse(db, page_id, int64(res.StatusCode), res.Header.Get("Content-Type"), res.ContentLength, string(body))
	if err != nil {
		return err
	}

	response_id, err := responseResult.LastInsertId()
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	err = res.Header.Write(&buffer)
	if err != nil {
		return err
	}

	err = InsertHeaders(db, response_id, buffer.String())
	if err != nil {
		return err
	}

	return err
}

func InsertPage(db *sql.DB, url string, host string) (sql.Result, error) {
	query := `
		INSERT INTO pages (url, host)
		VALUES (?, ?);
	`

	sqlResult, err := db.Exec(query, url, host)
	if err != nil {
		return nil, err
	}
	return sqlResult, err
}

func InsertResponse(db *sql.DB, page_id int64, status_code int64, content_type string, content_length int64, body string) (sql.Result, error) {
	query := `
		INSERT INTO responses (page_id, status_code, content_type, content_length, body)
		VALUES (?, ?, ?, ?, ?)
	`

	sqlResult, err := db.Exec(query, page_id, status_code, content_type, content_length, body)
	if err != nil {
		return nil, err
	}
	return sqlResult, err
}

func InsertHeaders(db *sql.DB, response_id int64, headers string) error {
	query := `
		INSERT INTO headers (response_id, headers)
		VALUES (?, ?)
	`

	_, err := db.Exec(query, response_id, headers)
	if err != nil {
		return err
	}
	return err
}

func GetAllPages(db *sql.DB) ([]Page, error) {
	query := `
		SELECT 
			p.id,
			p.url,
			p.host,
			p.crawled_at,
			r.id,
			r.status_code,
			r.content_type,
			r.content_length,
			r.body,
			h.id,
			h.headers
		FROM pages p
		JOIN responses r
		ON p.id = r.page_id
		JOIN headers h
		ON r.id = h.response_id;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []Page

	for rows.Next() {
		var page Page
		err := rows.Scan(
			&page.ID,
			&page.URL,
			&page.Host,
			&page.CrawledAt,
			&page.ResponseID,
			&page.StatusCode,
			&page.ContentType,
			&page.ContentLength,
			&page.Body,
			&page.HeaderID,
			&page.Headers)

		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	return pages, rows.Err()
}
