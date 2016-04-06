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
func GetConnectorDirectory(outputDirectory, uuid string) string {
	return path.Join(outputDirectory, uuid)
}

// GetBinDirectory gets the OS specific log directory
func GetBinDirectory(outputDirectory string) string {
	return path.Join(outputDirectory, "bin")
}

// GetLaunchFilePath gets the OS specific launch file directory
func GetLaunchFilePath(uuid string) string {
	return path.Join(os.Getenv("HOME"), "Library", "LaunchAgents", GetServiceFileName(uuid))
}

// GetLogDirectory gets the OS specific log directory
func GetLogDirectory(outputDirectory, uuid string) string {
	return path.Join(GetConnectorDirectory(outputDirectory, uuid), "log")
}

// GetServiceFileName gets the OS specific service file
func GetServiceFileName(uuid string) string {
	return fmt.Sprintf("com.octoblu.%s.plist", uuid)
}

// GetServiceFilePath gets the OS specific service path
func GetServiceFilePath(uuid, outputDirectory string) string {
	return path.Join(outputDirectory, GetServiceFileName(uuid))
}
