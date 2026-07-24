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

type Stats struct {
	TotalPages      int
	CrawledPages    int
	QueuedPages     int
	FailedPages     int
	ResponsesStored int
	HTMLPages       int
	NonHTMLPages    int
	TopHosts        map[string]int
}

func InsertPageRecord(db *sql.DB, pageID int, res *http.Response) ([]byte, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	responseResult, err := InsertResponse(db, int64(pageID), int64(res.StatusCode), res.Header.Get("Content-Type"), res.ContentLength, string(body))
	if err != nil {
		return nil, err
	}

	response_id, err := responseResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = res.Header.Write(&buffer)
	if err != nil {
		return nil, err
	}

	err = InsertHeaders(db, response_id, buffer.String())
	if err != nil {
		return nil, err
	}

	return body, err
}

func InsertPage(db *sql.DB, url string, host string) (sql.Result, error) {
	query := `
		INSERT OR IGNORE INTO pages (url, host)
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

func GetNextPage(db *sql.DB) (*Page, error) {
	query := `
		SELECT
			id,
			url
		FROM pages
		WHERE status = 'queued'
		LIMIT 1;
	`

	var page Page
	err := db.QueryRow(query).Scan(&page.ID, &page.URL)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

func UpdateStatus(db *sql.DB, id int, status string) error {
	query := `
		UPDATE pages
		SET status = ?
		WHERE id = ?
	`
	_, err := db.Exec(query, status, id)

	return err
}

func GetStatistics(db *sql.DB) (*Stats, error) {
	stats := &Stats{}

	err := db.QueryRow(`SELECT COUNT (*) FROM pages; `).Scan(&stats.TotalPages)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`SELECT COUNT(*) FROM pages WHERE status = 'done'; `).Scan(&stats.CrawledPages)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`SELECT COUNT(*) FROM pages WHERE status = 'queued'; `).Scan(&stats.QueuedPages)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`SELECT COUNT(*) FROM pages WHERE status = 'failed'; `).Scan(&stats.FailedPages)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`SELECT COUNT(*) FROM responses; `).Scan(&stats.ResponsesStored)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`SELECT COUNT(*) FROM responses WHERE content_type LIKE 'text/html%'; `).Scan(&stats.HTMLPages)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`SELECT COUNT(*) FROM responses WHERE content_type NOT LIKE 'text/html%'; `).Scan(&stats.NonHTMLPages)
	if err != nil {
		return nil, err
	}

	stats.TopHosts = make(map[string]int)

	rows, err := db.Query(`SELECT host, COUNT(*) FROM pages GROUP BY host ORDER BY COUNT(*) DESC LIMIT 20; `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var host string
		var count int

		if err := rows.Scan(&host, &count); err != nil {
			return nil, err
		}

		stats.TopHosts[host] = count
	}

	return stats, rows.Err()
}
