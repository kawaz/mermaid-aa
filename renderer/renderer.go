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
	Data   [][]string // Use string slices to support multi-byte characters
}

// NewCanvas creates a new canvas with the specified size
func NewCanvas(width, height int) *Canvas {
	data := make([][]string, height)
	for i := range data {
		data[i] = make([]string, width)
		for j := range data[i] {
			data[i][j] = " "
		}
	}
	return &Canvas{Width: width, Height: height, Data: data}
}

// Set sets a character at the specified position
func (c *Canvas) Set(x, y int, ch string) {
	if x >= 0 && x < c.Width && y >= 0 && y < c.Height {
		c.Data[y][x] = ch
	}
}

// SetRune sets a rune at the specified position
func (c *Canvas) SetRune(x, y int, ch rune) {
	c.Set(x, y, string(ch))
}

// SetString writes a string starting at the specified position
// Returns the number of cells used (accounting for wide characters)
func (c *Canvas) SetString(x, y int, s string) int {
	pos := x
	for _, ch := range s {
		if pos >= c.Width {
			break
		}
		c.Set(pos, y, string(ch))
		w := RuneWidth(ch)
		// For double-width characters, mark the next cell as empty
		if w == 2 && pos+1 < c.Width {
			c.Set(pos+1, y, "")
		}
		pos += w
	}
	return pos - x
}

// String converts the canvas to a string
func (c *Canvas) String() string {
	var sb strings.Builder
	for _, row := range c.Data {
		line := strings.Join(row, "")
		line = strings.TrimRight(line, " ")
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
	style      *BoxStyle
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
		style:      ASCIIStyle,
	}
}

