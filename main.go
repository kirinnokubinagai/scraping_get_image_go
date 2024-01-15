package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// コマンドライン引数からサイトのURLを取得
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <siteURL>")
		return
	}
	url := os.Args[1]

	// サイトのHTMLを取得
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println("Error fetching the URL:", err)
		return
	}

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
			imgURL := getFullURL(url, imgSrc)

			// 画像をダウンロード
			err := downloadImage(imgURL, downloadDir)
			if err != nil {
				fmt.Println("Error downloading image:", err)
			}
		}
	})
}

// 相対パスの画像URLをフルパスに変換
func getFullURL(baseURL, imgURL string) string {
	if strings.HasPrefix(imgURL, "http") {
		return imgURL
	}
	if strings.HasPrefix(imgURL, "//") {
		return "https:" + imgURL
	}
	if strings.HasPrefix(imgURL, "/") {
		return baseURL + imgURL
	}
	return baseURL + "/" + imgURL
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
