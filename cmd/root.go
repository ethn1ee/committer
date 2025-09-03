/*
Copyright © 2025 Ethan Lee <ethantlee21@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/ethn1ee/committer/internal/config"
	"github.com/spf13/cobra"
)

var version = "0.1.14"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "committer",
	Version: version,
	Short:   "committer - A simple CLI to generate git commit messages with AI",
	Long:    `committer - A simple CLI to generate git commit messages with AI`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&config.CfgFile, "config", "", "config file (default is $HOME/.config/committer/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
