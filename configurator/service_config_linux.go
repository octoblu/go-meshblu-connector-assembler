package configurator

// ServiceConfig interfaces with a remote meshblu server
type ServiceConfig struct {
	uuid, connector, workingDirectory string
}

// NewServiceConfig constructs a new Meshblu instance
func NewServiceConfig(uuid, connector, workingDirectory string) *ServiceConfig {
	return &ServiceConfig{uuid, connector, workingDirectory}
}

// Export the config
func (config *ServiceConfig) Export() ([]byte, error) {
	return nil, nil
}
