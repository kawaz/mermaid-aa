package renderer

import (
	"unicode"
)

// StringWidth returns the display width of a string
// considering that CJK characters and some symbols are double-width
func StringWidth(s string) int {
	width := 0
	for _, r := range s {
		width += RuneWidth(r)
	}
	return width
}

// RuneWidth returns the display width of a rune
func RuneWidth(r rune) int {
	// Check for common double-width ranges
	if isDoubleWidth(r) {
		return 2
	}
	return 1
}

// isDoubleWidth checks if a rune is typically displayed as double-width
func isDoubleWidth(r rune) bool {
	// CJK Unified Ideographs
	if r >= 0x4E00 && r <= 0x9FFF {
		return true
	}
	// CJK Unified Ideographs Extension A
	if r >= 0x3400 && r <= 0x4DBF {
		return true
	}
	// CJK Unified Ideographs Extension B-F
	if r >= 0x20000 && r <= 0x2FA1F {
		return true
	}
	// Hiragana
	if r >= 0x3040 && r <= 0x309F {
		return true
	}
	// Katakana
	if r >= 0x30A0 && r <= 0x30FF {
		return true
	}
	// Katakana Phonetic Extensions
	if r >= 0x31F0 && r <= 0x31FF {
		return true
	}
	// Hangul Syllables
	if r >= 0xAC00 && r <= 0xD7AF {
		return true
	}
	// Hangul Jamo
	if r >= 0x1100 && r <= 0x11FF {
		return true
	}
	// Full-width ASCII variants
	if r >= 0xFF01 && r <= 0xFF60 {
		return true
	}
	// Full-width brackets and symbols
	if r >= 0xFF5F && r <= 0xFFE6 {
		return true
	}
	// CJK Symbols and Punctuation
	if r >= 0x3000 && r <= 0x303F {
		return true
	}
	// Box Drawing characters (single-width in most terminals)
	if r >= 0x2500 && r <= 0x257F {
		return false
	}
	// Geometric shapes (varies, but typically single)
	if r >= 0x25A0 && r <= 0x25FF {
		return false
	}
	// Arrows (typically single-width)
	if r >= 0x2190 && r <= 0x21FF {
		return false
	}
	// Check if it's a wide character using unicode properties
	if unicode.Is(unicode.Han, r) {
		return true
	}

	return false
}

// PadRight pads a string to the specified display width
func PadRight(s string, width int) string {
	currentWidth := StringWidth(s)
	if currentWidth >= width {
		return s
	}
	padding := width - currentWidth
	for i := 0; i < padding; i++ {
		s += " "
	}
	return s
}

// PadCenter centers a string within the specified display width
func PadCenter(s string, width int) (string, int) {
	currentWidth := StringWidth(s)
	if currentWidth >= width {
		return s, 0
	}
	padding := width - currentWidth
	leftPad := padding / 2
	rightPad := padding - leftPad

	result := ""
	for i := 0; i < leftPad; i++ {
		result += " "
	}
	result += s
	for i := 0; i < rightPad; i++ {
		result += " "
	}
	return result, leftPad
}
