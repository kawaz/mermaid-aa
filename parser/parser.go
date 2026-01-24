// Package parser provides Mermaid flowchart parsing
package parser

import (
	"regexp"
	"strings"

	"github.com/kawaz/mermaid-aa/types"
)

// Parser parses Mermaid flowchart syntax
type Parser struct {
	errors     *types.ParseErrors
	lines      []string
	currentLine int
}

// New creates a new parser
func New() *Parser {
	return &Parser{
		errors: &types.ParseErrors{Filename: "input.mmd"},
	}
}

// Regular expressions for parsing
var (
	// Match graph/flowchart with optional direction
	headerRe = regexp.MustCompile(`^(?:graph|flowchart)\s*(?:(TB|TD|BT|LR|RL))?`)

	// Match node definitions with various shapes
	// A[text], A(text), A([text]), A{text}, A{{text}}, A((text)), A[/text/], A[/text\]
	nodePatterns = []struct {
		re    *regexp.Regexp
		shape types.NodeShape
	}{
		{regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\(\(\s*(.+?)\s*\)\)`), types.Circle},            // ((text))
		{regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\(\[\s*(.+?)\s*\]\)`), types.Stadium},           // ([text])
		{regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\{\{\s*(.+?)\s*\}\}`), types.Hexagon},           // {{text}}
		{regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\[\s*/\s*(.+?)\s*\\\s*\]`), types.Trapezoid},    // [/text\]
		{regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\[\s*/\s*(.+?)\s*/\s*\]`), types.Parallelogram}, // [/text/]
		{regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\(\s*(.+?)\s*\)`), types.Rounded},               // (text)
		{regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\{\s*(.+?)\s*\}`), types.Diamond},               // {text}
		{regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\[\s*(.+?)\s*\]`), types.Rectangle},             // [text]
	}

	// Match edges with various styles
	edgePatterns = []struct {
		re    *regexp.Regexp
		style types.EdgeStyle
	}{
		{regexp.MustCompile(`^\s*-\.->(?:\|([^|]*)\|)?`), types.DottedArrow}, // -.->
		{regexp.MustCompile(`^\s*==>(?:\|([^|]*)\|)?`), types.ThickArrow},    // ==>
		{regexp.MustCompile(`^\s*-->(?:\|([^|]*)\|)?`), types.SolidArrow},    // -->
		{regexp.MustCompile(`^\s*---(?:\|([^|]*)\|)?`), types.Open},          // ---
	}

	// Match incomplete edges (for error detection)
	incompleteEdgeRe = regexp.MustCompile(`^\s*(-->|-.->|==>|---)\s*$`)

	// Simple node ID
	nodeIDRe = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)`)

	// Unclosed brackets patterns (for error detection)
	unclosedBracketPatterns = []struct {
		re      *regexp.Regexp
		message string
		hint    string
	}{
		{regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\[([^\]]*$)`), "unclosed bracket '['", "add closing ']' to complete the node definition"},
		{regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\(([^\)]*$)`), "unclosed parenthesis '('", "add closing ')' to complete the node definition"},
		{regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)\{([^\}]*$)`), "unclosed brace '{'", "add closing '}' to complete the node definition"},
	}
)

// SetFilename sets the filename for error messages
func (p *Parser) SetFilename(filename string) {
	p.errors.Filename = filename
}

// Parse parses a Mermaid flowchart string
func (p *Parser) Parse(input string) (*types.Flowchart, error) {
	fc := types.NewFlowchart()
	p.errors = &types.ParseErrors{Filename: p.errors.Filename}

	// Normalize input: replace semicolons with newlines
	input = strings.ReplaceAll(input, ";", "\n")
	p.lines = strings.Split(input, "\n")

	hasHeader := false
	for lineNum, line := range p.lines {
		p.currentLine = lineNum + 1 // 1-indexed
		originalLine := line
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "%%") {
			continue
		}

		// Check for header
		if matches := headerRe.FindStringSubmatch(line); matches != nil {
			hasHeader = true
			if matches[1] != "" {
				fc.Direction = types.ParseDirection(matches[1])
			}
			continue
		}

		// Parse statements (nodes and edges)
		p.parseStatement(line, originalLine, fc)
	}

	// Warn if no header found (not an error, but good to know)
	if !hasHeader && len(p.lines) > 0 {
		// This is just a warning; Mermaid allows implicit graph
	}

	if p.errors.HasErrors() {
		return fc, p.errors
	}

	return fc, nil
}

func (p *Parser) parseStatement(line, originalLine string, fc *types.Flowchart) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// Check for unclosed brackets first
	for _, pat := range unclosedBracketPatterns {
		if matches := pat.re.FindStringSubmatch(line); matches != nil {
			col := strings.Index(originalLine, line) + len(matches[1]) + 1
			p.addError(col, len(line)-len(matches[1]), pat.message, pat.hint, originalLine)
			return
		}
	}

	// Try to parse as a chain of nodes and edges: A --> B --> C
	p.parseChain(line, originalLine, fc)
}

