// Package width provides Unicode character width calculation with ambiguous width support
package width

import (
	"github.com/mattn/go-runewidth"
)

// Calculator calculates string width with configurable ambiguous width
type Calculator struct {
	ambiguousWidth int
	cond           *runewidth.Condition
}

// NewCalculator creates a new width calculator
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
		ambiguousWidth: ambiguousWidth,
		cond:           cond,
	}
}

// StringWidth returns the display width of a string
func (c *Calculator) StringWidth(s string) int {
	return c.cond.StringWidth(s)
}

// RuneWidth returns the display width of a rune
func (c *Calculator) RuneWidth(r rune) int {
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
