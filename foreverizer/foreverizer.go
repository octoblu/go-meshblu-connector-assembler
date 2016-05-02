package foreverizer

import (
	"fmt"

	"github.com/kardianos/service"
	"github.com/octoblu/go-meshblu-connector-assembler/configurator"
)

// Foreverizer interfaces the long running services on the os
type Foreverizer interface {
	Do() error
}

var logger service.Logger

// Client defines the Foreverizer
type Client struct {
	opts configurator.Options
}

// New constructs a new Foreverizer
func New(opts configurator.Options) Foreverizer {
	return &Client{opts}
}

// Do will run the setup
func (client *Client) Do() error {
	fmt.Println("foreverizing...")
	opts := client.opts
	srvOptions := service.KeyValue{
		"UserService": true,
		"RunAtLoad":   true,
		"KeepAlive":   true,
	}
	svcConfig := &service.Config{
		Name:        opts.GetServiceName(),
		DisplayName: opts.GetDisplayName(),
		Description: opts.GetDescription(),
		Executable:  opts.GetExecutablePath(),
		Option:      srvOptions,
	}

	prg := &Program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}
	fmt.Println("installing service...")
	err = s.Install()
	if err != nil {
		return err
	}
	fmt.Println("starting...")
	err = s.Start()
	if err != nil {
		return err
	}
	return nil
}
