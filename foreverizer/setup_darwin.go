package foreverizer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/octoblu/go-meshblu-connector-installer/configurator"
)

// Setup configures the os to the device
func Setup(uuid, connector, outputDirectory string) error {
	err := setupStructure(outputDirectory)
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

func getLaunchFilePath(uuid string) string {
	return path.Join(os.Getenv("HOME"), "Library/LaunchAgents", getServiceFileName(uuid))
}
func getLogDirectory(outputDirectory string) string {
	return path.Join(outputDirectory, "log")
}

func getServiceFileName(uuid string) string {
	return fmt.Sprintf("com.octoblu.%s.plist", uuid)
}

func getServiceFilePath(uuid, outputDirectory string) string {
	return path.Join(outputDirectory, getServiceFileName(uuid))
}

func setupLaunchFile(uuid, outputDirectory string) error {
	launchFilePath := getLaunchFilePath(uuid)

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

	err = os.Symlink(getServiceFilePath(uuid, outputDirectory), launchFilePath)
	if err != nil {
		return fmt.Errorf("os.Symlink: %v", err.Error())
	}
	return nil
}

func setupStructure(outputDirectory string) error {
	err := os.MkdirAll(getLogDirectory(outputDirectory), 0777)
	if err != nil {
		return err
	}
	return err
}

func startService(uuid string) error {
	_, err := exec.Command("launchctl", "load", getLaunchFilePath(uuid)).Output()
	if err != nil {
		return err
	}
	return nil
}

func writeServiceFile(uuid, connector, outputDirectory string) error {
	serviceConfig := configurator.NewServiceConfig(uuid, connector, outputDirectory, getLogDirectory(outputDirectory))
	fileBytes, err := serviceConfig.Export()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(getServiceFilePath(uuid, outputDirectory), fileBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
