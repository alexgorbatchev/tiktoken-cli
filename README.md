# tiktoken-cli

A command-line interface for counting tokens using OpenAI's tiktoken tokenizer, powered by [tiktoken-go](https://github.com/pkoukk/tiktoken-go).

## Features

- Count tokens in text using various OpenAI tokenization encodings
- Encode text to token IDs
- Decode token IDs back to text
- Support for all major OpenAI models (GPT-4o, GPT-4, GPT-3.5-turbo, etc.)
- Read from stdin for piping
- Cross-platform (Linux, macOS, Windows)

## Installation

### Download Binary

Download the latest release from the [releases page](https://github.com/alexgorbatchev/tiktoken-cli/releases).

### Build from Source

```bash
go install github.com/alexgorbatchev/tiktoken-cli@latest
```

Or clone and build:

```bash
git clone https://github.com/alexgorbatchev/tiktoken-cli.git
cd tiktoken-cli
just build
# Compiled binary is output to ./bin/tiktoken
```

## Usage

### Count Tokens (Default)

Count is the default action - no subcommand needed. Uses `cl100k_base` encoding by default (GPT-4, GPT-3.5-turbo).

```bash
# Count tokens (uses cl100k_base encoding by default)
tiktoken "Hello, world!"

# Count tokens from stdin
echo "Hello, world!" | tiktoken

# Count tokens for a specific model
tiktoken -m gpt-4o "Hello, world!"

# Count tokens using a specific encoding
tiktoken -e o200k_base "Hello, world!"

# Count tokens from a file
cat myfile.txt | tiktoken

# Using explicit count subcommand (also works)
tiktoken count "Hello, world!"
```

### Encode Text

```bash
# Encode text to token IDs
tiktoken encode "Hello, world!"
# Output: 9906 11 1917 0

# Encode using a specific model
tiktoken encode -m gpt-4o "Hello, world!"
```

### Decode Tokens

```bash
# Decode token IDs back to text
tiktoken decode 15339 1917 0

# Decode from stdin
echo "15339 1917 0" | tiktoken decode

# Chain encode and decode
tiktoken encode "Hello, world!" | tiktoken decode
```

### List Models

```bash
# Show available models and encodings
tiktoken models
```

### Version

```bash
tiktoken version
```

## Available Encodings

| Encoding     | Models                                                              |
|--------------|---------------------------------------------------------------------|
| `o200k_base` | gpt-4o, gpt-4.1, gpt-4.5                                           |
| `cl100k_base`| gpt-4, gpt-3.5-turbo, text-embedding-ada-002, text-embedding-3-*   |
| `p50k_base`  | Codex models, text-davinci-002, text-davinci-003                   |
| `r50k_base`  | GPT-3 models (davinci, curie, babbage, ada)                        |

## Flags

| Flag            | Short | Description                                    |
|-----------------|-------|------------------------------------------------|
| `--model`       | `-m`  | OpenAI model name (e.g., gpt-4o, gpt-4)       |
| `--encoding`    | `-e`  | Encoding name (default: cl100k_base)          |
| `--help`        | `-h`  | Show help for command                          |

## Examples

### Counting tokens for a file

```bash
cat README.md | tiktoken count
```

### Using in shell scripts

```bash
TOKEN_COUNT=$(echo "Your text here" | tiktoken count)
echo "Token count: $TOKEN_COUNT"
```

### Comparing encodings

```bash
TEXT="Hello, world!"
echo "o200k_base: $(tiktoken count -e o200k_base "$TEXT")"
echo "cl100k_base: $(tiktoken count -e cl100k_base "$TEXT")"
echo "p50k_base: $(tiktoken count -e p50k_base "$TEXT")"
```

## Development & Testing

Build, run, and test with `just`:

```bash
# Build binary to ./bin/tiktoken
just build

# Run binary with arguments
just run "Hello, world!"

# Run unit and integration tests
just test

# Run linter
just lint
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- [tiktoken-go](https://github.com/pkoukk/tiktoken-go) - Go implementation of tiktoken
- [tiktoken](https://github.com/openai/tiktoken) - Original Python library by OpenAI
