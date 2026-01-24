/**
 * Tests for ASCII art renderer
 */

import { assertEquals, assertStringIncludes } from "@std/assert";
import { parse } from "./parser.ts";
import { render } from "./renderer.ts";

Deno.test("render - Simple flowchart", () => {
  const input = `flowchart TD
    A[Hello] --> B[World]`;

  const flowchart = parse(input);
  const output = render(flowchart, "unicode");

  // Should contain both labels
  assertStringIncludes(output, "Hello");
  assertStringIncludes(output, "World");

  // Should contain box characters
  assertStringIncludes(output, "┌");
  assertStringIncludes(output, "┐");
  assertStringIncludes(output, "└");
  assertStringIncludes(output, "┘");
});

Deno.test("render - ASCII charset", () => {
  const input = `flowchart TD
    A[Test]`;

  const flowchart = parse(input);
  const output = render(flowchart, "ascii");

  // Should contain ASCII box characters
  assertStringIncludes(output, "+");
  assertStringIncludes(output, "-");
  assertStringIncludes(output, "|");
});

Deno.test("render - Unicode round charset", () => {
  const input = `flowchart TD
    A[Test]`;

  const flowchart = parse(input);
  const output = render(flowchart, "unicode-round");

  // Should contain rounded corners
  assertStringIncludes(output, "╭");
  assertStringIncludes(output, "╮");
  assertStringIncludes(output, "╰");
  assertStringIncludes(output, "╯");
});

Deno.test("render - Unicode bold charset", () => {
  const input = `flowchart TD
    A[Test]`;

  const flowchart = parse(input);
  const output = render(flowchart, "unicode-bold");

  // Should contain bold box characters
  assertStringIncludes(output, "┏");
  assertStringIncludes(output, "┓");
  assertStringIncludes(output, "┗");
  assertStringIncludes(output, "┛");
});

Deno.test("render - Unicode double charset", () => {
  const input = `flowchart TD
    A[Test]`;

  const flowchart = parse(input);
  const output = render(flowchart, "unicode-double");

  // Should contain double box characters
  assertStringIncludes(output, "╔");
  assertStringIncludes(output, "╗");
  assertStringIncludes(output, "╚");
  assertStringIncludes(output, "╝");
});

Deno.test("render - LR direction", () => {
  const input = `flowchart LR
    A[Start] --> B[End]`;

  const flowchart = parse(input);
  const output = render(flowchart, "unicode");

  // Both nodes should be on same or similar lines
  const lines = output.split("\n");
  let startLine = -1;
  let endLine = -1;

  for (let i = 0; i < lines.length; i++) {
    if (lines[i].includes("Start")) startLine = i;
    if (lines[i].includes("End")) endLine = i;
  }

  // In LR mode, labels should be on the same row
  assertEquals(startLine, endLine);
});

Deno.test("render - Circle node shape", () => {
  const input = `flowchart TD
    A((Circle))`;

  const flowchart = parse(input);
  const output = render(flowchart, "unicode");

  // Should contain double parentheses
  assertStringIncludes(output, "((");
  assertStringIncludes(output, "))");
  assertStringIncludes(output, "Circle");
});

Deno.test("render - Diamond node shape", () => {
  const input = `flowchart TD
    A{Decision}`;

  const flowchart = parse(input);
  const output = render(flowchart, "unicode");

  assertStringIncludes(output, "Decision");
});

Deno.test("render - Hexagon node shape", () => {
  const input = `flowchart TD
    A{{Hexagon}}`;

  const flowchart = parse(input);
  const output = render(flowchart, "unicode");

  assertStringIncludes(output, "{{");
  assertStringIncludes(output, "}}");
  assertStringIncludes(output, "Hexagon");
});

Deno.test("render - CJK labels with ambiguous width 1", () => {
  const input = `flowchart TD
    A[日本語]`;

  const flowchart = parse(input);
  const output = render(flowchart, "unicode", 1);

  assertStringIncludes(output, "日本語");
});

Deno.test("render - CJK labels with ambiguous width 2", () => {
  const input = `flowchart TD
    A[日本語]`;

  const flowchart = parse(input);
  const output = render(flowchart, "unicode", 2);

  assertStringIncludes(output, "日本語");
});

Deno.test("render - Empty flowchart returns empty string", () => {
  const input = `flowchart TD`;

  const flowchart = parse(input);
  const output = render(flowchart, "unicode");

  assertEquals(output, "");
});

Deno.test("render - Multiple nodes", () => {
  const input = `flowchart TD
    A[First] --> B[Second]
    B --> C[Third]`;

  const flowchart = parse(input);
  const output = render(flowchart, "unicode");

  assertStringIncludes(output, "First");
  assertStringIncludes(output, "Second");
  assertStringIncludes(output, "Third");
});

Deno.test("render - Stadium shape", () => {
  const input = `flowchart TD
    A([Stadium])`;

  const flowchart = parse(input);
  const output = render(flowchart, "unicode");

  assertStringIncludes(output, "([");
  assertStringIncludes(output, "])");
  assertStringIncludes(output, "Stadium");
});
