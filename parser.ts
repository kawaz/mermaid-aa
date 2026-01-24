/**
 * Mermaid flowchart parser
 */

import type {
  Direction,
  Edge,
  EdgeStyle,
  Flowchart,
  Node,
  NodeShape,
} from "./types.ts";

/** Parse error with location information */
export class ParseError extends Error {
  constructor(
    message: string,
    public line: number,
    public column: number,
  ) {
    super(`Parse error at line ${line}, column ${column}: ${message}`);
    this.name = "ParseError";
  }
}

/** Parse a direction from string */
function parseDirection(dir: string): Direction {
  const normalized = dir.toUpperCase();
  switch (normalized) {
    case "TB":
    case "TD":
      return "TB";
    case "BT":
      return "BT";
    case "LR":
      return "LR";
    case "RL":
      return "RL";
    default:
      return "TB";
  }
}

/** Parse node shape from the node text */
function parseNode(
  nodeText: string,
): { id: string; label: string; shape: NodeShape } {
  // Strip whitespace
  nodeText = nodeText.trim();

  // Circle: A((text))
  let match = nodeText.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\(\((.+?)\)\)$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "circle" };
  }

  // Hexagon: A{{text}}
  match = nodeText.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\{\{(.+?)\}\}$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "hexagon" };
  }

  // Stadium: A([text])
  match = nodeText.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\(\[(.+?)\]\)$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "stadium" };
  }

  // Trapezoid: A[/text\]
  match = nodeText.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\[\/(.+?)\\\]$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "trapezoid" };
  }

  // Parallelogram: A[/text/]
  match = nodeText.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\[\/(.+?)\/\]$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "parallelogram" };
  }

  // Diamond: A{text}
  match = nodeText.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\{(.+?)\}$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "diamond" };
  }

  // Rounded: A(text)
  match = nodeText.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\((.+?)\)$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "rounded" };
  }

  // Rectangle: A[text]
  match = nodeText.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\[(.+?)\]$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "rectangle" };
  }

  // Just an ID (no shape specified)
  match = nodeText.match(/^([a-zA-Z_][a-zA-Z0-9_]*)$/);
  if (match) {
    return { id: match[1], label: match[1], shape: "rectangle" };
  }

  throw new Error(`Invalid node syntax: ${nodeText}`);
}

/** Parse edge style from arrow text */
function parseEdgeStyle(arrow: string): { style: EdgeStyle; label?: string } {
  // With label: --|text|-->
  let match = arrow.match(/^--\|(.+?)\|-->$/);
  if (match) {
    return { style: "solid", label: match[1] };
  }

  // Dotted with label: -.-|text|->
  match = arrow.match(/^-\.-\|(.+?)\|->$/);
  if (match) {
    return { style: "dotted", label: match[1] };
  }

  // Thick with label: ==|text|==>
  match = arrow.match(/^==\|(.+?)\|==>$/);
  if (match) {
    return { style: "thick", label: match[1] };
  }

  // Simple arrows
  if (arrow === "-->") {
    return { style: "solid" };
  }
  if (arrow === "-.->") {
    return { style: "dotted" };
  }
  if (arrow === "==>") {
    return { style: "thick" };
  }
  if (arrow === "---") {
    return { style: "open" };
  }

  // Default to solid
  return { style: "solid" };
}

/** Parse a single statement (node definition or edge) */
function parseStatement(
  statement: string,
  nodes: Map<string, Node>,
): Edge | null {
  statement = statement.trim();
  if (!statement) return null;

  // Try to find an edge arrow
  const arrowPatterns = [
    /--\|.+?\|-->/,
    /-\.-\|.+?\|->/,
    /==\|.+?\|==>/,
    /-->/,
    /-\.->|\.->/,
    /==>/,
    /---/,
  ];

  let arrowMatch: RegExpMatchArray | null = null;
  let arrowPattern: RegExp | null = null;

  for (const pattern of arrowPatterns) {
    arrowMatch = statement.match(pattern);
    if (arrowMatch) {
      arrowPattern = pattern;
      break;
    }
  }

  if (!arrowMatch || !arrowPattern) {
    // No edge, just a node definition
    try {
      const node = parseNode(statement);
      if (!nodes.has(node.id)) {
        nodes.set(node.id, node);
      }
    } catch {
      // Ignore invalid node syntax
    }
    return null;
  }

  // Split by arrow
  const arrowIndex = statement.indexOf(arrowMatch[0]);
  const leftPart = statement.substring(0, arrowIndex).trim();
  const rightPart = statement.substring(arrowIndex + arrowMatch[0].length)
    .trim();

  // Parse both sides
  const fromNode = parseNode(leftPart);
  const toNode = parseNode(rightPart);

  // Register nodes if not already present
  if (!nodes.has(fromNode.id)) {
    nodes.set(fromNode.id, fromNode);
  }
  if (!nodes.has(toNode.id)) {
    nodes.set(toNode.id, toNode);
  }

  // Parse edge style
  const { style, label } = parseEdgeStyle(arrowMatch[0]);

  return {
    from: fromNode.id,
    to: toNode.id,
    style,
    label,
  };
}

/** Parse a Mermaid flowchart */
export function parse(input: string): Flowchart {
  const lines = input.split("\n");
  let direction: Direction = "TB";
  const nodes = new Map<string, Node>();
  const edges: Edge[] = [];

  let inFlowchart = false;

  for (const line of lines) {
    const trimmed = line.trim();

    // Skip empty lines and comments
    if (!trimmed || trimmed.startsWith("%%")) {
      continue;
    }

    // Check for flowchart/graph declaration
    const flowchartMatch = trimmed.match(
      /^(flowchart|graph)\s*(TB|TD|BT|LR|RL)?/i,
    );

    let remainingContent = trimmed;

    if (flowchartMatch) {
      inFlowchart = true;
      if (flowchartMatch[2]) {
        direction = parseDirection(flowchartMatch[2]);
      }
      // Remove the flowchart/graph declaration from the line
      remainingContent = trimmed.substring(flowchartMatch[0].length).trim();
      // If nothing remains on this line, continue to next line
      if (!remainingContent) {
        continue;
      }
    }

    if (!inFlowchart) {
      continue;
    }

    // Split by semicolon for multiple statements on one line
    const statements = remainingContent.split(";");

    for (const stmt of statements) {
      const edge = parseStatement(stmt, nodes);
      if (edge) {
        edges.push(edge);
      }
    }
  }

  return {
    direction,
    nodes,
    edges,
  };
}
