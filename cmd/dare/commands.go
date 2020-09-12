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
	ui := &cli.PrefixedUi{
		OutputPrefix: "==> ",
		InfoPrefix:   "    ",
		ErrorPrefix:  "==> [ERROR]",
		Ui:           &cli.BasicUi{Writer: os.Stdout},
	}
	Commands = map[string]cli.CommandFactory{

		"keygen": func() (cli.Command, error) {
			return &command.KeygenCommand{
				Ui: ui,
			}, nil
		},
		"encrypt": func() (cli.Command, error) {
			return &command.EncryptCommand{
				Ui: ui,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Revision: version.Revision,
				Version:  version.Version,
				Ui:       ui,
			}, nil
		},
	}
}
