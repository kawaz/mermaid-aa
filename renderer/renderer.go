package renderer

import (
	"strings"

	"github.com/kawaz/mermaid-aa/types"
	"github.com/kawaz/mermaid-aa/width"
)

// Renderer renders flowcharts as ASCII art
type Renderer struct {
	charset *Charset
	calc    *width.Calculator
}

// Options for rendering
type Options struct {
	Charset        string
	AmbiguousWidth int
	Direction      types.Direction
}

// DefaultOptions returns default rendering options
func DefaultOptions() *Options {
	return &Options{
		Charset:        "unicode",
		AmbiguousWidth: 1,
		Direction:      types.TB,
	}
}

// New creates a new renderer
func New(opts *Options) *Renderer {
	if opts == nil {
		opts = DefaultOptions()
	}
	return &Renderer{
		charset: GetCharset(opts.Charset),
		calc:    width.NewCalculator(opts.AmbiguousWidth),
	}
}

// nodeLayout holds the position and size of a rendered node
type nodeLayout struct {
	x, y  int
	w, h  int
	node  *types.Node
	lines []string
}

// Render renders a flowchart to ASCII art
func (r *Renderer) Render(fc *types.Flowchart) string {
	if len(fc.Nodes) == 0 {
		return ""
	}

	// Calculate layout
	layouts := r.calculateLayout(fc)

	// Determine canvas size
	canvasWidth, canvasHeight := r.calculateCanvasSize(layouts)

	// Create canvas
	canvas := NewCanvas(canvasWidth, canvasHeight, r.calc)

	// Draw edges first (behind nodes)
	r.drawEdges(canvas, fc, layouts)

	// Draw nodes on top
	for _, layout := range layouts {
		r.drawNode(canvas, layout)
	}

	return canvas.String()
}

func (r *Renderer) calculateLayout(fc *types.Flowchart) map[string]*nodeLayout {
	layouts := make(map[string]*nodeLayout)

	// First pass: calculate node sizes
	maxWidth := 0
	maxHeight := 0
	for _, id := range fc.NodeOrder {
		node := fc.Nodes[id]
		lines := r.renderNodeContent(node)
		w, h := r.nodeSize(node, lines)
		layouts[id] = &nodeLayout{
			node:  node,
			w:     w,
			h:     h,
			lines: lines,
		}
		if w > maxWidth {
			maxWidth = w
		}
		if h > maxHeight {
			maxHeight = h
		}
	}

	// Build adjacency info for layout
	children := make(map[string][]string)
	parents := make(map[string][]string)
	for _, edge := range fc.Edges {
		children[edge.From] = append(children[edge.From], edge.To)
		parents[edge.To] = append(parents[edge.To], edge.From)
	}

	// Find roots (nodes with no parents)
	var roots []string
	for _, id := range fc.NodeOrder {
		if len(parents[id]) == 0 {
			roots = append(roots, id)
		}
	}
	if len(roots) == 0 && len(fc.NodeOrder) > 0 {
		roots = []string{fc.NodeOrder[0]}
	}

	// Calculate levels using BFS
	levels := make(map[string]int)
	visited := make(map[string]bool)
	queue := make([]string, 0)

	for _, root := range roots {
		levels[root] = 0
		queue = append(queue, root)
		visited[root] = true
	}

	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		for _, child := range children[id] {
			if !visited[child] {
				levels[child] = levels[id] + 1
				visited[child] = true
				queue = append(queue, child)
			}
		}
	}

	// Assign positions to unvisited nodes
	maxLevel := 0
	for _, level := range levels {
		if level > maxLevel {
			maxLevel = level
		}
	}
	for _, id := range fc.NodeOrder {
		if !visited[id] {
			maxLevel++
			levels[id] = maxLevel
		}
	}

	// Group nodes by level
	levelNodes := make(map[int][]string)
	for id, level := range levels {
		levelNodes[level] = append(levelNodes[level], id)
	}

	// Spacing
	hSpacing := 4
	vSpacing := 2

	// Position nodes based on direction
	isVertical := fc.Direction == types.TB || fc.Direction == types.TD || fc.Direction == types.BT
	isReversed := fc.Direction == types.BT || fc.Direction == types.RL

	if isVertical {
		y := 0
		levelOrder := make([]int, 0, len(levelNodes))
		for l := range levelNodes {
			levelOrder = append(levelOrder, l)
		}
		// Sort levels
		for i := 0; i < len(levelOrder)-1; i++ {
			for j := i + 1; j < len(levelOrder); j++ {
				if levelOrder[i] > levelOrder[j] {
					levelOrder[i], levelOrder[j] = levelOrder[j], levelOrder[i]
				}
			}
		}
		if isReversed {
			// Reverse level order for BT
			for i, j := 0, len(levelOrder)-1; i < j; i, j = i+1, j-1 {
				levelOrder[i], levelOrder[j] = levelOrder[j], levelOrder[i]
			}
		}

		for _, level := range levelOrder {
			nodes := levelNodes[level]
			x := 0
			levelHeight := 0
			for _, id := range nodes {
				layout := layouts[id]
				layout.x = x
				layout.y = y
				x += layout.w + hSpacing
				if layout.h > levelHeight {
					levelHeight = layout.h
				}
			}
			y += levelHeight + vSpacing
		}
	} else {
		// Horizontal layout (LR/RL)
		x := 0
		levelOrder := make([]int, 0, len(levelNodes))
		for l := range levelNodes {
			levelOrder = append(levelOrder, l)
		}
		// Sort levels
		for i := 0; i < len(levelOrder)-1; i++ {
			for j := i + 1; j < len(levelOrder); j++ {
				if levelOrder[i] > levelOrder[j] {
					levelOrder[i], levelOrder[j] = levelOrder[j], levelOrder[i]
				}
			}
		}
		if isReversed {
			// Reverse level order for RL
			for i, j := 0, len(levelOrder)-1; i < j; i, j = i+1, j-1 {
				levelOrder[i], levelOrder[j] = levelOrder[j], levelOrder[i]
			}
		}

		for _, level := range levelOrder {
			nodes := levelNodes[level]
			y := 0
			levelWidth := 0
			for _, id := range nodes {
				layout := layouts[id]
				layout.x = x
				layout.y = y
				y += layout.h + vSpacing
				if layout.w > levelWidth {
					levelWidth = layout.w
				}
			}
			x += levelWidth + hSpacing
		}
	}

	return layouts
}

