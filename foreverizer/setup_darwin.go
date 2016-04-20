package foreverizer

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/octoblu/go-meshblu-connector-installer/configurator"
)

// Setup configures the os to the device
func Setup(opts *configurator.Options) error {
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

func setupLaunchFile(opts *configurator.Options) error {
	fmt.Println("setting up launch file")
	launchFilePath := configurator.GetLaunchFilePath(opts)

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

	err = os.Symlink(configurator.GetServiceFilePath(opts), launchFilePath)
	if err != nil {
		fmt.Println("error symlinking service file", err.Error())
		return fmt.Errorf("os.Symlink: %v", err.Error())
	}
	return nil
}

func setupStructure(opts *configurator.Options) error {
	fmt.Println("setting up log directory")
	err := os.MkdirAll(opts.LogDirectory, 0777)
	if err != nil {
		fmt.Println("error creating log directory", err.Error())
		return err
	}
	return err
}

func startService(opts *configurator.Options) error {
	fmt.Println("starting service")
	_, err := exec.Command("launchctl", "load", configurator.GetLaunchFilePath(opts)).Output()
	if err != nil {
		fmt.Println("Error starting service", err.Error())
		return err
	}
	return nil
}

func writeServiceFile(opts *configurator.Options) error {
	fmt.Println("writing service file")
	serviceConfig := configurator.NewServiceConfig(opts)
	fileBytes, err := serviceConfig.Export()
	if err != nil {
		fmt.Println("error exporting service config", err.Error())
		return err
	}
	err = ioutil.WriteFile(configurator.GetServiceFilePath(opts), fileBytes, 0644)
	if err != nil {
		fmt.Println("error writing service file", err.Error())
		return err
	}
	return nil
}
