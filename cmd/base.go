package cmd

import (
	"github.com/urfave/cli"
)

var BaseCommands = []cli.Command{
	{
		Name:        BaseCommandGet,
		Subcommands: GetResourceCommands(),
	},
}
