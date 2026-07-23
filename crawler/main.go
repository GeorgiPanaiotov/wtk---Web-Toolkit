package main

import (
	"bytes"
	"crawler/database"
	"crawler/spider"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("Please provide a target url in the following format: 'https://example.com'\n")
		return
	}
	db, err := database.New()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = database.InitSchema(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Database Ready!\n")

	parsedURL, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.InsertPage(db, os.Args[1], parsedURL.Host)
	if err != nil {
		log.Fatal(err)
	}

	for {
		page, err := database.GetNextPage(db)
		if err == sql.ErrNoRows {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		err = database.UpdateStatus(db, page.ID, "fetching")
		if err != nil {
			log.Fatal(err)
		}

		res, err := spider.Fetch(page.URL)
		if err != nil {
			err = database.UpdateStatus(db, page.ID, "failed")
			if err != nil {
				log.Fatal(err)
			}

			continue
		}

		if res.StatusCode == 404 {
			err = database.UpdateStatus(db, page.ID, "done")
			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		var content_type = res.Header.Get("Content-Type")

		if !strings.Contains(content_type, "text/html") {
			fmt.Print("The page doesn't contain HTML! Nothing will be written in the database\n")
			err = database.UpdateStatus(db, page.ID, "done")
			if err != nil {
				log.Fatal(err)
			}
		} else {
			body, err := database.InsertPageRecord(db, page.ID, res)
			if err != nil {
				err = database.UpdateStatus(db, page.ID, "failed")
				if err != nil {
					log.Fatal(err)
				}
				continue
			}
			res.Body.Close()

			doc, err := html.Parse(bytes.NewReader(body))
			if err != nil {
				log.Fatal(err)
			}

			links, err := spider.ExtractLinks(res.Request.URL, doc)
			if err != nil {
				log.Fatal(err)
			}

			err = database.UpdateStatus(db, page.ID, "done")
			if err != nil {
				log.Fatal(err)
			}

			for _, link := range links {
				parsed, err := url.Parse(link)
				if err != nil {
					continue
				}

				if parsed.Host != parsedURL.Host {
					continue
				}

				if parsed.Scheme != "http" && parsed.Scheme != "https" {
					continue
				}

				_, err = database.InsertPage(db, parsed.String(), parsed.Host)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	db.Close()
}
