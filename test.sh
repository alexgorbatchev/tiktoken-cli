#!/bin/bash

# Bash tests for tiktoken-cli
# Run with: ./test.sh

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counter for tests
TESTS_PASSED=0
TESTS_FAILED=0

# Build the binary if it doesn't exist
if [ ! -f "./tiktoken" ]; then
    echo -e "${YELLOW}Building tiktoken...${NC}"
    go build -o tiktoken .
fi

# Helper function to run a test
run_test() {
    local name="$1"
    local command="$2"
    local expected="$3"
    
    echo -n "Testing: $name... "
    
    # Run the command and capture output
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

# Helper function to run a test with regex match
run_test_contains() {
    local name="$1"
    local command="$2"
    local expected_pattern="$3"
    
    echo -n "Testing: $name... "
    
    # Run the command and capture output
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
echo "Running tiktoken-cli tests"
echo "=========================================="
echo ""

# ===========================================
# Version command tests
# ===========================================
echo -e "${YELLOW}=== Version Command ===${NC}"

run_test_contains "version command" \
    "./tiktoken version" \
    "tiktoken version"

# ===========================================
# Help command tests
# ===========================================
echo ""
echo -e "${YELLOW}=== Help Command ===${NC}"

run_test_contains "help flag" \
    "./tiktoken --help" \
    "Usage:"

run_test_contains "count help" \
    "./tiktoken count --help" \
    "Count tokens"

run_test_contains "encode help" \
    "./tiktoken encode --help" \
    "Encode text"

run_test_contains "decode help" \
    "./tiktoken decode --help" \
    "Decode token IDs"

# ===========================================
# Count command tests
# ===========================================
echo ""
echo -e "${YELLOW}=== Count Command ===${NC}"

run_test "count from argument" \
    "./tiktoken count 'Hello, world'" \
    "3"

run_test "count from stdin" \
    "echo 'Hello, world' | ./tiktoken count" \
    "3"

run_test "count with gpt-4o model" \
    "./tiktoken count -m gpt-4o 'Hello, world'" \
    "3"

run_test "count with gpt-4 model" \
    "./tiktoken count -m gpt-4 'Hello, world'" \
    "3"

run_test "count with gpt-3.5-turbo model" \
    "./tiktoken count -m gpt-3.5-turbo 'Hello, world'" \
    "3"

# ===========================================
# Default count behavior (no subcommand)
# ===========================================
echo ""
echo -e "${YELLOW}=== Default Count (no subcommand) ===${NC}"

run_test "default count from argument" \
    "./tiktoken 'Hello, world'" \
    "3"

run_test "default count from stdin" \
    "echo 'Hello, world' | ./tiktoken" \
    "3"

run_test "default count with gpt-4o model" \
    "./tiktoken -m gpt-4o 'Hello, world'" \
    "3"

run_test "default count with encoding" \
    "./tiktoken -e o200k_base 'Hello, world'" \
    "3"

run_test "default count with o200k_base encoding" \
    "./tiktoken count -e o200k_base 'Hello, world'" \
    "3"

run_test "count with cl100k_base encoding" \
    "./tiktoken count -e cl100k_base 'Hello, world'" \
    "3"

run_test "count with p50k_base encoding" \
    "./tiktoken count -e p50k_base 'Hello, world'" \
    "3"

run_test "count with r50k_base encoding" \
    "./tiktoken count -e r50k_base 'Hello, world'" \
    "3"

run_test "count multiline from stdin" \
    "printf 'Hello\nWorld' | ./tiktoken count" \
    "3"

run_test "count empty string" \
    "./tiktoken count ''" \
    "0"

run_test "count unicode text" \
    "./tiktoken count '你好世界'" \
    "5"

run_test "count emoji" \
    "./tiktoken count '🎉'" \
    "3"

# ===========================================
# Encode command tests
# ===========================================
echo ""
echo -e "${YELLOW}=== Encode Command ===${NC}"

run_test "encode from argument" \
    "./tiktoken encode 'Hello world'" \
    "9906 1917"

run_test "encode from stdin" \
    "echo 'Hello world' | ./tiktoken encode" \
    "9906 1917"

run_test "encode with gpt-4o model" \
    "./tiktoken encode -m gpt-4o 'Hello world'" \
    "13225 2375"

run_test "encode with cl100k_base encoding" \
    "./tiktoken encode -e cl100k_base 'Hello world'" \
    "9906 1917"

run_test "encode with o200k_base encoding" \
    "./tiktoken encode -e o200k_base 'Hello world'" \
    "13225 2375"

# ===========================================
# Decode command tests
# ===========================================
echo ""
echo -e "${YELLOW}=== Decode Command ===${NC}"

run_test "decode from arguments" \
    "./tiktoken decode 9906 1917" \
    "Hello world"

run_test "decode from stdin" \
    "echo '9906 1917' | ./tiktoken decode" \
    "Hello world"

run_test "decode with gpt-4o model" \
    "./tiktoken decode -m gpt-4o 13225 2375" \
    "Hello world"

run_test "decode with cl100k_base encoding" \
    "./tiktoken decode -e cl100k_base 9906 1917" \
    "Hello world"

run_test "decode with o200k_base encoding" \
    "./tiktoken decode -e o200k_base 13225 2375" \
    "Hello world"

# ===========================================
# Encode/Decode round-trip tests
# ===========================================
echo ""
echo -e "${YELLOW}=== Round-trip Tests ===${NC}"

run_test "encode then decode (cl100k_base)" \
    "./tiktoken encode 'Hello world' | ./tiktoken decode" \
    "Hello world"

run_test "encode then decode (o200k_base)" \
    "./tiktoken encode -e o200k_base 'Hello world' | ./tiktoken decode -e o200k_base" \
    "Hello world"

run_test "encode then decode unicode" \
    "./tiktoken encode '你好世界' | ./tiktoken decode" \
    "你好世界"

run_test "encode then decode with spaces" \
    "./tiktoken encode 'Hello   world' | ./tiktoken decode" \
    "Hello   world"

# ===========================================
# Models command tests
# ===========================================
echo ""
echo -e "${YELLOW}=== Models Command ===${NC}"

run_test_contains "models lists o200k_base" \
    "./tiktoken models" \
    "o200k_base"

run_test_contains "models lists cl100k_base" \
    "./tiktoken models" \
    "cl100k_base"

run_test_contains "models lists gpt-4o" \
    "./tiktoken models" \
    "gpt-4o"

run_test_contains "models lists gpt-4" \
    "./tiktoken models" \
    "gpt-4"

# ===========================================
# Error handling tests
# ===========================================
echo ""
echo -e "${YELLOW}=== Error Handling ===${NC}"

run_test_contains "invalid model error" \
    "./tiktoken count -m invalid-model 'test' 2>&1 || true" \
    "no encoding for model"

run_test_contains "invalid encoding error" \
    "./tiktoken count -e invalid-encoding 'test' 2>&1 || true" \
    "Unknown encoding"

run_test_contains "invalid token ID error" \
    "./tiktoken decode abc 2>&1 || true" \
    "invalid token ID"

# ===========================================
# Summary
# ===========================================
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
