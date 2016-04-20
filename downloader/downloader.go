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
		fmt.Println("error formating url", err.Error())
		return "", err
	}
	fmt.Println("downloading connector: ", connector, tag, platform)

	downloadFile := path.Join(client.OutputDirectory, fmt.Sprintf("connector.%s", getExt()))
	fmt.Println("to: ", downloadFile)
	outputStream, err := os.Create(downloadFile)

	if err != nil {
		fmt.Println("error opening file to write to: ", err.Error())
		return "", err
	}

	defer outputStream.Close()

	response, err := http.Get(uri)

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

	fmt.Println("downloaded connector!")

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
