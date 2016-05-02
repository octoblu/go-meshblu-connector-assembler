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

// GetExecutablePath gets the OS specific service path
func (opts *OptionsConfig) GetExecutablePath() string {
	return path.Join(opts.GetConnectorDirectory(), "start")
}

// GetPathEnv gets the OS specific PATH env
func (opts *OptionsConfig) GetPathEnv() string {
	return fmt.Sprintf("PATH=%s:%s", opts.GetBinDirectory(), os.Getenv("PATH"))
}
