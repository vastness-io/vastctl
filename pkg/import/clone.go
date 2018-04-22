package importing

import (
	"github.com/masterminds/vcs"
	"io/ioutil"
)

func newRepository(remoteURL string) (vcs.Repo, error) {
	local, _ := ioutil.TempDir("", "import")

	repo, err := vcs.NewRepo(remoteURL, local)

	if err != nil {
		return nil, err
	}

	return repo, nil
}

func cloneRepository(remoteURL string) (vcs.Repo, error) {
	repo, err := newRepository(remoteURL)

	if err != nil {
		return nil, err
	}

	return repo, repo.Get()
}

func checkoutSpecificVersionWithFallback(repo vcs.Repo, version string) error {

	oldVer, err := repo.Current()

	if err != nil {
		return err
	}

	if version != "" {
		err = repo.UpdateVersion(version)
		if err != nil {
			return repo.UpdateVersion(oldVer)
		}
	}

	return nil
}
