package configurator

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/codegangsta/cli"
)

// OptionsOptions defines the service configurations
type OptionsOptions struct {
	IgnitionTag     string
	Connector       string
	GithubSlug      string
	Tag             string
	OutputDirectory string
	ServiceName     string
	Hostname        string
	Port            int
	UUID, Token     string
	Debug           bool
}

// Options defines the service configurations
type Options interface {
	GetIgnitionURI() string
	GetIgnitionTag() string
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
	GetHostname() string
	GetPort() int
	GetUUID() string
	GetToken() string
	GetDebug() bool
}

type options struct {
	ignitionTag     string
	connector       string
	githubSlug      string
	tag             string
	outputDirectory string
	serviceName     string
	hostname        string
	port            int
	uuid, token     string
	debug           bool
}

// NewOptionsFromContext should create an options interface from the context
func NewOptionsFromContext(context *cli.Context) (Options, error) {
	return NewOptions(OptionsOptions{
		IgnitionTag:     context.String("ignition"),
		Connector:       context.String("connector"),
		GithubSlug:      context.String("github-slug"),
		Tag:             context.String("tag"),
		OutputDirectory: context.String("output"),
		UUID:            context.String("uuid"),
		Token:           context.String("token"),
		Debug:           context.Bool("debug"),
	})
}

// NewOptions should create an options interface
func NewOptions(opts OptionsOptions) (Options, error) {
	if opts.Connector == "" {
		return nil, fmt.Errorf("Missing required opt: opts.Connector")
	}
	if opts.GithubSlug == "" {
		return nil, fmt.Errorf("Missing required opt: opts.GithubSlug")
	}
	if opts.Tag == "" {
		return nil, fmt.Errorf("Missing required opt: opts.Tag")
	}
	if opts.UUID == "" {
		return nil, fmt.Errorf("Missing required opt: opts.UUID")
	}
	if opts.Token == "" {
		return nil, fmt.Errorf("Missing required opt: opts.Token")
	}
	if opts.IgnitionTag == "" {
		return nil, fmt.Errorf("Missing required opt: opts.IgnitionTag")
	}

	hostname := opts.Hostname
	if hostname == "" {
		hostname = "meshblu.octoblu.com"
	}

	port := opts.Port
	if port == 0 {
		port = 443
	}

	outputDirectory := opts.OutputDirectory
	if outputDirectory == "" {
		outputDirectory = GetDefaultServiceDirectory()
	}

	outputDirectory, err := filepath.Abs(outputDirectory)
	if err != nil {
		return nil, err
	}

	return &options{
		hostname:        hostname,
		port:            port,
		outputDirectory: outputDirectory,
		connector:       opts.Connector,
		ignitionTag:     opts.IgnitionTag,
		githubSlug:      opts.GithubSlug,
		tag:             opts.Tag,
		uuid:            opts.UUID,
		token:           opts.Token,
		debug:           opts.Debug,
	}, nil
}

// GetIgnitionURI gets the OS specific connector path
func (opts *options) GetIgnitionURI() string {
	baseURI := "https://github.com/octoblu/go-meshblu-connector-ignition/releases/download"
	return fmt.Sprintf("%s/%s/meshblu-connector-ignition-%s-%s", baseURI, opts.ignitionTag, runtime.GOOS, runtime.GOARCH)
}

// GetConnectorDirectory gets the OS specific connector path
func (opts *options) GetConnectorDirectory() string {
	return filepath.Join(opts.outputDirectory, opts.uuid)
}

// GetLogDirectory gets the OS specific log directory
func (opts *options) GetLogDirectory() string {
	return filepath.Join(opts.GetConnectorDirectory(), "log")
}

// GetConnector get connector name
func (opts *options) GetConnector() string {
	return opts.connector
}

// GetIgnitionTag gets the ignition tag
func (opts *options) GetIgnitionTag() string {
	return opts.ignitionTag
}

// GetGithubSlug get connector name
func (opts *options) GetGithubSlug() string {
	return opts.githubSlug
}

// GetTag get connector name
func (opts *options) GetTag() string {
	return opts.tag
}

// GetDisplayName get service display name
func (opts *options) GetDisplayName() string {
	return fmt.Sprintf("MeshbluConnector %s", opts.GetUUID())
}

// GetDescription get service description
func (opts *options) GetDescription() string {
	return fmt.Sprintf("MeshbluConnector (%s) %s", opts.GetConnector(), opts.GetUUID())
}

// GetDownloadURI get download uri
func (opts *options) GetDownloadURI() string {
	tag := opts.GetTag()
	connector := opts.GetConnector()
	baseURI := fmt.Sprintf("https://github.com/%s/releases/download", opts.githubSlug)
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	fileName := fmt.Sprintf("%s-%s-%s.%s", connector, runtime.GOOS, runtime.GOARCH, ext)
	return fmt.Sprintf("%s/%s/%s", baseURI, tag, fileName)
}

// GetOutputDirectory get output directory
func (opts *options) GetOutputDirectory() string {
	return opts.outputDirectory
}

// GetHostname get meshblu hostname
func (opts *options) GetHostname() string {
	return opts.hostname
}

// GetPort get meshblu port
func (opts *options) GetPort() int {
	return opts.port
}

// GetUUID get meshblu uuid
func (opts *options) GetUUID() string {
	return opts.uuid
}

// GetToken get meshblu token
func (opts *options) GetToken() string {
	return opts.token
}

// GetDebug gets the debug flag
func (opts *options) GetDebug() bool {
	return opts.debug
}
