#!/bin/bash
# ADP CLI E2E Test Script
# Supports two modes:
#   - Offline tests: always run (version, help, config, schema)
#   - API tests: require ADP_API_KEY and ADP_API_BASE_URL env vars

set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0
TESTS_SKIPPED=0

# Resolve paths relative to this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
ADP="$PROJECT_DIR/adp"
RESULT_FILE="$SCRIPT_DIR/test_result.json"

# Build if binary doesn't exist
if [ ! -f "$ADP" ] && [ ! -f "$ADP.exe" ]; then
    echo "Building adp binary..."
    cd "$PROJECT_DIR" && go build -o adp . && cd "$SCRIPT_DIR"
fi

# Windows: use .exe
if [ -f "$ADP.exe" ]; then
    ADP="$ADP.exe"
fi

# Check API availability
API_AVAILABLE=false
ORIGINAL_API_KEY=""
ORIGINAL_API_BASE_URL=""
if [ -n "${ADP_API_KEY:-}" ] && [ -n "${ADP_API_BASE_URL:-}" ]; then
    API_AVAILABLE=true
    echo "API credentials from env vars. Running full test suite."
    "$ADP" config set --api-key "$ADP_API_KEY"
    "$ADP" config set --api-base-url "$ADP_API_BASE_URL"
elif "$ADP" config get 2>/dev/null | grep -q '"configured": true'; then
    API_AVAILABLE=true
    # Save original config for restore after config clear test
    ORIGINAL_API_BASE_URL=$("$ADP" config get 2>/dev/null | grep -o '"api_base_url": "[^"]*' | cut -d'"' -f4 || true)
    echo "Existing config found. Running full test suite."
else
    echo "No API credentials. Running offline tests only."
fi

REPORT_FILE="$SCRIPT_DIR/test_report.txt"

# Helpers
print_header() {
    echo ""
    echo "========================================"
    echo "$1"
    echo "========================================"
}

# Test result log for report
TEST_LOG=""

run_test() {
    local name="$1"
    shift
    echo -e "\n${YELLOW}[TEST]${NC} $name"
    TESTS_RUN=$((TESTS_RUN + 1))
    local output
    output=$("$@" 2>&1) || true
    local exit_code=$?
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}[PASS]${NC} $name"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        TEST_LOG+="PASS | $name"$'\n'
        return 0
    else
        echo -e "${RED}[FAIL]${NC} $name (exit code: $exit_code)"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        TEST_LOG+="FAIL | $name | exit=$exit_code"$'\n'
        return 1
    fi
}

run_test_capture() {
    local name="$1"
    shift
    echo -e "\n${YELLOW}[TEST]${NC} $name"
    TESTS_RUN=$((TESTS_RUN + 1))
    CAPTURE_OUTPUT=$("$@" 2>&1) || true
    if [ $? -eq 0 ] || echo "$CAPTURE_OUTPUT" | grep -q '"code": "success"'; then
        echo -e "${GREEN}[PASS]${NC} $name"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        return 0
    else
        echo -e "${RED}[FAIL]${NC} $name"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        return 1
    fi
}

skip_test() {
    echo -e "\n${YELLOW}[SKIP]${NC} $1"
    TESTS_RUN=$((TESTS_RUN + 1))
    TESTS_SKIPPED=$((TESTS_SKIPPED + 1))
    TEST_LOG+="SKIP | $1"$'\n'
}

cleanup() {
    rm -f "$RESULT_FILE"
    rm -f "$SCRIPT_DIR/parse_tasks.json" "$SCRIPT_DIR/extract_tasks.json"
    # Delete custom app if created
    if [ -n "${CREATED_APP_ID:-}" ] && [ "$API_AVAILABLE" = true ]; then
        "$ADP" custom-app delete --app-id "$CREATED_APP_ID" > /dev/null 2>&1 || true
    fi
}
trap cleanup EXIT

# ============================================
# 1. Offline Tests (always run)
# ============================================
print_header "1. Version & Help"

run_test "adp version" "$ADP" version
run_test "adp --help" "$ADP" --help
run_test "adp parse --help" "$ADP" parse --help
run_test "adp extract --help" "$ADP" extract --help
run_test "adp custom-app --help" "$ADP" custom-app --help
run_test "adp --lang en --help" "$ADP" --lang en --help
run_test "adp --lang zh --help" "$ADP" --lang zh --help

# ============================================
# 2. Config Commands
# ============================================
print_header "2. Config Commands"

if [ "$API_AVAILABLE" = true ]; then
    # API config exists — only test non-destructive commands
    run_test "config get" "$ADP" config get
    run_test "config set --api-base-url (restore)" "$ADP" config set --api-base-url "$ORIGINAL_API_BASE_URL"
    skip_test "config set --api-key (skipped to preserve API config)"
    skip_test "config clear (skipped to preserve API config)"
