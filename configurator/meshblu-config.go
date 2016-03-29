package configurator

import "encoding/json"

// Config interfaces with a remote meshblu server
type Config struct {
	UUID     string `json:"uuid"`
	Token    string `json:"token"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

// NewConfig constructs a new Meshblu instance
func NewConfig(UUID, Token, Hostname string, Port int) *Config {
	return &Config{UUID, Token, Hostname, Port}
}

// ToJSON serializes the object to the meshblu.json format
func (config *Config) ToJSON() ([]byte, error) {
	return json.Marshal(config)
}
