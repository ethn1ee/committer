package committer

import (
	"context"
	"fmt"

	"github.com/ethn1ee/committer/internal/config"
	"github.com/ethn1ee/committer/internal/llm"
	"github.com/ethn1ee/committer/internal/models"
	"github.com/ethn1ee/committer/internal/utils"
)

func Generate(cfg *config.Config, ctx context.Context) (string, error) {
	status, err := utils.GetStatus(cfg.WorkTree)
	if err != nil {
		return "", fmt.Errorf("failed to get git status: %w", err)
	}

	if status.IsClean() {
		return "", fmt.Errorf("no changes detected")
	}

	diffs, err := utils.GetDiffs(cfg.HeadTree, cfg.WorkTree)
	if err != nil {
		return "", fmt.Errorf("failed to get git diffs: %w", err)
	}

	prompt := &models.Prompt{
		Instruction: "Create a concise git commit message based on the provided status and diffs. Strictly follow the rules provided.",
		Status:      status.String(),
		Diffs:       diffs,
		Rules: []string{
			"The message contains two parts: header (first line of the commit) and body (optional, after the header)",
			"The header should be in <type>: <description> format",
			"The possible <type> includes: feat, chore, enhancement, fix, docs",
			"The <description> should be a short summary of the changes made",
			"The <description> MUST be all lowercase, except names and proper nouns",
			"The body should provide additional context and details about the changes made in a bullet list",
			"The body should be cased properly, with each sentence starting with a capital letter",
			"Do not surround message with any quotation or markdown syntax",
		},
	}

	msg, err := llm.Ask(cfg, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to get response from Gemini: %w", err)
	}

	return msg, nil
}
