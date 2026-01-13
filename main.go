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
)

var version = "0.1.0"

func main() {
	var (
		showVersion bool
		debugMode   bool
		inputFile   string
	)

	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showVersion, "v", false, "Show version (short)")
	flag.BoolVar(&debugMode, "debug", false, "Show debug output")
	flag.BoolVar(&debugMode, "d", false, "Show debug output (short)")
	flag.StringVar(&inputFile, "file", "", "Input file")
	flag.StringVar(&inputFile, "f", "", "Input file (short)")
	flag.Parse()

	if showVersion {
		fmt.Printf("mermaid-aa version %s\n", version)
		os.Exit(0)
	}

	// Read input
	var input string
	var err error

	// Priority: 1. file flag, 2. positional argument, 3. stdin
	if inputFile != "" {
		input, err = readFile(inputFile)
	} else if flag.NArg() > 0 {
		// Check if argument is a file
		arg := flag.Arg(0)
		if _, statErr := os.Stat(arg); statErr == nil {
			input, err = readFile(arg)
		} else {
			// Treat as mermaid text directly
			input = strings.Join(flag.Args(), " ")
		}
	} else {
		// Read from stdin
		input, err = readStdin()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(input) == "" {
		fmt.Fprintf(os.Stderr, "Error: no input provided\n")
		fmt.Fprintf(os.Stderr, "Usage: mermaid-aa [options] [file or mermaid text]\n")
		fmt.Fprintf(os.Stderr, "       echo 'graph TD; A-->B' | mermaid-aa\n")
		os.Exit(1)
	}

	// Parse mermaid text
	graph, err := parser.Parse(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	if debugMode {
		fmt.Print(renderer.RenderSimple(graph))
		fmt.Println("---")
	}

	// Render as ASCII art
	r := renderer.NewRenderer(graph)
	output := r.Render()
	fmt.Print(output)
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(data), nil
}

func readStdin() (string, error) {
	// Check if stdin has data
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// No pipe, return empty
		return "", nil
	}

	reader := bufio.NewReader(os.Stdin)
	var sb strings.Builder

	for {
		line, err := reader.ReadString('\n')
		sb.WriteString(line)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read stdin: %w", err)
		}
	}

	return sb.String(), nil
}
