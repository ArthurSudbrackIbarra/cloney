package git

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// GitRepository is a struct that represents a git repository.
type GitRepository struct {
	// URL is the URL of the git repository.
	URL string

	// Branch is the branch of the git repository.
	Branch string

	// Tag is the tag of the git repository.
	Tag string

	// Auth is the authentication method to use when cloning the repository.
	Auth transport.AuthMethod
}

// repositoryRegex is a regular expression to match a git repository URL.
var repositoryRegex = regexp.MustCompile(`^(?:https?|git|ssh):\/\/([^\/]+)\/([^\/]+)\/([^\/]+)(?:)?\.git$`)

// MatchesGitRepositoryURL returns true if a string matches a git repository URL.
func MatchesGitRepositoryURL(str string) bool {
	return repositoryRegex.MatchString(str)
}

// Validate validates if a git repository is valid.
func (r *GitRepository) Validate() error {
	if !MatchesGitRepositoryURL(r.URL) {
		return fmt.Errorf("invalid repository URL")
	}
	// Either branch or tag must be specified, but not both.
	if r.Branch == "" && r.Tag == "" {
		return fmt.Errorf("branch or tag must be specified")
	}
	if r.Branch != "" && r.Tag != "" {
		return fmt.Errorf("branch and tag cannot be specified at the same time")
	}
	return nil
}

// GetName returns the name of a git repository.
func (r *GitRepository) GetName() string {
	matches := repositoryRegex.FindStringSubmatch(r.URL)
	if len(matches) < 4 {
		return ""
	}
	// matches[0] is the whole string.
	// matches[1] is the host.
	// matches[2] is the owner.
	// matches[3] is the repository name.
	return matches[3]
}

// AuthenticateWithToken authenticates to the git repository with a token.
func (r *GitRepository) AuthenticateWithToken(token string) {
	r.Auth = &http.BasicAuth{
		Username: "token",
		Password: token,
	}
}

// Clone clones the git repository.
func (r *GitRepository) Clone(path string) error {
	var referenceName plumbing.ReferenceName
	// If the branch is specified, use it as the reference name.
	// If the tag is specified, use it as the reference name.
	if r.Branch != "" {
		referenceName = plumbing.NewBranchReferenceName(r.Branch)
	} else if r.Tag != "" {
		referenceName = plumbing.NewTagReferenceName(r.Tag)
	}

	// If the URL is HTTPS, use the token authentication method.
	var auth transport.AuthMethod
	if strings.HasPrefix(r.URL, "https://") {
		auth = r.Auth
	}

	// Clone the repository.
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           r.URL,
		ReferenceName: referenceName,
		Auth:          auth,
	})
	return err
}

// GetFileContent returns the content of a raw file in the git repository.
func (r *GitRepository) GetFileContent(filePath string) (string, error) {
	// Clone the repository in a temporary directory.
	temporaryDir := fmt.Sprintf("%s/cloney/%s", os.TempDir(), r.GetName())
	err := r.Clone(temporaryDir)
	if err != nil {
		return "", err
	}

	// Get the file content.
	fileContent, err := os.ReadFile(fmt.Sprintf("%s/%s", temporaryDir, filePath))
	if err != nil {
		// Clean up the temporary directory on error.
		os.RemoveAll(temporaryDir)
		return "", err
	}

	// Delete the temporary directory.
	os.RemoveAll(temporaryDir)
	if err != nil {
		return "", err
	}

	return string(fileContent), nil
}
