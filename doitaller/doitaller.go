package doitaller

import (
	"os"

	"github.com/octoblu/go-meshblu-connector-assembler/configurator"
	"github.com/octoblu/go-meshblu-connector-assembler/meshbluconfig"
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

	err = createDirectories(opts.GetOutputDirectory(), opts.GetLogDirectory(), opts.GetBinDirectory())
	if err != nil {
		return err
	}

	err = writeMeshbluConfig(opts.GetOutputDirectory(), opts.GetUUID(), opts.GetToken(), opts.GetHostname(), opts.GetPort())
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

func createDirectories(outputDir, logDir, getBinDir string) error {
	var err error

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
	err = os.MkdirAll(getBinDir, 0755)
	if err != nil {
		return err
	}

	return nil
}

func writeMeshbluConfig(dirPath, uuid, token, hostname string, port int) error {
	return meshbluconfig.Write(meshbluconfig.Options{
		DirPath:  dirPath,
		UUID:     uuid,
		Token:    token,
		Hostname: hostname,
		Port:     port,
	})
}