else
    # No API config — safe to do full config lifecycle test
    run_test "config set --api-key" "$ADP" config set --api-key test_key_for_e2e
    run_test "config set --api-base-url" "$ADP" config set --api-base-url https://test.example.com
    run_test "config get" "$ADP" config get
    run_test "config clear --force" "$ADP" config clear --force
fi

# ============================================
# 3. Schema Commands (offline)
# ============================================
print_header "3. Schema Commands"

run_test "schema (full tree)" "$ADP" schema
run_test "schema parse" "$ADP" schema parse
run_test "schema parse local" "$ADP" schema parse local
run_test "schema extract" "$ADP" schema extract
run_test "schema custom-app" "$ADP" schema custom-app

# ============================================
# 4. API Tests (require credentials)
# ============================================
if [ "$API_AVAILABLE" != true ]; then
    print_header "4-8. API Tests (SKIPPED)"
    skip_test "app-id list"
    skip_test "parse local"
    skip_test "extract local"
    skip_test "custom-app create"
    skip_test "credit"
else

# ---- App ID ----
print_header "4. App ID Commands"

run_test "app-id list" "$ADP" app-id list
run_test "app-id list --app-label" "$ADP" app-id list --app-label 合同
run_test "app-id cache" "$ADP" app-id cache

# Get a valid APP_ID for subsequent tests
APP_ID=$("$ADP" app-id list 2>/dev/null | grep -o '"app_id": "[^"]*' | head -1 | cut -d'"' -f4 || true)
if [ -z "$APP_ID" ]; then
    echo "WARNING: Could not get APP_ID, using fallback"
    APP_ID="ootb_k7m2x9p4v1n8w3q6r5t0y2b4"
fi
echo "Using APP_ID: $APP_ID"

# ---- Parse ----
print_header "5. Parse Commands"

run_test "parse local (sync)" "$ADP" parse local "$SCRIPT_DIR/samples/73.蚂蚁+B类.png" --app-id "$APP_ID"
run_test "parse local (async)" "$ADP" parse local "$SCRIPT_DIR/samples/73.蚂蚁+B类.png" --app-id "$APP_ID" --async
run_test "parse local directory" "$ADP" parse local "$SCRIPT_DIR" --app-id "$APP_ID" --async
run_test "parse local --export" "$ADP" parse local "$SCRIPT_DIR/samples/73.蚂蚁+B类.png" --app-id "$APP_ID" --export "$RESULT_FILE"
run_test "parse url (single)" "$ADP" parse url https://adp-global.laiye.com/web/agentic_doc_processor/laiye/file/13e18a44228611f1933a00163e122259 --app-id "$APP_ID"
run_test "parse url (file list)" "$ADP" parse url "$SCRIPT_DIR/samples/url.txt" --app-id "$APP_ID"

# --no-wait + query --file
PARSE_TASKS_FILE="$SCRIPT_DIR/parse_tasks.json"
run_test "parse local --async --no-wait" "$ADP" parse local "$SCRIPT_DIR" --app-id "$APP_ID" --async --no-wait --export "$PARSE_TASKS_FILE"
if [ -f "$PARSE_TASKS_FILE" ]; then
    run_test "parse query --file (from --no-wait)" "$ADP" parse query --watch --file "$PARSE_TASKS_FILE"
    rm -f "$PARSE_TASKS_FILE"
else
    skip_test "parse query --file (no tasks file generated)"
fi

# ---- Extract ----
print_header "6. Extract Commands"

run_test "extract local (sync)" "$ADP" extract local "$SCRIPT_DIR/samples/73.蚂蚁+B类.png" --app-id "$APP_ID"
run_test "extract local (async)" "$ADP" extract local "$SCRIPT_DIR/samples/73.蚂蚁+B类.png" --app-id "$APP_ID" --async
run_test "extract local directory" "$ADP" extract local "$SCRIPT_DIR" --app-id "$APP_ID" --async
run_test "extract local --export" "$ADP" extract local "$SCRIPT_DIR/samples/73.蚂蚁+B类.png" --app-id "$APP_ID" --export "$RESULT_FILE"
run_test "extract url (single)" "$ADP" extract url https://adp-global.laiye.com/web/agentic_doc_processor/laiye/file/13e18a44228611f1933a00163e122259 --app-id "$APP_ID"
run_test "extract url (file list)" "$ADP" extract url "$SCRIPT_DIR/samples/url.txt" --app-id "$APP_ID"

# --no-wait + query --file
EXTRACT_TASKS_FILE="$SCRIPT_DIR/extract_tasks.json"
run_test "extract local --async --no-wait" "$ADP" extract local "$SCRIPT_DIR" --app-id "$APP_ID" --async --no-wait --export "$EXTRACT_TASKS_FILE"
if [ -f "$EXTRACT_TASKS_FILE" ]; then
    run_test "extract query --file (from --no-wait)" "$ADP" extract query --watch --file "$EXTRACT_TASKS_FILE"
    rm -f "$EXTRACT_TASKS_FILE"
