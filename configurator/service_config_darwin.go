package configurator

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"github.com/DHowett/go-plist"
)

// ServiceData config for the launchagent
type ServiceData struct {
	Label                string            `plist:"Label"`
	ProgramArguments     []string          `plist:"ProgramArguments"`
	KeepAlive            bool              `plist:"KeepAlive"`
	StandardOutPath      string            `plist:"StandardOutPath"`
	StandardErrorPath    string            `plist:"StandardErrorPath"`
	WorkingDirectory     string            `plist:"WorkingDirectory"`
	EnvironmentVariables map[string]string `plist:"EnvironmentVariables"`
}

// ServiceConfig interfaces with a remote meshblu server
type ServiceConfig struct {
	opts *Options
}

// NewServiceConfig constructs a new Meshblu instance
func NewServiceConfig(opts *Options) *ServiceConfig {
	return &ServiceConfig{opts}
}

// Export the config
func (config *ServiceConfig) Export() ([]byte, error) {
	serviceData := config.getServiceData()
	return encodeServiceData(serviceData)
}

func convertBooleanToSelfClose(buf *bytes.Buffer) *bytes.Buffer {
	fileString := buf.String()
	fileString = strings.Replace(fileString, "<true></true>", "<true/>", -1)
	fileString = strings.Replace(fileString, "<false></false>", "<false/>", -1)
	return bytes.NewBufferString(fileString)
}

func encodeServiceData(data *ServiceData) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := plist.NewEncoder(buf)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	buf = convertBooleanToSelfClose(buf)
	return buf.Bytes(), nil
}

func (config *ServiceConfig) getLegacyFlag() string {
	if config.opts.Legacy {
		return "--legacy"
	}
	return ""
}

func (config *ServiceConfig) getServiceData() *ServiceData {
	opts := config.opts
	label := opts.ServiceName
	startCmd := path.Join(opts.ConnectorDirectory, "start")
	pArgs := []string{startCmd, config.getLegacyFlag()}
	keepAlive := true
	outPath := path.Join(opts.LogDirectory, fmt.Sprintf("%s.log", opts.UUID))
	errPath := path.Join(opts.LogDirectory, fmt.Sprintf("%s-error.log", opts.UUID))
	env := map[string]string{
		"PATH":              fmt.Sprintf("/sbin:/usr/sbin:/bin:/usr/bin:%s", opts.BinDirectory),
		"MESHBLU_CONNECTOR": opts.Connector,
	}
	return &ServiceData{label, pArgs, keepAlive, outPath, errPath, opts.ConnectorDirectory, env}
}
