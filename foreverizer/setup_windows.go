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
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(opts.GetServiceName())
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", opts.GetServiceName())
	}
	mrgConfig := mgr.Config{
		StartType:   mgr.StartAutomatic,
		DisplayName: fmt.Sprintf("Meshblu Connector: %s", opts.GetConnector()),
		Description: fmt.Sprintf("Meshblu Connector %s, %s", opts.GetConnector(), opts.GetUUID()),
	}
	s, err = m.CreateService(opts.GetServiceName(), opts.GetExecutablePath(), mrgConfig, opts.GetLegacyFlag())
	if err != nil {
		return err
	}
	defer s.Close()
	err = eventlog.InstallAsEventCreate(opts.GetServiceName(), eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}
	return nil
}
