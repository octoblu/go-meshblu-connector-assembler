package configurator

import (
	"encoding/json"
	"path"
)

// ServiceConfig is the runner connector config structure.
type ServiceConfig struct {
	ServiceName   string
	DisplayName   string
	Description   string
	ConnectorName string
	Legacy        bool
	Dir           string
	Env           []string

	Stderr, Stdout string
}

// NewServiceConfig constructs a new Meshblu instance
func NewServiceConfig(opts Options) *ServiceConfig {
	env := []string{
		opts.GetPathEnv(),
	}
	return &ServiceConfig{
		ServiceName:   opts.GetServiceName(),
		DisplayName:   opts.GetDisplayName(),
		Description:   opts.GetDescription(),
		ConnectorName: opts.GetConnector(),
		Legacy:        opts.GetLegacy(),
		Dir:           opts.GetConnectorDirectory(),
		Env:           env,
		Stderr:        path.Join(opts.GetLogDirectory(), "connector-error.log"),
		Stdout:        path.Join(opts.GetLogDirectory(), "connector.log"),
	}
}

// ToJSON serializes the object to the meshblu.json format
func (config *ServiceConfig) ToJSON() ([]byte, error) {
	return json.Marshal(config)
}
