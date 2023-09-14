package steps

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ArthurSudbrackIbarra/cloney/git"
	"github.com/ArthurSudbrackIbarra/cloney/metadata"
	"github.com/ArthurSudbrackIbarra/cloney/templates"
	"github.com/ArthurSudbrackIbarra/cloney/utils"
)

// This file defines common steps used by multiple commands.

// GetCurrentWorkingDirectory returns the current working directory.
func GetCurrentWorkingDirectory() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Could not get user's current directory:", err)
		return "", err
	}

	return currentDir, nil
}

// GetUserVariablesMap returns the template variables provided by the user.
func GetUserVariablesMap(currentDir, variablesJSON string, variablesFilePath string) (map[string]interface{}, error) {
	var err error
	var variablesMap map[string]interface{}

	if variablesJSON != "" {
		variablesMap, err = metadata.NewCloneyUserVariablesFromRawJSON(variablesJSON)
		if err != nil {
			fmt.Println("Could not parse your template variables raw JSON:", err)
			return nil, err
		}
	} else {
		variablesMap, err = metadata.NewCloneyUserVariablesFromFile(variablesFilePath)
		if err != nil {
			fmt.Println("Could not read your template variables file:", err)
			return nil, err
		}
	}
	fmt.Println("[OK] Your template variables were found and parsed.")

	return variablesMap, nil
}

// CreateAndValidateRepository creates the Git repository instance and validates it.
func CreateAndValidateRepository(repositoryURL, branch, tag, output string) (*git.GitRepository, error) {
	// Create the Git repository instance.
	repository := &git.GitRepository{
		URL:    repositoryURL,
		Branch: branch,
		Tag:    tag,
	}

	// Validate the repository.
	err := repository.Validate()
	if err != nil {
		fmt.Println("Error validating repository:", err)
		return nil, err
	}
	fmt.Println("[OK] The template repository reference is valid.")

	return repository, nil
}

// AuthenticateToRepository authenticates to the repository if a token is provided.
func AuthenticateToRepository(repository *git.GitRepository, gitToken string) {
	// If a token is provided, authenticate with it.
	if gitToken != "" {
		fmt.Println("[Remove Later] Authenticating with token...")
		repository.AuthenticateWithToken(gitToken)
	}
}

// CalculateClonePath calculates the clone path.
func CalculateClonePath(repository *git.GitRepository, currentDir, output string) string {
	repositoryName := repository.GetName()
	if output == "" {
		return filepath.Join(currentDir, repositoryName)
	}
	return filepath.Join(currentDir, output)
}

// CloneRepository clones the repository.
func CloneRepository(repository *git.GitRepository, clonePath string) error {
	err := repository.Clone(clonePath)
	if err != nil {
		fmt.Println("Could not clone repository:", err)
		return err
	}
	fmt.Println("[OK] The template repository was cloned.")

	return nil
}

// ReadRepositoryMetadata reads the repository metadata.
func ReadRepositoryMetadata(metadataFilePath string) (string, error) {
	metadataBytes, err := os.ReadFile(metadataFilePath)
	if err != nil {
		fmt.Println("Could not read the template repository metadata file:", err)
		return "", err
	}
	fmt.Println("[OK] The template repository metadata file was found.")

	return string(metadataBytes), nil
}

// ParseRepositoryMetadata parses the repository metadata.
func ParseRepositoryMetadata(metadataContent string) (*metadata.CloneyMetadata, error) {
	// Create the metadata struct from raw YAML data.
	cloneyMetadata, err := metadata.NewCloneyMetadataFromRawYAML(metadataContent)
	if err != nil {
		fmt.Println("Could not parse the template repository template metadata:", err)
		return nil, err
	}
	fmt.Println("[OK] The template repository metadata file is valid.")

	return cloneyMetadata, nil
}

// MatchUserVariables matches the user variables with the template variables.
func MatchUserVariables(cloneyMetadata *metadata.CloneyMetadata, variablesMap map[string]interface{}) error {
	// Validate if the user variables match the template variables.
	// Also fill default values of the variables if they are not defined.
	var err error
	variablesMap, err = cloneyMetadata.MatchUserVariables(variablesMap)
	if err != nil {
		fmt.Println("Error validating your template variables:", err)
		return err
	}
	fmt.Println("[OK] Your variables match the template repository variables.")

	return nil
}

// FillTemplateVariables fills the template variables in the cloned directory.
func FillTemplateVariables(
	templateOptions templates.TemplateFillOptions,
	ignoreOptions utils.IgnorePathOptions,
	variablesMap map[string]interface{},
) error {
	filler := templates.NewTemplateFiller(variablesMap)
	err := filler.FillDirectory(templateOptions, ignoreOptions)
	if err != nil {
		fmt.Println("Error filling template variables:", err)
		return err
	}
	if !templateOptions.TerminalMode {
		fmt.Println("[OK] The template variables were filled.")
	}

	return nil
}
