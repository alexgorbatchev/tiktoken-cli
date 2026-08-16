---
created_on: 2025-08-16 08:12
last_modified: 2025-08-16 08:12
status: current
---

# tiktoken-cli

CLI tool for tokenizing, counting, encoding, and decoding text with OpenAI tiktoken encodings in Go.

## Commands
- Build: `just build` (outputs to `bin/tiktoken`)
- Run: `just run "Hello, world!"` or `bin/tiktoken count "Hello, world!"`
- Test: `just test` (runs unit and integration test suite)
- Unit Test: `just test-unit` (`go test -v ./...`)
- Integration Test: `just test-integration` (executes CLI tests against `bin/tiktoken`)
- Lint: `just lint` (`go vet ./...`)
- Clean: `just clean` (removes `bin/` build directory)

## Setup
- Go 1.26+ installed (`go.mod` targets Go 1.26).
- `just` command-line runner installed.

## Conventions
- Structure Cobra subcommands as individual files inside `cmd/` (e.g., `cmd/count.go`, `cmd/encode.go`, `cmd/decode.go`, `cmd/models.go`).
- Input resolution: Use `getText(args)` or `getTokens(args)` helper functions in `cmd/` to accept positional arguments or fall back to reading `os.Stdin`.
- Encoding resolution: Check `--model` (`-m`) first via `tiktoken.EncodingForModel`, falling back to `--encoding` (`-e`, default `cl100k_base`) via `tiktoken.GetEncoding`.
- Keep output plain and copy-pasteable (number for `count`, space-separated IDs for `encode`, text string for `decode`).

## Gotchas
- Binary placement: Always output compiled binaries to `bin/tiktoken` (`just build`). Never output binaries into the repository root.
- Interactive `os.Stdin` check: `getText` and `getTokens` check `os.Stdin.Stat().Mode() & os.ModeCharDevice`. In non-interactive or piped contexts, pass input via stdin or supply arguments.
- Integration tests: `test.sh` / `just test-integration` relies on `bin/tiktoken`. Run `just build` before running individual integration test commands directly.

## Boundaries
- Always: automatically record all new instructions in the most appropriate `AGENTS.md` file immediately upon receipt (check with user if existing instructions conflict).
- Always: any time code is changed such that results from running that code are changed, a test file must be changed as well; 90% code coverage is required (scripts/ folder is excluded from this rule).
- Always: use `just` recipes (`just build`, `just test`, `just lint`) for build and test automation; place compiled binaries in `bin/`.
- Ask first: adding new third-party Go dependencies, changing default encoding/model flags, or altering CLI output formatting.
- Never: commit compiled Go binaries (`bin/`, root `tiktoken`), commit secrets or raw API tokens, or delete failing tests to force a passing run.

## References
- `justfile`
- `README.md`
- `cmd/root.go`
