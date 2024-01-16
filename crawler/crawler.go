package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"scraping_go/imageDownloader"
	"scraping_go/utils"

	"github.com/PuerkitoBio/goquery"
)

var visited = make(map[string]bool)

// CrawlSite is the entry point for crawling a site
func CrawlSite(urlStr string) {
	if visited[urlStr] {
		return
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	res, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("Error Not Found URL:", err)
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("Error fetching the URL:", err)
		return
	}

	fmt.Println("Visited:", urlStr)

	downloadDir := "./downloaded_images/"
	err = os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating download directory:", err)
		return
	}

	imageSelector := "img"

	downloadImages(u, doc, downloadDir, imageSelector)

	visited[urlStr] = true

	processLinks(u, doc)
}

// downloadImages downloads images from the page
func downloadImages(baseURL *url.URL, doc *goquery.Document, destDir, imageSelector string) {
	doc.Find(imageSelector).Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("src")
		if exists {
			imgURL := utils.GetAbsoluteURL(baseURL, imgSrc)

			err := imageDownloader.DownloadImage(imgURL, destDir)
			if err != nil {
				fmt.Println("Error downloading image:", err)
			}
		}
	})
}

// processLinks processes links on the page recursively
func processLinks(baseURL *url.URL, doc *goquery.Document) {
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			absURL := utils.GetAbsoluteURL(baseURL, link)

			if utils.IsSameSite(baseURL, absURL) {
				CrawlSite(absURL)
			}
		}
	})
}
