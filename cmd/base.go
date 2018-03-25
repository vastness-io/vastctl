package cmd

import (
	"github.com/urfave/cli"
	"github.com/vastness-io/vastctl/cmd/get"
)

var BaseCommands = []cli.Command{
	{
		Name:        "get",
		Subcommands: get.GetResourceCommands(),
	},
}
