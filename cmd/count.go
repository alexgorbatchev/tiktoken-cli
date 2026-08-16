package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	countModel    string
	countEncoding string
)

var countCmd = &cobra.Command{
	Use:   "count [text]",
	Short: "Count tokens in the given text",
	Long: `Count the number of tokens in the provided text using the specified model or encoding.

If no text is provided as an argument, it reads from stdin.

Available encodings:
  - o200k_base   (gpt-4o, gpt-4.1, gpt-4.5)
  - cl100k_base  (gpt-4, gpt-3.5-turbo, text-embedding-ada-002)
  - p50k_base    (Codex models, text-davinci-002, text-davinci-003)
  - r50k_base    (GPT-3 models like davinci)

Examples:
  # Count tokens from argument
  tiktoken count "Hello, world!"

  # Count tokens from stdin
  echo "Hello, world!" | tiktoken count

  # Count tokens for a specific model
  tiktoken count -m gpt-4o "Hello, world!"

  # Count tokens using a specific encoding
  tiktoken count -e o200k_base "Hello, world!"

  # Count tokens from a file
  cat myfile.txt | tiktoken count`,
	RunE: runCount,
}

func init() {
	rootCmd.AddCommand(countCmd)
	countCmd.Flags().StringVarP(&countModel, "model", "m", "", "OpenAI model name (e.g., gpt-4o, gpt-4, gpt-3.5-turbo)")
	countCmd.Flags().StringVarP(&countEncoding, "encoding", "e", "cl100k_base", "Encoding name (o200k_base, cl100k_base, p50k_base, r50k_base)")
}

func runCount(cmd *cobra.Command, args []string) error {
	text, err := getText(args)
	if err != nil {
		return err
	}

	enc, err := resolveEncoding(countModel, countEncoding)
	if err != nil {
		return err
	}

	tokens := enc.Encode(text, nil, nil)
	fmt.Fprintln(cmd.OutOrStdout(), len(tokens))

	return nil
}

func getText(args []string) (string, error) {
	stat, _ := os.Stdin.Stat()
	var isTerminal bool
	if stat != nil {
		isTerminal = (stat.Mode() & os.ModeCharDevice) != 0
	}
	return getTextFromReader(args, os.Stdin, isTerminal)
}

func getTextFromReader(args []string, r io.Reader, isTerminal bool) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	if isTerminal {
		return "", fmt.Errorf("no text provided. Either pass text as argument or pipe through stdin")
	}

	reader := bufio.NewReader(r)
	var builder strings.Builder

	for {
		line, err := reader.ReadString('\n')
		builder.WriteString(line)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error reading stdin: %w", err)
		}
	}

	return strings.TrimSuffix(builder.String(), "\n"), nil
}
