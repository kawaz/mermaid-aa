package renderer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kawaz/mermaid-aa/parser"
)

// Canvas represents a 2D character canvas for drawing
type Canvas struct {
	Width  int
	Height int
	Data   [][]rune
}

// NewCanvas creates a new canvas with the specified size
func NewCanvas(width, height int) *Canvas {
	data := make([][]rune, height)
	for i := range data {
		data[i] = make([]rune, width)
		for j := range data[i] {
			data[i][j] = ' '
		}
	}
	return &Canvas{Width: width, Height: height, Data: data}
}

// Set sets a character at the specified position
func (c *Canvas) Set(x, y int, ch rune) {
	if x >= 0 && x < c.Width && y >= 0 && y < c.Height {
		c.Data[y][x] = ch
	}
}

// SetString writes a string starting at the specified position
func (c *Canvas) SetString(x, y int, s string) {
	for i, ch := range s {
		c.Set(x+i, y, ch)
	}
}

// String converts the canvas to a string
func (c *Canvas) String() string {
	var sb strings.Builder
	for _, row := range c.Data {
		line := strings.TrimRight(string(row), " ")
		sb.WriteString(line)
		sb.WriteRune('\n')
	}
	return sb.String()
}

// NodePos represents the position of a node
type NodePos struct {
	X      int
	Y      int
	Width  int
	Height int
}

// Renderer renders a graph as ASCII art
type Renderer struct {
	graph      *parser.Graph
	nodePos    map[string]*NodePos
	nodeWidth  int
	nodeHeight int
	hSpacing   int
	vSpacing   int
}

// NewRenderer creates a new renderer for the given graph
func NewRenderer(graph *parser.Graph) *Renderer {
	return &Renderer{
		graph:      graph,
		nodePos:    make(map[string]*NodePos),
		nodeWidth:  14,
		nodeHeight: 3,
		hSpacing:   6,
		vSpacing:   2,
	}
}

// Render renders the graph as ASCII art
func (r *Renderer) Render() string {
	if len(r.graph.Nodes) == 0 {
		return ""
	}

	// Calculate layout
	r.calculateLayout()

	// Determine canvas size
	maxX, maxY := 0, 0
	for _, pos := range r.nodePos {
		if pos.X+pos.Width > maxX {
			maxX = pos.X + pos.Width
		}
		if pos.Y+pos.Height > maxY {
			maxY = pos.Y + pos.Height
		}
	}

	canvas := NewCanvas(maxX+4, maxY+4)

	// Draw edges first (so nodes are on top)
	r.drawEdges(canvas)

	// Draw nodes
	r.drawNodes(canvas)

	return canvas.String()
}

func (r *Renderer) calculateLayout() {
	// Build adjacency list and find root nodes
	outgoing := make(map[string][]string)
	incoming := make(map[string][]string)

	for _, edge := range r.graph.Edges {
		outgoing[edge.From] = append(outgoing[edge.From], edge.To)
		incoming[edge.To] = append(incoming[edge.To], edge.From)
	}

	// Find root nodes (nodes with no incoming edges)
	var roots []string
	for id := range r.graph.Nodes {
		if len(incoming[id]) == 0 {
			roots = append(roots, id)
		}
	}
	sort.Strings(roots)

	// If no roots found, use the first node
	if len(roots) == 0 {
		for id := range r.graph.Nodes {
			roots = append(roots, id)
			break
		}
	}

	// Assign levels using BFS
	levels := make(map[string]int)
	visited := make(map[string]bool)
	queue := make([]string, len(roots))
	copy(queue, roots)

	for _, root := range roots {
		levels[root] = 0
		visited[root] = true
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		children := outgoing[current]
		sort.Strings(children)
		for _, next := range children {
			if !visited[next] {
				visited[next] = true
				levels[next] = levels[current] + 1
				queue = append(queue, next)
			}
		}
	}

	// Handle unvisited nodes
	for id := range r.graph.Nodes {
		if !visited[id] {
			levels[id] = 0
		}
	}

	// Group nodes by level
	levelNodes := make(map[int][]string)
	maxLevel := 0
	for id, level := range levels {
		levelNodes[level] = append(levelNodes[level], id)
		if level > maxLevel {
			maxLevel = level
		}
	}

	// Sort nodes within each level
	for level := range levelNodes {
		sort.Strings(levelNodes[level])
	}

	// Calculate node widths
	nodeWidths := make(map[string]int)
	for id, node := range r.graph.Nodes {
		nodeWidths[id] = max(len(node.Label)+4, r.nodeWidth)
	}

	// Calculate positions based on direction
	isHorizontal := r.graph.Direction == parser.LeftRight || r.graph.Direction == parser.RightLeft

	if isHorizontal {
		// Horizontal layout
		xOffset := 0
		for level := 0; level <= maxLevel; level++ {
			nodes := levelNodes[level]
			maxWidth := 0
			for _, id := range nodes {
				if nodeWidths[id] > maxWidth {
					maxWidth = nodeWidths[id]
				}
			}

			yOffset := 0
			for _, id := range nodes {
				r.nodePos[id] = &NodePos{
					X:      xOffset,
					Y:      yOffset,
					Width:  nodeWidths[id],
					Height: r.nodeHeight,
				}
				yOffset += r.nodeHeight + r.vSpacing
			}
			xOffset += maxWidth + r.hSpacing
		}
	} else {
		// Vertical layout
		yOffset := 0
		for level := 0; level <= maxLevel; level++ {
			nodes := levelNodes[level]

			// Calculate total width of this level
			totalWidth := 0
			for i, id := range nodes {
				totalWidth += nodeWidths[id]
				if i > 0 {
					totalWidth += r.hSpacing / 2
				}
			}

			xOffset := 0
			for _, id := range nodes {
				r.nodePos[id] = &NodePos{
					X:      xOffset,
					Y:      yOffset,
					Width:  nodeWidths[id],
					Height: r.nodeHeight,
				}
				xOffset += nodeWidths[id] + r.hSpacing/2
			}
			yOffset += r.nodeHeight + r.vSpacing
		}
	}
}

