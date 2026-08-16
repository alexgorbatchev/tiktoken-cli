package cmd

import (
	"fmt"

	"github.com/pkoukk/tiktoken-go"
)

// resolveEncoding resolves a tiktoken encoding from either an explicit model name
// or an encoding name.
func resolveEncoding(model, encoding string) (*tiktoken.Tiktoken, error) {
	if model != "" {
		enc, err := tiktoken.EncodingForModel(model)
		if err != nil {
			return nil, fmt.Errorf("failed to get encoding for model %s: %w", model, err)
		}
		return enc, nil
	}
	enc, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return nil, fmt.Errorf("failed to get encoding %s: %w", encoding, err)
	}
	return enc, nil
}
