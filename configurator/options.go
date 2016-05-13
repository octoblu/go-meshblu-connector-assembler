package configurator

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/codegangsta/cli"
)

// OptionsConfig defines the service configurations
type OptionsConfig struct {
	IgnitionTag     string
	Connector       string
	GithubSlug      string
	Tag             string
	OutputDirectory string
	Legacy          bool
	ServiceName     string
	Hostname        string
	Port            int
	UUID, Token     string
}

// Options defines the service configurations
type Options interface {
	GetIgnitionURI() string
	GetConnectorDirectory() string
	GetBinDirectory() string
	GetLogDirectory() string
	GetServiceName() string
	GetExecutablePath() string
	GetConnector() string
	GetGithubSlug() string
	GetTag() string
	GetDisplayName() string
	GetUserName() (string, error)
	GetDescription() string
	GetDownloadURI() string
	GetOutputDirectory() string
	GetLegacy() bool
	GetLegacyFlag() string
	GetHostname() string
	GetPort() int
	GetUUID() string
	GetToken() string
}

// NewOptionsFromContext should create an options interface from the context
func NewOptionsFromContext(context *cli.Context) Options {
	outputDirectory := context.String("output")
	if outputDirectory == "" {
		outputDirectory = GetDefaultServiceDirectory()
	}
	outputDirectory, err := filepath.Abs(outputDirectory)
	if err != nil {
		log.Fatalln("Invalid output directory:", err.Error())
	}
	ignitionTag := context.String("ignition")
	if ignitionTag == "" {
		ignitionTag = "v3.0.4"
	}
	return &OptionsConfig{
		IgnitionTag:     ignitionTag,
		Connector:       context.String("connector"),
		GithubSlug:      context.String("github-slug"),
		Tag:             context.String("tag"),
		OutputDirectory: outputDirectory,
		Legacy:          context.Bool("legacy"),
		Hostname:        "meshblu.octoblu.com",
		Port:            443,
		UUID:            context.String("uuid"),
		Token:           context.String("token"),
	}
}

// NewOptions should create an options interface
func NewOptions(optConfig *OptionsConfig) Options {
	outputDirectory := optConfig.OutputDirectory
	if outputDirectory == "" {
		outputDirectory = GetDefaultServiceDirectory()
	}
	outputDirectory, err := filepath.Abs(outputDirectory)
	if err != nil {
		log.Fatalln("Invalid output directory:", err.Error())
	}
	ignitionTag := optConfig.IgnitionTag
	if ignitionTag == "" {
		ignitionTag = "v3.0.3"
	}
	optConfig.IgnitionTag = ignitionTag
	optConfig.OutputDirectory = outputDirectory
	return optConfig
}

// GetIgnitionURI gets the OS specific connector path
func (opts *OptionsConfig) GetIgnitionURI() string {
	baseURI := "https://github.com/octoblu/go-meshblu-connector-ignition/releases/download"
	return fmt.Sprintf("%s/%s/meshblu-connector-ignition-%s-%s", baseURI, opts.IgnitionTag, runtime.GOOS, runtime.GOARCH)
}

// GetConnectorDirectory gets the OS specific connector path
func (opts *OptionsConfig) GetConnectorDirectory() string {
	return filepath.Join(opts.OutputDirectory, opts.UUID)
}

// GetLogDirectory gets the OS specific log directory
func (opts *OptionsConfig) GetLogDirectory() string {
	return filepath.Join(opts.GetConnectorDirectory(), "log")
}

// GetConnector get connector name
func (opts *OptionsConfig) GetConnector() string {
	return opts.Connector
}

// GetGithubSlug get connector name
func (opts *OptionsConfig) GetGithubSlug() string {
	return opts.GithubSlug
}

// GetTag get connector name
func (opts *OptionsConfig) GetTag() string {
	return opts.Tag
}

// GetDisplayName get service display name
func (opts *OptionsConfig) GetDisplayName() string {
	return fmt.Sprintf("MeshbluConnector %s", opts.GetUUID())
}

// GetDescription get service description
func (opts *OptionsConfig) GetDescription() string {
	return fmt.Sprintf("MeshbluConnector (%s) %s", opts.GetConnector(), opts.GetUUID())
}

// GetDownloadURI get download uri
func (opts *OptionsConfig) GetDownloadURI() string {
	baseURI := fmt.Sprintf("https://github.com/%s/releases/download", opts.GithubSlug)
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	fileName := fmt.Sprintf("%s-%s-%s.%s", opts.Connector, runtime.GOOS, runtime.GOARCH, ext)
	return fmt.Sprintf("%s/%s/%s", baseURI, opts.Tag, fileName)
}

// GetOutputDirectory get output directory
func (opts *OptionsConfig) GetOutputDirectory() string {
	return opts.OutputDirectory
}

// GetLegacy get legacy bool
func (opts *OptionsConfig) GetLegacy() bool {
	return opts.Legacy
}

// GetLegacyFlag get legacy flag
func (opts *OptionsConfig) GetLegacyFlag() string {
	if opts.GetLegacy() {
		return "--legacy"
	}
	return ""
}

// GetHostname get meshblu hostname
func (opts *OptionsConfig) GetHostname() string {
	return opts.Hostname
}

// GetPort get meshblu port
func (opts *OptionsConfig) GetPort() int {
	return opts.Port
}

// GetUUID get meshblu uuid
func (opts *OptionsConfig) GetUUID() string {
	return opts.UUID
}

// GetToken get meshblu token
func (opts *OptionsConfig) GetToken() string {
	return opts.Token
}
