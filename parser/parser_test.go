package parser

import (
	"testing"

	"github.com/kawaz/mermaid-aa/types"
)

func TestParseDirection(t *testing.T) {
	tests := []struct {
		input    string
		expected types.Direction
	}{
		{"graph TD; A-->B", types.TD},
		{"graph TB; A-->B", types.TB},
		{"graph LR; A-->B", types.LR},
		{"graph RL; A-->B", types.RL},
		{"graph BT; A-->B", types.BT},
		{"flowchart TD; A-->B", types.TD},
		{"graph; A-->B", types.TB}, // default
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			fc, err := ParseFlowchart(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fc.Direction != tt.expected {
				t.Errorf("direction = %v, want %v", fc.Direction, tt.expected)
			}
		})
	}
}

func TestParseNodes(t *testing.T) {
	tests := []struct {
		input  string
		nodeID string
		label  string
		shape  types.NodeShape
	}{
		{"graph TD; A[Hello]", "A", "Hello", types.Rectangle},
		{"graph TD; B(World)", "B", "World", types.Rounded},
		{"graph TD; C([Stadium])", "C", "Stadium", types.Stadium},
		{"graph TD; D{Diamond}", "D", "Diamond", types.Diamond},
		{"graph TD; E{{Hexagon}}", "E", "Hexagon", types.Hexagon},
		{"graph TD; F((Circle))", "F", "Circle", types.Circle},
		{"graph TD; G[/Parallel/]", "G", "Parallel", types.Parallelogram},
		{`graph TD; H[/Trap\]`, "H", "Trap", types.Trapezoid},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			fc, err := ParseFlowchart(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			node, ok := fc.Nodes[tt.nodeID]
			if !ok {
				t.Fatalf("node %s not found", tt.nodeID)
			}
			if node.Label != tt.label {
				t.Errorf("label = %q, want %q", node.Label, tt.label)
			}
			if node.Shape != tt.shape {
				t.Errorf("shape = %v, want %v", node.Shape, tt.shape)
			}
		})
	}
}

func TestParseEdges(t *testing.T) {
	tests := []struct {
		input  string
		fromID string
		toID   string
		style  types.EdgeStyle
	}{
		{"graph TD; A-->B", "A", "B", types.SolidArrow},
		{"graph TD; A-.->B", "A", "B", types.DottedArrow},
		{"graph TD; A==>B", "A", "B", types.ThickArrow},
		{"graph TD; A---B", "A", "B", types.Open},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			fc, err := ParseFlowchart(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(fc.Edges) == 0 {
				t.Fatal("no edges found")
			}
			edge := fc.Edges[0]
			if edge.From != tt.fromID {
				t.Errorf("from = %q, want %q", edge.From, tt.fromID)
			}
			if edge.To != tt.toID {
				t.Errorf("to = %q, want %q", edge.To, tt.toID)
			}
			if edge.Style != tt.style {
				t.Errorf("style = %v, want %v", edge.Style, tt.style)
			}
		})
	}
}

func TestParseChain(t *testing.T) {
	input := "graph TD; A-->B-->C"
	fc, err := ParseFlowchart(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fc.Nodes) != 3 {
		t.Errorf("nodes count = %d, want 3", len(fc.Nodes))
	}
	if len(fc.Edges) != 2 {
		t.Errorf("edges count = %d, want 2", len(fc.Edges))
	}

	// Check edges
	if fc.Edges[0].From != "A" || fc.Edges[0].To != "B" {
		t.Errorf("edge 0 = %s-->%s, want A-->B", fc.Edges[0].From, fc.Edges[0].To)
	}
	if fc.Edges[1].From != "B" || fc.Edges[1].To != "C" {
		t.Errorf("edge 1 = %s-->%s, want B-->C", fc.Edges[1].From, fc.Edges[1].To)
	}
}

func TestParseWithLabels(t *testing.T) {
	input := "graph TD; A[Start]-->B[End]"
	fc, err := ParseFlowchart(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fc.Nodes["A"].Label != "Start" {
		t.Errorf("A label = %q, want %q", fc.Nodes["A"].Label, "Start")
	}
	if fc.Nodes["B"].Label != "End" {
		t.Errorf("B label = %q, want %q", fc.Nodes["B"].Label, "End")
	}
}

func TestParseMultiLine(t *testing.T) {
	input := `graph TD
    A[Start]
    B[Process]
    C[End]
    A-->B
    B-->C`

	fc, err := ParseFlowchart(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(fc.Nodes) != 3 {
		t.Errorf("nodes count = %d, want 3", len(fc.Nodes))
	}
	if len(fc.Edges) != 2 {
		t.Errorf("edges count = %d, want 2", len(fc.Edges))
	}
}

func TestParseErrorMissingNodeAfterArrow(t *testing.T) {
	input := "graph TD; A -->"
	_, err := ParseFlowchart(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	parseErrs, ok := err.(*types.ParseErrors)
	if !ok {
		t.Fatalf("expected *types.ParseErrors, got %T", err)
	}

	if len(parseErrs.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(parseErrs.Errors))
	}

	if parseErrs.Errors[0].Message != "missing node identifier after arrow" {
		t.Errorf("unexpected message: %s", parseErrs.Errors[0].Message)
	}
}

func TestParseErrorUnclosedBracket(t *testing.T) {
	input := "graph TD; A[Hello"
	_, err := ParseFlowchart(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	parseErrs, ok := err.(*types.ParseErrors)
	if !ok {
		t.Fatalf("expected *types.ParseErrors, got %T", err)
	}

	if len(parseErrs.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(parseErrs.Errors))
	}

	if parseErrs.Errors[0].Message != "unclosed bracket '['" {
		t.Errorf("unexpected message: %s", parseErrs.Errors[0].Message)
	}
}

func TestParseErrorInvalidEdgeSyntax(t *testing.T) {
	input := "graph TD; A->B"
	_, err := ParseFlowchart(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	parseErrs, ok := err.(*types.ParseErrors)
	if !ok {
		t.Fatalf("expected *types.ParseErrors, got %T", err)
	}

	if len(parseErrs.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(parseErrs.Errors))
	}

	if parseErrs.Errors[0].Message != "invalid edge syntax" {
		t.Errorf("unexpected message: %s", parseErrs.Errors[0].Message)
	}
}

func TestParseErrorMultiple(t *testing.T) {
	input := `graph TD
    A -->
    B[Unclosed`

	_, err := ParseFlowchart(input)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	parseErrs, ok := err.(*types.ParseErrors)
	if !ok {
		t.Fatalf("expected *types.ParseErrors, got %T", err)
	}

	if len(parseErrs.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(parseErrs.Errors))
	}
}

func TestParseErrorFormat(t *testing.T) {
	_, err := ParseFlowchartWithFilename("graph TD; A -->", "test.mmd")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	parseErrs, ok := err.(*types.ParseErrors)
	if !ok {
		t.Fatalf("expected *types.ParseErrors, got %T", err)
	}

	formatted := parseErrs.Format()
	if formatted == "" {
		t.Fatal("expected formatted error, got empty string")
	}

	// Check that the filename appears in the output
	if !contains(formatted, "test.mmd") {
		t.Errorf("expected filename in output, got: %s", formatted)
	}

	// Check that help hint appears
	if !contains(formatted, "= help:") {
		t.Errorf("expected help hint in output, got: %s", formatted)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
