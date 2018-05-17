package importing

import (
	event "github.com/vastness-io/vcs-webhook-svc/webhook"
)

type ImportProjectInfo struct {
	RemoteURL string
	Version   string
}

type RepoImporter interface {
	MapToPushEvent(vcsType string) (*event.VcsPushEvent, error)
}
