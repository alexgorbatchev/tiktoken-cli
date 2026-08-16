package cmd

import (
	"bytes"
	"testing"
)

func TestGetText(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{
			name:    "from arguments",
			args:    []string{"hello", "world"},
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "single argument",
			args:    []string{"hello world"},
			want:    "hello world",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getText(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getText(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("getText(%v) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func TestRunRootCount(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		model     string
		encoding  string
		wantErr   bool
		errSubstr string
	}{
		{
			name:     "valid argument default encoding",
			args:     []string{"hello world"},
			encoding: "cl100k_base",
			wantErr:  false,
		},
		{
			name:     "valid model",
			args:     []string{"hello world"},
			model:    "gpt-4o",
			wantErr:  false,
		},
		{
			name:      "invalid model",
			args:      []string{"hello world"},
			model:     "nonexistent-model",
			wantErr:   true,
			errSubstr: "failed to get encoding for model",
		},
		{
			name:      "invalid encoding",
			args:      []string{"hello world"},
			encoding:  "invalid-encoding",
			wantErr:   true,
			errSubstr: "failed to get encoding",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootModel = tt.model
			rootEncoding = tt.encoding
			defer func() {
				rootModel = ""
				rootEncoding = "cl100k_base"
			}()

			err := runRootCount(rootCmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("runRootCount() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && tt.errSubstr != "" {
				if !bytes.Contains([]byte(err.Error()), []byte(tt.errSubstr)) {
					t.Errorf("expected error containing %q, got %q", tt.errSubstr, err.Error())
				}
			}
		})
	}
}
