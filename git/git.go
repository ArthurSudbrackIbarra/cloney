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

	// Auth is the authentication method to use when cloning the repository.
	Auth transport.AuthMethod
}

// repositoryRegex is a regular expression to match a git repository URL.
var repositoryRegex = regexp.MustCompile(`^(?:https?|git|ssh):\/\/([^\/]+)\/([^\/]+)\/([^\/]+)(?:)?\.git$`)

// NewGitRepository creates a new GitRepository struct.
func NewGitRepository(url string, branch string) (*GitRepository, error) {
	// Validate parameters.
	if !repositoryRegex.MatchString(url) {
		return nil, fmt.Errorf("invalid repository URL")
	}
	return &GitRepository{
		URL:    url,
		Branch: branch,
	}, nil
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
	var auth transport.AuthMethod
	// If the URL is HTTPS, use the token authentication method.
	if strings.HasPrefix(r.URL, "https://") {
		auth = r.Auth
	}
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: r.URL,
		ReferenceName: plumbing.ReferenceName(
			fmt.Sprintf("refs/heads/%s", r.Branch),
		),
		Auth: auth,
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
