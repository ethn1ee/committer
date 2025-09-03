/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethn1ee/committer/internal/committer"
	"github.com/ethn1ee/committer/internal/config"
	"github.com/ethn1ee/committer/internal/utils"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
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
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Start()

		ctx := context.Background()

		cfg, err := config.Init()
		if err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}

		msg, err := committer.Generate(cfg, ctx)
		if err != nil {
			return fmt.Errorf("failed to generate commit message: %w", err)
		}

		s.Stop()

		headerColor := color.New(color.FgCyan, color.Bold)
		bodyColor := color.New(color.FgHiBlack)

		parts := strings.Split(msg, "\n")
		for i, part := range parts {
			if i == 0 {
				parts[i] = headerColor.Sprint(part)
			} else {
				parts[i] = bodyColor.Sprint(part)
			}
		}

		fmt.Fprintln(os.Stdout, "\n"+strings.Join(parts, "\n")+"\n")

		prompt := promptui.Prompt{
			Label:     "Accept and commit",
			IsConfirm: true,
			Templates: &promptui.PromptTemplates{
				Prompt: "{{ . }} | cyan",
			},
		}

		if commit || push {
			_, err := prompt.Run()
			if err != nil {
				fmt.Fprintln(os.Stdout, "Rejected")
				return nil
			}

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
