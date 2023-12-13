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

	// KnownIgnorePaths is a collection of paths that should be excluded in the template repository.
	// These paths typically include directories like .git and node_modules, which are not relevant to the template.
	KnownIgnorePaths []string

	// IgnorePrefix is the prefix used to ignore files and directories in the template repository.
	IgnorePrefix string

	// DefaultMetadataDescriptionValue is the default value for the description field in the metadata file when creating a new template repository.
	DefaultMetadataDescriptionValue string

	// DefaultMetadataLicenseValue is the default value for the license field in the metadata file when creating a new template repository.
	DefaultMetadataLicenseValue string

	// DefaultMetadataTemplateVersionValue is the default value for the template version field in the metadata file when creating a new template repository.
	DefaultMetadataTemplateVersionValue string

	// CloneyExampleRepositoryURL is the URL of the Cloney example repository used when creating a new template repository.
	CloneyExampleRepositoryURL string

	// CloneyDocumentationURL is the URL of the Cloney documentation.
	CloneyDocumentationURL string
}

// globalConfig is the global application configuration.
var globalConfig = &AppConfig{
	//! AppVersion is set automatically during the pipeline that tags the release (.github/workflows/auto_tag.yaml).
	//! Keep this value as it is.
	AppVersion: "X.X.X",

	MetadataFileName:        ".cloney.yaml",
	MetadataManifestVersion: "v1",
	SupportedManifestVersions: []string{
		"v1",
	},

	DefaultUserVariablesFileName: ".cloney-vars.yaml",
	DefaultDryRunDirectoryName:   "cloney-dry-run-results",
	DefaultCloneyProjectName:     "cloney-template",

	KnownIgnorePaths: []string{
		".cloney.yaml",      // Cloney metadata file.
		".cloney-vars.yaml", // Cloney default user variables file.
		".git",              // Git directory.
		"node_modules",      // Node.js modules directory.
		".venv",             // Python virtual environment directory.
	},
	IgnorePrefix: "__",

	DefaultMetadataDescriptionValue:     "A Cloney template repository",
	DefaultMetadataLicenseValue:         "MIT",
	DefaultMetadataTemplateVersionValue: "0.0.0",

	CloneyExampleRepositoryURL: "https://github.com/ArthurSudbrackIbarra/cloney-example.git",
	CloneyDocumentationURL:     "https://arthursudbrackibarra.github.io/cloney-documentation",
}

// GetAppConfig returns a copy of the global application configuration.
func GetAppConfig() AppConfig {
	return *globalConfig
}
