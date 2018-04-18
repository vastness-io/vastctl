package shared

import "errors"

var (
	EmptyRepositoryErr  = errors.New("empty repository")
	InvalidVcsRemoteURL = errors.New("invalid vcs remote url")
)
