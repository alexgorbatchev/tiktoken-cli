package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunEncode(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		model     string
		encoding  string
		wantOut   string
		wantErr   bool
		errSubstr string
	}{
		{
			name:     "encode argument cl100k_base",
			args:     []string{"Hello world"},
			encoding: "cl100k_base",
			wantOut:  "9906 1917\n",
			wantErr:  false,
		},
		{
			name:     "encode argument gpt-4o",
			args:     []string{"Hello world"},
			model:    "gpt-4o",
			wantOut:  "13225 2375\n",
			wantErr:  false,
		},
		{
			name:      "encode invalid model",
			args:      []string{"Hello world"},
			model:     "invalid-model",
			wantOut:   "",
			wantErr:   true,
			errSubstr: "failed to get encoding for model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encodeModel = tt.model
			encodeEncoding = tt.encoding
			defer func() {
				encodeModel = ""
				encodeEncoding = "cl100k_base"
			}()

			var outBuf bytes.Buffer
			cmd := encodeCmd
			cmd.SetOut(&outBuf)

			err := runEncode(cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("runEncode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if err == nil || (tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr)) {
					t.Errorf("expected error containing %q, got %v", tt.errSubstr, err)
				}
			} else {
				if outBuf.String() != tt.wantOut {
					t.Errorf("got output %q, want %q", outBuf.String(), tt.wantOut)
				}
			}
		})
	}
}
