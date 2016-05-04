package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Downloader interface with a way of downloading connector bundles
type Downloader interface {
	Download(downloadURI string) (string, error)
}

// Client interfaces with remote cdn
type Client struct {
	OutputDirectory string
}

// New constructs new Downloader instance
func New(OutputDirectory string) Downloader {
	return &Client{OutputDirectory}
}

// Download downloads the connector the local directory
func (client *Client) Download(downloadURI string) (string, error) {
	fmt.Println("downloading connector: ", downloadURI)

	downloadFile := client.getDownloadFile(downloadURI)
	fmt.Println("to: ", downloadFile)
	outputStream, err := os.Create(downloadFile)

	if err != nil {
		fmt.Println("error opening file to write to: ", err.Error())
		return "", err
	}

	defer outputStream.Close()

	response, err := http.Get(downloadURI)

	if err != nil {
		fmt.Println("http error downloading: ", err.Error())
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("download returned invalid response code: %v", response.StatusCode)
	}

	_, err = io.Copy(outputStream, response.Body)

	if err != nil {
		fmt.Println("error downloading to file", err.Error())
		return "", err
	}

	fmt.Println("downloaded!")

	return downloadFile, nil
}

func (client *Client) getDownloadFile(downloadURI string) string {
	fileName := getFileName(downloadURI)
	return filepath.Join(client.OutputDirectory, fileName)
}

func getFileName(source string) string {
	segments := strings.Split(source, "/")
	return segments[len(segments)-1]
}
