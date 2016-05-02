package configurator

import "encoding/json"

// MeshbluConfig interfaces with a remote meshblu server
type MeshbluConfig struct {
	UUID     string `json:"uuid"`
	Token    string `json:"token"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

// NewMeshbluConfig constructs a new Meshblu instance
func NewMeshbluConfig(opts Options) *MeshbluConfig {
	return &MeshbluConfig{
		UUID:     opts.GetUUID(),
		Token:    opts.GetToken(),
		Hostname: opts.GetHostname(),
		Port:     opts.GetPort(),
	}
}

// ToJSON serializes the object to the meshblu.json format
func (config *MeshbluConfig) ToJSON() ([]byte, error) {
	return json.Marshal(config)
}
