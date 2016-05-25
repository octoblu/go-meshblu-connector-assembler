package foreverizer

import (
	"runtime"

	"github.com/kardianos/service"
	"github.com/octoblu/go-meshblu-connector-assembler/configurator"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("meshblu-connector-assembler:foreverizer")

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
	debug("foreverizing...")
	opts := client.opts
	userService := true
	if runtime.GOOS == "linux" {
		userService = false
	}
	srvOptions := service.KeyValue{
		"UserService": userService,
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

	debug("maybe stop and removing service...")
	s.Uninstall()

	debug("installing service...")
	err = s.Install()
	if err != nil {
		return err
	}
	debug("starting...")
	err = s.Start()
	if err != nil {
		return err
	}
	return nil
}
