package importing

import (
	"fmt"
	"github.com/masterminds/vcs"
	"github.com/vastness-io/vastctl/pkg/shared"
	event "github.com/vastness-io/vcs-webhook-svc/webhook"
	"regexp"
	"strings"
	"time"
)

type gitVcs struct {
	vcs.Repo
}

func (g *gitVcs) MapToPushEvent(vcsType string) (*event.VcsPushEvent, error) {

	defer CleanupTemporaryImportDir(g.LocalPath())

	ev, err := WalkGitTree(g, vcsType)

	if err != nil {
		return nil, err
	}

	return ev, nil

}

func WalkGitTree(repo vcs.Repo, vcsType string) (*event.VcsPushEvent, error) {
	from, err := repo.RunFromDir("git", "rev-list", "--reverse", "HEAD")

	if err != nil {
		return nil, shared.EmptyRepositoryErr
	}

	ref, err := repo.Current()

	if err != nil {
		return nil, err
	}

	var (
		commits = strings.Split(strings.TrimSuffix(string(from), "\n"), "\n")
		out     = event.VcsPushEvent{}
	)

	owner, repository, err := SplitVcsRemoteUrl(repo.Remote())

	if err != nil {
		return nil, err
	}

	for _, commit := range commits {
		b, err := repo.RunFromDir("git", "show", "--name-status", "-r", "--pretty=format:%an,%ae,%aI,%cn,%ce,%cI", commit)

		if err != nil {
			return nil, err
		}

		out.Commits = append(out.Commits, createPushCommit(commit, string(b)))

	}

	out.HeadCommit = out.GetCommits()[len(out.GetCommits())-1]

	out.Ref = ref

	out.Created = true
	out.VcsType = event.VcsType(event.VcsType_value[vcsType])

	out.Organization = &event.User{
		Login: owner,
		Type:  vcsType,
		Name:  owner,
	}
	out.Repository = &event.Repository{
		Owner: &event.User{
			Login: owner,
			Type:  vcsType,
			Name:  owner,
		},
		Name:     repository,
		FullName: fmt.Sprintf("%s/%s", owner, repository),
		Organization: &event.User{
			Login: owner,
			Type:  vcsType,
			Name:  owner,
		},
	}

	return &out, nil

}

func createPushCommit(sha, fileString string) *event.PushCommit {
	const (
		status = iota
		name
		renamedName
	)

	out := event.PushCommit{
		Sha: sha,
	}

	commitInfoRaw := strings.Split(fileString, "\n")

	if len(commitInfoRaw) > 1 {

		const (
			authorNameIndex = iota
			authorEmailIndex
			authorDateIndex
			committerNameIndex
			committerEmailIndex
			committerDateIndex
			commitInfoTotalFields
			//TODO get commit message
		)

		commitInfo := strings.Split(commitInfoRaw[0], ",")

		if len(commitInfo) == commitInfoTotalFields {
			out.Author = &event.CommitAuthor{
				Name:     commitInfo[authorNameIndex],
				Email:    commitInfo[authorEmailIndex],
				Username: commitInfo[authorNameIndex],
				Date:     commitInfo[authorDateIndex],
			}

			out.Timestamp = commitInfo[authorDateIndex]

			out.Committer = &event.CommitAuthor{
				Name:     commitInfo[committerNameIndex],
				Email:    commitInfo[committerEmailIndex],
				Username: commitInfo[committerNameIndex],
				Date:     commitInfo[committerDateIndex],
			}

		} else {
			var (
				defaultName = "Merge commit"
				defaultDate = time.Time{}.UTC().Format(time.RFC3339) //default
			)

			out.Author = &event.CommitAuthor{
				Name:     defaultName,
				Username: defaultName,
				Date:     defaultDate,
			}

			out.Timestamp = defaultDate

			out.Committer = &event.CommitAuthor{
				Name:     defaultName,
				Username: defaultName,
				Date:     defaultDate,
			}
		}

		for _, file := range commitInfoRaw[1:] {

			if nameStatus := strings.Fields(file); len(nameStatus) != 0 {

				status := nameStatus[status]

				renamedRegex := regexp.MustCompile(`^R[0-9]+$`)

				switch {

				case status == "A":
					out.Added = append(out.Added, nameStatus[name])

				case status == "M":
					out.Modified = append(out.Modified, nameStatus[name])

				case renamedRegex.MatchString(status):
					out.Removed = append(out.Removed, nameStatus[name])
					out.Added = append(out.Added, nameStatus[renamedName])
				case status == "D" || status == "RM":
					out.Removed = append(out.Removed, nameStatus[name])
				}
			}
		}
	}

	return &out

}
