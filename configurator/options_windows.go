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

// GetExecutablePath gets the OS specific service path
func (opts *OptionsConfig) GetExecutablePath() string {
	return path.Join(opts.GetConnectorDirectory(), "start.exe")
}

// GetPathEnv gets the OS specific PATH env
func (opts *OptionsConfig) GetPathEnv() string {
	npmPath := path.Join(opts.GetBinDirectory(), "node_modules", "npm", "bin")
	return fmt.Sprintf("PATH=%s:%s:%s", opts.GetBinDirectory(), npmPath, os.Getenv("PATH"))
}
