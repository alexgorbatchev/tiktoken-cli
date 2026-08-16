package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintCatalog(t *testing.T) {
	var buf bytes.Buffer
	printCatalog(&buf, defaultCatalog)

	output := buf.String()

	expectedSubstrings := []string{
		"Available Models and Encodings:",
		"Encoding: o200k_base",
		"- gpt-4o",
		"Encoding: cl100k_base",
		"- gpt-4",
		"Encoding: p50k_base",
		"Encoding: r50k_base (gpt2)",
		"Available Encodings:",
		"o200k_base",
		"cl100k_base",
	}

	for _, substr := range expectedSubstrings {
		if !strings.Contains(output, substr) {
			t.Errorf("expected output to contain %q, but got:\n%s", substr, output)
		}
	}
}

func TestModelsCmdRun(t *testing.T) {
	var buf bytes.Buffer
	modelsCmd.SetOut(&buf)

	modelsCmd.Run(modelsCmd, []string{})

	output := buf.String()
	if !strings.Contains(output, "Available Models and Encodings:") {
		t.Errorf("modelsCmd output missing expected header, got:\n%s", output)
	}
}
