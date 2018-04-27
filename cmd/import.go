package cmd

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/urfave/cli"
	toolkit "github.com/vastness-io/toolkit/pkg/grpc"
	"github.com/vastness-io/vastctl/pkg/import"
	"github.com/vastness-io/vastctl/pkg/in"
	webhook "github.com/vastness-io/vcs-webhook-svc/webhook"
	"google.golang.org/grpc"
	"strings"
)

func GitAction(vcsType string) cli.ActionFunc {
	return func(ctx *cli.Context) error {

		var (
			tracer      = opentracing.GlobalTracer()
			coordinator = ctx.GlobalString(CoordinatorFlagName)
			useFile     = ctx.String(FileFlag)
		)

		var (
			importProjects []*importing.ImportProjectInfo
			err            error
		)

		if useFile != "" {
			importProjects, err = parseImportFile(useFile)
		} else {
			importProjects, err = validateImportArgs(ctx.Args(), nil)
		}

		if err != nil {
			return err
		}

		aggregatedError := &AggregatedError{}

		for _, project := range importProjects {
			err := func() error {

				s := NewSpinner(fmt.Sprintf("importing %s", project.RemoteURL), "")

				s.Start()

				defer s.Stop()

				repoImporter, err := importing.NewVcs(project.RemoteURL, project.RemoteURL)

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

				return err
			}()

			aggregatedError.Add(err)

		}

		return aggregatedError.ToError()
	}
}

func parseImportFile(file string) ([]*importing.ImportProjectInfo, error) {

	fileContents, err := in.ReadFile(file)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSuffix(fileContents, "\n"), "\n")

	if len(lines) == 0 || lines[0] == "" {
		return nil, InvalidFileFormat
	}

	out := make([]*importing.ImportProjectInfo, len(lines))

	for i, _ := range lines {
		const (
			remoteURL = iota
			version
		)

		fields := strings.Fields(lines[i])

		if len(fields) == 1 {
			out[i] = &importing.ImportProjectInfo{
				RemoteURL: fields[remoteURL],
			}
		} else if len(fields) >= 2 {
			out[i] = &importing.ImportProjectInfo{
				RemoteURL: fields[remoteURL],
				Version:   fields[version],
			}
		}
	}

	return out, nil
}

func validateImportArgs(args cli.Args, importProjects []*importing.ImportProjectInfo) ([]*importing.ImportProjectInfo, error) {

	const (
		remoteURL = iota
		version
	)
	out := make([]*importing.ImportProjectInfo, 0)

	if len(args) == 0 {
		if importProjects == nil {
			return nil, MissingRemoteURLErr
		} else {
			return importProjects, nil
		}
	} else if len(args) == 1 {
		out = append(out, &importing.ImportProjectInfo{
			RemoteURL: args[remoteURL],
		})
	} else if len(args) >= 2 {
		out = append(out, &importing.ImportProjectInfo{
			RemoteURL: args[remoteURL],
			Version:   args[version],
		})
	}
	return validateImportArgs(args.Tail(), out)
}
