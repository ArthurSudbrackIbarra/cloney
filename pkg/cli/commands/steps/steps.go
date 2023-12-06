package steps

import (
	"os"
	"path/filepath"

	"github.com/ArthurSudbrackIbarra/cloney/pkg/git"
	"github.com/ArthurSudbrackIbarra/cloney/pkg/metadata"
	"github.com/ArthurSudbrackIbarra/cloney/pkg/templates"
	"github.com/ArthurSudbrackIbarra/cloney/pkg/terminal"
)

// This file defines common steps used by multiple commands.

// Global variables:
// suppressPrints is a flag to determine if the functions in this package should print to the messages.
var suppressPrints bool

// SetSuppressPrints sets the suppressPrints flag.
func SetSuppressPrints(value bool) {
	suppressPrints = value
}

// GetCurrentWorkingDirectory returns the current working directory.
func GetCurrentWorkingDirectory() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		terminal.ErrorMessage("Could not get user's current directory", err)
		return "", err
	}

	return currentDir, nil
}

// GetUserVariablesMap returns the template variables provided by the user.
func GetUserVariablesMap(currentDir, variables string) (map[string]interface{}, error) {
	var err error
	var variablesMap map[string]interface{}

	// First, assume 'variables' is a raw YAML string.
	variablesMap, err = metadata.NewCloneyUserVariablesFromRawYAML(variables)
	if err != nil {
		// In case of error, assume 'variables' is a file path.
		variablesMap, err = metadata.NewCloneyUserVariablesFromFile(variables)
		if err != nil {
			// If it is not a file path, return an empty map.
			return map[string]interface{}{}, nil
		}
	}

	return variablesMap, nil
}

// CreateAndValidateRepository creates the Git repository instance and validates it.
func CreateAndValidateRepository(repositoryURL, branch, tag string) (*git.GitRepository, error) {
	// Create the Git repository instance.
	repository := &git.GitRepository{
		URL:    repositoryURL,
		Branch: branch,
		Tag:    tag,
	}

	// Validate the repository.
	err := repository.Validate()
	if err != nil {
		terminal.ErrorMessage("Error validating repository", err)
		return nil, err
	}
	if !suppressPrints {
		terminal.OKMessage("The template repository reference is valid")
	}

	return repository, nil
}

// AuthenticateToRepository authenticates to the repository if a token is provided.
func AuthenticateToRepository(repository *git.GitRepository, gitToken string) {
	// If the token is empty, try to get it from the environment variable.
	if gitToken == "" {
		gitToken = os.Getenv("CLONEY_GIT_TOKEN")
	}
	// Only if the token is not empty, authenticate to the repository.
	if gitToken != "" {
		repository.AuthenticateWithToken(gitToken)
	}
}

// CalculatePath calculates the absolute path for a given relative or absolute path string.
// If the path is already absolute, it is returned as-is.
// If the path is empty, the defaultName is appended to the current working directory.
// Otherwise, the path is appended to the current working directory.
func CalculatePath(path string, defaultName string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	currentDir, err := GetCurrentWorkingDirectory()
	if err != nil {
		return "", err
	}
	if path == "" {
		return filepath.Join(currentDir, defaultName), nil
	}
	newPath := filepath.Join(currentDir, path)
	if os.PathSeparator == '\\' {
		newPath = filepath.ToSlash(newPath)
	}

	return newPath, nil
}

// CloneRepository clones the repository.
func CloneRepository(repository *git.GitRepository, clonePath string) error {
	err := repository.Clone(clonePath)
	if err != nil {
		terminal.ErrorMessage("Could not clone repository", err)
		return err
	}
	if !suppressPrints {
		terminal.OKMessage("The template repository was cloned")
	}

	return nil
}

// ReadRepositoryMetadata reads the repository metadata.
func ReadRepositoryMetadata(metadataFilePath string) (string, error) {
	metadataBytes, err := os.ReadFile(metadataFilePath)
	if err != nil {
		terminal.ErrorMessage("Could not find the template repository metadata file", err)
		return "", err
	}
	if !suppressPrints {
		terminal.OKMessage("The template repository metadata file was found")
	}

	return string(metadataBytes), nil
}

// ParseRepositoryMetadata parses the repository metadata.
func ParseRepositoryMetadata(metadataContent string, supportedManifestVersions []string) (*metadata.CloneyMetadata, error) {
	// Create the metadata struct from raw YAML data.
	cloneyMetadata, err := metadata.NewCloneyMetadataFromRawYAML(metadataContent, supportedManifestVersions)
	if err != nil {
		terminal.ErrorMessage("Could not parse the template repository template metadata", err)
		return nil, err
	}
	if !suppressPrints {
		terminal.OKMessage("The template repository metadata file is valid")
	}

	return cloneyMetadata, nil
}

// DeleteIgnoredPaths removes files and directories from the specified 'directory' if their
// paths match any of the patterns listed in 'cloneyMetadata.Configuration.IgnorePaths'.
// It iterates through the ignore paths and deletes them recursively.
//
// Parameters:
//   - cloneyMetadata: The metadata containing the configuration of the template repository.
//   - directory: The base directory from which to start removing files and directories.
func DeleteIgnoredPaths(directory string, ignorePaths []string) {
	err := templates.DeleteIgnoredFiles(directory, ignorePaths)
	if err != nil {
		terminal.ErrorMessage("Failed to delete ignored files", err)
	}
}

// MatchUserVariables matches the user variables with the template variables.
func MatchUserVariables(cloneyMetadata *metadata.CloneyMetadata, variablesMap map[string]interface{}) error {
	// Validate if the user variables match the template variables.
	// Also fill default values of the variables if they are not defined.
	var err error
	_, err = cloneyMetadata.MatchUserVariables(variablesMap)
	if err != nil {
		terminal.ErrorMessage("Error validating your template variables", err)
		return err
	}
	if !suppressPrints {
		terminal.OKMessage("Your variables are valid and match the template repository variables")
	}

	return nil
}

// FillDirectory fills template variables in files within the source directory.
func FillDirectory(
	src string,
	ignorePaths []string,
	outputInTerminal bool,
	variablesMap map[string]interface{}) error {
	// Create a new template filler with the provided variables.
	filler := templates.NewTemplateFiller(variablesMap)

	// Fill the template variables in the source directory.
	err := filler.FillDirectory(src, ignorePaths, outputInTerminal)
	if err != nil {
		if outputInTerminal {
			terminal.ErrorMessage("Failed to print results to the terminal", err)
		} else {
			terminal.ErrorMessage("Failed to fill the template variables", err)
		}
		return err
	}

	if !suppressPrints && !outputInTerminal {
		terminal.OKMessage("Template variables successfully filled")
	}

	return nil
}
