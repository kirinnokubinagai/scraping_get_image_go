// main.go
package main

import (
	"fmt"
	"os"

	"scraping_go/crawler"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <siteURL>")
		return
	}
	startURL := os.Args[1]

	crawler.CrawlSite(startURL)
}
