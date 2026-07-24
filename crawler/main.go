package main

import "C"

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

//export crawler_main
func crawler_main() {
	args := os.Args
	if len(args) > 1 && args[1] == "crawler" {
		args = args[1:]
	}

	var (
		target  string
		verbose bool
	)

	for _, arg := range args[1:] {
		switch arg {
		case "-v":
			verbose = true
		case "--help":
			fmt.Printf("Usage: crawler [-v] <target_url>\n\n")
			fmt.Printf("\t -v: Verbose mode\n")
			fmt.Printf("Please provide a target url in the following format: 'https://example.com'\n")
			return
		default:
			if strings.HasPrefix(arg, "-") {
				fmt.Printf("Unknown option: %s\n", arg)
				return
			}

			if target == "" {
				target = arg
			} else {
				fmt.Printf("Unexpected argument: %s\n", arg)
				return
			}
		}
	}

	if target == "" || !strings.Contains(target, "http") {
		fmt.Println("Missing target URL")
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

	parsedURL, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.InsertPage(db, target, parsedURL.Host)
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

		if verbose {
			fmt.Printf("FETCHING: %s\n", page.URL)
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

	PrintStats(db)

	db.Close()
}

func PrintStats(db *sql.DB) {
	stats, err := database.GetStatistics(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	fmt.Println("Statistics")
	fmt.Println("------------------")
	fmt.Printf("Pages discovered : %d\n", stats.TotalPages)
	fmt.Printf("Pages crawled    : %d\n", stats.CrawledPages)
	fmt.Printf("Pages queued     : %d\n", stats.QueuedPages)
	fmt.Printf("Pages failed     : %d\n", stats.FailedPages)
	fmt.Println()
	fmt.Printf("Responses stored : %d\n", stats.ResponsesStored)
	fmt.Println()
	fmt.Printf("HTML pages       : %d\n", stats.HTMLPages)
	fmt.Printf("Non-HTML pages   : %d\n", stats.NonHTMLPages)
	fmt.Println()
	fmt.Println("Top Hosts")
	fmt.Println("---------")

	for host, count := range stats.TopHosts {
		fmt.Printf("%-30s %d\n", host, count)
	}
}

func main() {}
