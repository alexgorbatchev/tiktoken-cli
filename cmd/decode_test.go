package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestGetTokensFromReader(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		input      string
		isTerminal bool
		want       []int
		wantErr    bool
	}{
		{
			name:       "from args",
			args:       []string{"15339", "1917", "0"},
			input:      "",
			isTerminal: false,
			want:       []int{15339, 1917, 0},
			wantErr:    false,
		},
		{
			name:       "from stdin piped",
			args:       []string{},
			input:      "15339 1917 0\n",
			isTerminal: false,
			want:       []int{15339, 1917, 0},
			wantErr:    false,
		},
		{
			name:       "invalid token integer",
			args:       []string{"abc"},
			input:      "",
			isTerminal: false,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "empty input terminal",
			args:       []string{},
			input:      "",
			isTerminal: true,
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdr := strings.NewReader(tt.input)
			got, err := getTokensFromReader(tt.args, rdr, tt.isTerminal)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getTokensFromReader() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
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
			name:     "decode arguments cl100k_base",
			args:     []string{"9906", "1917"},
			encoding: "cl100k_base",
			wantOut:  "Hello world\n",
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
			args:      []string{"9906", "1917"},
			model:     "bad-model",
			wantOut:   "",
			wantErr:   true,
			errSubstr: "failed to get encoding for model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decodeModel = tt.model
			decodeEncoding = tt.encoding
			defer func() {
				decodeModel = ""
				decodeEncoding = "cl100k_base"
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
