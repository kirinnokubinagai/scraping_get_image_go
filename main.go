package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var visited = make(map[string]bool)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <siteURL>")
		return
	}
	startURL := os.Args[1]

	crawlSite(startURL)
}

func crawlSite(urlStr string) {
	// 既に訪れたURLは再帰的に処理しない
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

	// サイトのHTMLを取得
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("Error fetching the URL:", err)
		return
	}

	// このページでやりたい処理をここに追加
	fmt.Println("Visited:", urlStr)

	// 画像をダウンロードするためのディレクトリ
	downloadDir := "./downloaded_images/"
	err = os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating download directory:", err)
		return
	}

	// 画像のセレクタを変更してください
	imageSelector := "img"

	// 画像のURLを取得し、ダウンロード
	doc.Find(imageSelector).Each(func(i int, s *goquery.Selection) {
		imgSrc, exists := s.Attr("src")
		if exists {
			// 画像のフルパスを取得
			imgURL := getAbsoluteURL(u, imgSrc)

			// 画像をダウンロード
			err := downloadImage(imgURL, downloadDir)
			if err != nil {
				fmt.Println("Error downloading image:", err)
			}
		}
	})

	// このページを訪れたことを記録
	visited[urlStr] = true

	// ページ内のリンクを再帰的に処理
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists {
			// リンクの正規化
			absURL := getAbsoluteURL(u, link)

			// 同じサイトのページかどうか確認
			if isSameSite(u, absURL) {
				crawlSite(absURL)
			}
		}
	})
}

// 相対パスや不正な形式のURLを絶対URLに変換
func getAbsoluteURL(base *url.URL, href string) string {
	relURL, err := url.Parse(href)
	if err != nil {
		return ""
	}
	absURL := base.ResolveReference(relURL)
	return absURL.String()
}

// 同じサイトかどうかを確認
func isSameSite(base *url.URL, target string) bool {
	targetURL, err := url.Parse(target)
	if err != nil {
		return false
	}
	return base.Host == targetURL.Host
}

// 画像をダウンロード
func downloadImage(url, destDir string) error {
	// 画像の名前を取得
	tokens := strings.Split(url, "/")
	imageName := tokens[len(tokens)-1]

	// 画像の保存先ファイル
	filePath := destDir + imageName

	// 画像のダウンロード
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 画像の保存
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	fmt.Println("Downloaded:", imageName)
	return nil
}
