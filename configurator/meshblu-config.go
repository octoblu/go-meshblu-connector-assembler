package configurator

import "encoding/json"

// Config interfaces with a remote meshblu server
type Config struct {
	UUID     string `json:"uuid"`
	Token    string `json:"token"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

// NewMeshbluConfig constructs a new Meshblu instance
func NewMeshbluConfig(opts Options) *Config {
	UUID := opts.GetUUID()
	Token := opts.GetToken()
	Hostname := opts.GetHostname()
	Port := opts.GetPort()
	return &Config{UUID, Token, Hostname, Port}
}

// ToJSON serializes the object to the meshblu.json format
func (config *Config) ToJSON() ([]byte, error) {
	return json.Marshal(config)
}
