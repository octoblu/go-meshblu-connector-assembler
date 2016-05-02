package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/go-meshblu-connector-assembler/configurator"
	"github.com/octoblu/go-meshblu-connector-assembler/downloader"
	"github.com/octoblu/go-meshblu-connector-assembler/extractor"
	"github.com/octoblu/go-meshblu-connector-assembler/foreverizer"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("meshblu-connector-assembler:main")

func main() {
	app := cli.NewApp()
	app.Name = "meshblu-connector-assembler"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "connector, c",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_CONNECTOR",
			Usage:  "Connector name",
		},
		cli.StringFlag{
			Name:   "download-uri, d",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_DOWNLOAD_URI",
			Usage:  "Download URI",
		},
		cli.BoolFlag{
			Name:   "legacy, l",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_LEGACY",
			Usage:  "Run legacy meshblu connector",
		},
		cli.StringFlag{
			Name:   "uuid, u",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_UUID",
			Usage:  "Meshblu device uuid",
		},
		cli.StringFlag{
			Name:   "token, t",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_TOKEN",
			Usage:  "Meshblu device token",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	opts := getOpts(context)

	fmt.Println("creating directory", opts.GetConnectorDirectory())
	err := os.MkdirAll(opts.GetConnectorDirectory(), 0755)
	fatalIfError("error creating output directory", err)

	downloadClient := downloader.New(opts.GetConnectorDirectory())
	downloadFile, err := downloadClient.Download(opts.GetDownloadURI())
	fatalIfError("error downloading", err)

	extractorClient := extractor.New()
	err = extractorClient.Do(downloadFile, opts.GetConnectorDirectory())
	fatalIfError("error extracting:", err)

	configuratorClient := configurator.New(opts)
	err = configuratorClient.WriteMeshblu()
	fatalIfError("error writing meshblu config:", err)

	foreverizerClient := foreverizer.New(opts)
	err = foreverizerClient.Do()
	fatalIfError("error setuping device to run forever", err)

	fmt.Println("done installing")
}

func getOpts(context *cli.Context) configurator.Options {
	opts := configurator.NewOptions(context)

	if opts.GetConnector() == "" || opts.GetDownloadURI() == "" || opts.GetUUID() == "" || opts.GetToken() == "" {
		cli.ShowAppHelp(context)

		if opts.GetConnector() == "" {
			color.Red("  Missing required flag --connector, c or MESHBLU_CONNECTOR_ASSEMBLER_CONNECTOR")
		}

		if opts.GetDownloadURI() == "" {
			color.Red("  Missing required flag --download-uri, d or MESHBLU_CONNECTOR_ASSEMBLER_DOWNLOAD_URI")
		}

		if opts.GetUUID() == "" {
			color.Red("  Missing required flag --uuid, -u or MESHBLU_CONNECTOR_ASSEMBLER_UUID")
		}

		if opts.GetToken() == "" {
			color.Red("  Missing required flag --token, -t or MESHBLU_CONNECTOR_ASSEMBLER_TOKEN")
		}
		os.Exit(1)
	}

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
