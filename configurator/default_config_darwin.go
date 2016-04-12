package configurator

import (
	"fmt"
	"os"
	"path"
)

// GetDefaultServiceDirectory gets the OS specific install directory
func GetDefaultServiceDirectory() string {
	return path.Join(os.Getenv("HOME"), "Library", "Application Support", "MeshbluConnectors")
}

// GetConnectorDirectory gets the OS specific connector path
func GetConnectorDirectory(opts *Options) string {
	return path.Join(opts.OutputDirectory, opts.UUID)
}

// GetBinDirectory gets the OS specific log directory
func GetBinDirectory(opts *Options) string {
	return path.Join(opts.OutputDirectory, "bin")
}

// GetLaunchFilePath gets the OS specific launch file directory
func GetLaunchFilePath(opts *Options) string {
	return path.Join(os.Getenv("HOME"), "Library", "LaunchAgents", GetServiceFileName(opts))
}

// GetLogDirectory gets the OS specific log directory
func GetLogDirectory(opts *Options) string {
	return path.Join(opts.ConnectorDirectory, "log")
}

// GetServiceName gets the OS specific service name
func GetServiceName(opts *Options) string {
	return fmt.Sprintf("com.octoblu.%s", opts.UUID)
}

// GetServiceFileName gets the OS specific service file
func GetServiceFileName(opts *Options) string {
	return fmt.Sprintf("%s.plist", GetServiceName(opts))
}

// GetServiceFilePath gets the OS specific service path
func GetServiceFilePath(opts *Options) string {
	return path.Join(opts.ConnectorDirectory, GetServiceFileName(opts))
}
