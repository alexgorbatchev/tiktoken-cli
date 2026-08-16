# Default target: list available recipes
default:
    @just --list

# Build binary into ./bin/tiktoken
build:
    @mkdir -p bin
    go build -o bin/tiktoken .

# Run the CLI binary with optional arguments
run *ARGS: build
    ./bin/tiktoken {{ ARGS }}

# Run unit and integration tests
test: test-unit test-integration

# Run Go unit tests
test-unit:
    go test -v ./...

# Run integration tests against ./bin/tiktoken
test-integration: build
    #!/usr/bin/env bash
    set -euo pipefail

    BIN="./bin/tiktoken"
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    NC='\033[0m'

    TESTS_PASSED=0
    TESTS_FAILED=0

    run_test() {
        local name="$1"
        local command="$2"
        local expected="$3"
        
        echo -n "Testing: $name... "
        local output
        output=$(eval "$command" 2>&1) || true
        
        if [ "$output" = "$expected" ]; then
            echo -e "${GREEN}PASSED${NC}"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            echo -e "${RED}FAILED${NC}"
            echo "  Expected: '$expected'"
            echo "  Got:      '$output'"
            TESTS_FAILED=$((TESTS_FAILED + 1))
        fi
    }

    run_test_contains() {
        local name="$1"
        local command="$2"
        local expected_pattern="$3"
        
        echo -n "Testing: $name... "
        local output
        output=$(eval "$command" 2>&1) || true
        
        if echo "$output" | grep -q "$expected_pattern"; then
            echo -e "${GREEN}PASSED${NC}"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            echo -e "${RED}FAILED${NC}"
            echo "  Expected to contain: '$expected_pattern'"
            echo "  Got: '$output'"
            TESTS_FAILED=$((TESTS_FAILED + 1))
        fi
    }

    echo "=========================================="
    echo "Running tiktoken-cli integration tests"
    echo "=========================================="
    echo ""

    TMP_TEXT_FILE=$(mktemp /tmp/tiktoken_sample.XXXXXX.txt)
    TMP_TOKEN_FILE=$(mktemp /tmp/tiktoken_tokens.XXXXXX.txt)
    echo -n "Hello world" > "$TMP_TEXT_FILE"
    echo -n "9906 1917" > "$TMP_TOKEN_FILE"
    trap 'rm -f "$TMP_TEXT_FILE" "$TMP_TOKEN_FILE"' EXIT

    echo -e "${YELLOW}=== Version Command ===${NC}"
    run_test_contains "version command" "$BIN version" "tiktoken version"

    echo ""
    echo -e "${YELLOW}=== Help Command ===${NC}"
    run_test_contains "help flag" "$BIN --help" "Usage:"
    run_test_contains "count help" "$BIN count --help" "Count tokens"
    run_test_contains "encode help" "$BIN encode --help" "Encode text"
    run_test_contains "decode help" "$BIN decode --help" "Decode token IDs"

    echo ""
    echo -e "${YELLOW}=== Count Command ===${NC}"
    run_test "count from argument" "$BIN count 'Hello, world'" "3"
    run_test "count from stdin" "echo 'Hello, world' | $BIN count" "3"
    run_test "count with gpt-4o model" "$BIN count -m gpt-4o 'Hello, world'" "3"
    run_test "count with gpt-4 model" "$BIN count -m gpt-4 'Hello, world'" "3"
    run_test "count with gpt-3.5-turbo model" "$BIN count -m gpt-3.5-turbo 'Hello, world'" "3"

    echo ""
    echo -e "${YELLOW}=== Default Count (no subcommand) ===${NC}"
    run_test "default count from argument" "$BIN 'Hello, world'" "3"
    run_test "default count from stdin" "echo 'Hello, world' | $BIN" "3"
    run_test "default count with gpt-4o model" "$BIN -m gpt-4o 'Hello, world'" "3"
    run_test "default count with encoding" "$BIN -e o200k_base 'Hello, world'" "3"
    run_test "default count with o200k_base encoding" "$BIN count -e o200k_base 'Hello, world'" "3"
    run_test "count with cl100k_base encoding" "$BIN count -e cl100k_base 'Hello, world'" "3"
    run_test "count with p50k_base encoding" "$BIN count -e p50k_base 'Hello, world'" "3"
    run_test "count with r50k_base encoding" "$BIN count -e r50k_base 'Hello, world'" "3"
    run_test "count multiline from stdin" "printf 'Hello\nWorld' | $BIN count" "3"
    run_test "count empty string" "$BIN count ''" "0"
    run_test "count unicode text" "$BIN count '你好世界'" "5"
    run_test "count emoji" "$BIN count '🎉'" "3"

    echo ""
    echo -e "${YELLOW}=== Encode Command ===${NC}"
    run_test "encode from argument" "$BIN encode 'Hello world'" "9906 1917"
    run_test "encode from stdin" "echo 'Hello world' | $BIN encode" "9906 1917"
    run_test "encode with gpt-4o model" "$BIN encode -m gpt-4o 'Hello world'" "13225 2375"
    run_test "encode with cl100k_base encoding" "$BIN encode -e cl100k_base 'Hello world'" "9906 1917"
    run_test "encode with o200k_base encoding" "$BIN encode -e o200k_base 'Hello world'" "13225 2375"

    echo ""
    echo -e "${YELLOW}=== Decode Command ===${NC}"
    run_test "decode from arguments" "$BIN decode 9906 1917" "Hello world"
    run_test "decode from stdin" "echo '9906 1917' | $BIN decode" "Hello world"
    run_test "decode with gpt-4o model" "$BIN decode -m gpt-4o 13225 2375" "Hello world"
    run_test "decode with cl100k_base encoding" "$BIN decode -e cl100k_base 9906 1917" "Hello world"
    run_test "decode with o200k_base encoding" "$BIN decode -e o200k_base 13225 2375" "Hello world"

    echo ""
    echo -e "${YELLOW}=== File Input Tests ===${NC}"
    run_test "count from file path argument" "$BIN count '$TMP_TEXT_FILE'" "2"
    run_test "count using -f flag" "$BIN count -f '$TMP_TEXT_FILE'" "2"
    run_test "default count from file path" "$BIN '$TMP_TEXT_FILE'" "2"
    run_test "encode from file flag" "$BIN encode -f '$TMP_TEXT_FILE'" "9906 1917"
    run_test "decode from file flag" "$BIN decode -f '$TMP_TOKEN_FILE'" "Hello world"
    run_test "encode then decode (cl100k_base)" "$BIN encode 'Hello world' | $BIN decode" "Hello world"
    run_test "encode then decode (o200k_base)" "$BIN encode -e o200k_base 'Hello world' | $BIN decode -e o200k_base" "Hello world"
    run_test "encode then decode unicode" "$BIN encode '你好世界' | $BIN decode" "你好世界"
    run_test "encode then decode with spaces" "$BIN encode 'Hello   world' | $BIN decode" "Hello   world"

    echo ""
    echo -e "${YELLOW}=== Models Command ===${NC}"
    run_test_contains "models lists o200k_base" "$BIN models" "o200k_base"
    run_test_contains "models lists cl100k_base" "$BIN models" "cl100k_base"
    run_test_contains "models lists gpt-4o" "$BIN models" "gpt-4o"
    run_test_contains "models lists gpt-4" "$BIN models" "gpt-4"

    echo ""
    echo -e "${YELLOW}=== Error Handling ===${NC}"
    run_test_contains "invalid model error" "$BIN count -m invalid-model 'test' 2>&1 || true" "no encoding for model"
    run_test_contains "invalid encoding error" "$BIN count -e invalid-encoding 'test' 2>&1 || true" "Unknown encoding"
    run_test_contains "invalid token ID error" "$BIN decode abc 2>&1 || true" "invalid token ID"

    echo ""
    echo "=========================================="
    echo "Test Summary"
    echo "=========================================="
    echo -e "Passed: ${GREEN}$TESTS_PASSED${NC}"
    echo -e "Failed: ${RED}$TESTS_FAILED${NC}"
    echo ""

    if [ $TESTS_FAILED -gt 0 ]; then
        echo -e "${RED}Some tests failed!${NC}"
        exit 1
    else
        echo -e "${GREEN}All tests passed!${NC}"
        exit 0
    fi

# Run go vet
lint:
    go vet ./...

# Remove build artifacts
clean:
    rm -rf bin/ tiktoken
