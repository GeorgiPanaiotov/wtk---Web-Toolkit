package main

import (
	"crawler/database"
	"crawler/spider"
	"fmt"
	"io"
	"log"
	"os"
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

	res, err := spider.Fetch(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	database.InsertPage(db, res.Request.URL.String(), res.Request.URL.Host, string(body), res.ContentLength)

	pages, err := database.GetAllPages(db)

	if err != nil {
		log.Fatal(err)
	}

	for _, page := range pages {
		fmt.Println(page.URL)
	}
}
