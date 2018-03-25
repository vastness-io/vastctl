package get

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vastness-io/coordinator-svc/project"
	toolkit "github.com/vastness-io/toolkit/pkg/grpc"
	"github.com/vastness-io/vastctl/pkg/render"
	"google.golang.org/grpc"
	"time"
)

var (
	log                      = logrus.WithField("subcommand", "projects")
	projectResource Resource = &resource{
		singular:  "project",
		plural:    "projects",
		shortName: "prj",
		action:    getProjects(),
		flags: []cli.Flag{
			cli.StringFlag{
				Name:  "type,t",
				Usage: "VCS Type",
			},
			cli.BoolFlag{
				Name:  "all,a",
				Usage: "Get all projects regardless of type",
			},
		},
	}
)

func getProjects() func(*cli.Context) error {
	return func(ctx *cli.Context) error {

		var (
			all     = ctx.Bool("all")
			vcsType = ctx.String("type")
		)

		if !all && vcsType == "" {
			return cli.NewExitError("type should be set", 1)
		}

		var (
			address = ctx.GlobalString("coordinator")
			tracer  = opentracing.GlobalTracer()
		)

		projectsConn, err := toolkit.NewGRPCClient(tracer, log, grpc.WithInsecure())(address)

		if err != nil {
			log.Fatal(err)
		}

		client := project.NewProjectsClient(projectsConn)

		httpCtx, cancelFunc := context.WithTimeout(context.Background(), 10 * time.Second)

		defer cancelFunc()

		defer projectsConn.Close()

		if all {
			res, err := client.GetProjects(httpCtx, &empty.Empty{})

			if err != nil {
				return err
			}

			prettyJson, err := render.PrettyPrintJSON(res.GetProjects())

			if err != nil {
				return cli.NewExitError("Unable to render projects as json", 1)
			}

			fmt.Println(prettyJson)

			return nil
		}

		return nil
	}
}
