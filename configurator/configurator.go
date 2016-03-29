package configurator

import (
	"io/ioutil"
	"path"
)

// Configurator interfaces writing configuration files
type Configurator interface {
	WriteMeshblu(UUID, token, hostname string, port int) error
}

// Client interfaces the Configurator
type Client struct {
	outputDirectory string
}

// New constructs a new Configurator client
func New(outputDirectory string) Configurator {
	return &Client{outputDirectory}
}

// WriteMeshblu writes the configuration for meshblu
func (client *Client) WriteMeshblu(UUID, Token, Hostname string, Port int) error {
	config := NewConfig(UUID, Token, Hostname, Port)
	configJSON, err := config.ToJSON()
	if err != nil {
		return err
	}

	configFilePath := path.Join(client.outputDirectory, "meshblu.json")
	writeErr := ioutil.WriteFile(configFilePath, configJSON, 0644)
	if writeErr != nil {
		return writeErr
	}
	return nil
}