func (r *Renderer) drawNodes(canvas *Canvas) {
	for id, pos := range r.nodePos {
		node := r.graph.Nodes[id]
		r.drawNode(canvas, node, pos)
	}
}

func (r *Renderer) drawNode(canvas *Canvas, node *parser.Node, pos *NodePos) {
	x, y := pos.X, pos.Y
	width := pos.Width
	label := node.Label

	// Center the label
	labelX := x + (width-len(label))/2
	labelY := y + 1

	switch node.Shape {
	case parser.Rectangle:
		// Draw rectangle
		canvas.SetString(x, y, "+"+strings.Repeat("-", width-2)+"+")
		canvas.SetString(x, y+1, "|"+strings.Repeat(" ", width-2)+"|")
		canvas.SetString(x, y+2, "+"+strings.Repeat("-", width-2)+"+")
		canvas.SetString(labelX, labelY, label)

	case parser.RoundedRect:
		// Draw rounded rectangle
		canvas.SetString(x, y, "/"+strings.Repeat("-", width-2)+"\\")
		canvas.SetString(x, y+1, "|"+strings.Repeat(" ", width-2)+"|")
		canvas.SetString(x, y+2, "\\"+strings.Repeat("-", width-2)+"/")
		canvas.SetString(labelX, labelY, label)

	case parser.Rhombus:
		// Draw diamond/rhombus
		halfW := width / 2
		canvas.SetString(x+halfW-1, y, "/\\")
		canvas.SetString(x, y+1, "<"+strings.Repeat(" ", width-2)+">")
		canvas.SetString(x+halfW-1, y+2, "\\/")
		canvas.SetString(labelX, labelY, label)

	case parser.Circle:
		// Draw circle
		canvas.SetString(x, y, "("+strings.Repeat("-", width-2)+")")
		canvas.SetString(x, y+1, "("+strings.Repeat(" ", width-2)+")")
		canvas.SetString(x, y+2, "("+strings.Repeat("-", width-2)+")")
		canvas.SetString(labelX, labelY, label)

	case parser.Hexagon:
		// Draw hexagon
		canvas.SetString(x+1, y, "/"+strings.Repeat("-", width-4)+"\\")
		canvas.SetString(x, y+1, "<"+strings.Repeat(" ", width-2)+">")
		canvas.SetString(x+1, y+2, "\\"+strings.Repeat("-", width-4)+"/")
		canvas.SetString(labelX, labelY, label)

	case parser.Stadium:
		// Draw stadium shape
		canvas.SetString(x, y, "("+strings.Repeat("=", width-2)+")")
		canvas.SetString(x, y+1, "|"+strings.Repeat(" ", width-2)+"|")
		canvas.SetString(x, y+2, "("+strings.Repeat("=", width-2)+")")
		canvas.SetString(labelX, labelY, label)

	default:
		// Default to rectangle
		canvas.SetString(x, y, "+"+strings.Repeat("-", width-2)+"+")
		canvas.SetString(x, y+1, "|"+strings.Repeat(" ", width-2)+"|")
		canvas.SetString(x, y+2, "+"+strings.Repeat("-", width-2)+"+")
		canvas.SetString(labelX, labelY, label)
	}
}

