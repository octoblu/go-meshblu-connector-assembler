package configurator

// ServiceConfig interfaces with a remote meshblu server
type ServiceConfig struct {
	uuid, workingDirectory, logDirectory string
}

// NewServiceConfig constructs a new Meshblu instance
func NewServiceConfig(uuid, workingDirectory, logDirectory string) *ServiceConfig {
	return &ServiceConfig{uuid, workingDirectory, logDirectory}
}

// Export the config
func (config *ServiceConfig) Export() ([]byte, error) {
	return nil, nil
}