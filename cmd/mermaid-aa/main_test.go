package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"-h"}, nil, &stdout, &stderr)

	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Error("help output should contain 'Usage:'")
	}
}

func TestRunVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"-V"}, nil, &stdout, &stderr)

	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "mermaid-aa") {
		t.Error("version output should contain 'mermaid-aa'")
	}
}

func TestRunWithArgument(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"graph TD; A-->B"}, nil, &stdout, &stderr)

	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}
	output := stdout.String()
	if !strings.Contains(output, "A") || !strings.Contains(output, "B") {
		t.Errorf("output should contain nodes A and B: %s", output)
	}
}

func TestRunWithStdin(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("graph TD; A[Hello]-->B[World]")
	code := run([]string{}, stdin, &stdout, &stderr)

	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr: %s", code, stderr.String())
	}
	output := stdout.String()
	if !strings.Contains(output, "Hello") || !strings.Contains(output, "World") {
		t.Errorf("output should contain 'Hello' and 'World': %s", output)
	}
}

func TestRunWithCharset(t *testing.T) {
	charsets := []string{"ascii", "unicode", "unicode-round", "unicode-bold", "unicode-double"}

	for _, cs := range charsets {
		t.Run(cs, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run([]string{"-c", cs, "graph TD; A-->B"}, nil, &stdout, &stderr)

			if code != 0 {
				t.Errorf("exit code = %d, want 0 for charset %s", code, cs)
			}
		})
	}
}

func TestRunWithDirection(t *testing.T) {
	directions := []string{"TB", "TD", "BT", "LR", "RL"}

	for _, dir := range directions {
		t.Run(dir, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run([]string{"-d", dir, "graph TD; A-->B"}, nil, &stdout, &stderr)

			if code != 0 {
				t.Errorf("exit code = %d, want 0 for direction %s", code, dir)
			}
		})
	}
}

func TestRunWithAmbiguousWidth(t *testing.T) {
	var stdout1, stderr1 bytes.Buffer
	code1 := run([]string{"-a", "1", "graph TD; A-->B"}, nil, &stdout1, &stderr1)

	var stdout2, stderr2 bytes.Buffer
	code2 := run([]string{"-a", "2", "graph TD; A-->B"}, nil, &stdout2, &stderr2)

	if code1 != 0 {
		t.Errorf("exit code = %d, want 0 for -a 1", code1)
	}
	if code2 != 0 {
		t.Errorf("exit code = %d, want 0 for -a 2", code2)
	}
}

func TestRunNoInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := strings.NewReader("")
	code := run([]string{}, stdin, &stdout, &stderr)

	if code != 1 {
		t.Errorf("exit code = %d, want 1 for empty input", code)
	}
}

func TestRunInvalidOption(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--invalid-option"}, nil, &stdout, &stderr)

	if code != 3 {
		t.Errorf("exit code = %d, want 3 for invalid option", code)
	}
}
