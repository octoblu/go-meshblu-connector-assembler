package configurator

import (
	"fmt"
	"os"
	"path"
)

// GetDefaultServiceDirectory gets the OS specific install directory
func GetDefaultServiceDirectory() string {
	return path.Join(os.Getenv("HOME"), ".octoblu", "MeshbluConnectors")
}

// GetExecutablePath gets the OS specific service path
func (opts *OptionsConfig) GetExecutablePath() string {
	return path.Join(opts.GetConnectorDirectory(), "start")
}

// GetPathEnv gets the OS specific PATH env
func (opts *OptionsConfig) GetPathEnv() string {
	return fmt.Sprintf("PATH=%s:%s", opts.GetBinDirectory(), os.Getenv("PATH"))
}

// GetServiceName gets the OS specific log directory
func (opts *OptionsConfig) GetServiceName() string {
	return fmt.Sprintf("MeshbluConnector-%s", opts.UUID)
}

// GetUserName get service display name
func (opts *OptionsConfig) GetUserName() (string, error) {
	return "", nil
}
