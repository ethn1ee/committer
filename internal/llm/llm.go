package llm

import (
	"fmt"

	"github.com/ethn1ee/committer/internal/config"
	"github.com/ethn1ee/committer/internal/models"
)

func Ask(cfg *config.Config, prompt *models.Prompt) (string, error) {
	switch cfg.LLM {
	case config.LLM_GEMINI:
		return AskGemini(cfg, prompt)
	default:
		return "", fmt.Errorf("unknown llm type: %s", cfg.LLM)
	}
}
