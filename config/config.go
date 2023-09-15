package config

import (
	"github.com/spf13/viper"
)

// AppConfig represents the application configuration.
type AppConfig struct {
	// AppVersion is the current version of the application.
	AppVersion string

	// MetadataFileName is the name of the metadata file in the repository.
	// This file contains the information about the template repository.
	MetadataFileName string

	// MetadataManifestVersion is the version of the metadata manifest supported by this version of cloney.
	MetadataManifestVersion string

	// DefaultUserVariablesFileName is the default name of the file containing the user variables.
	// This file is used to fill the template variables in the cloned directory.'
	DefaultUserVariablesFileName string

	// DefaultDryRunDirectoryName is the default name of the directory created when running a template in dryrun mode.
	DefaultDryRunDirectoryName string

	// DefaultCloneyProjectName is the default name to use when creating a new cloney project.
	DefaultCloneyProjectName string

	// DefaultMetadataDescriptionValue is the default value for the description field in the metadata file.
	DefaultMetadataDescriptionValue string

	// DefaultMetadataLicenseValue is the default value for the license field in the metadata file.
	DefaultMetadataLicenseValue string

	// CloneyExampleRepositoryURL is the URL of the cloney example repository.
	CloneyExampleRepositoryURL string

	// GitToken is the token used to authenticate when dealing with private git repositories.
	GitToken string `mapstructure:"GIT_TOKEN"`
}

// globalConfig is the global application configuration.
var globalConfig = &AppConfig{
	// Default values.
	AppVersion: "0.0.0",

	MetadataFileName:        ".cloney.yaml",
	MetadataManifestVersion: "v1",

	DefaultUserVariablesFileName: ".cloney-vars.yaml",
	DefaultDryRunDirectoryName:   "cloney-dry-run-results",
	DefaultCloneyProjectName:     "cloney-template",

	DefaultMetadataDescriptionValue: "A cloney template repository",
	DefaultMetadataLicenseValue:     "MIT",

	CloneyExampleRepositoryURL: "https://github.com/ArthurSudbrackIbarra/cloney-test.git",
}

// LoadConfig loads the global application configuration.
func LoadConfig() error {
	// Enable reading environment variables with a prefix.
	viper.SetEnvPrefix("CLONEY")
	viper.AutomaticEnv()

	// Unmarshal the configuration into the globalConfig variable.
	err := viper.Unmarshal(globalConfig)
	if err != nil {
		return err
	}
	return nil
}

// GetAppConfig returns a copy of the global application configuration.
func GetAppConfig() *AppConfig {
	return &AppConfig{
		AppVersion: globalConfig.AppVersion,

		MetadataFileName:        globalConfig.MetadataFileName,
		MetadataManifestVersion: globalConfig.MetadataManifestVersion,

		DefaultUserVariablesFileName: globalConfig.DefaultUserVariablesFileName,
		DefaultDryRunDirectoryName:   globalConfig.DefaultDryRunDirectoryName,
		DefaultCloneyProjectName:     globalConfig.DefaultCloneyProjectName,

		DefaultMetadataDescriptionValue: globalConfig.DefaultMetadataDescriptionValue,
		DefaultMetadataLicenseValue:     globalConfig.DefaultMetadataLicenseValue,

		CloneyExampleRepositoryURL: globalConfig.CloneyExampleRepositoryURL,

		GitToken: globalConfig.GitToken,
	}
}
