package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethn1ee/committer/internal/config"
	"github.com/ethn1ee/committer/internal/models"
	"google.golang.org/genai"
)

func AskGemini(cfg *config.Config, prompt *models.Prompt) (string, error) {
	ctx := context.Background()

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

	promptStr, err := prompt.String()

	res, err := chat.Send(ctx, &genai.Part{
		Text: promptStr,
	})
	if err != nil {
		return "", fmt.Errorf("failed to send message to Gemini: %w", err)
	}

	msg := strings.Trim(res.Text(), "\n")

	return msg, nil
}
