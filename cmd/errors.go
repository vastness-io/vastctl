package cmd

import (
	"fmt"
	"github.com/urfave/cli"
	"strings"
)

var (
	MissingProjectNameArgErr = cli.NewExitError("Missing project name.", 1)
	MissingVCSTypeArgErr     = cli.NewExitError("Missing vcs type.", 1)
	MissingRemoteURLErr      = cli.NewExitError("Missing vcs remote url.", 1)
	UnsupportedVCSTypeErr    = cli.NewExitError("Unsupported vcs type.", 1)
	ProjectNotFoundErr       = cli.NewExitError("No Project(s) can be found.", 1)
	MalformedRemoteVcsURLErr = cli.NewExitError("Unable to detect vcs from remote url.", 1)
	APIServerUnavailableErr  = cli.NewExitError("Can't connect to the coordinator.", 1)
	RenderAsJSONErr          = cli.NewExitError("Unable to render projects as json.", 1)
	InvalidFileFormat        = cli.NewExitError("Invalid file format.", 1)
	GenericExitErr           = func(err error) error {
		return cli.NewExitError(fmt.Sprintf("something went wrong, %s", err.Error()), 1)
	}
)

type AggregatedError struct {
	errors []error
}

func (e *AggregatedError) Add(err error) *AggregatedError {
	if err != nil {
		e.errors = append(e.errors, err)
	}
	return e
}

func (e *AggregatedError) ToError() error {
	if len(e.errors) == 0 {
		return nil
	}
	return e
}

func (e *AggregatedError) Error() string {
	b := strings.Builder{}
	for _, err := range e.errors {
		b.WriteString(fmt.Sprintf("<%s>\n", err.Error()))
	}
	return b.String()
}
