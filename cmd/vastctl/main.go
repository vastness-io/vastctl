package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vastness-io/vastctl/cmd"
	"os"
)

const (
	name        = "vastctl"
	description = "CLI to interact with vastness"
)

var (
	commit  string
	version string
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = description
	app.Version = fmt.Sprintf("%s (%s)", version, commit)
	app.Commands = cmd.BaseCommands
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "coordinator, c",
			Usage:  "Coordinator address",
			Value:  "127.0.0.1:8080",
			EnvVar: "COORDINATOR",
		},
	}
	app.Run(os.Args)
}
