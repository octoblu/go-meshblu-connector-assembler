package configurator

import (
	"io/ioutil"
	"path"
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

// WriteConfig writes the configuration for meshblu
func (client *Client) WriteConfig() error {
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

	configFilePath := path.Join(client.opts.GetConnectorDirectory(), "meshblu.json")
	writeErr := ioutil.WriteFile(configFilePath, configJSON, 0644)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

// WriteServiceConfig writes the configuration for the service
func (client *Client) WriteServiceConfig() error {
	config := NewServiceConfig(client.opts)
	configJSON, err := config.ToJSON()
	if err != nil {
		return err
	}

	configFilePath := path.Join(client.opts.GetConnectorDirectory(), "service.json")
	writeErr := ioutil.WriteFile(configFilePath, configJSON, 0644)
	if writeErr != nil {
		return writeErr
	}
	return nil
}
