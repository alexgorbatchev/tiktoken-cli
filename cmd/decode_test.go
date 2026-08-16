package cmd

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGetTokensWithFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/tokens.txt"
	if err := os.WriteFile(tmpFile, []byte("15339 1917 0\n"), 0644); err != nil {
		t.Fatalf("failed to write tmp file: %v", err)
	}

	// Test positional file path
	got, err := getTokens([]string{tmpFile}, "")
	if err != nil {
		t.Fatalf("getTokens() error = %v", err)
	}
	want := []int{15339, 1917, 0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("getTokens() = %v, want %v", got, want)
	}

	// Test explicit file flag
	got, err = getTokens([]string{}, tmpFile)
	if err != nil {
		t.Fatalf("getTokens() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("getTokens() = %v, want %v", got, want)
	}

	// Test non-existent file error
	_, err = getTokens([]string{}, tmpDir+"/nonexistent.txt")
	if err == nil {
		t.Errorf("expected error for non-existent file, got nil")
	}
}

func TestGetTokensFromReader(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		input      string
		isTerminal bool
		want       []int
		wantErr    bool
		errSubstr  string
		useErrRdr  bool
	}{
		{
			name:       "from arguments",
			args:       []string{"15339", "1917", "0"},
			input:      "",
			isTerminal: false,
			want:       []int{15339, 1917, 0},
			wantErr:    false,
		},
		{
			name:       "from stdin whitespace separated",
			args:       []string{},
			input:      "15339 1917\n0\n",
			isTerminal: false,
			want:       []int{15339, 1917, 0},
			wantErr:    false,
		},
		{
			name:       "terminal without args error",
			args:       []string{},
			input:      "",
			isTerminal: true,
			want:       nil,
			wantErr:    true,
			errSubstr:  "no token IDs provided",
		},
		{
			name:       "invalid token integer arg",
			args:       []string{"15339", "invalid", "0"},
			input:      "",
			isTerminal: false,
			want:       nil,
			wantErr:    true,
			errSubstr:  "invalid token ID \"invalid\"",
		},
		{
			name:       "reader error",
			args:       []string{},
			input:      "",
			isTerminal: false,
			want:       nil,
			wantErr:    true,
			useErrRdr:  true,
		},
		{
			name:       "empty input result error",
			args:       []string{},
			input:      "   \n  \t ",
			isTerminal: false,
			want:       nil,
			wantErr:    true,
			errSubstr:  "no valid token IDs provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.useErrRdr {
				_, err := getTokensFromReader(tt.args, &errorReader{}, tt.isTerminal)
				if (err != nil) != tt.wantErr {
					t.Fatalf("getTokensFromReader() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			r := strings.NewReader(tt.input)
			got, err := getTokensFromReader(tt.args, r, tt.isTerminal)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getTokensFromReader() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if err == nil || (tt.errSubstr != "" && !strings.Contains(err.Error(), tt.errSubstr)) {
					t.Errorf("expected error containing %q, got %v", tt.errSubstr, err)
				}
			} else {
				if len(got) != len(tt.want) {
					t.Fatalf("got len %d, want len %d", len(got), len(tt.want))
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("got[%d] = %d, want %d", i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestRunDecode(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/decode_tokens.txt"
	if err := os.WriteFile(tmpFile, []byte("15339 1917 0\n"), 0644); err != nil {
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
			name:     "decode arguments cl100k_base",
			args:     []string{"15339", "1917", "0"},
			encoding: "cl100k_base",
			wantOut:  "hello world!\n",
			wantErr:  false,
		},
		{
			name:     "decode file argument cl100k_base",
			args:     []string{tmpFile},
			encoding: "cl100k_base",
			wantOut:  "hello world!\n",
			wantErr:  false,
		},
		{
			name:     "decode file flag cl100k_base",
			file:     tmpFile,
			encoding: "cl100k_base",
			wantOut:  "hello world!\n",
			wantErr:  false,
		},
		{
			name:     "decode arguments gpt-4o",
			args:     []string{"13225", "2375"},
			model:    "gpt-4o",
			wantOut:  "Hello world\n",
			wantErr:  false,
		},
		{
			name:      "decode invalid model",
			args:      []string{"15339"},
			model:     "invalid-model",
			wantOut:   "",
			wantErr:   true,
			errSubstr: "failed to get encoding for model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decodeModel = tt.model
			decodeEncoding = tt.encoding
			decodeFile = tt.file
			defer func() {
				decodeModel = ""
				decodeEncoding = "cl100k_base"
				decodeFile = ""
			}()

			var outBuf bytes.Buffer
			cmd := decodeCmd
			cmd.SetOut(&outBuf)

			err := runDecode(cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("runDecode() error = %v, wantErr %v", err, tt.wantErr)
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
