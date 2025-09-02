package committer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethn1ee/committer/internal/config"
	"github.com/ethn1ee/committer/internal/utils"
	"github.com/go-git/go-git/v6"
	"google.golang.org/genai"
)

type Diff struct {
	Path       string
	StatusCode string
	Before     string
	After      string
}

type Prompt struct {
	Instruction string
	Diffs       []*Diff
	Rules       []string
}

func Generate(cfg *config.Config) (string, error) {
	ctx := context.Background()

	worktree, headTree, err := utils.GetTrees()

	_, err = worktree.Add(".")
	if err != nil {
		return "", fmt.Errorf("failed to add changes to staging: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree status: %w", err)
	}

	if status.IsClean() {
		return "No changes to commit", nil
	}

	diffs := make([]*Diff, 0, len(status))

	for path, fs := range status {
		d := &Diff{
			Path:       path,
			StatusCode: string(fs.Staging),
		}

		switch fs.Staging {
		case git.Modified:
			before, err := utils.GetBefore(headTree, path)
			if err != nil {
				return "", fmt.Errorf("failed to get before contents for file %s: %w", path, err)
			}

			after, err := utils.GetAfter(worktree, path)
			if err != nil {
				return "", fmt.Errorf("failed to get after contents for file %s: %w", path, err)
			}

			d.Before = before
			d.After = after

		case git.Added:
			after, err := utils.GetAfter(worktree, path)
			if err != nil {
				return "", fmt.Errorf("failed to get after contents for file %s: %w", path, err)
			}

			d.Before = ""
			d.After = after

		case git.Deleted:
			before, err := utils.GetBefore(headTree, path)
			if err != nil {
				return "", fmt.Errorf("failed to get before contents for file %s: %w", path, err)
			}

			d.Before = before
			d.After = ""

		case git.Renamed:
			// fmt.Fprintf(os.Stdout, "Renamed: %s\n", path)

		case git.Untracked:
			// fmt.Fprintf(os.Stdout, "Untracked: %s\n", path)
		default:
			// fmt.Fprintf(os.Stdout, "Other change (%s): %s\n", string(fs.Worktree), path)
		}

		diffs = append(diffs, d)
	}

	prompt := &Prompt{
		Instruction: "Create a concise git commit message based on the following information",
		Diffs:       diffs,
		Rules: []string{
			"The message contains two parts: header (first line of the commit) and body (optional, after the header)",
			"The header  should be in `<type>: <description>` format",
			"The possible <type> includes: feat, chore, enhancement, fix, docs",
			"The <description> should be a short summary of the changes made",
			"Do not capitalize the first letter of the <description>",
			"The body should provide additional context and details about the changes made in a bullet list",
		},
	}

	promptStr, err := json.MarshalIndent(prompt, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal prompt: %w", err)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.GeminiApiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %w", err)
	}

	chat, err := client.Chats.Create(ctx, "gemini-2.0-flash", nil, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini chat: %w", err)
	}

	res, err := chat.Send(ctx, &genai.Part{
		Text: string(promptStr),
	})
	if err != nil {
		return "", fmt.Errorf("failed to send message to Gemini: %w", err)
	}

	return res.Text(), nil
}
