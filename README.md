# mermaid-aa Test Cases

Cross-implementation test cases for mermaid-aa (Go, Rust, TypeScript, MoonBit).

## Directory Structure

```
tests/
├── inputs/           # Input Mermaid files (.mmd)
├── expected/         # Expected output by charset
│   ├── ascii/
│   ├── unicode/
│   ├── unicode-round/
│   ├── unicode-bold/
│   └── unicode-double/
├── output/           # Actual output (generated during tests)
└── run_tests.sh      # Test runner script
```

## Test Cases

### Basic Direction (01-04)
- `01_basic_td.mmd` - Top-Down (TD) flowchart
- `02_basic_lr.mmd` - Left-Right (LR) flowchart
- `03_basic_bt.mmd` - Bottom-Top (BT) flowchart
- `04_basic_rl.mmd` - Right-Left (RL) flowchart

### Node Shapes (05-13)
- `05_shape_rectangle.mmd` - `A[text]`
- `06_shape_rounded.mmd` - `A(text)`
- `07_shape_stadium.mmd` - `A([text])`
- `08_shape_diamond.mmd` - `A{text}`
- `09_shape_hexagon.mmd` - `A{{text}}`
- `10_shape_circle.mmd` - `A((text))`
- `11_shape_parallelogram.mmd` - `A[/text/]`
- `12_shape_trapezoid.mmd` - `A[/text\]`
- `13_shapes_mixed.mmd` - Multiple shapes combined

### Edge Styles (14-19)
- `14_edge_solid.mmd` - `-->`
- `15_edge_dotted.mmd` - `-.->`
- `16_edge_thick.mmd` - `==>`
- `17_edge_open.mmd` - `---`
- `18_edge_labeled.mmd` - `-->|label|`
- `19_edges_chain.mmd` - `A-->B-->C-->D`

### CJK Support (20-22)
- `20_cjk_japanese.mmd` - Japanese labels
- `21_cjk_chinese.mmd` - Chinese labels
- `22_cjk_mixed.mmd` - Mixed CJK and ASCII

### Complex Graphs (23-30)
- `23_complex_diamond.mmd` - Branch and merge
- `24_complex_parallel.mmd` - Parallel paths
- `25_complex_cycle.mmd` - Cyclic graph
- `26_flowchart_keyword.mmd` - `flowchart` keyword
- `27_edge_styles_mixed.mmd` - Mixed edge styles
- `28_multiline_simple.mmd` - Multiple disconnected edges
- `29_deep_hierarchy.mmd` - Deep node chain
- `30_wide_graph.mmd` - Wide branching

## Usage

### Run All Tests

```bash
./tests/run_tests.sh
```

### Test Specific Implementation

```bash
./tests/run_tests.sh -i go
./tests/run_tests.sh -i rust
./tests/run_tests.sh -i typescript
./tests/run_tests.sh -i moonbit
```

### Test Specific Charset

```bash
./tests/run_tests.sh -c unicode
./tests/run_tests.sh -c ascii
```

### Generate Expected Output

```bash
# Generate from Go implementation (default reference)
./tests/run_tests.sh -g

# Generate from specific implementation
./tests/run_tests.sh -g -r rust
```

### Verbose Output

```bash
./tests/run_tests.sh -v
```

## Environment Variables

Configure implementation paths:

```bash
export MERMAID_AA_GO="path/to/mermaid-aa-go"
export MERMAID_AA_RUST="path/to/mermaid-aa"
export MERMAID_AA_TS="deno run --allow-read path/to/main.ts"
export MERMAID_AA_MOONBIT="path/to/mermaid-aa-moonbit"
```

## Adding New Test Cases

1. Add input file to `tests/inputs/` with `.mmd` extension
2. Generate expected output: `./tests/run_tests.sh -g -r <reference-impl>`
3. Run tests to verify: `./tests/run_tests.sh -v`

## License

MIT License
