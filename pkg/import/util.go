package importing

import (
	"fmt"
	"github.com/vastness-io/vastctl/pkg/shared"
	"net/url"
	"os"
	"strings"
)

func SplitVcsRemoteUrl(remoteURL string) (owner string, repository string, err error) {

	rurl, err := url.Parse(strings.TrimSuffix(remoteURL, ".git"))

	if err != nil {
		err = shared.InvalidVcsRemoteURL
	}

	var (
		sanitizedPath = strings.TrimPrefix(rurl.Path, "/")
	)

	chunkedPath := strings.Split(sanitizedPath, "/")

	if len(chunkedPath)%2 == 0 {
		return chunkedPath[len(chunkedPath)-2], chunkedPath[len(chunkedPath)-1], nil
	}

	return "", "", shared.InvalidVcsRemoteURL

}

func CleanupTemporaryImportDir(filePath string) {
	if rerr := os.RemoveAll(filePath); rerr != nil {
		fmt.Printf("failed to remove temp path: %s\n", filePath)
	}
}
