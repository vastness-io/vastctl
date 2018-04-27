package cmd

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
	"github.com/urfave/cli"
	"github.com/vastness-io/coordinator-svc/project"
	toolkit "github.com/vastness-io/toolkit/pkg/grpc"
	"github.com/vastness-io/vastctl/pkg/render"
	"google.golang.org/grpc"
)

func getProjects() func(*cli.Context) error {
	return func(ctx *cli.Context) error {

		var (
			address = ctx.GlobalString(CoordinatorFlagName)
			tracer  = opentracing.GlobalTracer()
		)

		clientConn, err := toolkit.NewGRPCClient(tracer, nil, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(-1)))(address)

		if err != nil {
			return err
		}

		defer clientConn.Close()

		client := project.NewProjectsClient(clientConn)

		httpCtx, cancelFunc := context.WithTimeout(context.Background(), TimeOut)

		defer cancelFunc()

		res, err := client.GetProjects(httpCtx, &empty.Empty{})

		if err != nil {
			return err
		}

		prettyJson, err := render.PrettyPrintJSON(ExtractNamesFromProjects(res.GetProjects()))

		if err != nil {
			return RenderAsJSONErr
		}

		fmt.Println(prettyJson)

		return nil
	}

}

func DescribeProject() func(*cli.Context) error {
	return func(ctx *cli.Context) error {

		if ctx.NArg() == 0 {
			return MissingProjectNameArgErr
		} else if ctx.NArg() == 1 {
			return MissingVCSTypeArgErr
		}

		var (
			name    = ctx.Args().First()
			address = ctx.GlobalString(CoordinatorFlagName)
			tracer  = opentracing.GlobalTracer()
		)

		vcsType, err := MapVcsTypesToVcsMessage(ctx.Args().Get(1))

		if err != nil {
			return err
		}

		clientConn, err := toolkit.NewGRPCClient(tracer, nil, grpc.WithInsecure())(address)

		if err != nil {
			return err
		}

		defer clientConn.Close()

		client := project.NewProjectsClient(clientConn)

		httpCtx, cancelFunc := context.WithTimeout(context.Background(), TimeOut)

		defer cancelFunc()

		prj, err := client.GetProject(httpCtx, &project.GetProjectMessage{
			Name: name,
			Type: vcsType,
		})

		if err != nil {
			return err
		}

		prettyJson, err := render.PrettyPrintJSON(prj)

		if err != nil {
			return RenderAsJSONErr
		}

		fmt.Println(prettyJson)

		return nil
	}
}
