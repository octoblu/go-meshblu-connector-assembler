package configurator

import (
	"bytes"
	"fmt"
	"path"

	"github.com/DHowett/go-plist"
)

type DarwinServiceData struct {
	Label             string   `plist:"Label"`
	ProgramArguments  []string `plist:"ProgramArguments"`
	KeepAlive         bool     `plist:"KeepAlive"`
	StandardOutPath   string   `plist:"StandardOutPath"`
	StandardErrorPath string   `plist:"StandardErrorPath"`
	WorkingDirectory  string   `plist:"WorkingDirectory"`
}

// DarwinServiceConfig interfaces with a remote meshblu server
type DarwinServiceConfig struct {
	UUID     string `json:"uuid"`
	Token    string `json:"token"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

// NewDarwinServiceConfig constructs a new Meshblu instance
func NewDarwinServiceConfig(uuid, workingDirectory string) *DarwinServiceConfig {
	return &DarwinServiceConfig{uuid, workingDirectory}
}

// String serializes the object to the meshblu.json format
func (config *DarwinServiceConfig) String() ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := plist.NewEncoder(buf)
	label := fmt.Sprintf("com.octoblu.%s", config.uuid)
	pArgs := []string{config.workingDirectory}
	keepAlive := true
	outPath := path.Join("/var/log/octoblu", fmt.Sprintf("%s.log", config.uuid))
	errPath := path.Join("/var/log/octoblu", fmt.Sprintf("%s-error.log", config.uuid))
	data := &DarwinServiceData{label, pArgs, keepAlive, outPath, errPath}
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
