package doitaller

import (
	"os"

	"github.com/octoblu/go-meshblu-connector-assembler/configurator"
	"github.com/octoblu/go-meshblu-connector-assembler/meshbluconfig"
	"github.com/octoblu/go-meshblu-connector-assembler/serviceconfig"
)
import De "github.com/tj/go-debug"

var debug = De.Debug("meshblu-connector-assembler:doitaller")

// DoItAll does all the things a connector assembler does
// including:
//   creating directories
//   writing the meshblu config
//   writing the service config
//   installing the ignition
//   foreverize
func DoItAll(opts configurator.Options) error {
	var err error

	err = createDirectories(opts)
	if err != nil {
		return err
	}

	err = writeMeshbluConfig(opts)
	if err != nil {
		return err
	}

	err = writeServiceConfig(opts)
	if err != nil {
		return err
	}
	//
	// err = installIgnition(opts)
	// if err != nil {
	// 	return err
	// }
	//
	// err = foreverize(opts)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func createDirectories(opts configurator.Options) error {
	var err error

	outputDir := opts.GetOutputDirectory()
	logDir := opts.GetLogDirectory()
	binDir := opts.GetBinDirectory()

	debug("creating directories")
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	debug("creating log directory")
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		return err
	}

	debug("creating bin directory")
	err = os.MkdirAll(binDir, 0755)
	if err != nil {
		return err
	}

	return nil
}

func writeMeshbluConfig(opts configurator.Options) error {
	return meshbluconfig.Write(meshbluconfig.Options{
		DirPath:  opts.GetConnectorDirectory(),
		UUID:     opts.GetUUID(),
		Token:    opts.GetToken(),
		Hostname: opts.GetHostname(),
		Port:     opts.GetPort(),
	})
}

func writeServiceConfig(opts configurator.Options) error {
	return serviceconfig.Write(serviceconfig.Options{
		ServiceName:   opts.GetServiceName(),
		DisplayName:   opts.GetDisplayName(),
		Description:   opts.GetDescription(),
		ConnectorName: opts.GetConnector(),
		GithubSlug:    opts.GetGithubSlug(),
		Tag:           opts.GetTag(),
		BinPath:       opts.GetBinDirectory(),
		Dir:           opts.GetConnectorDirectory(),
		LogDir:        opts.GetLogDirectory(),
	})
}
