package foreverizer

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"github.com/DHowett/go-plist"
	"github.com/octoblu/go-meshblu-connector-assembler/configurator"
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
	opts configurator.Options
}

// NewServiceConfig constructs a new Meshblu instance
func NewServiceConfig(opts configurator.Options) *ServiceConfig {
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

func (config *ServiceConfig) getServiceData() *ServiceData {
	opts := config.opts
	label := opts.GetServiceName()
	startCmd := opts.GetExecutablePath()
	pArgs := []string{startCmd, opts.GetLegacyFlag()}
	keepAlive := true
	outPath := path.Join(opts.GetLogDirectory(), "connector.log")
	errPath := path.Join(opts.GetLogDirectory(), "connector-error.log")
	env := map[string]string{
		"PATH": opts.GetPathEnv(),
		"MESHBLU_CONNECTOR_NAME":   opts.GetConnector(),
		"MESHBLU_CONNECTOR_LEGACY": fmt.Sprintf("%v", opts.GetLegacy()),
	}
	return &ServiceData{label, pArgs, keepAlive, outPath, errPath, opts.GetConnectorDirectory(), env}
}
