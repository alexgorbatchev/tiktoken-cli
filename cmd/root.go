package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"

	// Root command flags (for default count behavior)
	rootModel    string
	rootEncoding string
)

var rootCmd = &cobra.Command{
	Use:   "tiktoken [text]",
	Short: "A CLI tool for counting tokens using OpenAI's tiktoken",
	Long: `tiktoken-cli is a command line interface for the tiktoken-go library.
It allows you to count tokens in text using various OpenAI tokenization encodings.

When called without a subcommand, it defaults to counting tokens using the
cl100k_base encoding (used by GPT-4 and GPT-3.5-turbo).

Examples:
  # Count tokens (uses cl100k_base encoding by default)
  echo "Hello, world!" | tiktoken

  # Count tokens from argument
  tiktoken "Hello, world!"

  # Count tokens for a specific model
  tiktoken -m gpt-4o "Hello, world!"

  # Count tokens using a specific encoding
  tiktoken -e o200k_base "Hello, world!"

  # Encode text to token IDs
  tiktoken encode "Hello, world!"

  # Decode token IDs back to text
  tiktoken decode 15339 1917 0`,
	Args: cobra.ArbitraryArgs,
	RunE: runRootCount,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "tiktoken version %s\n", version)
		fmt.Fprintf(cmd.OutOrStdout(), "  commit: %s\n", commit)
		fmt.Fprintf(cmd.OutOrStdout(), "  built:  %s\n", date)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().StringVarP(&rootModel, "model", "m", "", "OpenAI model name (e.g., gpt-4o, gpt-4, gpt-3.5-turbo)")
	rootCmd.Flags().StringVarP(&rootEncoding, "encoding", "e", "cl100k_base", "Encoding name (o200k_base, cl100k_base, p50k_base, r50k_base)")
}

// runRootCount is the default behavior when no subcommand is provided
func runRootCount(cmd *cobra.Command, args []string) error {
	stat, _ := os.Stdin.Stat()
	if len(args) == 0 && stat != nil && (stat.Mode()&os.ModeCharDevice) != 0 {
		return cmd.Help()
	}

	text, err := getText(args)
	if err != nil {
		return err
	}

	enc, err := resolveEncoding(rootModel, rootEncoding)
	if err != nil {
		return err
	}

	tokens := enc.Encode(text, nil, nil)
	fmt.Fprintln(cmd.OutOrStdout(), len(tokens))

	return nil
}
