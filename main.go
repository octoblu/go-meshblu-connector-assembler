package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
			Usage:  "Connector name",
		},
		cli.StringFlag{
			Name:   "download-uri, d",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_DOWNLOAD_URI",
			Usage:  "Download URI",
		},
		cli.BoolFlag{
			Name:   "legacy, l",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_LEGACY",
			Usage:  "Run legacy meshblu connector",
		},
		cli.StringFlag{
			Name:   "uuid, u",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_UUID",
			Usage:  "Meshblu device uuid",
		},
		cli.StringFlag{
			Name:   "token, t",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_TOKEN",
			Usage:  "Meshblu device token",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	opts := getOpts(context)

	fmt.Println("creating directory", opts.ConnectorDirectory)
	err := os.MkdirAll(opts.ConnectorDirectory, 755)
	fatalIfError("error creating output directory", err)

	downloadClient := downloader.New(opts.ConnectorDirectory)
	downloadFile, err := downloadClient.Download(opts.DownloadURI)
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

func getOpts(context *cli.Context) *configurator.Options {
	opts := configurator.NewOptions(context)

	if opts.Connector == "" || opts.DownloadURI == "" || opts.UUID == "" || opts.Token == "" {
		cli.ShowAppHelp(context)

		if opts.Connector == "" {
			color.Red("  Missing required flag --connector, c or MESHBLU_CONNECTOR_INSTALLER_CONNECTOR")
		}

		if opts.DownloadURI == "" {
			color.Red("  Missing required flag --download-uri, d or MESHBLU_CONNECTOR_INSTALLER_DOWNLOAD_URI")
		}

		if opts.UUID == "" {
			color.Red("  Missing required flag --uuid, -u or MESHBLU_CONNECTOR_INSTALLER_UUID")
		}

		if opts.Token == "" {
			color.Red("  Missing required flag --token, -t or MESHBLU_CONNECTOR_INSTALLER_TOKEN")
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
