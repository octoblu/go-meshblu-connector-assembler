package configurator

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetDefaultServiceDirectory gets the OS specific install directory
func GetDefaultServiceDirectory() string {
	return filepath.Join(os.Getenv("HOME"), ".octoblu", "MeshbluConnectors")
}

// GetBinDirectory gets the OS specific log directory
func (opts *OptionsConfig) GetBinDirectory() string {
	return filepath.Join(opts.GetOutputDirectory(), "bin")
}

// GetExecutablePath gets the OS specific service path
func (opts *OptionsConfig) GetExecutablePath() string {
	return filepath.Join(opts.GetConnectorDirectory(), "start")
}

// GetServiceName gets the OS specific log directory
func (opts *OptionsConfig) GetServiceName() string {
	return fmt.Sprintf("com.octoblu.%s", opts.GetUUID())
}

// GetUserName get service display name
func (opts *OptionsConfig) GetUserName() (string, error) {
	return "", nil
}
