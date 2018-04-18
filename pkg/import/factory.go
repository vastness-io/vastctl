package importing

import (
	"github.com/masterminds/vcs"
)

func NewVcs(remoteUrl string, version string) (RepoImporter, error) {

	repo, err := cloneRepository(remoteUrl)

	if err != nil {
		return nil, err
	}

	if err := checkoutSpecificVersionWithFallback(repo, version); err != nil {
		return nil, err
	}

	switch repo.Vcs() {

	case vcs.Git:
		return &gitVcs{
			repo,
		}, nil
	default:
		return nil, vcs.ErrCannotDetectVCS
	}
}
