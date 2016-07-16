package serviceconfig

import (
	"encoding/json"
	"path/filepath"

	"github.com/spf13/afero"
)

type Options struct {
	ServiceName   string
	DisplayName   string
	Description   string
	ConnectorName string
	GithubSlug    string
	Tag           string
	BinPath       string
	Dir           string

	Stderr, Stdout string
}

// Write a ServiceConfig JSON to the file system
func Write(options Options) error {
	return WriteWithFS(options, afero.NewOsFs())
}

// WriteWithFS does everything Write does, only you get to supply
// your own FileSystem!
func WriteWithFS(options Options, fs afero.Fs) error {
	data, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return err
	}

	serviceConfigPath := filepath.Join(options.Dir, "service.json")
	err = afero.WriteFile(fs, serviceConfigPath, data, 0644)

	if err != nil {
		return err
	}

	return nil
}
