package cmd

import (
	"strings"
	"testing"
)

func TestResolveEncoding(t *testing.T) {
	tests := []struct {
		name         string
		model        string
		encoding     string
		wantErr      bool
		errSubstring string
		wantTokens   int
	}{
		{
			name:       "valid model gpt-4o",
			model:      "gpt-4o",
			encoding:   "",
			wantErr:    false,
			wantTokens: 1,
		},
		{
			name:       "valid model gpt-4",
			model:      "gpt-4",
			encoding:   "",
			wantErr:    false,
			wantTokens: 1,
		},
		{
			name:         "invalid model",
			model:        "invalid-model-xyz",
			encoding:     "",
			wantErr:      true,
			errSubstring: "failed to get encoding for model invalid-model-xyz",
		},
		{
			name:       "valid encoding o200k_base",
			model:      "",
			encoding:   "o200k_base",
			wantErr:    false,
			wantTokens: 1,
		},
		{
			name:       "valid encoding cl100k_base",
			model:      "",
			encoding:   "cl100k_base",
			wantErr:    false,
			wantTokens: 1,
		},
		{
			name:       "valid encoding p50k_base",
			model:      "",
			encoding:   "p50k_base",
			wantErr:    false,
			wantTokens: 1,
		},
		{
			name:       "valid encoding r50k_base",
			model:      "",
			encoding:   "r50k_base",
			wantErr:    false,
			wantTokens: 1,
		},
		{
			name:         "invalid encoding name",
			model:        "",
			encoding:     "invalid_encoding_xyz",
			wantErr:      true,
			errSubstring: "failed to get encoding invalid_encoding_xyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := resolveEncoding(tt.model, tt.encoding)
			if (err != nil) != tt.wantErr {
				t.Fatalf("resolveEncoding(%q, %q) error = %v, wantErr %v", tt.model, tt.encoding, err, tt.wantErr)
			}
			if tt.wantErr {
				if err == nil || (tt.errSubstring != "" && !strings.Contains(err.Error(), tt.errSubstring)) {
					t.Errorf("expected error containing %q, got %v", tt.errSubstring, err)
				}
			} else {
				if enc == nil {
					t.Fatalf("expected non-nil tiktoken.Tiktoken instance")
				}
				tokens := enc.Encode("hello", nil, nil)
				if len(tokens) != tt.wantTokens {
					t.Errorf("len(enc.Encode(\"hello\")) = %d, want %d", len(tokens), tt.wantTokens)
				}
			}
		})
	}
}
