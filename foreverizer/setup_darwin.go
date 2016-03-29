package foreverizer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/octoblu/meshblu-connector-installer-go/configurator"
)

// Setup configures the os to the device
func Setup(uuid, outputDirectory string) error {
	logDirectory := path.Join(outputDirectory, "log")
	err := os.MkdirAll(logDirectory, 0777)
	if err != nil {
		return err
	}
	serviceConfig := configurator.NewServiceConfig(uuid, outputDirectory, logDirectory)
	fileBytes, err := serviceConfig.Export()
	if err != nil {
		return err
	}

	serviceFileName := fmt.Sprintf("com.octoblu.%s.plist", uuid)
	serviceFilePath := path.Join(outputDirectory, serviceFileName)
	err = ioutil.WriteFile(serviceFilePath, fileBytes, 0644)
	if err != nil {
		return err
	}

	launchFilePath := path.Join(os.Getenv("HOME"), "Library/LaunchAgents", serviceFileName)
	_, err = os.Stat(launchFilePath)

	fileExists, err := filePathExists(launchFilePath)
	if err != nil {
		return err
	}

	if fileExists {
		err = os.Remove(launchFilePath)
		if err != nil {
			return err
		}
	}

	err = os.Symlink(serviceFilePath, launchFilePath)
	if err != nil {
		return fmt.Errorf("os.Symlink: %v", err.Error())
	}

	err = startService(launchFilePath)
	if err != nil {
		return err
	}
	return nil
}

func filePathExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func startService(launchFilePath string) error {
	_, err := exec.Command("launchctl", "load", launchFilePath).Output()
	if err != nil {
		return err
	}
	return nil
}
