// Package types defines common types for mermaid-aa
package types

import (
	"fmt"
	"strings"
)

// Position represents a position in the source code
type Position struct {
	Line   int // 1-indexed
	Column int // 1-indexed
}

// ParseError represents a parse error with location information
type ParseError struct {
	Position Position
	Message  string
	Source   string // the line content where error occurred
	Hint     string // optional help message
	Length   int    // length of the error span for underlining
}

// Error implements the error interface
func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d:%d: %s", e.Position.Line, e.Position.Column, e.Message)
}

// Format returns a formatted error message with visual display
func (e *ParseError) Format(filename string) string {
	var sb strings.Builder

	// Error header
	sb.WriteString(fmt.Sprintf("Error: %s\n", e.Message))
	sb.WriteString(fmt.Sprintf(" --> %s:%d:%d\n", filename, e.Position.Line, e.Position.Column))

	// Line number width
	lineNumStr := fmt.Sprintf("%d", e.Position.Line)
	padding := strings.Repeat(" ", len(lineNumStr))

	// Empty line with bar
	sb.WriteString(fmt.Sprintf("%s |\n", padding))

	// Source line
	sb.WriteString(fmt.Sprintf("%d | %s\n", e.Position.Line, e.Source))

	// Error pointer line
	if e.Length <= 0 {
		e.Length = 1
	}
	spaces := strings.Repeat(" ", e.Position.Column-1)
	carets := strings.Repeat("^", e.Length)
	sb.WriteString(fmt.Sprintf("%s | %s%s\n", padding, spaces, carets))

	// Help hint
	if e.Hint != "" {
		sb.WriteString(fmt.Sprintf("%s |\n", padding))
		sb.WriteString(fmt.Sprintf("%s = help: %s\n", padding, e.Hint))
	}

	return sb.String()
}

// ParseErrors holds multiple parse errors
type ParseErrors struct {
	Errors   []*ParseError
	Filename string
}

// Error implements the error interface
func (e *ParseErrors) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	}
	return fmt.Sprintf("%d errors found", len(e.Errors))
}

// Format returns all formatted error messages
func (e *ParseErrors) Format() string {
	var sb strings.Builder
	for i, err := range e.Errors {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(err.Format(e.Filename))
	}
	return sb.String()
}

// HasErrors returns true if there are any errors
func (e *ParseErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// Add adds a new error
func (e *ParseErrors) Add(err *ParseError) {
	e.Errors = append(e.Errors, err)
}

// Direction represents the flowchart direction
type Direction int

const (
	TB Direction = iota // Top to Bottom (default)
	TD                  // Top Down (same as TB)
	BT                  // Bottom to Top
	LR                  // Left to Right
	RL                  // Right to Left
)

// ParseDirection parses a direction string
func ParseDirection(s string) Direction {
	switch s {
	case "TB", "":
		return TB
	case "TD":
		return TD
	case "BT":
		return BT
	case "LR":
		return LR
	case "RL":
		return RL
	default:
		return TB
	}
}

func (d Direction) String() string {
	switch d {
	case TB:
		return "TB"
	case TD:
		return "TD"
	case BT:
		return "BT"
	case LR:
		return "LR"
	case RL:
		return "RL"
	default:
		return "TB"
	}
}

// IsVertical returns true if direction is vertical
func (d Direction) IsVertical() bool {
	return d == TB || d == TD || d == BT
}

// NodeShape represents the shape of a node
type NodeShape int

const (
	Rectangle     NodeShape = iota // A[text]
	Rounded                        // A(text)
	Stadium                        // A([text])
	Diamond                        // A{text}
	Hexagon                        // A{{text}}
	Circle                         // A((text))
	Parallelogram                  // A[/text/]
	Trapezoid                      // A[/text\]
)

func (s NodeShape) String() string {
	names := []string{
		"rectangle", "rounded", "stadium", "diamond",
		"hexagon", "circle", "parallelogram", "trapezoid",
	}
	if int(s) < len(names) {
		return names[s]
	}
	return "unknown"
}

// EdgeStyle represents the style of an edge
type EdgeStyle int

const (
	SolidArrow  EdgeStyle = iota // -->
	DottedArrow                  // -.->
	ThickArrow                   // ==>
	Open                         // ---
)

// Node represents a flowchart node
type Node struct {
	ID    string
	Label string
	Shape NodeShape
}

// Edge represents a connection between nodes
type Edge struct {
	From  string
	To    string
	Label string
	Style EdgeStyle
}

// Flowchart represents a parsed flowchart
type Flowchart struct {
	Direction Direction
	Nodes     map[string]*Node
	Edges     []*Edge
	NodeOrder []string // maintain insertion order
}

// NewFlowchart creates a new empty flowchart
func NewFlowchart() *Flowchart {
	return &Flowchart{
		Direction: TB,
		Nodes:     make(map[string]*Node),
		Edges:     make([]*Edge, 0),
		NodeOrder: make([]string, 0),
	}
}

// AddNode adds a node to the flowchart
func (f *Flowchart) AddNode(node *Node) {
	if _, exists := f.Nodes[node.ID]; !exists {
		f.NodeOrder = append(f.NodeOrder, node.ID)
	}
	f.Nodes[node.ID] = node
}

// AddEdge adds an edge to the flowchart
func (f *Flowchart) AddEdge(edge *Edge) {
	f.Edges = append(f.Edges, edge)
	// Ensure nodes exist
	if _, exists := f.Nodes[edge.From]; !exists {
		f.AddNode(&Node{ID: edge.From, Label: edge.From, Shape: Rectangle})
	}
	if _, exists := f.Nodes[edge.To]; !exists {
		f.AddNode(&Node{ID: edge.To, Label: edge.To, Shape: Rectangle})
	}
}
