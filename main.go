package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
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
	if opts.GetDebug() {
		tellWindowsToWait()
	}
}

func writeConfiguration(opts configurator.Options) {
	debug("writing configuration files")
	configuratorClient := configurator.New(opts)
	err := configuratorClient.WriteConfigs()
	fatalIfError(opts.GetDebug(), "error writing configs:", err)
}

func shouldDownloadConnector(opts configurator.Options) bool {
	tag := strings.Replace(opts.GetIgnitionTag(), "v", "", 1)
	parts := strings.Split(tag, ".")
	major, _ := strconv.ParseInt(parts[0], 10, 0)
	minor, _ := strconv.ParseInt(parts[1], 10, 0)
	if major < 6 {
		return true
	}
	if major == 6 && minor < 1 {
		return true
	}
	return false
}

func downloadConnector(opts configurator.Options) {
	downloadIt := shouldDownloadConnector(opts)
	if !downloadIt {
		debug("skipping downloading the connector because the ignition script will do it")
		return
	}
	debug("downloading the connector")
	client := extractor.New()
	err := client.DoWithURI(opts.GetDownloadURI(), opts.GetConnectorDirectory())
	fatalIfError(opts.GetDebug(), "error downloading:", err)
}

func createDirectories(opts configurator.Options) {
	debug("creating directories")
	err := os.MkdirAll(opts.GetOutputDirectory(), 0755)
	fatalIfError(opts.GetDebug(), "create output directory", err)

	debug("creating log directory")
	err = os.MkdirAll(opts.GetLogDirectory(), 0755)
	fatalIfError(opts.GetDebug(), "create log directory", err)

	debug("creating bin directory")
	err = os.MkdirAll(opts.GetBinDirectory(), 0755)
	fatalIfError(opts.GetDebug(), "create bin directory", err)
}

func installIgnition(opts configurator.Options) {
	client := downloader.New(opts.GetConnectorDirectory())
	ignitionFile, err := client.Download(opts.GetIgnitionURI())
	fatalIfError(opts.GetDebug(), "error downloading ignition", err)

	err = os.Rename(ignitionFile, opts.GetExecutablePath())
	fatalIfError(opts.GetDebug(), "error moving ignition", err)

	err = os.Chmod(opts.GetExecutablePath(), os.FileMode(int(0777)))
	fatalIfError(opts.GetDebug(), "error making exectuable", err)
}

func foreverize(opts configurator.Options) {
	foreverizerClient := foreverizer.New(opts)
	err := foreverizerClient.Do()
	fatalIfError(opts.GetDebug(), "error setuping device to run forever", err)
}

func getOpts(context *cli.Context) configurator.Options {
	opts, err := configurator.NewOptionsFromContext(context)
	fatalIfError(true, "Failed to create configurationion", err)
	return opts
}

func fatalIfError(windowsShouldWait bool, msg string, err error) {
	if err == nil {
		return
	}

	log.Println(msg, err.Error())
	if windowsShouldWait {
		tellWindowsToWait()
	}
	log.Fatalln("Exiting...")
}

func tellWindowsToWait() {
	if runtime.GOOS != "windows" {
		return
	}

	fmt.Println("Press any key to continue >>>")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
