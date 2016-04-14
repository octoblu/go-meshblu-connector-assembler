package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
)

// Downloader interface with a way of downloading connector bundles
type Downloader interface {
	DownloadConnector(connector string, tag string, platform string) (string, error)
	buildURI(connector string, tag string, platform string) (string, error)
}

// Client interfaces with remote cdn
type Client struct {
	OutputDirectory string
	baseURI         string
}

// New constructs new Downloader instance
func New(OutputDirectory string, baseURI string) Downloader {
	return &Client{OutputDirectory, baseURI}
}

// DownloadConnector downloads the connector the local directory
func (client *Client) DownloadConnector(connector string, tag string, platform string) (string, error) {
	uri, err := client.buildURI(connector, tag, platform)
	if err != nil {
		fmt.Println("Error on client.buildURI", err.Error())
		return "", err
	}
	fmt.Println("Downloading connector...", uri)

	downloadFile := path.Join(client.OutputDirectory, fmt.Sprintf("connector.%s", getExt()))
	outputStream, err := os.Create(downloadFile)

	if err != nil {
		fmt.Println("Error on os.Create", err.Error())
		return "", err
	}

	defer outputStream.Close()

	response, err := http.Get(uri)

	if err != nil {
		fmt.Println("Error on http.Get", err.Error())
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Download returned invalid response code: %v", response.StatusCode)
	}

	_, err = io.Copy(outputStream, response.Body)

	if err != nil {
		fmt.Println("Error on io.Copy", err.Error())
		return "", err
	}

	return downloadFile, nil
}

func (client *Client) buildURI(connector string, tag string, platform string) (string, error) {
	uri, err := url.Parse(client.baseURI)
	if err != nil {
		return "", err
	}

	uri.Path = fmt.Sprintf("/connectors/%v/%v/%v.bundle.%s", connector, tag, platform, getExt())
	return uri.String(), nil
}

func getExt() string {
	if runtime.GOOS == "windows" {
		return "zip"
	}
	return "tar.gz"
}
