package config

import (
	"fmt"
	"os"
	"path"

	"github.com/ethn1ee/committer/internal/utils"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/spf13/viper"
)

var CfgFile string

type Config struct {
	Remotes  []*git.Remote `mapstructure:"remotes"`
	HeadTree *object.Tree  `mapstructure:"headTree"`
	WorkTree *git.Worktree `mapstructure:"workTree"`

	LLM          string `mapstructure:"llm"`
	GeminiApiKey string `mapstructure:"geminiApiKey"`
}

const (
	LLM_GEMINI = "gemini"
)

func Init() (*Config, error) {
	if err := loadGit(); err != nil {
		return nil, fmt.Errorf("failed to load git info: %w", err)
	}
	loadEnv()
	setDefaults()

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

func setDefaults() {
	viper.SetDefault("LLM", LLM_GEMINI)
}

func loadGit() error {
	remotes, workTree, headTree, err := utils.GetTrees()
	if err != nil {
		return fmt.Errorf("failed to get git trees: %w", err)
	}

	viper.Set("remotes", remotes)
	viper.Set("headTree", headTree)
	viper.Set("workTree", workTree)

	return nil
}
