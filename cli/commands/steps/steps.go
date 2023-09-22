package steps

import (
	"os"
	"path/filepath"

	"github.com/ArthurSudbrackIbarra/cloney/git"
	"github.com/ArthurSudbrackIbarra/cloney/metadata"
	"github.com/ArthurSudbrackIbarra/cloney/templates"
	"github.com/ArthurSudbrackIbarra/cloney/terminal"
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
			terminal.ErrorMessage("Error parsing template variables", err)
			return nil, err
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
	// If a token is provided, authenticate with it.
	if gitToken != "" {
		repository.AuthenticateWithToken(gitToken)
	}
}

// CalculatePath calculates the absolute path for a given relative or absolute path string.
// If the provided 'path' is empty, it returns the 'defaultName' as the result.
// If 'path' is an absolute path, it returns 'path' itself.
// If 'path' is a relative path, it joins it with the current working directory to create an absolute path.
func CalculatePath(path string, defaultName string) (string, error) {
	if path == "" {
		return defaultName, nil
	}

	if filepath.IsAbs(path) {
		return path, nil
	}

	currentDir, err := GetCurrentWorkingDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(currentDir, path), nil
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

// FillTemplateVariables fills the template variables in the cloned directory.
func FillTemplateVariables(
	templateOptions templates.TemplateFillOptions,
	ignoreOptions templates.IgnorePathOptions,
	variablesMap map[string]interface{},
) error {
	filler := templates.NewTemplateFiller(variablesMap)
	err := filler.FillDirectory(templateOptions, ignoreOptions)
	if err != nil {
		terminal.ErrorMessage("Error filling template variables", err)
		return err
	}
	if !suppressPrints && !templateOptions.PrintMode {
		terminal.OKMessage("The template variables were filled")
	}

	return nil
}
