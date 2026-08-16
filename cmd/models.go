package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// ModelGroup associates an encoding name with supported OpenAI models.
type ModelGroup struct {
	Encoding string
	Models   []string
}

// EncodingSummary describes an encoding and its typical usage.
type EncodingSummary struct {
	Name        string
	Description string
}

// ModelCatalog holds structured information about models and encodings.
type ModelCatalog struct {
	Groups    []ModelGroup
	Encodings []EncodingSummary
}

var defaultCatalog = ModelCatalog{
	Groups: []ModelGroup{
		{
			Encoding: "o200k_base",
			Models:   []string{"gpt-4o", "gpt-4.1", "gpt-4.5"},
		},
		{
			Encoding: "cl100k_base",
			Models: []string{
				"gpt-4",
				"gpt-3.5-turbo",
				"text-embedding-ada-002",
				"text-embedding-3-small",
				"text-embedding-3-large",
			},
		},
		{
			Encoding: "p50k_base",
			Models: []string{
				"text-davinci-002",
				"text-davinci-003",
				"code-davinci-002",
				"code-cushman-001",
			},
		},
		{
			Encoding: "r50k_base (gpt2)",
			Models:   []string{"davinci", "curie", "babbage", "ada"},
		},
	},
	Encodings: []EncodingSummary{
		{Name: "o200k_base", Description: "newest, used by GPT-4o models"},
		{Name: "cl100k_base", Description: "used by GPT-4 and GPT-3.5-turbo"},
		{Name: "p50k_base", Description: "used by Codex models"},
		{Name: "p50k_edit", Description: "used by edit models"},
		{Name: "r50k_base", Description: "used by GPT-3 models"},
	},
}

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List available models and their encodings",
	Long:  `Display a list of available OpenAI models and their corresponding tokenization encodings.`,
	Run: func(cmd *cobra.Command, args []string) {
		printCatalog(cmd.OutOrStdout(), defaultCatalog)
	},
}

func init() {
	rootCmd.AddCommand(modelsCmd)
}

func printCatalog(out io.Writer, catalog ModelCatalog) {
	fmt.Fprintln(out, "Available Models and Encodings:")
	fmt.Fprintln(out)

	for _, group := range catalog.Groups {
		fmt.Fprintf(out, "Encoding: %s\n", group.Encoding)
		for _, model := range group.Models {
			fmt.Fprintf(out, "  - %s\n", model)
		}
		fmt.Fprintln(out)
	}

	fmt.Fprintln(out, "Available Encodings:")
	for _, enc := range catalog.Encodings {
		fmt.Fprintf(out, "  - %-12s (%s)\n", enc.Name, enc.Description)
	}
}
