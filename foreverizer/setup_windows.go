package foreverizer

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/octoblu/go-meshblu-connector-installer/configurator"
)

// Setup configures the os to the device
func Setup(opts *configurator.Options) error {
	err := setupStructure(opts)
	if err != nil {
		return err
	}

	err = setEnvInService(opts)
	if err != nil {
		return err
	}

	err = startService(opts)
	if err != nil {
		return err
	}

	return nil
}

func setupStructure(opts *configurator.Options) error {
	err := os.MkdirAll(opts.LogDirectory, 0777)
	if err != nil {
		return err
	}
	return err
}

func startService(opts *configurator.Options) error {
	startExe := path.Join(opts.ConnectorDirectory, "start")
	_, err := exec.Command(getNSSMExe(opts), "install", opts.ServiceName, startExe).Output()
	if err != nil {
		return err
	}
	return nil
}

func setEnvInService(opts *configurator.Options) error {
	_, err := exec.Command(
		getNSSMExe(opts),
		"set",
		opts.ServiceName,
		"AppEnvironmentExtra",
		getPath(opts),
	).Output()
	if err != nil {
		return err
	}
	return nil
}

func getNSSMExe(opts *configurator.Options) string {
	return path.Join(opts.BinDirectory, "nssm.exe")
}

func getPath(opts *configurator.Options) string {
	npmPath := path.Join(opts.BinDirectory, "node_modules", "npm", "bin")
	return fmt.Sprintf("PATH=%s:%s:%s", os.Getenv("PATH"), opts.BinDirectory, npmPath)
}
