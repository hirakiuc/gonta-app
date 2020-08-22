package release

import (
	"github.com/hirakiuc/gonta-app/usecase"
)

type VersionFetcher struct {
	usecase.Base
}

func NewVersionFetcher(u *Release) *VersionFetcher {
	return &VersionFetcher{
		Base: u.Base,
	}
}

func (f *VersionFetcher) Fetch(repo string, prefix string) ([]string, error) {
	versions := []string{"v1.0.0", "v1.1.0", "v1.1.1"}

	return versions, nil
}
