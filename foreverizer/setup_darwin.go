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
	logDirectory := "/var/log/octoblu"
	err := os.MkdirAll(logDirectory, 0777)
	if err != nil {
		return err
	}
	serviceConfig := configurator.NewServiceConfig(uuid, outputDirectory, logDirectory)
	fileBytes, err := serviceConfig.Export()
	if err != nil {
		return nil
	}
	serviceFileName := fmt.Sprintf("com.octoblu.%s.plist", uuid)
	serviceFilePath := path.Join(outputDirectory, serviceFileName)
	err = ioutil.WriteFile(serviceFilePath, fileBytes, 06444)
	if err != nil {
		return err
	}

	launchFilePath := path.Join(os.Getenv("HOME"), "Library/LaunchAgents", serviceFileName)
	err = os.Symlink(serviceFilePath, launchFilePath)
	if err != nil {
		return err
	}

	err = startService(launchFilePath)
	if err != nil {
		return err
	}

	return nil
}

func startService(launchFilePath string) error {
	_, err := exec.Command("launchctl", "load", launchFilePath).Output()
	if err != nil {
		return err
	}
	return nil
}
