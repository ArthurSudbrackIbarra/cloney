package config

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
}

// globalConfig is the global application configuration.
var globalConfig = &AppConfig{
	AppVersion: "0.1.1",

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

// GetAppConfig returns a copy of the global application configuration.
func GetAppConfig() AppConfig {
	return *globalConfig
}
