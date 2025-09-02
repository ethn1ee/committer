/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/ethn1ee/committer/internal/committer"
	"github.com/ethn1ee/committer/internal/config"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate a commit message based on git diffs",
	Long:    `Generate a commit message based on git diffs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Init()
		if err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}

		msg, err := committer.Generate(cfg)
		if err != nil {
			return fmt.Errorf("failed to generate commit message: %w", err)
		}
		fmt.Println(msg)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
