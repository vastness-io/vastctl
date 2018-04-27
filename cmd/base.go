package cmd

import (
	"github.com/urfave/cli"
	"github.com/vastness-io/vcs-webhook-svc/webhook"
)

var BaseCommands = []cli.Command{
	{
		Name:        BaseCommandGet,
		Description: "Retrieves a list of resource names",
		Subcommands: cli.Commands{
			{
				Name:      "projects",
				ShortName: "prj",
				Aliases: []string{
					"project",
				},
				Usage:       "Retrieves a list of projects",
				Description: "Retrieves a list of project names",
				Action:      HandleErr(getProjects()),
			},
		},
	},
	{
		Name:        BaseCommandDescribe,
		Description: "Prints detailed information about various resources in Vastness",
		Subcommands: cli.Commands{
			{
				Name:      "projects",
				ShortName: "prj",
				Aliases: []string{
					"project",
				},
				Usage:       "Describes a project",
				Description: "Prints detailed information of a particular vcs project",
				ArgsUsage:   "NAME [github|bitbucket-server]",
				Action:      HandleErr(DescribeProject()),
			},
		},
	},
	{
		Name:        BaseCommandImport,
		Description: "Imports a vcs project to Vastness",
		Subcommands: cli.Commands{
			{
				Name:        "github",
				Action:      HandleErr(GitAction(vcs.VcsType_GITHUB.String())),
				Usage:       "Imports a github project",
				Description: "Imports a github project to Vastness",
				ArgsUsage:   "REMOTE_URL BRANCH",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  CreateFlagName(FileFlag, FileShortFlag),
						Usage: FileFlagUsage,
					},
				},
			},
			{
				Name:        "bitbucket-server",
				Action:      HandleErr(GitAction(vcs.VcsType_BITBUCKET_SERVER.String())),
				Usage:       "Imports a bitbucket server project",
				Description: "Imports a bitbucket server project to Vastness",
				ArgsUsage:   "REMOTE_URL BRANCH",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  CreateFlagName(FileFlag, FileShortFlag),
						Usage: FileFlagUsage,
					},
				},
			},
		},
	},
}
