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
	"sync"
)

func GitAction(vcsType string) cli.ActionFunc {
	return func(ctx *cli.Context) error {

		var (
			tracer               = opentracing.GlobalTracer()
			coordinator          = ctx.GlobalString(CoordinatorFlagName)
			useFile              = ctx.String(FileFlag)
			maxConcurrentImports = ctx.Int(MaxConcurrentFlag)
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

		return BatchProcess(importProjects, maxConcurrentImports, vcsType, coordinator, tracer)

	}
}
func BatchProcess(importProjects []*importing.ImportProjectInfo, maxConcurrent int, vcsType, coordinatorAddr string, tracer opentracing.Tracer) error {
	var (
		chunkedProjects = ChunkImportProjects(importProjects, maxConcurrent)
		wg              sync.WaitGroup
		errCh           = make(chan error, maxConcurrent)
		aggregatedError = &AggregatedError{}
	)

	for _, group := range chunkedProjects {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ImportBatch(group, vcsType, coordinatorAddr, tracer, errCh)
		}()
		wg.Wait()
	}

	for range chunkedProjects {
		aggregatedError.Add(<-errCh)
	}

	return aggregatedError.ToError()
}

func ImportBatch(batch []*importing.ImportProjectInfo, vcsType, coordinatorAddr string, tracer opentracing.Tracer, errCh chan<- error) {

	var (
		tmpErrCh = make(chan error)

		wg sync.WaitGroup
	)

	go func() {
		aggregatedError := &AggregatedError{}

		for err := range tmpErrCh {
			aggregatedError.Add(err)
		}
		errCh <- aggregatedError.ToError()
	}()

	for _, project := range batch {
		wg.Add(1)
		go func(prj *importing.ImportProjectInfo) {
			defer wg.Done()
			fmt.Printf("importing %s\n", prj.RemoteURL)

			repoImporter, err := importing.NewVcs(prj.RemoteURL, prj.Version)

			if err != nil {
				tmpErrCh <- err
				return
			}

			event, err := repoImporter.MapToPushEvent(vcsType)

			if err != nil {
				tmpErrCh <- err
				return
			}

			coordinatorConn, err := toolkit.NewGRPCClient(tracer, nil, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)))(coordinatorAddr)

			if err != nil {
				tmpErrCh <- err
				return
			}

			defer coordinatorConn.Close()

			cc := webhook.NewVcsEventClient(coordinatorConn)

			httpCtx, cancelFunc := context.WithTimeout(context.Background(), TimeOut)

			defer cancelFunc()

			_, err = cc.OnPush(httpCtx, event)

			tmpErrCh <- err
		}(project)
	}

	wg.Wait()

	close(tmpErrCh)

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
