// Package types defines common types for mermaid-aa
package types

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
