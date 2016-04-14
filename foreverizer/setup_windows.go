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

	err = setEnvInService(opts, getPath(opts))
	if err != nil {
		return err
	}

	err = setEnvInService(opts, fmt.Sprintf("MESHBLU_CONNECTOR=%s", opts.Connector))
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
	return os.MkdirAll(opts.LogDirectory, 0777)
}

func getLegacyFlag(opts *configurator.Options) string {
	if opts.Legacy {
		return "--legacy"
	}
	return ""
}

func getStartExe(opts *configurator.Options) string {
	return path.Join(opts.ConnectorDirectory, "start.exe")
}

func startService(opts *configurator.Options) error {
	return exec.Command(
		getNSSMExe(opts),
		"install",
		opts.ServiceName,
		getStartExe(opts),
		getLegacyFlag(opts),
	).Run()
}

func setEnvInService(opts *configurator.Options, env string) error {
	return exec.Command(
		getNSSMExe(opts),
		"set",
		opts.ServiceName,
		"AppEnvironmentExtra",
		env,
	).Run()
}

func getNSSMExe(opts *configurator.Options) string {
	return path.Join(opts.BinDirectory, "nssm.exe")
}

func getPath(opts *configurator.Options) string {
	npmPath := path.Join(opts.BinDirectory, "node_modules", "npm", "bin")
	return fmt.Sprintf("PATH=%s:%s:%s", os.Getenv("PATH"), opts.BinDirectory, npmPath)
}
