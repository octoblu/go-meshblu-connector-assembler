package configurator

import (
	"encoding/json"
	"path"
)

// ServiceConfig is the runner connector config structure.
type ServiceConfig struct {
	ConnectorName string
	Legacy        bool
	Dir           string
	Args          []string
	Env           []string

	Stderr, Stdout string
}

// NewServiceConfig constructs a new Meshblu instance
func NewServiceConfig(opts Options) *ServiceConfig {
	args := []string{
		opts.GetLegacyFlag(),
	}
	env := []string{
		opts.GetPathEnv(),
	}
	return &ServiceConfig{
		ConnectorName: opts.GetConnector(),
		Legacy:        opts.GetLegacy(),
		Dir:           opts.GetConnectorDirectory(),
		Args:          args,
		Env:           env,
		Stderr:        path.Join(opts.GetLogDirectory(), "connector-error.log"),
		Stdout:        path.Join(opts.GetLogDirectory(), "connector.log"),
	}
}

// ToJSON serializes the object to the meshblu.json format
func (config *ServiceConfig) ToJSON() ([]byte, error) {
	return json.Marshal(config)
}
