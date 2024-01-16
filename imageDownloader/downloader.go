package imageDownloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// DownloadImage downloads an image from the given URL to the specified directory
// and categorizes the image based on its size.
func DownloadImage(url, destDir string) error {
	tokens := strings.Split(url, "/")
	imageName := tokens[len(tokens)-1]

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Get the image size category based on Content-Length
	imageSize := getImageSize(response.ContentLength)

	// Choose the subdirectory based on the image size
	subDir := getSubdirectory(imageSize)
	subDirPath := destDir + subDir + "/"

	// Create the subdirectory if it doesn't exist
	err = os.MkdirAll(subDirPath, os.ModePerm)
	if err != nil {
		return err
	}

	filePath := subDirPath + imageName

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded %s to %s\n", imageName, subDir)
	return nil
}

// getImageSize returns a size category based on the image size
func getImageSize(contentLength int64) string {
	// This is a simplified example; you may need to adjust the thresholds based on your needs
	const smallSize = 1024 * 100  // 100 KB
	const mediumSize = 1024 * 500 // 500 KB

	if contentLength <= smallSize {
		return "small"
	} else if contentLength <= mediumSize {
		return "medium"
	} else {
		return "large"
	}
}

// getSubdirectory returns the subdirectory based on the image size
func getSubdirectory(imageSize string) string {
	return imageSize
}
