package foreverizer

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/octoblu/go-meshblu-connector-installer/configurator"
)

// Setup configures the os to the device
func Setup(opts *Options) error {
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

func setupStructure(opts *Options) error {
	err := os.MkdirAll(opts.LogDirectory, 0777)
	if err != nil {
		return err
	}
	return err
}

func startService(opts *Options) error {
	startExe := path.Join(opts.ConnectorDirectory, "start.exe")
	_, err := exec.Command(getNSSM(opts), "install", configurator.GetServiceFileName(opts), startExe).Output()
	if err != nil {
		return err
	}
	return nil
}

func setEnvInService(opts *Options) error {
	_, err := exec.Command(
		getNSSM(opts),
		"set",
		configurator.ServiceName,
		"AppEnvironmentExtra",
		fmt.Sprintf("PATH=%s", getPath(opts),
	).Output()
	if err != nil {
		return err
	}
	return nil
}

func getNSSM(opts *Options) string {
	return path.Join(opts.BinDirectory, "nssm.exe")
}

func getPath(opts *Options) string {
	npmPath := path.Join(opts.BinDirectory, "node_modules", "npm", "bin")
	return fmt.Sprintf("%s:%s:%s", os.Getenv("PATH"), opts.BinDirectory, npmPath)
}
