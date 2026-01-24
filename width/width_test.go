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
	// Test with ambiguous width = 2 (full-width)
	calc2 := NewCalculator(2)

	// Box-drawing characters are in the Ambiguous category
	boxChar := "─"

	w1 := calc1.StringWidth(boxChar)
	w2 := calc2.StringWidth(boxChar)

	if w1 != 1 {
		t.Errorf("with ambiguous=1, width of %q = %d, want 1", boxChar, w1)
	}
	if w2 != 2 {
		t.Errorf("with ambiguous=2, width of %q = %d, want 2", boxChar, w2)
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
