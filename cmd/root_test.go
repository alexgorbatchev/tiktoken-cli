package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestVersionCmd(t *testing.T) {
	var buf bytes.Buffer
	versionCmd.SetOut(&buf)

	versionCmd.Run(versionCmd, []string{})

	output := buf.String()
	if !strings.Contains(output, "tiktoken version") {
		t.Errorf("versionCmd output missing expected header, got:\n%s", output)
	}
}

func TestRootCmdExecution(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/root.txt"
	if err := os.WriteFile(tmpFile, []byte("hello world\n"), 0644); err != nil {
		t.Fatalf("failed to write tmp file: %v", err)
	}

	tests := []struct {
		name      string
		args      []string
		model     string
		encoding  string
		file      string
		wantOut   string
		contains  string
		wantErr   bool
		errSubstr string
	}{
		{
			name:     "root count from arg",
			args:     []string{"hello world"},
			encoding: "cl100k_base",
			wantOut:  "2\n",
			wantErr:  false,
		},
		{
			name:     "root count from file arg",
			args:     []string{tmpFile},
			encoding: "cl100k_base",
			wantOut:  "2\n",
			wantErr:  false,
		},
		{
			name:     "root count from file flag",
			file:     tmpFile,
			encoding: "cl100k_base",
			wantOut:  "2\n",
			wantErr:  false,
		},
		{
			name:     "root count zero args help",
			args:     []string{},
			contains: "Usage:",
			wantErr:  false,
		},
		{
			name:     "root count with model",
			args:     []string{"hello world"},
			model:    "gpt-4o",
			wantOut:  "2\n",
			wantErr:  false,
		},
		{
			name:      "root count invalid model",
			args:      []string{"hello world"},
			model:     "invalid-model",
			wantOut:   "",
			wantErr:   true,
			errSubstr: "failed to get encoding for model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootModel = tt.model
			rootEncoding = tt.encoding
			rootFile = tt.file
			defer func() {
				rootModel = ""
				rootEncoding = "cl100k_base"
				rootFile = ""
			}()

			var outBuf bytes.Buffer
			cmd := rootCmd
			cmd.SetOut(&outBuf)

			err := runRootCount(cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("runRootCount() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if err == nil || (tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr)) {
					t.Errorf("expected error containing %q, got %v", tt.errSubstr, err)
				}
			} else {
				if tt.contains != "" {
					if !strings.Contains(outBuf.String(), tt.contains) {
						t.Errorf("got output %q, expected to contain %q", outBuf.String(), tt.contains)
					}
				} else if outBuf.String() != tt.wantOut {
					t.Errorf("got output %q, want %q", outBuf.String(), tt.wantOut)
				}
			}
		})
	}
}

func TestRootCmdHelpText(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	_ = rootCmd.Help()

	output := buf.String()
	if strings.Contains(output, "tiktoken-go") {
		t.Errorf("rootCmd help description contains unwanted library implementation details ('tiktoken-go')")
	}
	if !strings.Contains(output, "fast command-line tool for tokenizing text") {
		t.Errorf("rootCmd help description missing expected long description, got:\n%s", output)
	}
	if rootCmd.Short != "Count, encode, and decode tokens for OpenAI models" {
		t.Errorf("rootCmd.Short = %q, want 'Count, encode, and decode tokens for OpenAI models'", rootCmd.Short)
	}
}

func TestExecuteRootCmd(t *testing.T) {
	rootCmd.SetArgs([]string{"version"})
	var buf bytes.Buffer
	versionCmd.SetOut(&buf)
	rootCmd.SetOut(&buf)
	Execute()
	if !strings.Contains(buf.String(), "tiktoken version") {
		t.Errorf("Execute() with version subcommand failed, output:\n%s", buf.String())
	}
}
