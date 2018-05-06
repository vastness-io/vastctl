package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/masterminds/vcs"
	"github.com/urfave/cli"
	"github.com/vastness-io/coordinator-svc/project"
	"github.com/vastness-io/vastctl/pkg/import"
	"github.com/vastness-io/vastctl/pkg/shared"
	vcs2 "github.com/vastness-io/vcs-webhook-svc/webhook"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

func HandleErr(actionFunc cli.ActionFunc) cli.ActionFunc {
	return func(context *cli.Context) error {
		err := actionFunc(context)

		if err != nil {
			_, ok := err.(*cli.ExitError)
			if ok {
				return err
			}
			switch err {
			case vcs.ErrCannotDetectVCS:
				return MalformedRemoteVcsURLErr

			case shared.InvalidVcsRemoteURL:
				return MalformedRemoteVcsURLErr
			}
			switch status.Code(err) {
			case codes.Unavailable:
				return APIServerUnavailableErr
			case codes.NotFound:
				return ProjectNotFoundErr
			default:
				return GenericExitErr(err)
			}
		}

		return nil

	}
}

func ExtractNamesFromProjects(projects []*project.Project) (out []string) {
	for _, prj := range projects {
		out = append(out, prj.GetName())
	}
	return
}

func MapVcsTypesToVcsMessage(vcsType string) (string, error) {
	switch strings.ToLower(vcsType) {
	case "github":
		return vcs2.VcsType_GITHUB.String(), nil
	case "bitbucket-server":
		return vcs2.VcsType_BITBUCKET_SERVER.String(), nil

	default:
		return "", UnsupportedVCSTypeErr
	}
}

func CreateFlagName(full, short string) string {
	return fmt.Sprintf("%s, %s", full, short)
}

func NewSpinner(prefix, finalMsg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = prefix

	if finalMsg != "" {
		s.FinalMSG = finalMsg
	}
	return s

}

func ChunkImportProjects(importProjects []*importing.ImportProjectInfo, n int) [][]*importing.ImportProjectInfo {

	var groups [][]*importing.ImportProjectInfo

	chunkSize := n

	for i := 0; i < len(importProjects); i += chunkSize {
		end := i + chunkSize

		if end > len(importProjects) {
			end = len(importProjects)
		}

		groups = append(groups, importProjects[i:end])
	}

	return groups
}
