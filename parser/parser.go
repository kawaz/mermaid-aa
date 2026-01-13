package parser

import (
	"regexp"
	"strings"
)

// Direction represents the graph direction
type Direction string

const (
	TopDown   Direction = "TD"
	TopBottom Direction = "TB"
	BottomTop Direction = "BT"
	LeftRight Direction = "LR"
	RightLeft Direction = "RL"
)

// NodeShape represents the shape of a node
type NodeShape int

const (
	Rectangle NodeShape = iota
	RoundedRect
	Stadium
	Cylinder
	Circle
	Rhombus
	Hexagon
	Parallelogram
	Trapezoid
)

// Node represents a node in the graph
type Node struct {
	ID    string
	Label string
	Shape NodeShape
}

// EdgeStyle represents the style of an edge
type EdgeStyle int

const (
	Arrow EdgeStyle = iota
	DottedArrow
	ThickArrow
	NoArrow
)

// Edge represents a connection between nodes
type Edge struct {
	From  string
	To    string
	Label string
	Style EdgeStyle
}

// Graph represents a parsed mermaid graph
type Graph struct {
	Direction Direction
	Nodes     map[string]*Node
	Edges     []*Edge
}

// Parse parses mermaid text and returns a Graph
func Parse(text string) (*Graph, error) {
	graph := &Graph{
		Direction: TopDown,
		Nodes:     make(map[string]*Node),
		Edges:     make([]*Edge, 0),
	}

	// Split by newlines first, then by semicolons
	lines := strings.Split(text, "\n")
	var statements []string
	for _, line := range lines {
		// Split by semicolons to handle "graph TD; A-->B; B-->C"
		parts := strings.Split(line, ";")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				statements = append(statements, part)
			}
		}
	}

	// Parse graph declaration
	graphDeclRegex := regexp.MustCompile(`(?i)^\s*(graph|flowchart)\s+(TD|TB|BT|LR|RL)\s*$`)

	for _, stmt := range statements {
		if stmt == "" || strings.HasPrefix(stmt, "%%") {
			continue
		}

		// Check for graph declaration
		if matches := graphDeclRegex.FindStringSubmatch(stmt); matches != nil {
			graph.Direction = Direction(strings.ToUpper(matches[2]))
			continue
		}

		// Parse edges and nodes - handle chained edges like A --> B --> C
		parseChainedEdges(graph, stmt)
	}

	return graph, nil
}

// parseChainedEdges handles chained edge definitions like A --> B --> C
func parseChainedEdges(graph *Graph, line string) {
	// Split by arrow patterns while keeping track of edge styles
	// This regex captures: node (arrow) node (arrow) node ...

	// First, let's tokenize the line
	tokens := tokenizeLine(line)

	if len(tokens) < 3 {
		// Not enough tokens for an edge
		return
	}

	// Process pairs of nodes connected by arrows
	i := 0
	for i < len(tokens)-2 {
		fromToken := tokens[i]
		arrowToken := tokens[i+1]
		toToken := tokens[i+2]

		style, label := parseArrow(arrowToken)
		if style == -1 {
			i++
			continue
		}

		fromNode := parseNode(fromToken)
		toNode := parseNode(toToken)

		if graph.Nodes[fromNode.ID] == nil {
			graph.Nodes[fromNode.ID] = fromNode
		}
		if graph.Nodes[toNode.ID] == nil {
			graph.Nodes[toNode.ID] = toNode
		}

		graph.Edges = append(graph.Edges, &Edge{
			From:  fromNode.ID,
			To:    toNode.ID,
			Label: label,
			Style: style,
		})

		i += 2 // Move to next potential edge (toToken becomes fromToken)
	}
}

