package configurator

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/DHowett/go-plist"
)

// ServiceConfig interfaces with a remote meshblu server
type ServiceConfig struct {
	uuid, workingDirectory, logDirectory string
}

// NewServiceConfig constructs a new Meshblu instance
func NewServiceConfig(uuid, workingDirectory, logDirectory string) *ServiceConfig {
	return &ServiceConfig{uuid, workingDirectory, logDirectory}
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
	label := fmt.Sprintf("com.octoblu.%s", config.uuid)

	startCmd := path.Join(config.workingDirectory, "start")
	pArgs := []string{startCmd}
	keepAlive := true
	outPath := path.Join(config.logDirectory, fmt.Sprintf("%s.log", config.uuid))
	errPath := path.Join(config.logDirectory, fmt.Sprintf("%s-error.log", config.uuid))
	env := map[string]string{"Path": os.Getenv("PATH")}
	data := &ServiceData{label, pArgs, keepAlive, outPath, errPath, config.workingDirectory, env}
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
