package extractor

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extractor ungzips and untars the source to the target
type Extractor interface {
	Do(downloadFile, target string) error
	Ungzip(source, target string) error
	Untar(tarball, target string) error
}

// Client interfaces with the Extractor
type Client struct {
}

// New constructs a new Extractor
func New() Extractor {
	return &Client{}
}

// Do extracts the tar.gz file
func (client *Client) Do(downloadFile, target string) error {
	tarFile := strings.Replace(downloadFile, "tar.gz", "tar", 1)
	ungzipErr := client.Ungzip(downloadFile, target)
	if ungzipErr != nil {
		return ungzipErr
	}
	untarErr := client.Untar(tarFile, target)
	if untarErr != nil {
		return untarErr
	}
	removeDownloadErr := os.Remove(downloadFile)
	if removeDownloadErr != nil {
		return removeDownloadErr
	}
	removeTarFile := os.Remove(tarFile)
	if removeTarFile != nil {
		return removeTarFile
	}
	return nil
}

// Ungzip the source to the target
func (client *Client) Ungzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("Ungzip, os.Open: %v", err.Error())
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("Ungzip, gzip.NewReader: %v", err.Error())
	}
	defer archive.Close()

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("Ungzip, os.Create: %v", err.Error())
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}

// Untar the source to the target
func (client *Client) Untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return fmt.Errorf("Untar, os.Open: %v", err.Error())
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Untar, tarReader.Next: %v", err.Error())
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return fmt.Errorf("Untar, os.MkdirAll: %v", err.Error())
			}
			continue
		}

		if header.Typeflag == tar.TypeSymlink {
			os.Remove(path)
			err = os.Symlink(header.Linkname, path)
			if err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return fmt.Errorf("Untar, os.OpenFile: %v", err.Error())
		}

		_, err = io.Copy(file, tarReader)
		if err != nil {
			file.Close()
			return fmt.Errorf("Untar, io.Copy: %v", err.Error())
		}
		file.Close()
	}
	return nil
}
