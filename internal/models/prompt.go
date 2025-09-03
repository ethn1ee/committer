package models

import (
	"encoding/json"
	"fmt"
)

type Prompt struct {
	Instruction string
	Status      string
	Diffs       []*Diff
	Rules       []string
}

func (p *Prompt) String() (string, error) {
	str, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal prompt: %w", err)
	}

	return string(str), nil
}
