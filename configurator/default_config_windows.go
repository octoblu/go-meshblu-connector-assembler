package configurator

import (
	"fmt"
	"os"
	"path"
)

// GetDefaultServiceDirectory gets the OS specific install directory
func GetDefaultServiceDirectory() string {
	return path.Join(os.Getenv("LOCALAPPDATA"), "MeshbluConnectors")
}

// GetConnectorDirectory gets the OS specific connector path
func GetConnectorDirectory(opts *Options) string {
	return path.Join(opts.OutputDirectory, opts.UUID)
}

// GetLogDirectory gets the OS specific log directory
func GetLogDirectory(opts *Options) string {
	return path.Join(GetConnectorDirectory(opts), "log")
}

// GetBinDirectory gets the OS specific log directory
func GetBinDirectory(opts *Options) string {
	return path.Join(opts.OutputDirectory, "bin")
}

// GetServiceName gets the OS specific service file
func GetServiceName(opts *Options) string {
	return fmt.Sprintf("MeshbluConnector-%s", opts.UUID)
}
