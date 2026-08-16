package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read error simulated")
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
			var rdr = strings.NewReader(tt.input)
			var errRdr *errorReader
			var r = rdr
			if tt.useErrRdr {
				errRdr = &errorReader{}
				got, err := getTextFromReader(tt.args, errRdr, tt.isTerminal)
				if (err != nil) != tt.wantErr {
					t.Fatalf("getTextFromReader() error = %v, wantErr %v", err, tt.wantErr)
				}
				if got != tt.want {
					t.Errorf("got %q, want %q", got, tt.want)
				}
				return
			}

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
			name:     "valid args count",
			args:     []string{"Hello, world!"},
			encoding: "cl100k_base",
			wantOut:  "3\n",
			wantErr:  false,
		},
		{
			name:     "valid model count",
			args:     []string{"Hello, world!"},
			model:    "gpt-4o",
			wantOut:  "3\n",
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
			defer func() {
				countModel = ""
				countEncoding = "cl100k_base"
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
