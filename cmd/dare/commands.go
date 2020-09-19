package main

import (
	"os"

	command "github.com/da-moon/terraform-provider-dare/cmd/dare/command"
	version "github.com/da-moon/version"
	cli "github.com/mitchellh/cli"
)

// Commands is the mapping of all the available Serf commands.
var Commands map[string]cli.CommandFactory

func init() {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	Commands = map[string]cli.CommandFactory{
		"dd": func() (cli.Command, error) {
			return &command.DDCommand{
				UI: ui,
			}, nil
		},
		"keygen": func() (cli.Command, error) {
			return &command.KeygenCommand{
				UI: ui,
			}, nil
		},
		"encrypt": func() (cli.Command, error) {
			return &command.EncryptCommand{
				UI: ui,
			}, nil
		},
		"decrypt": func() (cli.Command, error) {
			return &command.DecryptCommand{
				UI: ui,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Revision: version.Revision,
				Version:  version.Version,
				UI:       ui,
			}, nil
		},
	}
}
