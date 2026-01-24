// Package width provides Unicode character width calculation with ambiguous width support
package width

import (
	"github.com/mattn/go-runewidth"
)

// WidthMode represents the ambiguous width mode
type WidthMode string

const (
	// ModeHalf sets ambiguous characters to half-width (1 cell)
	ModeHalf WidthMode = "half"
	// ModeFull sets ambiguous characters to full-width (2 cells)
	ModeFull WidthMode = "full"
	// ModeConsole is an alias for half (ambiguous=1, box drawing=1)
	ModeConsole WidthMode = "console"
	// ModeLegacy sets all ambiguous including box drawing to full-width (2 cells)
	ModeLegacy WidthMode = "legacy"
)

// ParseWidthMode parses a width mode string
// Returns the mode and whether box drawing should be treated as full-width
func ParseWidthMode(s string) (ambiguousWidth int, boxDrawingFullWidth bool) {
	switch s {
	case "1", "half", "console":
		return 1, false
	case "2", "full":
		return 2, false
	case "legacy":
		return 2, true
	default:
		return 1, false
	}
}

// Calculator calculates string width with configurable ambiguous width
type Calculator struct {
	ambiguousWidth      int
	boxDrawingFullWidth bool
	cond                *runewidth.Condition
}

// NewCalculator creates a new width calculator
// Deprecated: Use NewCalculatorWithMode instead
func NewCalculator(ambiguousWidth int) *Calculator {
	if ambiguousWidth < 1 {
		ambiguousWidth = 1
	}
	if ambiguousWidth > 2 {
		ambiguousWidth = 2
	}

	cond := runewidth.NewCondition()
	cond.EastAsianWidth = ambiguousWidth == 2

	return &Calculator{
		ambiguousWidth:      ambiguousWidth,
		boxDrawingFullWidth: ambiguousWidth == 2,
		cond:                cond,
	}
}

// NewCalculatorWithMode creates a new width calculator with a specific mode
func NewCalculatorWithMode(mode string) *Calculator {
	ambiguousWidth, boxDrawingFullWidth := ParseWidthMode(mode)

	cond := runewidth.NewCondition()
	cond.EastAsianWidth = ambiguousWidth == 2

	return &Calculator{
		ambiguousWidth:      ambiguousWidth,
		boxDrawingFullWidth: boxDrawingFullWidth,
		cond:                cond,
	}
}

// isBoxDrawing returns true if r is a Box Drawing character (U+2500-U+257F)
func isBoxDrawing(r rune) bool {
	return r >= 0x2500 && r <= 0x257F
}

// StringWidth returns the display width of a string
func (c *Calculator) StringWidth(s string) int {
	// If box drawing has special handling, calculate manually
	if c.ambiguousWidth == 2 && !c.boxDrawingFullWidth {
		width := 0
		for _, r := range s {
			width += c.RuneWidth(r)
		}
		return width
	}
	return c.cond.StringWidth(s)
}

// RuneWidth returns the display width of a rune
func (c *Calculator) RuneWidth(r rune) int {
	// Box drawing characters: treat as half-width unless in legacy mode
	if isBoxDrawing(r) && !c.boxDrawingFullWidth {
		return 1
	}
	return c.cond.RuneWidth(r)
}

// Truncate truncates a string to fit in the given width
func (c *Calculator) Truncate(s string, w int, tail string) string {
	return c.cond.Truncate(s, w, tail)
}

// PadRight pads a string with spaces on the right to reach the target width
func (c *Calculator) PadRight(s string, targetWidth int) string {
	currentWidth := c.StringWidth(s)
	if currentWidth >= targetWidth {
		return s
	}
	padding := targetWidth - currentWidth
	for i := 0; i < padding; i++ {
		s += " "
	}
	return s
}

// PadLeft pads a string with spaces on the left to reach the target width
func (c *Calculator) PadLeft(s string, targetWidth int) string {
	currentWidth := c.StringWidth(s)
	if currentWidth >= targetWidth {
		return s
	}
	padding := targetWidth - currentWidth
	result := ""
	for i := 0; i < padding; i++ {
		result += " "
	}
	return result + s
}

// PadCenter pads a string with spaces on both sides to center it
func (c *Calculator) PadCenter(s string, targetWidth int) string {
	currentWidth := c.StringWidth(s)
	if currentWidth >= targetWidth {
		return s
	}
	total := targetWidth - currentWidth
	left := total / 2
	right := total - left
	result := ""
	for i := 0; i < left; i++ {
		result += " "
	}
	result += s
	for i := 0; i < right; i++ {
		result += " "
	}
	return result
}

// Default calculator with ambiguous width = 1
var Default = NewCalculator(1)

// StringWidth returns the display width of a string using default calculator
func StringWidth(s string) int {
	return Default.StringWidth(s)
}

// RuneWidth returns the display width of a rune using default calculator
func RuneWidth(r rune) int {
	return Default.RuneWidth(r)
}
