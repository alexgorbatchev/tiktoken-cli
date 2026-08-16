package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read error simulated")
}

func TestGetTextWithFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/test.txt"
	if err := os.WriteFile(tmpFile, []byte("Hello file world\n"), 0644); err != nil {
		t.Fatalf("failed to write tmp file: %v", err)
	}

	// Test positional file path
	got, err := getText([]string{tmpFile}, "")
	if err != nil {
		t.Fatalf("getText() error = %v", err)
	}
	if got != "Hello file world" {
		t.Errorf("getText() = %q, want %q", got, "Hello file world")
	}

	// Test explicit file flag
	got, err = getText([]string{}, tmpFile)
	if err != nil {
		t.Fatalf("getText() error = %v", err)
	}
	if got != "Hello file world" {
		t.Errorf("getText() = %q, want %q", got, "Hello file world")
	}

	// Test non-existent file path flag error
	_, err = getText([]string{}, tmpDir+"/nonexistent.txt")
	if err == nil {
		t.Errorf("expected error for non-existent file, got nil")
	}
}

func TestGetTextFromReader(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		input      string
		isTerminal bool
		want       string
		wantErr    bool
		useErrRdr  bool
	}{
		{
			name:       "from args",
			args:       []string{"hello", "world"},
			input:      "",
			isTerminal: false,
			want:       "hello world",
			wantErr:    false,
		},
		{
			name:       "from stdin piped",
			args:       []string{},
			input:      "hello from stdin\n",
			isTerminal: false,
			want:       "hello from stdin",
			wantErr:    false,
		},
		{
			name:       "terminal without args error",
			args:       []string{},
			input:      "",
			isTerminal: true,
			want:       "",
			wantErr:    true,
		},
		{
			name:       "reader error",
			args:       []string{},
			input:      "",
			isTerminal: false,
			want:       "",
			wantErr:    true,
			useErrRdr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useErrRdr {
				_, err := getTextFromReader(tt.args, &errorReader{}, tt.isTerminal)
				if (err != nil) != tt.wantErr {
					t.Fatalf("getTextFromReader() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			r := strings.NewReader(tt.input)
			got, err := getTextFromReader(tt.args, r, tt.isTerminal)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getTextFromReader() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("getTextFromReader() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRunCount(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/count.txt"
	if err := os.WriteFile(tmpFile, []byte("Hello, world!\n"), 0644); err != nil {
		t.Fatalf("failed to write tmp file: %v", err)
	}

	tests := []struct {
		name      string
		args      []string
		model     string
		encoding  string
		file      string
		wantOut   string
		wantErr   bool
		errSubstr string
	}{
		{
			name:     "valid args count",
			args:     []string{"Hello, world!"},
			encoding: "cl100k_base",
			wantOut:  "4\n",
			wantErr:  false,
		},
		{
			name:     "valid file count",
			file:     tmpFile,
			encoding: "cl100k_base",
			wantOut:  "4\n",
			wantErr:  false,
		},
		{
			name:     "valid model count",
			args:     []string{"Hello, world!"},
			model:    "gpt-4o",
			wantOut:  "4\n",
			wantErr:  false,
		},
		{
			name:      "invalid model count",
			args:      []string{"Hello, world!"},
			model:     "bad-model",
			wantOut:   "",
			wantErr:   true,
			errSubstr: "failed to get encoding for model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			countModel = tt.model
			countEncoding = tt.encoding
			countFile = tt.file
			defer func() {
				countModel = ""
				countEncoding = "cl100k_base"
				countFile = ""
			}()

			var outBuf bytes.Buffer
			cmd := countCmd
			cmd.SetOut(&outBuf)

			err := runCount(cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("runCount() error = %v, wantErr %v", err, tt.wantErr)
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
