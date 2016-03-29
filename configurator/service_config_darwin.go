package configurator

import (
	"bytes"
	"fmt"
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
		Label             string   `plist:"Label"`
		ProgramArguments  []string `plist:"ProgramArguments"`
		KeepAlive         bool     `plist:"KeepAlive"`
		StandardOutPath   string   `plist:"StandardOutPath"`
		StandardErrorPath string   `plist:"StandardErrorPath"`
		WorkingDirectory  string   `plist:"WorkingDirectory"`
	}

	buf := &bytes.Buffer{}
	encoder := plist.NewEncoder(buf)
	label := fmt.Sprintf("com.octoblu.%s", config.uuid)
	pArgs := []string{config.workingDirectory}
	keepAlive := true
	outPath := path.Join(config.logDirectory, fmt.Sprintf("%s.log", config.uuid))
	errPath := path.Join(config.logDirectory, fmt.Sprintf("%s-error.log", config.uuid))
	data := &ServiceData{label, pArgs, keepAlive, outPath, errPath, config.workingDirectory}
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