func (r *Renderer) calculateCanvasSize(layouts map[string]*nodeLayout) (int, int) {
	maxX, maxY := 0, 0
	for _, layout := range layouts {
		if layout.x+layout.w > maxX {
			maxX = layout.x + layout.w
		}
		if layout.y+layout.h > maxY {
			maxY = layout.y + layout.h
		}
	}
	return maxX + 1, maxY + 1
}

func (r *Renderer) renderNodeContent(node *types.Node) []string {
	// For now, single line label
	return []string{node.Label}
}

func (r *Renderer) nodeSize(node *types.Node, lines []string) (int, int) {
	// Calculate content width
	contentWidth := 0
	for _, line := range lines {
		w := r.calc.StringWidth(line)
		if w > contentWidth {
			contentWidth = w
		}
	}

	// Add padding for borders based on shape
	padding := 4 // 2 chars on each side for borders
	switch node.Shape {
	case types.Circle:
		padding = 6 // (( )) = 4 + extra space
	case types.Stadium:
		padding = 6 // ([ ]) = 4 + extra space
	case types.Hexagon:
		padding = 6 // {{ }} = 4 + extra space
	}

	width := contentWidth + padding
	height := len(lines) + 2 // +2 for top and bottom borders

	return width, height
}

func (r *Renderer) drawNode(canvas *Canvas, layout *nodeLayout) {
	x, y := layout.x, layout.y
	w, h := layout.w, layout.h
	cs := r.charset

	// Get the label centered
	label := layout.node.Label
	contentWidth := w - 4 // minus borders
	centeredLabel := r.calc.PadCenter(label, contentWidth)

	switch layout.node.Shape {
	case types.Rectangle:
		// Top border
		canvas.Set(x, y, []rune(cs.TopLeft)[0])
		for i := 1; i < w-1; i++ {
			canvas.Set(x+i, y, []rune(cs.HLine)[0])
		}
		canvas.Set(x+w-1, y, []rune(cs.TopRight)[0])

		// Middle with label
		for row := 1; row < h-1; row++ {
			canvas.Set(x, y+row, []rune(cs.VLine)[0])
			if row == 1 {
				canvas.SetString(x+1, y+row, " "+centeredLabel+" ")
			} else {
				for i := 1; i < w-1; i++ {
					canvas.Set(x+i, y+row, ' ')
				}
			}
			canvas.Set(x+w-1, y+row, []rune(cs.VLine)[0])
		}

		// Bottom border
		canvas.Set(x, y+h-1, []rune(cs.BotLeft)[0])
		for i := 1; i < w-1; i++ {
			canvas.Set(x+i, y+h-1, []rune(cs.HLine)[0])
		}
		canvas.Set(x+w-1, y+h-1, []rune(cs.BotRight)[0])

	case types.Rounded:
		// Top border with rounded corners
		canvas.Set(x, y, []rune(cs.RoundTopLeft)[0])
		for i := 1; i < w-1; i++ {
			canvas.Set(x+i, y, []rune(cs.HLine)[0])
		}
		canvas.Set(x+w-1, y, []rune(cs.RoundTopRight)[0])

		// Middle with label
		for row := 1; row < h-1; row++ {
			canvas.Set(x, y+row, []rune(cs.VLine)[0])
			if row == 1 {
				canvas.SetString(x+1, y+row, " "+centeredLabel+" ")
			}
			canvas.Set(x+w-1, y+row, []rune(cs.VLine)[0])
		}

		// Bottom border
		canvas.Set(x, y+h-1, []rune(cs.RoundBotLeft)[0])
		for i := 1; i < w-1; i++ {
			canvas.Set(x+i, y+h-1, []rune(cs.HLine)[0])
		}
		canvas.Set(x+w-1, y+h-1, []rune(cs.RoundBotRight)[0])

	case types.Stadium:
		// Stadium: ([ label ])
		topLine := "([" + strings.Repeat(string([]rune(cs.HLine)[0]), w-4) + "])"
		botLine := "([" + strings.Repeat(string([]rune(cs.HLine)[0]), w-4) + "])"
		canvas.SetString(x, y, topLine)
		canvas.Set(x, y+1, '(')
		canvas.SetString(x+1, y+1, " "+centeredLabel+" ")
		canvas.Set(x+w-1, y+1, ')')
		canvas.SetString(x, y+h-1, botLine)

	case types.Diamond:
		// Diamond: < label >
		midY := y + h/2
		canvas.Set(x, midY, []rune(cs.DiamondLeft)[0])
		canvas.SetString(x+1, midY, " "+centeredLabel+" ")
		canvas.Set(x+w-1, midY, []rune(cs.DiamondRight)[0])
		// Draw borders
		for i := 1; i < w-1; i++ {
			canvas.Set(x+i, y, []rune(cs.HLine)[0])
			canvas.Set(x+i, y+h-1, []rune(cs.HLine)[0])
		}

	case types.Hexagon:
		// Hexagon: {{ label }}
		canvas.SetString(x, y, "{{"+strings.Repeat(string([]rune(cs.HLine)[0]), w-4)+"}}")
		canvas.Set(x, y+1, '{')
		canvas.SetString(x+1, y+1, " "+centeredLabel+" ")
		canvas.Set(x+w-1, y+1, '}')
		canvas.SetString(x, y+h-1, "{{"+strings.Repeat(string([]rune(cs.HLine)[0]), w-4)+"}}")

	case types.Circle:
		// Circle: (( label ))
		canvas.SetString(x, y, "(("+strings.Repeat(string([]rune(cs.HLine)[0]), w-4)+"))")
		canvas.Set(x, y+1, '(')
		canvas.Set(x+1, y+1, '(')
		canvas.SetString(x+2, y+1, centeredLabel)
		canvas.Set(x+w-2, y+1, ')')
		canvas.Set(x+w-1, y+1, ')')
		canvas.SetString(x, y+h-1, "(("+strings.Repeat(string([]rune(cs.HLine)[0]), w-4)+"))")

	case types.Parallelogram:
		// Parallelogram: / label /
		canvas.Set(x, y, '/')
		for i := 1; i < w-1; i++ {
			canvas.Set(x+i, y, []rune(cs.HLine)[0])
		}
		canvas.Set(x+w-1, y, '/')
		canvas.Set(x, y+1, '/')
		canvas.SetString(x+1, y+1, " "+centeredLabel+" ")
		canvas.Set(x+w-1, y+1, '/')
		canvas.Set(x, y+h-1, '/')
		for i := 1; i < w-1; i++ {
			canvas.Set(x+i, y+h-1, []rune(cs.HLine)[0])
		}
		canvas.Set(x+w-1, y+h-1, '/')

	case types.Trapezoid:
		// Trapezoid: / label \
		canvas.Set(x, y, '/')
		for i := 1; i < w-1; i++ {
			canvas.Set(x+i, y, []rune(cs.HLine)[0])
		}
		canvas.Set(x+w-1, y, '\\')
		canvas.Set(x, y+1, '/')
		canvas.SetString(x+1, y+1, " "+centeredLabel+" ")
		canvas.Set(x+w-1, y+1, '\\')
		canvas.Set(x, y+h-1, '\\')
		for i := 1; i < w-1; i++ {
			canvas.Set(x+i, y+h-1, []rune(cs.HLine)[0])
		}
		canvas.Set(x+w-1, y+h-1, '/')
	}
}

