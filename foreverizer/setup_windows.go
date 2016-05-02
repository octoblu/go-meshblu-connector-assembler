package foreverizer

import (
	"fmt"
	"os"

	"github.com/octoblu/go-meshblu-connector-assembler/configurator"
)

// Setup configures the os to the device
func Setup(opts configurator.Options) error {
	err := setupStructure(opts)
	if err != nil {
		return err
	}

	err = startService(opts)
	if err != nil {
		return err
	}

	return nil
}

func setupStructure(opts configurator.Options) error {
	fmt.Println("setting up log directory")
	return os.MkdirAll(opts.GetLogDirectory(), 0777)
}

func startService(opts configurator.Options) error {
	return nil
}
