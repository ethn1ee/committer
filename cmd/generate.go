/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ethn1ee/committer/internal/committer"
	"github.com/ethn1ee/committer/internal/config"
	"github.com/ethn1ee/committer/internal/utils"
	"github.com/spf13/cobra"
)

var (
	commit bool
	push   bool
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate a commit message based on git diffs",
	Long:    `Generate a commit message based on git diffs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		cfg, err := config.Init()
		if err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}

		msg, err := committer.Generate(cfg, ctx)
		if err != nil {
			return fmt.Errorf("failed to generate commit message: %w", err)
		}

		if commit || push {
			hash, err := utils.Commit(cfg.WorkTree, msg)
			if err != nil {
				return fmt.Errorf("failed to commit changes: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Committed successfully: %s\n", hash)
		}

		if push {
			err := utils.Push(cfg.Remotes)
			if err != nil {
				return fmt.Errorf("failed to push changes: %w", err)
			}
			fmt.Fprintln(os.Stdout, "Pushed successfully")
		}

		return nil
	},
}

func init() {
	generateCmd.SilenceUsage = true
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	generateCmd.Flags().BoolVarP(&commit, "commit", "c", false, "commit with the generated message, without pushing")
	generateCmd.Flags().BoolVarP(&push, "push", "p", false, "commit and push with the generated message")
}
