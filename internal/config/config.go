package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"ssh-manager/internal/models"
)

var (
	appConfig *models.AppConfig
	once      sync.Once
)

// GetConfig returns a singleton instance of the AppConfig.
// It loads the configuration from the file on its first call.
func GetConfig() (*models.AppConfig, error) {
	var loadErr error
	once.Do(func() {
		appConfig = &models.AppConfig{
			Connections: make(map[string]models.Connection),
			SSHKeys:     make(map[string]models.SSHKey),
		}

		configFile := viper.ConfigFileUsed()
		if configFile == "" {
			// Config file not found, proceed with empty config.
			return
		}

		// Bypass viper.Unmarshal due to issues with time.Time/int64 fields.
		// Read the file manually and unmarshal with yaml.v3.
		bytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			loadErr = fmt.Errorf("could not read config file: %w", err)
			return
		}

		err = yaml.Unmarshal(bytes, &appConfig)
		if err != nil {
			loadErr = fmt.Errorf("unable to decode into struct: %w", err)
			return
		}
	})
	return appConfig, loadErr
}

// SaveConfig saves the current configuration back to the file.
func SaveConfig(config *models.AppConfig) error {
	configFile := viper.ConfigFileUsed()
	if configFile == "" {
		// If no config file is being used, create one in the default location.
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not get home directory: %w", err)
		}
		configDir := fmt.Sprintf("%s/.ssh-manager", home)
		if err := os.MkdirAll(configDir, 0700); err != nil {
			return fmt.Errorf("could not create config directory: %w", err)
		}
		configFile = fmt.Sprintf("%s/config.yaml", configDir)
		viper.SetConfigFile(configFile)
	}

	bytes, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configFile, bytes, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}