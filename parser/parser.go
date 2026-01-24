// Package parser provides Mermaid flowchart parsing
package parser

import (
	"regexp"
	"strings"

	"github.com/kawaz/mermaid-aa/types"
)

// Parser parses Mermaid flowchart syntax
type Parser struct{}

// New creates a new parser
func New() *Parser {
	return &Parser{}
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

	// Simple node ID
	nodeIDRe = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*)`)
)

// Parse parses a Mermaid flowchart string
func (p *Parser) Parse(input string) (*types.Flowchart, error) {
	fc := types.NewFlowchart()

	// Normalize input: replace semicolons with newlines
	input = strings.ReplaceAll(input, ";", "\n")
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "%%") {
			continue
		}

		// Check for header
		if matches := headerRe.FindStringSubmatch(line); matches != nil {
			if matches[1] != "" {
				fc.Direction = types.ParseDirection(matches[1])
			}
			continue
		}

		// Parse statements (nodes and edges)
		p.parseStatement(line, fc)
	}

	return fc, nil
}

func (p *Parser) parseStatement(line string, fc *types.Flowchart) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// Try to parse as a chain of nodes and edges: A --> B --> C
	p.parseChain(line, fc)
}

func (p *Parser) parseChain(line string, fc *types.Flowchart) {
	var nodes []*types.Node
	var edges []struct {
		style types.EdgeStyle
		label string
	}

	remaining := line
	for remaining != "" {
		remaining = strings.TrimSpace(remaining)
		if remaining == "" {
			break
		}

		// Try to parse a node
		node, consumed := p.parseNode(remaining)
		if node != nil {
			nodes = append(nodes, node)
			remaining = remaining[consumed:]
			remaining = strings.TrimSpace(remaining)

			// Try to parse an edge after the node
			edge, edgeLabel, edgeConsumed := p.parseEdge(remaining)
			if edgeConsumed > 0 {
				edges = append(edges, struct {
					style types.EdgeStyle
					label string
				}{edge, edgeLabel})
				remaining = remaining[edgeConsumed:]
			} else if remaining != "" {
				// If no edge found and there's remaining text, break
				break
			}
		} else {
			// Try to parse just a node ID
			if matches := nodeIDRe.FindStringSubmatch(remaining); matches != nil {
				id := matches[1]
				nodes = append(nodes, &types.Node{ID: id, Label: id, Shape: types.Rectangle})
				remaining = remaining[len(matches[0]):]
				remaining = strings.TrimSpace(remaining)

				// Try to parse an edge
				edge, edgeLabel, edgeConsumed := p.parseEdge(remaining)
				if edgeConsumed > 0 {
					edges = append(edges, struct {
						style types.EdgeStyle
						label string
					}{edge, edgeLabel})
					remaining = remaining[edgeConsumed:]
				} else if remaining != "" {
					break
				}
			} else {
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

// ParseFlowchart is a convenience function to parse a flowchart
func ParseFlowchart(input string) (*types.Flowchart, error) {
	return New().Parse(input)
}