// SetStyle sets the rendering style
func (r *Renderer) SetStyle(style *BoxStyle) {
	r.style = style
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

	// Calculate node widths using StringWidth for proper multi-byte support
	nodeWidths := make(map[string]int)
	for id, node := range r.graph.Nodes {
		labelWidth := StringWidth(node.Label)
		nodeWidths[id] = max(labelWidth+4, r.nodeWidth)
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
	s := r.style

	// Calculate label position (centered)
	labelWidth := StringWidth(label)
	labelX := x + (width-labelWidth)/2
	labelY := y + 1

	innerWidth := width - 2

	switch node.Shape {
	case parser.Rectangle:
		r.drawBox(canvas, x, y, width, s.TopLeft, s.TopRight, s.BottomLeft, s.BottomRight, s.Horizontal, s.Vertical)
		canvas.SetString(labelX, labelY, label)

	case parser.RoundedRect:
		r.drawBox(canvas, x, y, width, s.RoundTopLeft, s.RoundTopRight, s.RoundBottomLeft, s.RoundBottomRight, s.Horizontal, s.Vertical)
		canvas.SetString(labelX, labelY, label)

	case parser.Rhombus:
		// Draw diamond/rhombus
		halfW := width / 2
		canvas.SetString(x+halfW-1, y, "/\\")
		canvas.Set(x, y+1, s.DiamondLeft)
		for i := 1; i < innerWidth+1; i++ {
			canvas.Set(x+i, y+1, " ")
		}
		canvas.Set(x+width-1, y+1, s.DiamondRight)
		canvas.SetString(x+halfW-1, y+2, "\\/")
		canvas.SetString(labelX, labelY, label)

	case parser.Circle:
		// Draw circle using parentheses
		canvas.Set(x, y, "(")
		for i := 0; i < innerWidth; i++ {
			canvas.Set(x+1+i, y, s.Horizontal)
		}
		canvas.Set(x+width-1, y, ")")

		canvas.Set(x, y+1, "(")
		for i := 0; i < innerWidth; i++ {
			canvas.Set(x+1+i, y+1, " ")
		}
		canvas.Set(x+width-1, y+1, ")")

		canvas.Set(x, y+2, "(")
		for i := 0; i < innerWidth; i++ {
			canvas.Set(x+1+i, y+2, s.Horizontal)
		}
		canvas.Set(x+width-1, y+2, ")")
		canvas.SetString(labelX, labelY, label)

	case parser.Hexagon:
		// Draw hexagon
		canvas.Set(x+1, y, "/")
		for i := 0; i < innerWidth-2; i++ {
			canvas.Set(x+2+i, y, s.Horizontal)
		}
		canvas.Set(x+width-2, y, "\\")

		canvas.Set(x, y+1, "<")
		for i := 0; i < innerWidth; i++ {
			canvas.Set(x+1+i, y+1, " ")
		}
		canvas.Set(x+width-1, y+1, ">")

		canvas.Set(x+1, y+2, "\\")
		for i := 0; i < innerWidth-2; i++ {
			canvas.Set(x+2+i, y+2, s.Horizontal)
		}
		canvas.Set(x+width-2, y+2, "/")
		canvas.SetString(labelX, labelY, label)

	case parser.Stadium:
		// Draw stadium shape
		canvas.Set(x, y, "(")
		for i := 0; i < innerWidth; i++ {
			canvas.Set(x+1+i, y, "=")
		}
		canvas.Set(x+width-1, y, ")")

		canvas.Set(x, y+1, s.Vertical)
		for i := 0; i < innerWidth; i++ {
			canvas.Set(x+1+i, y+1, " ")
		}
		canvas.Set(x+width-1, y+1, s.Vertical)

		canvas.Set(x, y+2, "(")
		for i := 0; i < innerWidth; i++ {
			canvas.Set(x+1+i, y+2, "=")
		}
		canvas.Set(x+width-1, y+2, ")")
		canvas.SetString(labelX, labelY, label)

	default:
		// Default to rectangle
		r.drawBox(canvas, x, y, width, s.TopLeft, s.TopRight, s.BottomLeft, s.BottomRight, s.Horizontal, s.Vertical)
		canvas.SetString(labelX, labelY, label)
	}
}

func (r *Renderer) drawBox(canvas *Canvas, x, y, width int, tl, tr, bl, br, h, v string) {
	innerWidth := width - 2

	// Top line
	canvas.Set(x, y, tl)
	for i := 0; i < innerWidth; i++ {
		canvas.Set(x+1+i, y, h)
	}
	canvas.Set(x+width-1, y, tr)

	// Middle line
	canvas.Set(x, y+1, v)
	for i := 0; i < innerWidth; i++ {
		canvas.Set(x+1+i, y+1, " ")
	}
	canvas.Set(x+width-1, y+1, v)

	// Bottom line
	canvas.Set(x, y+2, bl)
	for i := 0; i < innerWidth; i++ {
		canvas.Set(x+1+i, y+2, h)
	}
	canvas.Set(x+width-1, y+2, br)
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
	s := r.style
	startX := from.X + from.Width
	startY := from.Y + from.Height/2
	endX := to.X
	endY := to.Y + to.Height/2

	if startY == endY {
		// Simple horizontal line
		for x := startX; x < endX-1; x++ {
			canvas.Set(x, startY, s.Horizontal)
		}
		canvas.Set(endX-1, endY, s.ArrowRight)
	} else {
		// Need to route around
		midX := (startX + endX) / 2

		// Draw from start to mid
		for x := startX; x < midX; x++ {
			canvas.Set(x, startY, s.Horizontal)
		}

		// Draw vertical part
		minY, maxY := startY, endY
		if startY > endY {
			minY, maxY = endY, startY
		}
		for y := minY + 1; y < maxY; y++ {
			canvas.Set(midX, y, s.Vertical)
		}

		// Draw corner characters
		if startY < endY {
			canvas.Set(midX, startY, s.TeeDown)
			canvas.Set(midX, endY, s.TeeUp)
		} else {
			canvas.Set(midX, startY, s.TeeUp)
			canvas.Set(midX, endY, s.TeeDown)
		}

		// Draw from mid to end
		for x := midX + 1; x < endX-1; x++ {
			canvas.Set(x, endY, s.Horizontal)
		}
		canvas.Set(endX-1, endY, s.ArrowRight)
	}

	// Draw edge label if present
	if edge.Label != "" {
		midX := (startX + endX) / 2
		labelWidth := StringWidth(edge.Label)
		canvas.SetString(midX-labelWidth/2, startY-1, edge.Label)
	}
}

func (r *Renderer) drawVerticalEdge(canvas *Canvas, edge *parser.Edge, from, to *NodePos) {
	s := r.style
	startX := from.X + from.Width/2
	startY := from.Y + from.Height
	endX := to.X + to.Width/2
	endY := to.Y

	if startX == endX {
		// Simple vertical line
		for y := startY; y < endY-1; y++ {
			canvas.Set(startX, y, s.Vertical)
		}
		canvas.Set(endX, endY-1, s.ArrowDown)

		// Draw edge label if present (for straight lines)
		if edge.Label != "" {
			midY := (startY + endY) / 2
			canvas.SetString(startX+2, midY, edge.Label)
		}
	} else {
		// Need to route around
		midY := (startY + endY) / 2

		// Draw from start down to mid
		for y := startY; y < midY; y++ {
			canvas.Set(startX, y, s.Vertical)
		}

		// Draw horizontal part
		minX, maxX := startX, endX
		if startX > endX {
			minX, maxX = endX, startX
		}
		for x := minX + 1; x < maxX; x++ {
			canvas.Set(x, midY, s.Horizontal)
		}

		// Draw corner characters
		if startX < endX {
			canvas.Set(startX, midY, s.TeeRight)
			canvas.Set(endX, midY, s.TeeLeft)
		} else {
			canvas.Set(startX, midY, s.TeeLeft)
			canvas.Set(endX, midY, s.TeeRight)
		}

		// Draw from mid down to end
		for y := midY + 1; y < endY-1; y++ {
			canvas.Set(endX, y, s.Vertical)
		}
		canvas.Set(endX, endY-1, s.ArrowDown)

		// Draw edge label on horizontal segment
		if edge.Label != "" {
			labelWidth := StringWidth(edge.Label)
			labelX := (minX + maxX) / 2 - labelWidth/2
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
