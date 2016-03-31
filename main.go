package main

import (
	"fmt"
	"log"
	"os"
	"path"
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

// CommandOpts defines the command line arguments
type CommandOpts struct {
	connector, hostname string
	legacy              bool
	outputDirectory     string
	port                int
	uuid, tag, token    string
}

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
	platform := runtime.GOOS
	err := os.MkdirAll(opts.outputDirectory, 0755)
	fatalIfError("Error creating output directory", err)

	baseURI := "https://meshblu-connector.octoblu.com"
	downloadClient := downloader.New(opts.outputDirectory, baseURI)
	downloadFile, err := downloadClient.DownloadConnector(getConnector(opts), opts.tag, platform)
	fatalIfError("Error downloading", err)

	extractorClient := extractor.New()
	err = extractorClient.Do(downloadFile, opts.outputDirectory)
	fatalIfError("Error extracting:", err)

	configuratorClient := configurator.New(opts.outputDirectory)
	err = configuratorClient.WriteMeshblu(opts.uuid, opts.token, opts.hostname, opts.port)
	fatalIfError("Error writing meshblu config:", err)

	foreverizerClient := foreverizer.New()
	err = foreverizerClient.Do(opts.uuid, opts.connector, opts.outputDirectory)
	fatalIfError("Error setuping device to run forever", err)
}

func getConnector(opts *CommandOpts) string {
	if opts.legacy {
		return "run-legacy"
	}
	return opts.connector
}

func getOpts(context *cli.Context) *CommandOpts {
	commandOpts := &CommandOpts{
		context.String("connector"),
		context.String("hostname"),
		context.Bool("legacy"),
		context.String("output"),
		context.Int("port"),
		context.String("uuid"),
		context.String("tag"),
		context.String("token"),
	}

	if commandOpts.connector == "" || commandOpts.uuid == "" || commandOpts.token == "" {
		cli.ShowAppHelp(context)

		if commandOpts.connector == "" {
			color.Red("  Missing required flag --connector or MESHBLU_CONNECTOR_INSTALLER_CONNECTOR")
		}

		if commandOpts.uuid == "" {
			color.Red("  Missing required flag --uuid or MESHBLU_CONNECTOR_INSTALLER_OUTPUT")
		}

		if commandOpts.token == "" {
			color.Red("  Missing required flag --token or MESHBLU_CONNECTOR_INSTALLER_OUTPUT")
		}
		os.Exit(1)
	}

	if commandOpts.outputDirectory == "" {
		commandOpts.outputDirectory = path.Join(os.Getenv("HOME"), "Library", "Application Support", "Octoblu", commandOpts.uuid)
	}

	outputDirectory, err := filepath.Abs(commandOpts.outputDirectory)
	if err != nil {
		log.Fatalln("Invalid output directory:", err.Error())
	}
	commandOpts.outputDirectory = outputDirectory

	if commandOpts.hostname == "" {
		commandOpts.hostname = "meshblu.octoblu.com"
	}

	if commandOpts.port == 0 {
		commandOpts.port = 443
	}

	if commandOpts.tag == "" {
		commandOpts.tag = "latest"
	}

	return commandOpts
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
