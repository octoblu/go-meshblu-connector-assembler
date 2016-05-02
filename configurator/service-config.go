package configurator

import "encoding/json"

// Config is the runner connector config structure.
type Config struct {
	Name, DisplayName, Description string

	Dir  string
	Exec string
	Args []string
	Env  []string

	Stderr, Stdout string
}

// NewServiceConfig constructs a new Meshblu instance
func NewServiceConfig(opts Options) *Config {
  args := []string{
    opts.GetLegacyFlag(),
  }
  env := []string{
    opts.GetPathEnv(),
    fmt.Sprintf("MESHBLU_CONNECTOR_NAME=%s", opts.GetConnector()),
    fmt.Sprintf("MESHBLU_CONNECTOR_LEGACY=%v", opts.GetLegacy()),
  }
	return &Config{
    Name: opts.GetServiceName(),
    DisplayName: opts.GetDisplayName(),
    Dir: opts.GetConnectorDir(),
    Exec: opts.GetExecutablePath(),
    Args: args,
    Env: env,
    Stderr: path.Join(opts.GetLogDirectory(), "connector-error.log"),
    Stdout: path.Join(opts.GetLogDirectory(), "connector.log"),
  }
}

// ToJSON serializes the object to the meshblu.json format
func (config *Config) ToJSON() ([]byte, error) {
	return json.Marshal(config)
}
