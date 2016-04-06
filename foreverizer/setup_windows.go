package foreverizer

import (
	"fmt"
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

	err = setEnvInService(uuid, outputDirectory)
	if err != nil {
		return err
	}
	
	err = startService(uuid, outputDirectory)
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
	startExe := path.Join(configurator.GetConnectorDirectory(outputDirectory, uuid), "start.exe")
	_, err := exec.Command(getNSSM(outputDirectory), "install", getServiceName(uuid), startExe).Output()
	if err != nil {
		return err
	}
	return nil
}

func setEnvInService(uuid, outputDirectory string) error {
	_, err := exec.Command(
		getNSSM(outputDirectory),
		"set",
		getServiceName(uuid),
		"AppEnvironmentExtra",
		fmt.Sprintf("PATH=%s", getPath(outputDirectory),
	).Output()
	if err != nil {
		return err
	}
	return nil
}

func getNSSM(outputDirectory string) string {
	return path.Join(configurator.GetBinDirectory(outputDirectory), "nssm.exe")
}

func getPath(outputDirectory string) string {
	binPath := configurator.GetBinDirectory(outputDirectory)
	npmPath := path.Join(binPath, "node_modules", "npm", "bin")
	return fmt.Sprintf("%s:%s:%s", os.Getenv("PATH"), binPath, npmPath)
}

func getServiceName(uuid string) string {
	return fmt.Sprintf("MeshbluConnector-%s", uuid)
}
