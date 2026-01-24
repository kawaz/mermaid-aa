// mermaid-aa converts Mermaid diagrams to ASCII art
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kawaz/mermaid-aa/parser"
	"github.com/kawaz/mermaid-aa/renderer"
	"github.com/kawaz/mermaid-aa/types"
)

const version = "0.1.0"

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("mermaid-aa", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		fileFlag      = fs.String("f", "", "Read from file")
		charsetFlag   = fs.String("c", "", "Character set [ascii|unicode|unicode-round|unicode-bold|unicode-double]")
		ambWidthFlag  = fs.String("a", "", "Width for ambiguous chars")
		directionFlag = fs.String("d", "", "Override direction [TB|TD|BT|LR|RL]")
		helpFlag      = fs.Bool("h", false, "Print help")
		versionFlag   = fs.Bool("V", false, "Print version")
	)

	// Also support long flags
	fs.StringVar(fileFlag, "file", "", "Read from file")
	fs.StringVar(charsetFlag, "charset", "", "Character set")
	fs.StringVar(ambWidthFlag, "ambiguous-width", "", "Width for ambiguous chars")
	fs.StringVar(directionFlag, "direction", "", "Override direction")
	fs.BoolVar(helpFlag, "help", false, "Print help")
	fs.BoolVar(versionFlag, "version", false, "Print version")

	if err := fs.Parse(args); err != nil {
		return 3
	}

	if *helpFlag {
		printHelp(stdout)
		return 0
	}

	if *versionFlag {
		fmt.Fprintf(stdout, "mermaid-aa %s\n", version)
		return 0
	}

	// Determine input source (priority: file > argument > stdin)
	var input string
	var err error

	if *fileFlag != "" {
		data, err := os.ReadFile(*fileFlag)
		if err != nil {
			fmt.Fprintf(stderr, "Error: cannot read file %s: %v\n", *fileFlag, err)
			return 2
		}
		input = string(data)
	} else if fs.NArg() > 0 {
		input = fs.Arg(0)
	} else {
		// Read from stdin
		reader := bufio.NewReader(stdin)
		var sb strings.Builder
		for {
			line, err := reader.ReadString('\n')
			sb.WriteString(line)
			if err != nil {
				break
			}
		}
		input = sb.String()
	}

	if strings.TrimSpace(input) == "" {
		fmt.Fprintln(stderr, "Error: no input provided")
		return 1
	}

	// Parse the flowchart
	filename := "input"
	if *fileFlag != "" {
		filename = *fileFlag
	}
	fc, err := parser.ParseFlowchartWithFilename(input, filename)
	if err != nil {
		// Check if it's a ParseErrors for formatted output
		if parseErrs, ok := err.(*types.ParseErrors); ok {
			fmt.Fprint(stderr, parseErrs.Format())
		} else {
			fmt.Fprintf(stderr, "Error: %v\n", err)
		}
		return 1
	}

	// Apply options
	opts := renderer.DefaultOptions()

	// Charset from environment or flag
	if env := os.Getenv("MERMAID_AA_CHARSET"); env != "" {
		opts.Charset = env
	}
	if *charsetFlag != "" {
		opts.Charset = *charsetFlag
	}

	// Ambiguous width from environment or flag
	if env := os.Getenv("MERMAID_AA_AMBIGUOUS_WIDTH"); env != "" {
		opts.AmbiguousWidthMode = env
	}
	if *ambWidthFlag != "" {
		opts.AmbiguousWidthMode = *ambWidthFlag
	}

	// Direction override
	if *directionFlag != "" {
		opts.Direction = types.ParseDirection(*directionFlag)
		fc.Direction = opts.Direction
	}

	// Render
	result := renderer.RenderFlowchart(fc, opts)
	fmt.Fprintln(stdout, result)

	return 0
}

func printHelp(w io.Writer) {
	help := `mermaid-aa - Convert Mermaid diagrams to ASCII art

Usage:
  mermaid-aa [OPTIONS] [INPUT]

Arguments:
  [INPUT]  Mermaid text or file path

Options:
  -f, --file <FILE>           Read from file
  -c, --charset <CHARSET>     Character set [default: unicode]
                              [ascii|unicode|unicode-round|unicode-bold|unicode-double]
  -a, --ambiguous-width <W>   Width for East Asian Ambiguous characters [default: 1]

                              Values:
                                1, half     Half-width (1 cell) - for Western terminals
                                2, full     Full-width (2 cells) - for CJK terminals
                                console     Ambiguous=1, Box Drawing=1 (recommended)
                                legacy      All ambiguous including Box Drawing=2
  -d, --direction <DIR>       Override direction [TB|TD|BT|LR|RL]
  -h, --help                  Print help
  -V, --version               Print version

Input methods (priority order):
  1. -f/--file option
  2. [INPUT] argument
  3. stdin (pipe or interactive)

Environment variables:
  MERMAID_AA_CHARSET          Default character set
  MERMAID_AA_AMBIGUOUS_WIDTH  Ambiguous width setting

Examples:
  mermaid-aa 'graph TD; A-->B'
  mermaid-aa -f diagram.mmd
  echo 'graph LR; A-->B' | mermaid-aa
  mermaid-aa -c ascii -a full 'graph TD; A[日本語]-->B'
`
	fmt.Fprint(w, help)
}
