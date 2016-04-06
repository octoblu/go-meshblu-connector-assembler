package foreverizer

import (
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

	err = startService(uuid)
	if err != nil {
		return err
	}

	return nil
}

func setupStructure(outputDirectory string) error {
	err := os.MkdirAll(configurator.GetLogDirectory(outputDirectory), 0777)
	if err != nil {
		return err
	}
	return err
}

func startService(uuid, outputDirectory string) error {
	nssmExe := path.Join(configurator.GetBinDirectory(outputDirectory), "nssm.exe")
	startExe := path.Join(configurator.GetConnectorDirectory(outputDirectory, uuid), "start.exe")
	_, err := exec.Command(nssmExe, "install", startExe).Output()
	if err != nil {
		return err
	}
	return nil
}
