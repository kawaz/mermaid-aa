#!/usr/bin/env bash
# mermaid-aa cross-implementation test runner

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INPUTS_DIR="${SCRIPT_DIR}/inputs"
EXPECTED_DIR="${SCRIPT_DIR}/expected"
OUTPUT_DIR="${SCRIPT_DIR}/output"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default implementations to test (can be overridden via environment)
IMPLEMENTATIONS="${IMPLEMENTATIONS:-go rust typescript moonbit}"

# Charsets to test
CHARSETS="${CHARSETS:-ascii unicode}"

# Ambiguous width options (-a option for CJK character width)
AMBIGUOUS_WIDTHS="${AMBIGUOUS_WIDTHS:-}"

usage() {
    cat <<EOF
Usage: $(basename "$0") [OPTIONS]

Run mermaid-aa tests across implementations.

OPTIONS:
    -i, --impl IMPL       Test specific implementation (go|rust|typescript|moonbit)
    -c, --charset SET     Test specific charset (ascii|unicode|unicode-round|unicode-bold|unicode-double)
    -a, --ambiguous W     Test specific ambiguous width (1|2|legacy)
    -A, --all-ambiguous   Test all ambiguous width options (1, 2, legacy)
    -g, --generate        Generate expected output from reference implementation
    -r, --reference IMPL  Reference implementation for generating expected output (default: go)
    -u, --update          Update expected output with actual output
    -v, --verbose         Verbose output
    -h, --help            Show this help

EXAMPLES:
    $(basename "$0")                           # Run all tests
    $(basename "$0") -i go                     # Test only Go implementation
    $(basename "$0") -c unicode                # Test only unicode charset
    $(basename "$0") -a 2                      # Test with ambiguous width = 2
    $(basename "$0") -A                        # Test all ambiguous width options
    $(basename "$0") -g -r rust                # Generate expected from Rust impl
    $(basename "$0") -g -a 2                   # Generate expected with -a 2
    $(basename "$0") -v                        # Verbose output

ENVIRONMENT:
    MERMAID_AA_GO         Path to Go implementation binary
    MERMAID_AA_RUST       Path to Rust implementation binary
    MERMAID_AA_TS         Path to TypeScript implementation (deno run ...)
    MERMAID_AA_MOONBIT    Path to MoonBit implementation

EOF
}

# Get binary path for implementation
get_binary() {
    local impl="$1"
    case "$impl" in
        go)
            echo "${MERMAID_AA_GO:-mermaid-aa-go}"
            ;;
        rust)
            echo "${MERMAID_AA_RUST:-mermaid-aa}"
            ;;
        typescript)
            echo "${MERMAID_AA_TS:-deno run --allow-read mermaid-aa.ts}"
            ;;
        moonbit)
            echo "${MERMAID_AA_MOONBIT:-mermaid-aa-moonbit}"
            ;;
        *)
            echo ""
            ;;
    esac
}

# Check if implementation is available
check_impl() {
    local impl="$1"
    local bin
    bin=$(get_binary "$impl")
    if [[ -z "$bin" ]]; then
        return 1
    fi
    # Check if command exists (first word of bin)
    local cmd
    cmd=$(echo "$bin" | awk '{print $1}')
    command -v "$cmd" >/dev/null 2>&1
}

# Run mermaid-aa with given implementation
run_mermaid_aa() {
    local impl="$1"
    local input_file="$2"
    local charset="$3"
    local ambiguous="${4:-}"
    local bin
    bin=$(get_binary "$impl")

    # Build command with optional -a flag
    local cmd="$bin -c $charset"
    if [[ -n "$ambiguous" ]]; then
        cmd="$cmd -a $ambiguous"
    fi
    cmd="$cmd -f $input_file"

    # Run and capture output
    $cmd 2>/dev/null || true
}

# Compare outputs
compare_outputs() {
    local expected="$1"
    local actual="$2"

    if [[ ! -f "$expected" ]]; then
        echo "SKIP (no expected)"
        return 2
    fi

    if diff -q "$expected" "$actual" >/dev/null 2>&1; then
        echo "PASS"
        return 0
    else
        echo "FAIL"
        return 1
    fi
}

