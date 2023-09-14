package config

import (
	"github.com/spf13/viper"
)

// AppConfig represents the application configuration.
type AppConfig struct {
	// Version is the current version of the application.
	Version string

	// MetadataFileName is the name of the metadata file in the repository.
	// This file contains the information about the template repository.
	MetadataFileName string

	// DefaultUserVariablesFileName is the default name of the file containing the user variables.
	// This file is used to fill the template variables in the cloned directory.'
	DefaultUserVariablesFileName string

	// DefaultDryRunDirectoryName is the default name of the directory created when running a template in dryrun mode.
	DefaultDryRunDirectoryName string

	// GitToken is the token used to authenticate when dealing with private git repositories.
	GitToken string `mapstructure:"GIT_TOKEN"`
}

// globalConfig is the global application configuration.
var globalConfig = &AppConfig{
	// Default values.
	Version:                      "v0.0.0",
	MetadataFileName:             ".cloney.yaml",
	DefaultUserVariablesFileName: ".cloney-vars.yaml",
	DefaultDryRunDirectoryName:   "cloney-dry-run-results",
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
		Version:                      globalConfig.Version,
		MetadataFileName:             globalConfig.MetadataFileName,
		DefaultUserVariablesFileName: globalConfig.DefaultUserVariablesFileName,
		DefaultDryRunDirectoryName:   globalConfig.DefaultDryRunDirectoryName,
		GitToken:                     globalConfig.GitToken,
	}
}
