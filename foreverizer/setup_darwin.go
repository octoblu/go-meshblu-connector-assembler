package foreverizer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/octoblu/go-meshblu-connector-assembler/configurator"
)

// Setup configures the os to the device
func Setup(opts configurator.Options) error {
	err := setupStructure(opts)
	if err != nil {
		return err
	}

	err = writeServiceFile(opts)
	if err != nil {
		return err
	}

	err = setupLaunchFile(opts)
	if err != nil {
		return err
	}

	err = startService(opts)
	if err != nil {
		return err
	}
	return nil
}

func getLaunchFilePath(opts configurator.Options) string {
	return path.Join(os.Getenv("HOME"), "Library", "LaunchAgents", opts.GetServiceName())
}

func getServiceFilePath(opts configurator.Options) string {
	launchFilePath := getServiceFilePath(opts)
	return fmt.Sprintf("%s.plist", launchFilePath)
}

func setupLaunchFile(opts configurator.Options) error {
	fmt.Println("setting up launch file")
	launchFilePath := getLaunchFilePath(opts)

	fileExists, err := FilePathExists(launchFilePath)
	if err != nil {
		return err
	}

	if fileExists {
		err = os.Remove(launchFilePath)
		if err != nil {
			fmt.Println("error removing old launch file", err.Error())
			return err
		}
	}

	err = os.Symlink(opts.GetExecutablePath(), launchFilePath)
	if err != nil {
		fmt.Println("error symlinking service file", err.Error())
		return fmt.Errorf("os.Symlink: %v", err.Error())
	}
	return nil
}

func setupStructure(opts configurator.Options) error {
	fmt.Println("setting up log directory")
	err := os.MkdirAll(opts.GetLogDirectory(), 0777)
	if err != nil {
		fmt.Println("error creating log directory", err.Error())
		return err
	}
	return err
}

func startService(opts configurator.Options) error {
	fmt.Println("starting service")
	launchFilePath := getLaunchFilePath(opts)
	_, err := exec.Command("launchctl", "load", launchFilePath).Output()
	if err != nil {
		fmt.Println("Error starting service", err.Error())
		return err
	}
	return nil
}

func writeServiceFile(opts configurator.Options) error {
	fmt.Println("writing service file")
	serviceConfig := NewServiceConfig(opts)
	fileBytes, err := serviceConfig.Export()
	if err != nil {
		fmt.Println("error exporting service config", err.Error())
		return err
	}
	err = ioutil.WriteFile(getServiceFilePath(opts), fileBytes, 0644)
	if err != nil {
		fmt.Println("error writing service file", err.Error())
		return err
	}
	return nil
}
