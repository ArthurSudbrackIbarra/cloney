package git

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GitRepository is a struct that represents a git repository.
type GitRepository struct {
	URL    string
	Branch string
}

// CloneRepository clones a git repository.
func CloneRepository(repository *GitRepository, path string) error {
	if repository == nil {
		return fmt.Errorf("git repository cannot be nil")
	}
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: repository.URL,
		ReferenceName: plumbing.ReferenceName(
			fmt.Sprintf("refs/heads/%s", repository.Branch),
		),
	})
	return err
}