func (r *Renderer) drawEdges(canvas *Canvas, fc *types.Flowchart, layouts map[string]*nodeLayout) {
	cs := r.charset
	isVertical := fc.Direction == types.TB || fc.Direction == types.TD || fc.Direction == types.BT

	for _, edge := range fc.Edges {
		fromLayout := layouts[edge.From]
		toLayout := layouts[edge.To]
		if fromLayout == nil || toLayout == nil {
			continue
		}

		// Calculate connection points
		var fromX, fromY, toX, toY int

		if isVertical {
			// Vertical layout: connect bottom of from to top of to
			fromX = fromLayout.x + fromLayout.w/2
			fromY = fromLayout.y + fromLayout.h
			toX = toLayout.x + toLayout.w/2
			toY = toLayout.y - 1
		} else {
			// Horizontal layout: connect right of from to left of to
			fromX = fromLayout.x + fromLayout.w
			fromY = fromLayout.y + fromLayout.h/2
			toX = toLayout.x - 1
			toY = toLayout.y + toLayout.h/2
		}

		// Select line character based on edge style
		hChar := []rune(cs.HLine)[0]
		vChar := []rune(cs.VLine)[0]
		switch edge.Style {
		case types.DottedArrow:
			hChar = []rune(cs.DottedH)[0]
			vChar = []rune(cs.DottedV)[0]
		case types.ThickArrow:
			hChar = []rune(cs.ThickH)[0]
			vChar = []rune(cs.ThickV)[0]
		}

		// Draw the edge
		if isVertical {
			// Draw vertical line
			midY := (fromY + toY) / 2
			canvas.VLine(fromX, fromY, midY, vChar)
			if fromX != toX {
				canvas.HLine(fromX, toX, midY, hChar)
			}
			canvas.VLine(toX, midY, toY, vChar)

			// Draw arrow
			if edge.Style != types.Open {
				arrowChar := []rune(cs.ArrowDown)[0]
				if fc.Direction == types.BT {
					arrowChar = []rune(cs.ArrowUp)[0]
				}
				canvas.Set(toX, toY, arrowChar)
			}
		} else {
			// Draw horizontal line
			midX := (fromX + toX) / 2
			canvas.HLine(fromX, midX, fromY, hChar)
			if fromY != toY {
				canvas.VLine(midX, fromY, toY, vChar)
			}
			canvas.HLine(midX, toX, toY, hChar)

			// Draw arrow
			if edge.Style != types.Open {
				arrowChar := []rune(cs.ArrowRight)[0]
				if fc.Direction == types.RL {
					arrowChar = []rune(cs.ArrowLeft)[0]
				}
				canvas.Set(toX, toY, arrowChar)
			}
		}

		// Draw edge label if present
		if edge.Label != "" {
			labelX := (fromX + toX) / 2
			labelY := (fromY + toY) / 2
			canvas.SetString(labelX-r.calc.StringWidth(edge.Label)/2, labelY, edge.Label)
		}
	}
}

// RenderFlowchart is a convenience function to render a flowchart
func RenderFlowchart(fc *types.Flowchart, opts *Options) string {
	return New(opts).Render(fc)
}
