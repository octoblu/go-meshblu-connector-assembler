package foreverizer

import (
	"fmt"

	"github.com/octoblu/go-meshblu-connector-assembler/configurator"
)

// Foreverizer interfaces the long running services on the os
type Foreverizer interface {
	Do() error
}

// Client defines the Foreverizer
type Client struct {
	opts *configurator.Options
}

// New constructs a new Foreverizer
func New(opts *configurator.Options) Foreverizer {
	return &Client{opts}
}

// Do will run the setup
func (client *Client) Do() error {
	fmt.Println("foreverizing")
	return Setup(client.opts)
}
