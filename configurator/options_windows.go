package configurator

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// GetDefaultServiceDirectory gets the OS specific install directory
func GetDefaultServiceDirectory() string {
	return filepath.Join(os.Getenv("LOCALAPPDATA"), "MeshbluConnectors")
}

// GetExecutablePath gets the OS specific service path
func (opts *OptionsConfig) GetExecutablePath() string {
	return filepath.Join(opts.GetConnectorDirectory(), "start.exe")
}

// GetPathEnv gets the OS specific PATH env
func (opts *OptionsConfig) GetPathEnv() string {
	npmPath := filepath.Join(opts.GetBinDirectory(), "node_modules", "npm", "bin")
	return fmt.Sprintf("PATH=%s:%s:%s", opts.GetBinDirectory(), npmPath, os.Getenv("PATH"))
}

// GetServiceName gets the OS specific log directory
func (opts *OptionsConfig) GetServiceName() string {
	return fmt.Sprintf("MeshbluConnector-%s", opts.UUID)
}

// GetUserName get service display name
func (opts *OptionsConfig) GetUserName() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return currentUser.Username, nil
}
