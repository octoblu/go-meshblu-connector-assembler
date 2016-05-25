package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"

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
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Debug mode, will prompt for user to continue on Windows",
		},
		cli.StringFlag{
			Name:   "github-slug, g",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_GITHUB_SLUG",
			Usage:  "Github Slug",
		},
		cli.StringFlag{
			Name:   "tag, T",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_TAG",
			Usage:  "Tag or Version",
		},
		cli.StringFlag{
			Name:   "legacy, l",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_LEGACY",
			Usage:  "Legacy Version",
		},
		cli.StringFlag{
			Name:   "ignition, i",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_IGNITION_TAG",
			Usage:  "Ignition Tag",
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

	createDirectories(opts)

	downloadConnector(opts)

	writeConfiguration(opts)

	installIgnition(opts)

	foreverize(opts)

	debug("done installing")
	windowsMustWait(opts)
}

func writeConfiguration(opts configurator.Options) {
	debug("writing conifuration files")
	configuratorClient := configurator.New(opts)
	err := configuratorClient.WriteConfigs()
	fatalIfError(opts, "error writing configs:", err)
}

func downloadConnector(opts configurator.Options) {
	debug("downloading the connector")
	client := extractor.New()
	err := client.DoWithURI(opts.GetDownloadURI(), opts.GetConnectorDirectory())
	fatalIfError(opts, "error downloading:", err)
}

func createDirectories(opts configurator.Options) {
	debug("creating directories")
	err := os.MkdirAll(opts.GetOutputDirectory(), 0755)
	fatalIfError(opts, "create output directory", err)

	debug("creating log directory")
	err = os.MkdirAll(opts.GetLogDirectory(), 0755)
	fatalIfError(opts, "create log directory", err)

	debug("creating bin directory")
	err = os.MkdirAll(opts.GetBinDirectory(), 0755)
	fatalIfError(opts, "create bin directory", err)
}

func installIgnition(opts configurator.Options) {
	client := downloader.New(opts.GetConnectorDirectory())
	ignitionFile, err := client.Download(opts.GetIgnitionURI())
	fatalIfError(opts, "error downloading ignition", err)

	err = os.Rename(ignitionFile, opts.GetExecutablePath())
	fatalIfError(opts, "error moving ignition", err)

	err = os.Chmod(opts.GetExecutablePath(), os.FileMode(int(0777)))
	fatalIfError(opts, "error making exectuable", err)
}

func foreverize(opts configurator.Options) {
	foreverizerClient := foreverizer.New(opts)
	err := foreverizerClient.Do()
	fatalIfError(opts, "error setuping device to run forever", err)
}

func getOpts(context *cli.Context) configurator.Options {
	opts := configurator.NewOptionsFromContext(context)

	errStr := opts.Validate()

	if errStr != "" {
		cli.ShowAppHelp(context)
		color.Red(errStr)
		os.Exit(1)
	}

	return opts
}

func fatalIfError(opts configurator.Options, msg string, err error) {
	if err == nil {
		return
	}

	log.Println(msg, err.Error())
	windowsMustWait(opts)
	log.Fatalln("Exiting...")
}

func windowsMustWait(opts configurator.Options) {
	if opts.GetDebug() == false {
		return
	}
	if runtime.GOOS == "windows" {
		fmt.Println("Press any key to continue >>>")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
