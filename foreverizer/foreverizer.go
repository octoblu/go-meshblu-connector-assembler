package foreverizer

// Foreverizer interfaces the long running services on the os
type Foreverizer interface {
	Do(uuid, outputDirectory string) error
}

// Client defines the Foreverizer
type Client struct {
}

// New constructs a new Foreverizer
func New() Foreverizer {
	return &Client{}
}

// Do will run the setup
func (client *Client) Do(uuid, outputDirectory string) error {
	return Setup(uuid, outputDirectory)
}
