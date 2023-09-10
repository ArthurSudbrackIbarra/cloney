package git

import (
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GitRepository is a struct that represents a git repository.
type GitRepository struct {
	URL    string
	Branch string
}

// Regex to match a git repository URL. not only for github.
var repositoryRegex = regexp.MustCompile(`^(?:https?|git|ssh):\/\/([^\/]+)\/([^\/]+)\/([^\/]+)(?:)?\.git$`)

// ValidateRepositoryURL validates a git repository URL.
func ValidateRepositoryURL(repositoryURL string) error {
	if !repositoryRegex.MatchString(repositoryURL) {
		return fmt.Errorf("invalid repository URL")
	}
	return nil
}

// GetRepositoryName gets the name of a git repository from its URL.
func GetRepositoryName(repositoryURL string) string {
	matches := repositoryRegex.FindStringSubmatch(repositoryURL)
	if len(matches) < 4 {
		return ""
	}
	return matches[3]
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
