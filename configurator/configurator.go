package configurator

import (
	"io/ioutil"
	"path/filepath"
)

// Configurator interfaces writing configuration files
type Configurator interface {
	WriteConfigs() error
	WriteMeshbluConfig() error
	WriteServiceConfig() error
}

// Client interfaces the Configurator
type Client struct {
	opts Options
}

// New constructs a new Configurator client
func New(opts Options) Configurator {
	return &Client{opts}
}

// WriteConfigs writes the configuration for meshblu
func (client *Client) WriteConfigs() error {
	err := client.WriteMeshbluConfig()
	if err != nil {
		return err
	}
	err = client.WriteServiceConfig()
	if err != nil {
		return err
	}
	return nil
}

// WriteMeshbluConfig writes the configuration for meshblu
func (client *Client) WriteMeshbluConfig() error {
	config := NewMeshbluConfig(client.opts)
	configJSON, err := config.ToJSON()
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(client.opts.GetConnectorDirectory(), "meshblu.json")
	return ioutil.WriteFile(configFilePath, configJSON, 0644)
}

// WriteServiceConfig writes the configuration for the service
func (client *Client) WriteServiceConfig() error {
	config := NewServiceConfig(client.opts)
	configJSON, err := config.ToJSON()
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(client.opts.GetConnectorDirectory(), "service.json")
	writeErr := ioutil.WriteFile(configFilePath, configJSON, 0644)
	if writeErr != nil {
		return writeErr
	}
	return nil
}
