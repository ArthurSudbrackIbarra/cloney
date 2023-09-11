package config

import (
	"github.com/spf13/viper"
)

// AppConfig represents the application configuration.
type AppConfig struct {
	// MetadataFileName is the name of the metadata file in the repository.
	// This file contains the information about the template.
	// Default is ".cloney.yaml".
	MetadataFileName string

	// GitToken is the token used to authenticate when dealing with private git repositories.
	GitToken string `mapstructure:"git_token"`
}

// globalConfig is the global application configuration.
var globalConfig = &AppConfig{
	// Default values.
	MetadataFileName: ".cloney.yaml",
}

// LoadConfig loads the global application configuration.
func LoadConfig() error {
	// Enable reading environment variables.
	viper.SetEnvPrefix("CLONEY")
	viper.AutomaticEnv()

	// Unmarshal the configuration into the globalConfig variable
	err := viper.Unmarshal(globalConfig)
	if err != nil {
		return err
	}

	return nil
}

// GetAppConfig returns a copy of the global application configuration.
func GetAppConfig() *AppConfig {
	return &AppConfig{
		MetadataFileName: globalConfig.MetadataFileName,
		GitToken:         globalConfig.GitToken,
	}
}
