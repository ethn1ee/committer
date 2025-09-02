package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

var CfgFile string

type Config struct {
	GeminiApiKey string `mapstructure:"geminiApiKey"`
}

func Init() (*Config, error) {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home dir: %w", err)
		}

		configPath := path.Join(home, ".config", "committer") // ~/.config/committer/config.yaml

		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	loadEnv()

	if err := viper.ReadInConfig(); err == nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", viper.ConfigFileUsed(), err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func loadEnv() {
	viper.BindEnv("geminiApiKey", "GEMINI_API_KEY")
}
