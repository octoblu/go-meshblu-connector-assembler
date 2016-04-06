package foreverizer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/octoblu/go-meshblu-connector-installer/configurator"
)

// Setup configures the os to the device
func Setup(uuid, connector, outputDirectory string) error {
	err := setupStructure(uuid, outputDirectory)
	if err != nil {
		return err
	}

	err = writeServiceFile(uuid, connector, outputDirectory)
	if err != nil {
		return err
	}

	err = setupLaunchFile(uuid, outputDirectory)
	if err != nil {
		return err
	}

	err = startService(uuid)
	if err != nil {
		return err
	}
	return nil
}

func setupLaunchFile(uuid, outputDirectory string) error {
	launchFilePath := configurator.GetLaunchFilePath(uuid)

	fileExists, err := FilePathExists(launchFilePath)
	if err != nil {
		return err
	}

	if fileExists {
		err = os.Remove(launchFilePath)
		if err != nil {
			return err
		}
	}

	err = os.Symlink(configurator.GetServiceFilePath(uuid, outputDirectory), launchFilePath)
	if err != nil {
		return fmt.Errorf("os.Symlink: %v", err.Error())
	}
	return nil
}

func setupStructure(uuid, outputDirectory string) error {
	err := os.MkdirAll(configurator.GetLogDirectory(outputDirectory, uuid), 0777)
	if err != nil {
		return err
	}
	return err
}

func startService(uuid string) error {
	_, err := exec.Command("launchctl", "load", configurator.GetLaunchFilePath(uuid)).Output()
	if err != nil {
		return err
	}
	return nil
}

func writeServiceFile(uuid, connector, outputDirectory string) error {
	serviceConfig := configurator.NewServiceConfig(uuid, connector, outputDirectory)
	fileBytes, err := serviceConfig.Export()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configurator.GetServiceFilePath(uuid, outputDirectory), fileBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
