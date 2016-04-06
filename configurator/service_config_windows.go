package configurator

// ServiceConfig interfaces with a remote meshblu server
type ServiceConfig struct {
	uuid, workingDirectory string
}

// NewServiceConfig constructs a new Meshblu instance
func NewServiceConfig(uuid, workingDirectory string) *ServiceConfig {
	return &ServiceConfig{uuid, workingDirectory}
}

// Export the config
func (config *ServiceConfig) Export() ([]byte, error) {
	return nil, nil
}
