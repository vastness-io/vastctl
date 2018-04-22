package cmd

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/urfave/cli"
	toolkit "github.com/vastness-io/toolkit/pkg/grpc"
	"github.com/vastness-io/vastctl/pkg/import"
	webhook "github.com/vastness-io/vcs-webhook-svc/webhook"
	"google.golang.org/grpc"
)

func GitAction(vcsType string) cli.ActionFunc {
	return func(ctx *cli.Context) error {

		var (
			tracer      = opentracing.GlobalTracer()
			coordinator = ctx.GlobalString(CoordinatorFlagName)
		)

		remoteURl, version, err := validateImportArgs(ctx.Args())

		if err != nil {
			return err
		}

		s := NewSpinner(fmt.Sprintf("importing %s", remoteURl), "")

		s.Start()

		repoImporter, err := importing.NewVcs(remoteURl, version)

		if err != nil {
			return err
		}

		event, err := repoImporter.MapToPushEvent(vcsType)

		if err != nil {
			return err
		}

		coordinatorConn, err := toolkit.NewGRPCClient(tracer, nil, grpc.WithInsecure())(coordinator)

		if err != nil {
			return err
		}

		defer coordinatorConn.Close()

		cc := webhook.NewVcsEventClient(coordinatorConn)

		httpCtx, cancelFunc := context.WithTimeout(context.Background(), TimeOut)

		defer cancelFunc()

		_, err = cc.OnPush(httpCtx, event)

		s.Stop()

		return err
	}
}
func validateImportArgs(args cli.Args) (remoteURL, version string, err error) {
	if len(args) == 1 {
		remoteURL = args[0]
	} else if len(args) >= 2 {
		remoteURL = args[0]
		version = args[1]
	} else {
		err = MissingRemoteURLErr
	}

	return

}