func (p *Parser) parseChain(line, originalLine string, fc *types.Flowchart) {
	var nodes []*types.Node
	var edges []struct {
		style types.EdgeStyle
		label string
	}

	remaining := line
	startOffset := strings.Index(originalLine, line)
	if startOffset < 0 {
		startOffset = 0
	}
	currentOffset := startOffset

	for remaining != "" {
		trimmed := strings.TrimLeft(remaining, " \t")
		currentOffset += len(remaining) - len(trimmed)
		remaining = trimmed
		if remaining == "" {
			break
		}

		// Try to parse a node
		node, consumed := p.parseNode(remaining)
		if node != nil {
			nodes = append(nodes, node)
			remaining = remaining[consumed:]
			currentOffset += consumed
			trimmed = strings.TrimLeft(remaining, " \t")
			currentOffset += len(remaining) - len(trimmed)
			remaining = trimmed

			// Try to parse an edge after the node
			edge, edgeLabel, edgeConsumed := p.parseEdge(remaining)
			if edgeConsumed > 0 {
				// Check if edge is incomplete (no target node)
				afterEdge := strings.TrimSpace(remaining[edgeConsumed:])
				if afterEdge == "" {
					// Edge at end of line without target
					edgeCol := currentOffset + 1
					edgeLen := edgeConsumed
					p.addError(edgeCol, edgeLen, "missing node identifier after arrow",
						"add a target node after the arrow, e.g., 'A --> B'", originalLine)
					return
				}

				edges = append(edges, struct {
					style types.EdgeStyle
					label string
				}{edge, edgeLabel})
				remaining = remaining[edgeConsumed:]
				currentOffset += edgeConsumed
			} else if remaining != "" {
				// Check if this looks like an invalid edge
				if strings.HasPrefix(remaining, "-") || strings.HasPrefix(remaining, "=") {
					p.addError(currentOffset+1, len(remaining), "invalid edge syntax",
						"use '-->', '-.->',  '==>', or '---' for edges", originalLine)
					return
				}
				// If no edge found and there's remaining text, break
				break
			}
		} else {
			// Try to parse just a node ID
			if matches := nodeIDRe.FindStringSubmatch(remaining); matches != nil {
				id := matches[1]
				nodes = append(nodes, &types.Node{ID: id, Label: id, Shape: types.Rectangle})
				remaining = remaining[len(matches[0]):]
				currentOffset += len(matches[0])
				trimmed = strings.TrimLeft(remaining, " \t")
				currentOffset += len(remaining) - len(trimmed)
				remaining = trimmed

				// Try to parse an edge
				edge, edgeLabel, edgeConsumed := p.parseEdge(remaining)
				if edgeConsumed > 0 {
					// Check if edge is incomplete
					afterEdge := strings.TrimSpace(remaining[edgeConsumed:])
					if afterEdge == "" {
						edgeCol := currentOffset + 1
						edgeLen := edgeConsumed
						p.addError(edgeCol, edgeLen, "missing node identifier after arrow",
							"add a target node after the arrow, e.g., 'A --> B'", originalLine)
						return
					}

					edges = append(edges, struct {
						style types.EdgeStyle
						label string
					}{edge, edgeLabel})
					remaining = remaining[edgeConsumed:]
					currentOffset += edgeConsumed
				} else if remaining != "" {
					// Check for invalid syntax after node ID
					if strings.HasPrefix(remaining, "-") || strings.HasPrefix(remaining, "=") {
						p.addError(currentOffset+1, len(remaining), "invalid edge syntax",
							"use '-->', '-.->',  '==>', or '---' for edges", originalLine)
						return
					}
					break
				}
			} else {
				// Unable to parse - report error
				if remaining != "" {
					p.addError(currentOffset+1, len(remaining), "unexpected token",
						"expected a node identifier (e.g., 'A', 'Node1')", originalLine)
				}
				break
			}
		}
	}

	// Add nodes and edges to flowchart
	for _, node := range nodes {
		fc.AddNode(node)
	}

	for i := 0; i < len(edges) && i+1 < len(nodes); i++ {
		fc.AddEdge(&types.Edge{
			From:  nodes[i].ID,
			To:    nodes[i+1].ID,
			Label: edges[i].label,
			Style: edges[i].style,
		})
	}
}

func (p *Parser) parseNode(s string) (*types.Node, int) {
	for _, pat := range nodePatterns {
		if matches := pat.re.FindStringSubmatch(s); matches != nil {
			label := matches[2]
			// Remove quotes if present
			label = strings.Trim(label, `"'`)
			return &types.Node{
				ID:    matches[1],
				Label: label,
				Shape: pat.shape,
			}, len(matches[0])
		}
	}
	return nil, 0
}

func (p *Parser) parseEdge(s string) (types.EdgeStyle, string, int) {
	for _, pat := range edgePatterns {
		if matches := pat.re.FindStringSubmatch(s); matches != nil {
			label := ""
			if len(matches) > 1 {
				label = matches[1]
			}
			return pat.style, label, len(matches[0])
		}
	}
	return types.SolidArrow, "", 0
}

func (p *Parser) addError(col, length int, message, hint, sourceLine string) {
	p.errors.Add(&types.ParseError{
		Position: types.Position{Line: p.currentLine, Column: col},
		Message:  message,
		Source:   sourceLine,
		Hint:     hint,
		Length:   length,
	})
}

// ParseFlowchart is a convenience function to parse a flowchart
func ParseFlowchart(input string) (*types.Flowchart, error) {
	return New().Parse(input)
}

// ParseFlowchartWithFilename parses a flowchart with a specific filename for error messages
func ParseFlowchartWithFilename(input, filename string) (*types.Flowchart, error) {
	parser := New()
	parser.SetFilename(filename)
	return parser.Parse(input)
}
