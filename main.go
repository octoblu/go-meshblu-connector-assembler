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
		cli.StringFlag{
			Name:   "github-slug, g",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_GITHUB_SLUG",
			Usage:  "Github Slug",
		},
		cli.StringFlag{
			Name:   "tag",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_TAG",
			Usage:  "Tag or Version",
		},
		cli.BoolFlag{
			Name:   "legacy, l",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_LEGACY",
			Usage:  "Run legacy meshblu connector",
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
			Name:   "token",
			EnvVar: "MESHBLU_CONNECTOR_ASSEMBLER_TOKEN",
			Usage:  "Meshblu device token",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	opts := getOpts(context)

	fmt.Println("creating output directory")
	err := os.MkdirAll(opts.GetOutputDirectory(), 0755)
	fatalIfError("create output directory", err)

	fmt.Println("creating log directory")
	err = os.MkdirAll(opts.GetLogDirectory(), 0755)
	fatalIfError("create log directory", err)

	fmt.Println("creating bin directory")
	err = os.MkdirAll(opts.GetBinDirectory(), 0755)
	fatalIfError("create bin directory", err)

	downloadClient := downloader.New(opts.GetConnectorDirectory())
	downloadFile, err := downloadClient.Download(opts.GetDownloadURI())
	fatalIfError("error downloading", err)

	extractorClient := extractor.New()
	err = extractorClient.Do(downloadFile, opts.GetConnectorDirectory())
	fatalIfError("error extracting:", err)

	configuratorClient := configurator.New(opts)
	err = configuratorClient.WriteConfigs()
	fatalIfError("error writing configs:", err)

	ignitionFile, err := downloadClient.Download(opts.GetIgnitionURI())
	fatalIfError("error downloading ignition", err)

	err = os.Rename(ignitionFile, opts.GetExecutablePath())
	fatalIfError("error moving ignition", err)

	err = os.Chmod(opts.GetExecutablePath(), os.FileMode(int(0777)))
	fatalIfError("error making exectuable", err)

	foreverizerClient := foreverizer.New(opts)
	err = foreverizerClient.Do()
	fatalIfError("error setuping device to run forever", err)

	fmt.Println("done installing")
	windowsMustWait()
}

func getOpts(context *cli.Context) configurator.Options {
	opts := configurator.NewOptionsFromContext(context)

	if opts.GetConnector() == "" ||
		opts.GetGithubSlug() == "" ||
		opts.GetTag() == "" ||
		opts.GetUUID() == "" ||
		opts.GetToken() == "" {

		cli.ShowAppHelp(context)

		if opts.GetConnector() == "" {
			color.Red("  Missing required flag --connector, c or MESHBLU_CONNECTOR_ASSEMBLER_CONNECTOR")
		}

		if opts.GetGithubSlug() == "" {
			color.Red("  Missing required flag --github-slug, g or MESHBLU_CONNECTOR_ASSEMBLER_GITHUB_SLUG")
		}

		if opts.GetTag() == "" {
			color.Red("  Missing required flag --tag or MESHBLU_CONNECTOR_ASSEMBLER_TAG")
		}

		if opts.GetUUID() == "" {
			color.Red("  Missing required flag --uuid, -u or MESHBLU_CONNECTOR_ASSEMBLER_UUID")
		}

		if opts.GetToken() == "" {
			color.Red("  Missing required flag --token or MESHBLU_CONNECTOR_ASSEMBLER_TOKEN")
		}
		os.Exit(1)
	}

	return opts
}

func fatalIfError(msg string, err error) {
	if err == nil {
		return
	}

	log.Println(msg, err.Error())
	windowsMustWait()
	log.Fatalln("Exiting...")
}

func windowsMustWait() {
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
