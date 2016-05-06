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

// GetBinDirectory gets the OS specific log directory
func (opts *OptionsConfig) GetBinDirectory() string {
	return filepath.Join(opts.GetOutputDirectory(), "bin")
}

// GetExecutablePath gets the OS specific service path
func (opts *OptionsConfig) GetExecutablePath() string {
	return filepath.Join(opts.GetConnectorDirectory(), "start.exe")
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
