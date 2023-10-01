package config

import (
	"github.com/spf13/viper"
)

// AppConfig represents the application configuration.
type AppConfig struct {
	// AppVersion is the current version of the application.
	AppVersion string

	// MetadataFileName is the name of the metadata file in the template repository.
	// This file contains information about the template repository.
	MetadataFileName string

	// MetadataManifestVersion is used when creating a new template repository.
	// It will always be the latest version supported by the current version of Cloney.
	MetadataManifestVersion string

	// SupportedManifestVersions is a list of the supported metadata manifest versions by this version of Cloney.
	SupportedManifestVersions []string

	// DefaultUserVariablesFileName is the default name of the file containing the user variables.
	// This file is used to fill the template variables in the cloned directory.
	DefaultUserVariablesFileName string

	// DefaultDryRunDirectoryName is the default name of the directory created when running a template repository in dryrun mode.
	DefaultDryRunDirectoryName string

	// DefaultCloneyProjectName is the default name to use when creating a new Cloney project.
	DefaultCloneyProjectName string

	// DefaultMetadataDescriptionValue is the default value for the description field in the metadata file when creating a new template repository.
	DefaultMetadataDescriptionValue string

	// DefaultMetadataLicenseValue is the default value for the license field in the metadata file when creating a new template repository.
	DefaultMetadataLicenseValue string

	// DefaultMetadataTemplateVersionValue is the default value for the template version field in the metadata file when creating a new template repository.
	DefaultMetadataTemplateVersionValue string

	// CloneyExampleRepositoryURL is the URL of the Cloney example repository used when creating a new template repository.
	CloneyExampleRepositoryURL string

	// GitToken is the token used to authenticate when dealing with private git repositories.
	// This variable is configured using the CLONEY_GIT_TOKEN environment variable.
	GitToken string `mapstructure:"GIT_TOKEN"`
}

// globalConfig is the global application configuration.
var globalConfig = &AppConfig{
	// Default values.
	AppVersion: "0.1.0",

	MetadataFileName:        ".cloney.yaml",
	MetadataManifestVersion: "v1",
	SupportedManifestVersions: []string{
		"v1",
	},

	DefaultUserVariablesFileName: ".cloney-vars.yaml",
	DefaultDryRunDirectoryName:   "cloney-dry-run-results",
	DefaultCloneyProjectName:     "cloney-template",

	DefaultMetadataDescriptionValue:     "A Cloney template repository",
	DefaultMetadataLicenseValue:         "MIT",
	DefaultMetadataTemplateVersionValue: "0.0.0",

	CloneyExampleRepositoryURL: "https://github.com/ArthurSudbrackIbarra/cloney-example.git",
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
		SupportedManifestVersions: []string{
			globalConfig.SupportedManifestVersions[0],
		},

		DefaultUserVariablesFileName: globalConfig.DefaultUserVariablesFileName,
		DefaultDryRunDirectoryName:   globalConfig.DefaultDryRunDirectoryName,
		DefaultCloneyProjectName:     globalConfig.DefaultCloneyProjectName,

		DefaultMetadataDescriptionValue:     globalConfig.DefaultMetadataDescriptionValue,
		DefaultMetadataLicenseValue:         globalConfig.DefaultMetadataLicenseValue,
		DefaultMetadataTemplateVersionValue: globalConfig.DefaultMetadataTemplateVersionValue,

		CloneyExampleRepositoryURL: globalConfig.CloneyExampleRepositoryURL,

		GitToken: globalConfig.GitToken,
	}
}