func (r *Renderer) drawEdges(canvas *Canvas) {
	for _, edge := range r.graph.Edges {
		fromPos := r.nodePos[edge.From]
		toPos := r.nodePos[edge.To]
		if fromPos == nil || toPos == nil {
			continue
		}

		r.drawEdge(canvas, edge, fromPos, toPos)
	}
}

func (r *Renderer) drawEdge(canvas *Canvas, edge *parser.Edge, from, to *NodePos) {
	isHorizontal := r.graph.Direction == parser.LeftRight || r.graph.Direction == parser.RightLeft

	if isHorizontal {
		r.drawHorizontalEdge(canvas, edge, from, to)
	} else {
		r.drawVerticalEdge(canvas, edge, from, to)
	}
}

func (r *Renderer) drawHorizontalEdge(canvas *Canvas, edge *parser.Edge, from, to *NodePos) {
	startX := from.X + from.Width
	startY := from.Y + from.Height/2
	endX := to.X
	endY := to.Y + to.Height/2

	if startY == endY {
		// Simple horizontal line
		for x := startX; x < endX; x++ {
			canvas.Set(x, startY, '-')
		}
		canvas.Set(endX-1, endY, '>')
	} else {
		// Need to route around
		midX := (startX + endX) / 2

		// Draw from start to mid
		for x := startX; x <= midX; x++ {
			canvas.Set(x, startY, '-')
		}

		// Draw vertical part
		minY, maxY := startY, endY
		if startY > endY {
			minY, maxY = endY, startY
		}
		for y := minY; y <= maxY; y++ {
			canvas.Set(midX, y, '|')
		}

		// Draw corner characters
		if startY < endY {
			canvas.Set(midX, startY, '+')
			canvas.Set(midX, endY, '+')
		} else {
			canvas.Set(midX, startY, '+')
			canvas.Set(midX, endY, '+')
		}

		// Draw from mid to end
		for x := midX; x < endX; x++ {
			canvas.Set(x, endY, '-')
		}
		canvas.Set(endX-1, endY, '>')
	}

	// Draw edge label if present
	if edge.Label != "" {
		midX := (startX + endX) / 2
		canvas.SetString(midX-len(edge.Label)/2, startY-1, edge.Label)
	}
}

func (r *Renderer) drawVerticalEdge(canvas *Canvas, edge *parser.Edge, from, to *NodePos) {
	startX := from.X + from.Width/2
	startY := from.Y + from.Height
	endX := to.X + to.Width/2
	endY := to.Y

	if startX == endX {
		// Simple vertical line
		for y := startY; y < endY; y++ {
			canvas.Set(startX, y, '|')
		}
		canvas.Set(endX, endY-1, 'v')

		// Draw edge label if present (for straight lines)
		if edge.Label != "" {
			midY := (startY + endY) / 2
			canvas.SetString(startX+2, midY, edge.Label)
		}
	} else {
		// Need to route around
		midY := (startY + endY) / 2

		// Draw from start down to mid
		for y := startY; y <= midY; y++ {
			canvas.Set(startX, y, '|')
		}

		// Draw horizontal part
		minX, maxX := startX, endX
		if startX > endX {
			minX, maxX = endX, startX
		}
		for x := minX; x <= maxX; x++ {
			canvas.Set(x, midY, '-')
		}

		// Draw corner characters
		canvas.Set(startX, midY, '+')
		canvas.Set(endX, midY, '+')

		// Draw from mid down to end
		for y := midY; y < endY; y++ {
			canvas.Set(endX, y, '|')
		}
		canvas.Set(endX, endY-1, 'v')

		// Draw edge label on horizontal segment
		if edge.Label != "" {
			labelX := (minX + maxX) / 2 - len(edge.Label)/2
			canvas.SetString(labelX, midY-1, edge.Label)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// RenderSimple provides a simple rendering for debugging
func RenderSimple(graph *parser.Graph) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Graph Direction: %s\n", graph.Direction))
	sb.WriteString(fmt.Sprintf("Nodes: %d\n", len(graph.Nodes)))
	for id, node := range graph.Nodes {
		sb.WriteString(fmt.Sprintf("  - %s: %q (shape: %d)\n", id, node.Label, node.Shape))
	}
	sb.WriteString(fmt.Sprintf("Edges: %d\n", len(graph.Edges)))
	for _, edge := range graph.Edges {
		label := ""
		if edge.Label != "" {
			label = fmt.Sprintf(" |%s|", edge.Label)
		}
		sb.WriteString(fmt.Sprintf("  - %s --> %s%s\n", edge.From, edge.To, label))
	}

	return sb.String()
}
