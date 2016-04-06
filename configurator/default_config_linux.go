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
func GetConnectorDirectory(outputDirectory, uuid string) string {
	return path.Join(outputDirectory, uuid)
}

// GetLogDirectory gets the OS specific log directory
func GetLogDirectory(outputDirectory, uuid string) string {
	return path.Join(GetConnectorDirectory(outputDirectory, uuid), "log")
}

// GetBinDirectory gets the OS specific log directory
func GetBinDirectory(outputDirectory string) string {
	return path.Join(outputDirectory, "bin")
}
