#!/bin/bash

echo "=== Testing retry program ==="
echo

# Build the program first
echo "Building retry..."
go build -o bin/retry cmd/retry/main.go
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi
echo "✓ Build successful"
echo

# Track test results
declare -a PASSED_TESTS
declare -a FAILED_TESTS

# Helper function to run a test
run_test() {
    local test_num="$1"
    local test_name="$2"
    local expected_exit="$3"
    shift 3
    local cmd=("$@")
    
    echo "Test $test_num: $test_name"
    
    # Run the command
    "${cmd[@]}"
    local actual_exit=$?
    
    # Check if it matches expected exit code
    if [ "$actual_exit" -eq "$expected_exit" ]; then
        echo "✓ PASSED (exit code: $actual_exit)"
        PASSED_TESTS+=("Test $test_num: $test_name")
    else
        echo "✗ FAILED (expected exit: $expected_exit, got: $actual_exit)"
        FAILED_TESTS+=("Test $test_num: $test_name")
    fi
    
    echo
    echo "---"
    echo
}

# Test 1: Simple success (should succeed on first try)
run_test 1 "Command that succeeds immediately" 0 \
    ./bin/retry 3 echo "Hello, World!"

# Test 2: Command that always fails
run_test 2 "Command that always fails" 1 \
    ./bin/retry 3 ls /nonexistent/directory

# Test 3: With delay flag
run_test 3 "With 2 second delay between retries" 1 \
    ./bin/retry --delay 2s 3 ls /fake/path

# Test 4: With verbose flag
run_test 4 "With verbose output" 0 \
    ./bin/retry -v 3 echo "Testing verbose mode"

# Test 5: Quiet mode
run_test 5 "Quiet mode (should only show command output)" 0 \
    ./bin/retry -q 3 echo "Testing quiet mode"

# Test 6: Create a flaky command (fails twice, then succeeds)
echo "Test 6: Flaky command (fails 2 times, succeeds on 3rd)"
cat > /tmp/flaky.sh << 'EOF'
#!/bin/bash
FILE=/tmp/retry_counter
if [ ! -f "$FILE" ]; then
    echo "0" > "$FILE"
fi
COUNT=$(cat "$FILE")
COUNT=$((COUNT + 1))
echo "$COUNT" > "$FILE"
echo "Attempt number: $COUNT"
if [ $COUNT -lt 3 ]; then
    echo "Failing..." >&2
    exit 1
fi
echo "Success!"
rm "$FILE"
exit 0
EOF
chmod +x /tmp/flaky.sh

run_test 6 "Flaky command (succeeds on 3rd attempt)" 0 \
    ./bin/retry 5 /tmp/flaky.sh

# Test 7: Command timeout - fast command (should succeed)
echo "Test 7: Command timeout - fast command (should succeed)"
run_test 7 "Fast command with timeout" 0 \
    ./bin/retry --command-timeout 5s 3 echo "Quick command"

# Test 8: Command timeout - slow command (should timeout)
echo "Test 8: Command timeout - slow command (should timeout and retry)"
cat > /tmp/slow.sh << 'EOF'
#!/bin/bash
echo "Starting slow command..."
sleep 10
echo "Finished!"
EOF
chmod +x /tmp/slow.sh

run_test 8 "Slow command with short timeout" 1 \
    ./bin/retry --command-timeout 2s 3 /tmp/slow.sh

# Test 9: Command timeout with verbose
run_test 9 "Timeout with verbose output" 1 \
    ./bin/retry -v --command-timeout 1s 2 sleep 5

# Test 10: Overall timeout (should stop retrying after overall timeout)
run_test 10 "Overall timeout" 1 \
    ./bin/retry --overall-timeout 5s 10 sleep 3


# Test 10: Real HTTP request (optional - needs internet)
run_test 10 "Real HTTP request" 0 \
    ./bin/retry 3 -- curl -s https://httpbin.org/status/200

# Test 11: Multiple retries with different delays
run_test 11 "Multiple retries with 0.5s delay" 1 \
    ./bin/retry --delay 0.5s 4 false

# Test 12: Success after timeout on earlier attempts
echo "Test 12: Command that times out first, then succeeds"
cat > /tmp/timeout_then_success.sh << 'EOF'
#!/bin/bash
FILE=/tmp/timeout_counter
if [ ! -f "$FILE" ]; then
    echo "0" > "$FILE"
fi
COUNT=$(cat "$FILE")
COUNT=$((COUNT + 1))
echo "$COUNT" > "$FILE"
echo "Attempt: $COUNT"

if [ $COUNT -lt 2 ]; then
    sleep 10  # Will timeout
else
    echo "Success on attempt $COUNT!"
    rm "$FILE"
fi
EOF
chmod +x /tmp/timeout_then_success.sh

run_test 12 "Timeout then success" 0 \
    ./bin/retry --command-timeout 2s 3 /tmp/timeout_then_success.sh

# Clean up
rm -f /tmp/flaky.sh /tmp/retry_counter /tmp/slow.sh /tmp/timeout_then_success.sh /tmp/timeout_counter

echo
echo "========================================"
echo "           TEST SUMMARY"
echo "========================================"
echo

echo "PASSED TESTS (${#PASSED_TESTS[@]}):"
if [ ${#PASSED_TESTS[@]} -eq 0 ]; then
    echo "  (none)"
else
    for test in "${PASSED_TESTS[@]}"; do
        echo "  ✓ $test"
    done
fi

echo
echo "FAILED TESTS (${#FAILED_TESTS[@]}):"
if [ ${#FAILED_TESTS[@]} -eq 0 ]; then
    echo "  (none)"
else
    for test in "${FAILED_TESTS[@]}"; do
        echo "  ✗ $test"
    done
fi

echo
echo "========================================"
echo "Total: $((${#PASSED_TESTS[@]} + ${#FAILED_TESTS[@]})) tests"
echo "Passed: ${#PASSED_TESTS[@]}"
echo "Failed: ${#FAILED_TESTS[@]}"
echo "========================================"

# Exit with failure if any tests failed
if [ ${#FAILED_TESTS[@]} -gt 0 ]; then
    exit 1
else
    exit 0
fi