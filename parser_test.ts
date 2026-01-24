/**
 * Tests for Mermaid parser
 */

import { assertEquals } from "@std/assert";
import { parse } from "./parser.ts";

Deno.test("parse - Simple flowchart with direction", () => {
  const input = `flowchart TD
    A --> B`;

  const result = parse(input);

  assertEquals(result.direction, "TB");
  assertEquals(result.nodes.size, 2);
  assertEquals(result.edges.length, 1);
  assertEquals(result.edges[0].from, "A");
  assertEquals(result.edges[0].to, "B");
});

Deno.test("parse - Graph keyword works like flowchart", () => {
  const input = `graph LR
    A --> B`;

  const result = parse(input);

  assertEquals(result.direction, "LR");
  assertEquals(result.nodes.size, 2);
});

Deno.test("parse - Rectangle node", () => {
  const input = `flowchart TD
    A[Hello World]`;

  const result = parse(input);

  assertEquals(result.nodes.size, 1);
  const node = result.nodes.get("A");
  assertEquals(node?.shape, "rectangle");
  assertEquals(node?.label, "Hello World");
});

Deno.test("parse - Rounded node", () => {
  const input = `flowchart TD
    A(Hello)`;

  const result = parse(input);

  const node = result.nodes.get("A");
  assertEquals(node?.shape, "rounded");
  assertEquals(node?.label, "Hello");
});

Deno.test("parse - Stadium node", () => {
  const input = `flowchart TD
    A([Stadium])`;

  const result = parse(input);

  const node = result.nodes.get("A");
  assertEquals(node?.shape, "stadium");
  assertEquals(node?.label, "Stadium");
});

Deno.test("parse - Diamond node", () => {
  const input = `flowchart TD
    A{Decision}`;

  const result = parse(input);

  const node = result.nodes.get("A");
  assertEquals(node?.shape, "diamond");
  assertEquals(node?.label, "Decision");
});

Deno.test("parse - Hexagon node", () => {
  const input = `flowchart TD
    A{{Hexagon}}`;

  const result = parse(input);

  const node = result.nodes.get("A");
  assertEquals(node?.shape, "hexagon");
  assertEquals(node?.label, "Hexagon");
});

Deno.test("parse - Circle node", () => {
  const input = `flowchart TD
    A((Circle))`;

  const result = parse(input);

  const node = result.nodes.get("A");
  assertEquals(node?.shape, "circle");
  assertEquals(node?.label, "Circle");
});

Deno.test("parse - Parallelogram node", () => {
  const input = `flowchart TD
    A[/Parallelogram/]`;

  const result = parse(input);

  const node = result.nodes.get("A");
  assertEquals(node?.shape, "parallelogram");
  assertEquals(node?.label, "Parallelogram");
});

Deno.test("parse - Trapezoid node", () => {
  const input = String.raw`flowchart TD
    A[/Trapezoid\]`;

  const result = parse(input);

  const node = result.nodes.get("A");
  assertEquals(node?.shape, "trapezoid");
  assertEquals(node?.label, "Trapezoid");
});

Deno.test("parse - Solid arrow", () => {
  const input = `flowchart TD
    A --> B`;

  const result = parse(input);

  assertEquals(result.edges[0].style, "solid");
});

Deno.test("parse - Dotted arrow", () => {
  const input = `flowchart TD
    A -.-> B`;

  const result = parse(input);

  assertEquals(result.edges[0].style, "dotted");
});

Deno.test("parse - Thick arrow", () => {
  const input = `flowchart TD
    A ==> B`;

  const result = parse(input);

  assertEquals(result.edges[0].style, "thick");
});

Deno.test("parse - Open connection", () => {
  const input = `flowchart TD
    A --- B`;

  const result = parse(input);

  assertEquals(result.edges[0].style, "open");
});

Deno.test("parse - Arrow with label", () => {
  const input = `flowchart TD
    A --|label|--> B`;

  const result = parse(input);

  assertEquals(result.edges[0].label, "label");
});

Deno.test("parse - Multiple statements on one line", () => {
  const input = `flowchart TD
    A --> B; B --> C`;

  const result = parse(input);

  assertEquals(result.edges.length, 2);
  assertEquals(result.edges[0].from, "A");
  assertEquals(result.edges[0].to, "B");
  assertEquals(result.edges[1].from, "B");
  assertEquals(result.edges[1].to, "C");
});

Deno.test("parse - Different directions", () => {
  assertEquals(parse("flowchart TB\n A-->B").direction, "TB");
  assertEquals(parse("flowchart TD\n A-->B").direction, "TB");
  assertEquals(parse("flowchart BT\n A-->B").direction, "BT");
  assertEquals(parse("flowchart LR\n A-->B").direction, "LR");
  assertEquals(parse("flowchart RL\n A-->B").direction, "RL");
});

Deno.test("parse - Comments are ignored", () => {
  const input = `flowchart TD
    %% This is a comment
    A --> B`;

  const result = parse(input);

  assertEquals(result.nodes.size, 2);
  assertEquals(result.edges.length, 1);
});

Deno.test("parse - CJK labels", () => {
  const input = `flowchart TD
    A[日本語] --> B[中国語]`;

  const result = parse(input);

  assertEquals(result.nodes.get("A")?.label, "日本語");
  assertEquals(result.nodes.get("B")?.label, "中国語");
});

Deno.test("parse - Default direction is TB", () => {
  const input = `flowchart
    A --> B`;

  const result = parse(input);

  assertEquals(result.direction, "TB");
});
