package configurator

// ServiceConfig interfaces with a remote meshblu server
type ServiceConfig struct {
	opts *Options
}

// NewServiceConfig constructs a new Meshblu instance
func NewServiceConfig(opts *Options) *ServiceConfig {
	return &ServiceConfig{opts}
}

// Export the config
func (config *ServiceConfig) Export() ([]byte, error) {
	return nil, nil
}
