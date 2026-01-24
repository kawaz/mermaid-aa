package width

import (
	"testing"
)

func TestStringWidth(t *testing.T) {
	calc := NewCalculator(1)

	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 5},
		{"日本語", 6},     // 3 full-width chars
		{"hello日本", 9}, // 5 + 4
		{"", 0},
		{"a", 1},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := calc.StringWidth(tt.input)
			if got != tt.expected {
				t.Errorf("StringWidth(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestAmbiguousWidth(t *testing.T) {
	// Test with ambiguous width = 1 (half-width)
	calc1 := NewCalculator(1)
	// Test with ambiguous width = 2 (full-width) - legacy behavior
	calc2 := NewCalculator(2)

	// Box-drawing characters are in the Ambiguous category
	boxChar := "─"

	w1 := calc1.StringWidth(boxChar)
	w2 := calc2.StringWidth(boxChar)

	if w1 != 1 {
		t.Errorf("with ambiguous=1, width of %q = %d, want 1", boxChar, w1)
	}
	if w2 != 2 {
		t.Errorf("with ambiguous=2 (legacy), width of %q = %d, want 2", boxChar, w2)
	}
}

func TestWidthModes(t *testing.T) {
	boxChar := "─"      // Box Drawing: U+2500
	ambChar := "α"      // Greek alpha: ambiguous
	fullChar := "日"    // CJK: always full-width

	tests := []struct {
		mode          string
		boxWidth      int
		ambWidth      int
		fullWidth     int
	}{
		{"1", 1, 1, 2},
		{"half", 1, 1, 2},
		{"console", 1, 1, 2},
		{"2", 1, 2, 2},       // full: ambiguous=2, box=1
		{"full", 1, 2, 2},    // full: ambiguous=2, box=1
		{"legacy", 2, 2, 2},  // legacy: all=2
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			calc := NewCalculatorWithMode(tt.mode)

			if got := calc.RuneWidth([]rune(boxChar)[0]); got != tt.boxWidth {
				t.Errorf("mode=%s, box drawing width = %d, want %d", tt.mode, got, tt.boxWidth)
			}
			if got := calc.RuneWidth([]rune(ambChar)[0]); got != tt.ambWidth {
				t.Errorf("mode=%s, ambiguous width = %d, want %d", tt.mode, got, tt.ambWidth)
			}
			if got := calc.RuneWidth([]rune(fullChar)[0]); got != tt.fullWidth {
				t.Errorf("mode=%s, full-width = %d, want %d", tt.mode, got, tt.fullWidth)
			}
		})
	}
}

func TestBoxDrawingRange(t *testing.T) {
	// Verify box drawing detection for U+2500-U+257F
	testCases := []struct {
		char     rune
		expected bool
	}{
		{'─', true},   // U+2500
		{'│', true},   // U+2502
		{'┌', true},   // U+250C
		{'┐', true},   // U+2510
		{'└', true},   // U+2514
		{'┘', true},   // U+2518
		{'├', true},   // U+251C
		{'┤', true},   // U+2524
		{'┬', true},   // U+252C
		{'┴', true},   // U+2534
		{'┼', true},   // U+253C
		{'╔', true},   // U+2554
		{'╗', true},   // U+2557
		{'╚', true},   // U+255A
		{'╝', true},   // U+255D
		{'\u257F', true},  // U+257F (end of range)
		{'\u2580', false}, // U+2580 (outside range)
		{'a', false},
		{'日', false},
	}

	for _, tc := range testCases {
		if got := isBoxDrawing(tc.char); got != tc.expected {
			t.Errorf("isBoxDrawing(%q U+%04X) = %v, want %v", tc.char, tc.char, got, tc.expected)
		}
	}
}

func TestPadRight(t *testing.T) {
	calc := NewCalculator(1)

	tests := []struct {
		input    string
		target   int
		expected string
	}{
		{"hello", 10, "hello     "},
		{"hello", 5, "hello"},
		{"hello", 3, "hello"}, // already longer
		{"日本", 6, "日本  "},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := calc.PadRight(tt.input, tt.target)
			if got != tt.expected {
				t.Errorf("PadRight(%q, %d) = %q, want %q", tt.input, tt.target, got, tt.expected)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	calc := NewCalculator(1)

	tests := []struct {
		input    string
		target   int
		expected string
	}{
		{"hello", 10, "     hello"},
		{"hello", 5, "hello"},
		{"hello", 3, "hello"}, // already longer
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := calc.PadLeft(tt.input, tt.target)
			if got != tt.expected {
				t.Errorf("PadLeft(%q, %d) = %q, want %q", tt.input, tt.target, got, tt.expected)
			}
		})
	}
}

func TestPadCenter(t *testing.T) {
	calc := NewCalculator(1)

	tests := []struct {
		input    string
		target   int
		expected string
	}{
		{"hi", 6, "  hi  "},
		{"hello", 5, "hello"},
		{"a", 4, " a  "}, // odd padding: left gets less
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := calc.PadCenter(tt.input, tt.target)
			if got != tt.expected {
				t.Errorf("PadCenter(%q, %d) = %q, want %q", tt.input, tt.target, got, tt.expected)
			}
		})
	}
}
