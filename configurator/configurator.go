package configurator

import (
	"io/ioutil"
	"path"
)

// Configurator interfaces writing configuration files
type Configurator interface {
	WriteMeshblu() error
}

// Client interfaces the Configurator
type Client struct {
	opts Options
}

// New constructs a new Configurator client
func New(opts Options) Configurator {
	return &Client{opts}
}

// WriteMeshblu writes the configuration for meshblu
func (client *Client) WriteMeshblu() error {
	config := NewMeshbluConfig(client.opts)
	configJSON, err := config.ToJSON()
	if err != nil {
		return err
	}

	configFilePath := path.Join(client.opts.GetConnectorDirectory(), "meshblu.json")
	writeErr := ioutil.WriteFile(configFilePath, configJSON, 0644)
	if writeErr != nil {
		return writeErr
	}
	return nil
}
