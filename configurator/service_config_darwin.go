package configurator

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"github.com/DHowett/go-plist"
)

// ServiceConfig interfaces with a remote meshblu server
type ServiceConfig struct {
	uuid, connector, outputDirectory string
}

// NewServiceConfig constructs a new Meshblu instance
func NewServiceConfig(uuid, connector, outputDirectory string) *ServiceConfig {
	return &ServiceConfig{uuid, connector, outputDirectory}
}

// Export the config
func (config *ServiceConfig) Export() ([]byte, error) {
	type ServiceData struct {
		Label                string            `plist:"Label"`
		ProgramArguments     []string          `plist:"ProgramArguments"`
		KeepAlive            bool              `plist:"KeepAlive"`
		StandardOutPath      string            `plist:"StandardOutPath"`
		StandardErrorPath    string            `plist:"StandardErrorPath"`
		WorkingDirectory     string            `plist:"WorkingDirectory"`
		EnvironmentVariables map[string]string `plist:"EnvironmentVariables"`
	}

	buf := &bytes.Buffer{}
	encoder := plist.NewEncoder(buf)
	label := GetServiceFileName(config.uuid)

	startCmd := path.Join(GetConnectorDirectory(config.outputDirectory, config.uuid), "start")
	pArgs := []string{startCmd}
	keepAlive := true
	logDirectory := GetLogDirectory(config.outputDirectory, config.uuid)
	outPath := path.Join(logDirectory, fmt.Sprintf("%s.log", config.uuid))
	errPath := path.Join(logDirectory, fmt.Sprintf("%s-error.log", config.uuid))
	binPath := GetBinDirectory(config.outputDirectory)
	env := map[string]string{
		"PATH":              fmt.Sprintf("/sbin:/usr/sbin:/bin:/usr/bin:%s", binPath),
		"MESHBLU_CONNECTOR": config.connector,
	}
	data := &ServiceData{label, pArgs, keepAlive, outPath, errPath, config.outputDirectory, env}
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	buf = convertBooleanToSelfClose(buf)
	return buf.Bytes(), nil
}

func convertBooleanToSelfClose(buf *bytes.Buffer) *bytes.Buffer {
	fileString := buf.String()
	fileString = strings.Replace(fileString, "<true></true>", "<true/>", -1)
	fileString = strings.Replace(fileString, "<false></false>", "<false/>", -1)
	return bytes.NewBufferString(fileString)
}
