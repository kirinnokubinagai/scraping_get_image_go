package imageDownloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// DownloadImage downloads an image from the given URL to the specified directory
func DownloadImage(url, destDir string) error {
	tokens := strings.Split(url, "/")
	imageName := tokens[len(tokens)-1]

	filePath := destDir + imageName

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

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
