package renderer

import (
	"strings"

	"github.com/kawaz/mermaid-aa/width"
)

// Canvas represents a 2D drawing surface
type Canvas struct {
	cells  [][]rune
	width  int
	height int
	calc   *width.Calculator
}

// NewCanvas creates a new canvas with the given dimensions
func NewCanvas(w, h int, calc *width.Calculator) *Canvas {
	cells := make([][]rune, h)
	for i := range cells {
		cells[i] = make([]rune, w)
		for j := range cells[i] {
			cells[i][j] = ' '
		}
	}
	return &Canvas{
		cells:  cells,
		width:  w,
		height: h,
		calc:   calc,
	}
}

// Set sets a character at the given position
func (c *Canvas) Set(x, y int, r rune) {
	if x >= 0 && x < c.width && y >= 0 && y < c.height {
		c.cells[y][x] = r
	}
}

// Get gets a character at the given position
func (c *Canvas) Get(x, y int) rune {
	if x >= 0 && x < c.width && y >= 0 && y < c.height {
		return c.cells[y][x]
	}
	return ' '
}

// SetString sets a string at the given position
func (c *Canvas) SetString(x, y int, s string) {
	for _, r := range s {
		if x >= 0 && x < c.width && y >= 0 && y < c.height {
			c.cells[y][x] = r
			rw := c.calc.RuneWidth(r)
			// For full-width characters, mark the next cell as a continuation
			// by setting it to a zero-width marker (we use 0)
			if rw == 2 && x+1 < c.width {
				c.cells[y][x+1] = 0
			}
			// Advance x by the display width of the rune
			x += rw
		} else {
			break
		}
	}
}

// HLine draws a horizontal line
func (c *Canvas) HLine(x1, x2, y int, r rune) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	rw := c.calc.RuneWidth(r)
	for x := x1; x <= x2; x += rw {
		c.Set(x, y, r)
	}
}

// VLine draws a vertical line
func (c *Canvas) VLine(x, y1, y2 int, r rune) {
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for y := y1; y <= y2; y++ {
		c.Set(x, y, r)
	}
}

// String returns the canvas as a string
func (c *Canvas) String() string {
	var sb strings.Builder
	for _, row := range c.cells {
		var lineBuf strings.Builder
		for _, r := range row {
			// Skip zero-width markers (continuation of full-width chars)
			if r == 0 {
				continue
			}
			lineBuf.WriteRune(r)
		}
		line := strings.TrimRight(lineBuf.String(), " ")
		sb.WriteString(line)
		sb.WriteRune('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Width returns the canvas width
func (c *Canvas) Width() int {
	return c.width
}

// Height returns the canvas height
func (c *Canvas) Height() int {
	return c.height
}