// tokenizeLine splits a line into nodes and arrows
func tokenizeLine(line string) []string {
	var tokens []string

	// Pattern to match arrows with optional labels
	// Order matters: more specific patterns first
	arrowPattern := regexp.MustCompile(`\s*(-->\s*\|[^|]+\||--\s+[^-]+\s+-->|==>|---|-.->|-\.->|-->)\s*`)

	parts := arrowPattern.Split(line, -1)
	arrows := arrowPattern.FindAllString(line, -1)

	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			tokens = append(tokens, part)
		}
		if i < len(arrows) {
			tokens = append(tokens, strings.TrimSpace(arrows[i]))
		}
	}

	return tokens
}

// parseArrow determines the edge style and label from an arrow token
func parseArrow(arrow string) (EdgeStyle, string) {
	arrow = strings.TrimSpace(arrow)

	// --> |label|
	if matches := regexp.MustCompile(`-->\s*\|([^|]+)\|`).FindStringSubmatch(arrow); matches != nil {
		return Arrow, strings.TrimSpace(matches[1])
	}

	// -- label -->
	if matches := regexp.MustCompile(`--\s*([^-]+)\s*-->`).FindStringSubmatch(arrow); matches != nil {
		return Arrow, strings.TrimSpace(matches[1])
	}

	// Simple arrows
	switch {
	case strings.Contains(arrow, "==>"):
		return ThickArrow, ""
	case strings.Contains(arrow, "-.->") || strings.Contains(arrow, "-..->"):
		return DottedArrow, ""
	case arrow == "---":
		return NoArrow, ""
	case strings.Contains(arrow, "-->"):
		return Arrow, ""
	}

	return -1, ""
}

func parseNode(text string) *Node {
	text = strings.TrimSpace(text)
	node := &Node{
		ID:    text,
		Label: text,
		Shape: Rectangle,
	}

	// Pattern: ID[Label] - Rectangle
	if matches := regexp.MustCompile(`^(\w+)\[([^\]]+)\]$`).FindStringSubmatch(text); matches != nil {
		node.ID = matches[1]
		node.Label = matches[2]
		node.Shape = Rectangle
		return node
	}

	// Pattern: ID([Label]) - Stadium (must check before rounded rect)
	if matches := regexp.MustCompile(`^(\w+)\(\[([^\]]+)\]\)$`).FindStringSubmatch(text); matches != nil {
		node.ID = matches[1]
		node.Label = matches[2]
		node.Shape = Stadium
		return node
	}

	// Pattern: ID((Label)) - Circle (must check before rounded rect)
	if matches := regexp.MustCompile(`^(\w+)\(\(([^)]+)\)\)$`).FindStringSubmatch(text); matches != nil {
		node.ID = matches[1]
		node.Label = matches[2]
		node.Shape = Circle
		return node
	}

	// Pattern: ID(Label) - Rounded rectangle
	if matches := regexp.MustCompile(`^(\w+)\(([^)]+)\)$`).FindStringSubmatch(text); matches != nil {
		node.ID = matches[1]
		node.Label = matches[2]
		node.Shape = RoundedRect
		return node
	}

	// Pattern: ID{{Label}} - Hexagon (must check before rhombus)
	if matches := regexp.MustCompile(`^(\w+)\{\{([^}]+)\}\}$`).FindStringSubmatch(text); matches != nil {
		node.ID = matches[1]
		node.Label = matches[2]
		node.Shape = Hexagon
		return node
	}

	// Pattern: ID{Label} - Rhombus/Diamond
	if matches := regexp.MustCompile(`^(\w+)\{([^}]+)\}$`).FindStringSubmatch(text); matches != nil {
		node.ID = matches[1]
		node.Label = matches[2]
		node.Shape = Rhombus
		return node
	}

	// Pattern: ID>Label] - Asymmetric/Flag
	if matches := regexp.MustCompile(`^(\w+)>([^\]]+)\]$`).FindStringSubmatch(text); matches != nil {
		node.ID = matches[1]
		node.Label = matches[2]
		node.Shape = Parallelogram
		return node
	}

	// Plain ID - just use the text as both ID and label
	if matches := regexp.MustCompile(`^(\w+)$`).FindStringSubmatch(text); matches != nil {
		node.ID = matches[1]
		node.Label = matches[1]
		return node
	}

	return node
}