# Main test function
run_tests() {
    local impl="$1"
    local charset="$2"
    local verbose="${3:-false}"
    local ambiguous="${4:-}"

    local pass=0
    local fail=0
    local skip=0

    # Build output directory path
    local output_subdir="${charset}"
    local expected_subdir="${charset}"
    if [[ -n "$ambiguous" ]]; then
        output_subdir="${charset}/ambiguous-${ambiguous}"
        expected_subdir="${charset}/ambiguous-${ambiguous}"
    fi

    mkdir -p "${OUTPUT_DIR}/${impl}/${output_subdir}"

    for input_file in "${INPUTS_DIR}"/*.mmd; do
        local basename
        basename=$(basename "$input_file" .mmd)
        local expected_file="${EXPECTED_DIR}/${expected_subdir}/${basename}.txt"
        local output_file="${OUTPUT_DIR}/${impl}/${output_subdir}/${basename}.txt"

        # Run implementation
        run_mermaid_aa "$impl" "$input_file" "$charset" "$ambiguous" > "$output_file"

        # Compare with expected
        local result
        result=$(compare_outputs "$expected_file" "$output_file")
        local status=$?

        if [[ "$verbose" == "true" ]]; then
            printf "  %-40s %s\n" "$basename" "$result"
        fi

        case $status in
            0) ((pass++)) ;;
            1) ((fail++)) ;;
            2) ((skip++)) ;;
        esac
    done

    echo "Results: $pass passed, $fail failed, $skip skipped"

    if [[ $fail -gt 0 ]]; then
        return 1
    fi
    return 0
}

# Generate expected output from reference implementation
generate_expected() {
    local ref_impl="$1"
    local charset="$2"
    local ambiguous="${3:-}"

    if ! check_impl "$ref_impl"; then
        echo -e "${RED}Error: Reference implementation '$ref_impl' not available${NC}" >&2
        return 1
    fi

    # Build output directory path
    local output_subdir="${charset}"
    if [[ -n "$ambiguous" ]]; then
        output_subdir="${charset}/ambiguous-${ambiguous}"
    fi

    mkdir -p "${EXPECTED_DIR}/${output_subdir}"

    local msg="Generating expected output using $ref_impl for charset $charset"
    if [[ -n "$ambiguous" ]]; then
        msg="$msg (ambiguous=$ambiguous)"
    fi
    echo "$msg..."

    for input_file in "${INPUTS_DIR}"/*.mmd; do
        local basename
        basename=$(basename "$input_file" .mmd)
        local output_file="${EXPECTED_DIR}/${output_subdir}/${basename}.txt"

        run_mermaid_aa "$ref_impl" "$input_file" "$charset" "$ambiguous" > "$output_file"
        echo "  Generated: $basename.txt"
    done

    echo "Done."
}

# Main
main() {
    local specific_impl=""
    local specific_charset=""
    local specific_ambiguous=""
    local all_ambiguous=false
    local generate=false
    local reference_impl="go"
    local update=false
    local verbose=false

    while [[ $# -gt 0 ]]; do
        case "$1" in
            -i|--impl)
                specific_impl="$2"
                shift 2
                ;;
            -c|--charset)
                specific_charset="$2"
                shift 2
                ;;
            -a|--ambiguous)
                specific_ambiguous="$2"
                shift 2
                ;;
            -A|--all-ambiguous)
                all_ambiguous=true
                shift
                ;;
            -g|--generate)
                generate=true
                shift
                ;;
            -r|--reference)
                reference_impl="$2"
                shift 2
                ;;
            -u|--update)
                update=true
                shift
                ;;
            -v|--verbose)
                verbose=true
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                echo "Unknown option: $1" >&2
                usage >&2
                exit 1
                ;;
        esac
    done

    # Determine implementations to test
    local impls
    if [[ -n "$specific_impl" ]]; then
        impls="$specific_impl"
    else
        impls="$IMPLEMENTATIONS"
    fi

    # Determine charsets to test
    local charsets
    if [[ -n "$specific_charset" ]]; then
        charsets="$specific_charset"
    else
        charsets="$CHARSETS"
    fi

    # Determine ambiguous widths to test
    local ambiguous_widths=""
    if [[ "$all_ambiguous" == "true" ]]; then
        ambiguous_widths="1 2 legacy"
    elif [[ -n "$specific_ambiguous" ]]; then
        ambiguous_widths="$specific_ambiguous"
    fi

    # Generate mode
    if [[ "$generate" == "true" ]]; then
        for charset in $charsets; do
            if [[ -n "$ambiguous_widths" ]]; then
                for aw in $ambiguous_widths; do
                    generate_expected "$reference_impl" "$charset" "$aw"
                done
            else
                generate_expected "$reference_impl" "$charset"
            fi
        done
        exit 0
    fi

    # Test mode
    local total_pass=0
    local total_fail=0

    for impl in $impls; do
        if ! check_impl "$impl"; then
            echo -e "${YELLOW}Skipping $impl (not available)${NC}"
            continue
        fi

        echo -e "${GREEN}Testing: $impl${NC}"

        for charset in $charsets; do
            if [[ -n "$ambiguous_widths" ]]; then
                for aw in $ambiguous_widths; do
                    echo "  Charset: $charset (ambiguous=$aw)"
                    if run_tests "$impl" "$charset" "$verbose" "$aw"; then
                        ((total_pass++))
                    else
                        ((total_fail++))
                    fi
                done
            else
                echo "  Charset: $charset"
                if run_tests "$impl" "$charset" "$verbose"; then
                    ((total_pass++))
                else
                    ((total_fail++))
                fi
            fi
        done
    done

    echo ""
    if [[ $total_fail -gt 0 ]]; then
        echo -e "${RED}Some tests failed${NC}"
        exit 1
    else
        echo -e "${GREEN}All tests passed${NC}"
        exit 0
    fi
}

main "$@"
