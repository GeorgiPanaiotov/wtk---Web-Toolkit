package main

import (
	"crawler/database"
	"crawler/spider"
	"fmt"
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

	err = database.InsertPageRecord(db, res)
	if err != nil {
		fmt.Println(err)
	}

	pages, err := database.GetAllPages(db)

	if err != nil {
		log.Fatal(err)
	}

	for _, page := range pages {
		fmt.Println(page.URL)
	}
}
