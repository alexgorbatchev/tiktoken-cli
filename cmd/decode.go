package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	decodeModel    string
	decodeEncoding string
	decodeFile     string
)

var decodeCmd = &cobra.Command{
	Use:   "decode [token_ids...|file]",
	Short: "Decode token IDs back to text",
	Long: `Decode token IDs back to text using the specified model or encoding.

Token IDs can be provided as arguments, in a file, or piped through stdin (space or newline separated).

Examples:
  # Decode token IDs from arguments
  tiktoken decode 15339 1917 0

  # Decode token IDs from a file
  tiktoken decode tokens.txt

  # Decode token IDs using file flag
  tiktoken decode -f tokens.txt

  # Decode using a specific model
  tiktoken decode -m gpt-4o 15339 1917 0

  # Decode from stdin
  echo "15339 1917 0" | tiktoken decode

  # Chain encode and decode
  tiktoken encode "Hello, world!" | tiktoken decode`,
	RunE: runDecode,
}

func init() {
	rootCmd.AddCommand(decodeCmd)
	decodeCmd.Flags().StringVarP(&decodeModel, "model", "m", "", "OpenAI model name (e.g., gpt-4o, gpt-4, gpt-3.5-turbo)")
	decodeCmd.Flags().StringVarP(&decodeEncoding, "encoding", "e", "cl100k_base", "Encoding name (o200k_base, cl100k_base, p50k_base, r50k_base)")
	decodeCmd.Flags().StringVarP(&decodeFile, "file", "f", "", "Path to input file")
}

func runDecode(cmd *cobra.Command, args []string) error {
	tokens, err := getTokens(args, decodeFile)
	if err != nil {
		return err
	}

	enc, err := resolveEncoding(decodeModel, decodeEncoding)
	if err != nil {
		return err
	}

	text := enc.Decode(tokens)
	fmt.Fprintln(cmd.OutOrStdout(), text)

	return nil
}

func getTokens(args []string, filePath string) ([]int, error) {
	if filePath != "" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
		}
		return parseTokensFromString(string(content))
	}

	if len(args) == 1 {
		if info, err := os.Stat(args[0]); err == nil && !info.IsDir() {
			content, err := os.ReadFile(args[0])
			if err == nil {
				return parseTokensFromString(string(content))
			}
		}
	}

	stat, _ := os.Stdin.Stat()
	var isTerminal bool
	if stat != nil {
		isTerminal = (stat.Mode() & os.ModeCharDevice) != 0
	}
	return getTokensFromReader(args, os.Stdin, isTerminal)
}

func getTokensFromReader(args []string, r io.Reader, isTerminal bool) ([]int, error) {
	if len(args) > 0 {
		return parseTokensFromString(strings.Join(args, " "))
	}

	if isTerminal {
		return nil, fmt.Errorf("no token IDs provided. Either pass token IDs as arguments or pipe through stdin")
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
			return nil, fmt.Errorf("error reading stdin: %w", err)
		}
	}

	return parseTokensFromString(builder.String())
}

func parseTokensFromString(str string) ([]int, error) {
	tokenStrs := strings.Fields(str)
	tokens := make([]int, 0, len(tokenStrs))
	for _, s := range tokenStrs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		t, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid token ID %q: %w", s, err)
		}
		tokens = append(tokens, t)
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("no valid token IDs provided")
	}

	return tokens, nil
}
