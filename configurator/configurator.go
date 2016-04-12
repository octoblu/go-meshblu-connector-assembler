package configurator

import (
	"io/ioutil"
	"path"

	"github.com/codegangsta/cli"
)

// Options defines the command line arguments
type Options struct {
	Connector          string
	Tag                string
	Legacy             bool
	OutputDirectory    string
	ConnectorDirectory string
	LogDirectory       string
	BinDirectory       string
	ServiceName        string
	Hostname           string
	Port               int
	UUID, Token        string
}

// Configurator interfaces writing configuration files
type Configurator interface {
	WriteMeshblu() error
}

// Client interfaces the Configurator
type Client struct {
	opts *Options
}

// New constructs a new Configurator client
func New(opts *Options) Configurator {
	return &Client{opts}
}

// NewOptions should create an options points
func NewOptions(context *cli.Context) *Options {
	return &Options{
		context.String("connector"),
		context.String("tag"),
		context.Bool("legacy"),
		context.String("output"),
		"",
		"",
		"",
		"",
		context.String("hostname"),
		context.Int("port"),
		context.String("uuid"),
		context.String("token"),
	}
}

// WriteMeshblu writes the configuration for meshblu
func (client *Client) WriteMeshblu() error {
	config := NewConfig(client.opts)
	configJSON, err := config.ToJSON()
	if err != nil {
		return err
	}

	configFilePath := path.Join(client.opts.ConnectorDirectory, "meshblu.json")
	writeErr := ioutil.WriteFile(configFilePath, configJSON, 0644)
	if writeErr != nil {
		return writeErr
	}
	return nil
}
