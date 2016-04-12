package configurator

import (
	"os"
	"path"
)

// GetDefaultServiceDirectory gets the OS specific install directory
func GetDefaultServiceDirectory() string {
	return path.Join(os.Getenv("HOME"), ".octoblu", "MeshbluConnectors")
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

// GetServiceFileName gets the OS specific log directory
func GetServiceFileName(opts *Options) string {
	return ""
}

// GetServiceFileName gets the OS specific log directory
func GetServiceName(opts *Options) string {
	return ""
}