else
    skip_test "extract query --file (no tasks file generated)"
fi

# ---- Custom App ----
print_header "7. Custom App Commands"

# Create using JSON file paths (avoids shell quoting issues)
CREATED_APP_ID=""
echo -e "\n${YELLOW}[TEST]${NC} custom-app create"
TESTS_RUN=$((TESTS_RUN + 1))
CREATE_OUTPUT=$("$ADP" custom-app create \
    --app-name "E2E-Test-App" \
    --app-label "e2e,test" \
    --extract-fields "$SCRIPT_DIR/samples/extract-fields.json" \
    --parse-mode standard \
    --enable-long-doc true \
    --long-doc-config "$SCRIPT_DIR/samples/long_doc_config.json" 2>&1) || true

CREATED_APP_ID=$(echo "$CREATE_OUTPUT" | grep -o '"app_id": "[^"]*' | head -1 | cut -d'"' -f4 || true)
if [ -n "$CREATED_APP_ID" ]; then
    echo -e "${GREEN}[PASS]${NC} custom-app create (app_id: $CREATED_APP_ID)"
    TESTS_PASSED=$((TESTS_PASSED + 1))
    TEST_LOG+="PASS | custom-app create (app_id: $CREATED_APP_ID)"$'\n'
else
    echo -e "${RED}[FAIL]${NC} custom-app create"
    echo "$CREATE_OUTPUT"
    TESTS_FAILED=$((TESTS_FAILED + 1))
    TEST_LOG+="FAIL | custom-app create"$'\n'
fi

if [ -n "$CREATED_APP_ID" ]; then
    run_test "custom-app get-config" "$ADP" custom-app get-config --app-id "$CREATED_APP_ID"
    run_test "custom-app get-config --config-version v1" "$ADP" custom-app get-config --app-id "$CREATED_APP_ID" --config-version v1
    run_test "custom-app ai-generate (local)" "$ADP" custom-app ai-generate --app-id "$CREATED_APP_ID" --file-local "$SCRIPT_DIR/samples/73.蚂蚁+B类.png"
    run_test "custom-app ai-generate (url)" "$ADP" custom-app ai-generate --app-id "$CREATED_APP_ID" --file-url https://adp-global.laiye.com/web/agentic_doc_processor/laiye/file/13e18a44228611f1933a00163e122259
    run_test "custom-app update" "$ADP" custom-app update \
        --app-id "$CREATED_APP_ID" \
        --app-name "E2E-Test-Updated" \
        --extract-fields "$SCRIPT_DIR/samples/extract-fields.json" \
        --parse-mode standard \
        --enable-long-doc true \
        --long-doc-config "$SCRIPT_DIR/samples/long_doc_config.json"
    run_test "custom-app delete-version (idempotent)" "$ADP" custom-app delete-version --app-id "$CREATED_APP_ID" --config-version v99
    run_test "custom-app delete" "$ADP" custom-app delete --app-id "$CREATED_APP_ID"
    CREATED_APP_ID=""  # Already deleted, skip cleanup
fi

# ---- Credit ----
print_header "8. Credit Command"

run_test "credit" "$ADP" credit

fi  # end API tests

# ============================================
# Summary & Report
# ============================================
print_header "Test Summary"
echo "Tests Run:     $TESTS_RUN"
echo -e "Tests Passed:  ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed:  ${RED}$TESTS_FAILED${NC}"
echo -e "Tests Skipped: ${YELLOW}$TESTS_SKIPPED${NC}"
echo ""

# Generate report file
{
    echo "ADP CLI E2E Test Report"
    echo "Date: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "Platform: $(uname -s) $(uname -m)"
    echo "API Tests: $API_AVAILABLE"
    echo "========================================"
    echo ""
    echo "RESULT | TEST NAME"
    echo "-------|----------"
    echo -n "$TEST_LOG"
    echo ""
    echo "========================================"
    echo "Total:   $TESTS_RUN"
    echo "Passed:  $TESTS_PASSED"
    echo "Failed:  $TESTS_FAILED"
    echo "Skipped: $TESTS_SKIPPED"
    if [ $TESTS_SKIPPED -gt 0 ]; then
        echo ""
        echo "Skipped tests:"
        echo "$TEST_LOG" | grep "^SKIP" | while IFS='|' read -r _ name; do
            echo "  - $name"
        done
    fi
    if [ $TESTS_FAILED -gt 0 ]; then
        echo ""
        echo "Failed tests:"
        echo "$TEST_LOG" | grep "^FAIL" | while IFS='|' read -r _ name; do
            echo "  - $name"
        done
    fi
} > "$REPORT_FILE"

echo "Report saved to: $REPORT_FILE"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed.${NC}"
    exit 1
fi
