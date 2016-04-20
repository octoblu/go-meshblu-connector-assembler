package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/go-meshblu-connector-installer/configurator"
	"github.com/octoblu/go-meshblu-connector-installer/downloader"
	"github.com/octoblu/go-meshblu-connector-installer/extractor"
	"github.com/octoblu/go-meshblu-connector-installer/foreverizer"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("meshblu-connector-installer:main")

func main() {
	app := cli.NewApp()
	app.Name = "meshblu-connector-installer"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "connector, c",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_CONNECTOR",
			Usage:  "Meshblu connector name",
		},
		cli.StringFlag{
			Name:   "hostname",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_HOSTNAME",
			Usage:  "Meshblu device hostname",
		},
		cli.BoolFlag{
			Name:   "legacy, l",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_LEGACY",
			Usage:  "Run legacy meshblu connector",
		},
		cli.StringFlag{
			Name:   "output, o",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_OUTPUT",
			Usage:  "Output directory",
		},
		cli.IntFlag{
			Name:   "port, -p",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_PORT",
			Usage:  "Meshblu device port",
		},
		cli.StringFlag{
			Name:   "uuid, -u",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_UUID",
			Usage:  "Meshblu device uuid",
		},
		cli.StringFlag{
			Name:   "tag, t",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_TAG",
			Usage:  "Tag version. Defaults to 'latest'",
		},
		cli.StringFlag{
			Name:   "token",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_TOKEN",
			Usage:  "Meshblu device token",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	opts := getOpts(context)
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	fmt.Println("creating directory", opts.ConnectorDirectory)
	err := os.MkdirAll(opts.ConnectorDirectory, 755)
	fatalIfError("error creating output directory", err)

	baseURI := "https://meshblu-connector.octoblu.com"
	downloadClient := downloader.New(opts.ConnectorDirectory, baseURI)
	downloadFile, err := downloadClient.DownloadConnector(getConnector(opts), opts.Tag, platform)
	fatalIfError("error downloading", err)

	extractorClient := extractor.New()
	err = extractorClient.Do(downloadFile, opts.ConnectorDirectory)
	fatalIfError("error extracting:", err)

	configuratorClient := configurator.New(opts)
	err = configuratorClient.WriteMeshblu()
	fatalIfError("error writing meshblu config:", err)

	foreverizerClient := foreverizer.New(opts)
	err = foreverizerClient.Do()
	fatalIfError("error setuping device to run forever", err)

	fmt.Println("done installing")
}

func getConnector(opts *configurator.Options) string {
	if opts.Legacy {
		return "run-legacy"
	}
	return opts.Connector
}

func getOpts(context *cli.Context) *configurator.Options {
	opts := configurator.NewOptions(context)

	if opts.Connector == "" || opts.UUID == "" || opts.Token == "" {
		cli.ShowAppHelp(context)

		if opts.Connector == "" {
			color.Red("  Missing required flag --connector or MESHBLU_CONNECTOR_INSTALLER_CONNECTOR")
		}

		if opts.UUID == "" {
			color.Red("  Missing required flag --uuid or MESHBLU_CONNECTOR_INSTALLER_OUTPUT")
		}

		if opts.Token == "" {
			color.Red("  Missing required flag --token or MESHBLU_CONNECTOR_INSTALLER_TOKEN")
		}
		os.Exit(1)
	}

	if opts.OutputDirectory == "" {
		opts.OutputDirectory = configurator.GetDefaultServiceDirectory()
	}

	outputDirectory, err := filepath.Abs(opts.OutputDirectory)
	if err != nil {
		log.Fatalln("Invalid output directory:", err.Error())
	}
	opts.OutputDirectory = outputDirectory
	opts.ConnectorDirectory = configurator.GetConnectorDirectory(opts)

	if opts.Hostname == "" {
		opts.Hostname = "meshblu.octoblu.com"
	}

	if opts.Port == 0 {
		opts.Port = 443
	}

	if opts.Tag == "" {
		opts.Tag = "latest"
	}

	opts.LogDirectory = configurator.GetLogDirectory(opts)
	opts.BinDirectory = configurator.GetBinDirectory(opts)
	opts.ServiceName = configurator.GetServiceName(opts)

	return opts
}

func fatalIfError(msg string, err error) {
	if err == nil {
		return
	}

	log.Fatalln(msg, err.Error())
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
