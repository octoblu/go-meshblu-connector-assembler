package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/meshblu-connector-installer-go/downloader"
	"github.com/octoblu/meshblu-connector-installer-go/extractor"
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
			Name:   "output, o",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_OUTPUT",
			Usage:  "Output directory",
		},
		cli.StringFlag{
			Name:   "platform, p",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_PLATFORM",
			Usage:  "Platform name, 'osx', 'linux', 'win32', 'win64'. Defualts to 'osx'",
		},
		cli.StringFlag{
			Name:   "tag, t",
			EnvVar: "MESHBLU_CONNECTOR_INSTALLER_TAG",
			Usage:  "Tag version. Defaults to 'latest'",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	connector, outputDirectory, platform, tag := getOpts(context)
	baseURI := "https://meshblu-connector.octoblu.com"
	downloadClient := downloader.New(outputDirectory, baseURI)
	downloadFile, err := downloadClient.DownloadConnector(connector, tag, platform)
	if err != nil {
		log.Fatalln("Error downloading:", err.Error())
	}
	extractorClient := extractor.New()
	tarFile := strings.Replace(downloadFile, "tar.gz", "tar", 1)
	ungzipErr := extractorClient.Ungzip(downloadFile, tarFile)
	if ungzipErr != nil {
		log.Fatalln("Error ungziping:", ungzipErr.Error())
	}
	untarErr := extractorClient.Untar(tarFile, outputDirectory)
	if untarErr != nil {
		log.Fatalln("Error untaring:", untarErr.Error())
	}
}

func getOpts(context *cli.Context) (string, string, string, string) {
	connector := context.String("connector")
	outputDirectory := context.String("output")
	platform := context.String("platform")
	tag := context.String("tag")

	if connector == "" {
		cli.ShowAppHelp(context)

		if connector == "" {
			color.Red("  Missing required flag --connector or MESHBLU_CONNECTOR_INSTALLER_CONNECTOR")
		}
		os.Exit(1)
	}
	if outputDirectory == "" {
		cli.ShowAppHelp(context)

		if outputDirectory == "" {
			color.Red("  Missing required flag --output or MESHBLU_CONNECTOR_INSTALLER_OUTPUT")
		}
		os.Exit(1)
	}

	if platform == "" {
		platform = "osx"
	}

	if tag == "" {
		tag = "latest"
	}

	return connector, outputDirectory, platform, tag
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
