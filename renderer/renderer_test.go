package renderer

import (
	"strings"
	"testing"

	"github.com/kawaz/mermaid-aa/parser"
	"github.com/kawaz/mermaid-aa/types"
	"github.com/kawaz/mermaid-aa/width"
)

func TestRenderSimple(t *testing.T) {
	fc, _ := parser.ParseFlowchart("graph TD; A-->B")
	opts := &Options{
		Charset:        "ascii",
		AmbiguousWidth: 1,
	}
	result := RenderFlowchart(fc, opts)

	// Check that result is not empty
	if result == "" {
		t.Error("render result is empty")
	}

	// Check that it contains node labels
	if !strings.Contains(result, "A") {
		t.Error("result should contain node A")
	}
	if !strings.Contains(result, "B") {
		t.Error("result should contain node B")
	}
}

func TestRenderWithLabels(t *testing.T) {
	fc, _ := parser.ParseFlowchart("graph TD; A[Start]-->B[End]")
	opts := &Options{
		Charset:        "ascii",
		AmbiguousWidth: 1,
	}
	result := RenderFlowchart(fc, opts)

	if !strings.Contains(result, "Start") {
		t.Error("result should contain 'Start'")
	}
	if !strings.Contains(result, "End") {
		t.Error("result should contain 'End'")
	}
}

func TestRenderHorizontal(t *testing.T) {
	fc, _ := parser.ParseFlowchart("graph LR; A-->B-->C")
	opts := &Options{
		Charset:        "ascii",
		AmbiguousWidth: 1,
	}
	result := RenderFlowchart(fc, opts)

	// Check that result is not empty
	if result == "" {
		t.Error("render result is empty")
	}

	// In LR layout, nodes should be on the same or similar lines
	lines := strings.Split(result, "\n")
	if len(lines) < 1 {
		t.Error("result should have at least one line")
	}
}

func TestRenderCharsets(t *testing.T) {
	charsets := []string{"ascii", "unicode", "unicode-round", "unicode-bold", "unicode-double"}

	for _, cs := range charsets {
		t.Run(cs, func(t *testing.T) {
			fc, _ := parser.ParseFlowchart("graph TD; A[Test]-->B[Node]")
			opts := &Options{
				Charset:        cs,
				AmbiguousWidth: 1,
			}
			result := RenderFlowchart(fc, opts)

			if result == "" {
				t.Errorf("render result for charset %s is empty", cs)
			}
		})
	}
}

func TestRenderNodeShapes(t *testing.T) {
	shapes := []struct {
		syntax string
		shape  types.NodeShape
	}{
		{"A[text]", types.Rectangle},
		{"A(text)", types.Rounded},
		{"A([text])", types.Stadium},
		{"A{text}", types.Diamond},
		{"A{{text}}", types.Hexagon},
		{"A((text))", types.Circle},
		{"A[/text/]", types.Parallelogram},
		{`A[/text\]`, types.Trapezoid},
	}

	for _, tt := range shapes {
		t.Run(tt.shape.String(), func(t *testing.T) {
			fc, _ := parser.ParseFlowchart("graph TD; " + tt.syntax)
			opts := &Options{
				Charset:        "ascii",
				AmbiguousWidth: 1,
			}
			result := RenderFlowchart(fc, opts)

			if result == "" {
				t.Errorf("render result for shape %s is empty", tt.shape)
			}
			if !strings.Contains(result, "text") {
				t.Errorf("result should contain 'text' for shape %s", tt.shape)
			}
		})
	}
}

func TestRenderDirections(t *testing.T) {
	directions := []string{"TB", "TD", "BT", "LR", "RL"}

	for _, dir := range directions {
		t.Run(dir, func(t *testing.T) {
			fc, _ := parser.ParseFlowchart("graph " + dir + "; A-->B-->C")
			opts := &Options{
				Charset:        "ascii",
				AmbiguousWidth: 1,
			}
			result := RenderFlowchart(fc, opts)

			if result == "" {
				t.Errorf("render result for direction %s is empty", dir)
			}
		})
	}
}

func TestRenderUnicodeLabels(t *testing.T) {
	fc, _ := parser.ParseFlowchart("graph TD; A[日本語]-->B[テスト]")
	opts := &Options{
		Charset:        "unicode",
		AmbiguousWidth: 1,
	}
	result := RenderFlowchart(fc, opts)

	if !strings.Contains(result, "日本語") {
		t.Error("result should contain '日本語'")
	}
	if !strings.Contains(result, "テスト") {
		t.Error("result should contain 'テスト'")
	}
}

func TestRenderAmbiguousWidth(t *testing.T) {
	fc, _ := parser.ParseFlowchart("graph TD; A[Test]-->B[Node]")

	// Render with ambiguous width = 1
	opts1 := &Options{
		Charset:        "unicode",
		AmbiguousWidth: 1,
	}
	result1 := RenderFlowchart(fc, opts1)

	// Render with ambiguous width = 2
	opts2 := &Options{
		Charset:        "unicode",
		AmbiguousWidth: 2,
	}
	result2 := RenderFlowchart(fc, opts2)

	// Both should produce non-empty output
	if result1 == "" {
		t.Error("result with ambiguous=1 is empty")
	}
	if result2 == "" {
		t.Error("result with ambiguous=2 is empty")
	}
}

func TestCanvasOperations(t *testing.T) {
	calc := width.NewCalculator(1)
	canvas := NewCanvas(10, 5, calc)

	// Test Set and Get
	canvas.Set(0, 0, 'X')
	if canvas.Get(0, 0) != 'X' {
		t.Error("Set/Get failed")
	}

	// Test SetString
	canvas.SetString(2, 1, "Hi")
	if canvas.Get(2, 1) != 'H' {
		t.Error("SetString failed for 'H'")
	}
	if canvas.Get(3, 1) != 'i' {
		t.Error("SetString failed for 'i'")
	}

	// Test bounds
	canvas.Set(-1, 0, 'Y')  // should not panic
	canvas.Set(100, 0, 'Y') // should not panic

	// Test String output
	result := canvas.String()
	if !strings.Contains(result, "X") {
		t.Error("String output should contain 'X'")
	}
}
