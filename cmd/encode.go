package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	encodeModel    string
	encodeEncoding string
)

var encodeCmd = &cobra.Command{
	Use:   "encode [text]",
	Short: "Encode text to token IDs",
	Long: `Encode the provided text into token IDs using the specified model or encoding.

If no text is provided as an argument, it reads from stdin.

Examples:
  # Encode text
  tiktoken encode "Hello, world!"

  # Encode using a specific model
  tiktoken encode -m gpt-4o "Hello, world!"

  # Encode using a specific encoding
  tiktoken encode -e o200k_base "Hello, world!"

  # Encode from stdin
  echo "Hello, world!" | tiktoken encode`,
	RunE: runEncode,
}

func init() {
	rootCmd.AddCommand(encodeCmd)
	encodeCmd.Flags().StringVarP(&encodeModel, "model", "m", "", "OpenAI model name (e.g., gpt-4o, gpt-4, gpt-3.5-turbo)")
	encodeCmd.Flags().StringVarP(&encodeEncoding, "encoding", "e", "cl100k_base", "Encoding name (o200k_base, cl100k_base, p50k_base, r50k_base)")
}

func runEncode(cmd *cobra.Command, args []string) error {
	text, err := getText(args)
	if err != nil {
		return err
	}

	enc, err := resolveEncoding(encodeModel, encodeEncoding)
	if err != nil {
		return err
	}

	tokens := enc.Encode(text, nil, nil)

	// Print tokens as space-separated values
	tokenStrs := make([]string, len(tokens))
	for i, t := range tokens {
		tokenStrs[i] = fmt.Sprintf("%d", t)
	}
	fmt.Fprintln(cmd.OutOrStdout(), strings.Join(tokenStrs, " "))

	return nil
}
