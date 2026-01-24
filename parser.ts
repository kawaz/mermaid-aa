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
    public sourceText?: string,
    public help?: string,
  ) {
    super(message);
    this.name = "ParseError";
  }

  /** Format error with visual indicator */
  format(filename?: string): string {
    const lines: string[] = [];
    const location = filename
      ? ` --> ${filename}:${this.line}:${this.column}`
      : ` --> input:${this.line}:${this.column}`;

    lines.push(`Error: ${this.message}`);
    lines.push(location);
    lines.push("  |");

    if (this.sourceText !== undefined) {
      const lineNumStr = String(this.line);
      const linePrefix = `${lineNumStr} | `;
      lines.push(`${linePrefix}${this.sourceText}`);

      // Create pointer line
      const padding = " ".repeat(linePrefix.length + this.column - 1);
      lines.push(`${padding}^^^^^`);
    }

    if (this.help) {
      lines.push(`  = help: ${this.help}`);
    }

    return lines.join("\n");
  }
}

/** Collection of parse errors */
export class ParseErrors extends Error {
  constructor(public errors: ParseError[]) {
    super(`Parse failed with ${errors.length} error(s)`);
    this.name = "ParseErrors";
  }

  /** Format all errors with visual indicators */
  format(filename?: string): string {
    return this.errors.map((e) => e.format(filename)).join("\n\n");
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

/** Context for parsing with position tracking */
interface ParseContext {
  errors: ParseError[];
  currentLine: number;
  currentLineText: string;
  lines: string[];
}

/** Parse node shape from the node text */
function parseNode(
  nodeText: string,
  ctx?: ParseContext,
  columnOffset?: number,
): { id: string; label: string; shape: NodeShape } | null {
  // Strip whitespace
  const trimmed = nodeText.trim();
  const leadingSpaces = nodeText.length - nodeText.trimStart().length;

  // Circle: A((text))
  let match = trimmed.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\(\((.+?)\)\)$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "circle" };
  }

  // Hexagon: A{{text}}
  match = trimmed.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\{\{(.+?)\}\}$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "hexagon" };
  }

  // Stadium: A([text])
  match = trimmed.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\(\[(.+?)\]\)$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "stadium" };
  }

  // Trapezoid: A[/text\]
  match = trimmed.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\[\/(.+?)\\\]$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "trapezoid" };
  }

  // Parallelogram: A[/text/]
  match = trimmed.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\[\/(.+?)\/\]$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "parallelogram" };
  }

  // Diamond: A{text}
  match = trimmed.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\{(.+?)\}$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "diamond" };
  }

  // Rounded: A(text)
  match = trimmed.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\((.+?)\)$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "rounded" };
  }

  // Rectangle: A[text]
  match = trimmed.match(/^([a-zA-Z_][a-zA-Z0-9_]*)\[(.+?)\]$/);
  if (match) {
    return { id: match[1], label: match[2], shape: "rectangle" };
  }

  // Just an ID (no shape specified)
  match = trimmed.match(/^([a-zA-Z_][a-zA-Z0-9_]*)$/);
  if (match) {
    return { id: match[1], label: match[1], shape: "rectangle" };
  }

  // Error handling with context
  if (ctx) {
    const column = (columnOffset ?? 0) + leadingSpaces + 1;
    let message: string;
    let help: string;

    if (trimmed === "") {
      message = "Missing node identifier";
      help = "Add a node identifier like 'A' or 'A[label]'";
    } else if (!/^[a-zA-Z_]/.test(trimmed)) {
      message = "Invalid node identifier start";
      help = "Node identifiers must start with a letter or underscore";
    } else if (trimmed.includes("[") && !trimmed.includes("]")) {
      message = "Unclosed bracket in node definition";
      help = "Close the bracket: A[text] or A[/text/]";
    } else if (trimmed.includes("(") && !trimmed.includes(")")) {
      message = "Unclosed parenthesis in node definition";
      help = "Close the parenthesis: A(text) or A((text))";
    } else if (trimmed.includes("{") && !trimmed.includes("}")) {
      message = "Unclosed brace in node definition";
      help = "Close the brace: A{text} or A{{text}}";
    } else {
      message = `Invalid node syntax: '${trimmed}'`;
      help = "Valid shapes: A[text], A(text), A{text}, A((text)), A{{text}}, A([text])";
    }

    ctx.errors.push(
      new ParseError(
        message,
        ctx.currentLine,
        column,
        ctx.currentLineText,
        help,
      ),
    );
  }

  return null;
}

/** Parse edge style from arrow text (arrow + optional label after arrow) */
function parseEdgeStyle(
  arrow: string,
  afterArrow?: string,
): { style: EdgeStyle; label?: string; consumedAfter?: string } {
  let label: string | undefined;
  let consumedAfter: string | undefined;

  // Check for label embedded in arrow: --|text|-->
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

  // Check for label after arrow: -->|text|
  if (afterArrow) {
    const labelMatch = afterArrow.match(/^\|([^|]+)\|(.*)$/);
    if (labelMatch) {
      label = labelMatch[1];
      consumedAfter = labelMatch[2].trim();
    }
  }

  // Simple arrows
  if (arrow === "-->") {
    return { style: "solid", label, consumedAfter };
  }
  if (arrow === "-.->") {
    return { style: "dotted", label, consumedAfter };
  }
  if (arrow === "==>") {
    return { style: "thick", label, consumedAfter };
  }
  if (arrow === "---") {
    return { style: "open", label, consumedAfter };
  }

  // Default to solid
  return { style: "solid", label, consumedAfter };
}

/** Arrow patterns for edge detection (order matters - more specific first) */
const ARROW_PATTERNS = [
  /--\|.+?\|-->/,
  /-\.-\|.+?\|->/,
  /==\|.+?\|==>/,
  /-->/,
  /-\.->|\.->/,
  /==>/,
  /---/,
];

/** Find the first arrow in the statement */
function findFirstArrow(statement: string): { match: string; index: number } | null {
  let earliest: { match: string; index: number; pattern: RegExp } | null = null;

  for (const pattern of ARROW_PATTERNS) {
    const match = statement.match(pattern);
    if (match && match.index !== undefined) {
      if (earliest === null || match.index < earliest.index) {
        earliest = { match: match[0], index: match.index, pattern };
      }
    }
  }

  return earliest ? { match: earliest.match, index: earliest.index } : null;
}

/** Extract the first node from a string that may contain chained edges */
function extractFirstNode(text: string): { nodeText: string; remaining: string } {
  text = text.trim();

  // Find if there's an arrow in the text
  const arrowInfo = findFirstArrow(text);
  if (arrowInfo) {
    const nodeText = text.substring(0, arrowInfo.index).trim();
    const remaining = text.substring(arrowInfo.index).trim();
    return { nodeText, remaining };
  }

  // No arrow, the entire text is the node
  return { nodeText: text, remaining: "" };
}

/** Parse a single statement (node definition or edge), returns multiple edges for chains */
function parseStatement(
  statement: string,
  nodes: Map<string, Node>,
  ctx: ParseContext,
  statementOffset: number,
): Edge[] {
  statement = statement.trim();
  if (!statement) return [];

  const edges: Edge[] = [];

  // Try to find an edge arrow
  const arrowInfo = findFirstArrow(statement);

  if (!arrowInfo) {
    // No edge, just a node definition
    const node = parseNode(statement, ctx, statementOffset);
    if (node && !nodes.has(node.id)) {
      nodes.set(node.id, node);
    }
    return [];
  }

  // Split by the first arrow
  const leftPart = statement.substring(0, arrowInfo.index).trim();
  let rightPart = statement.substring(arrowInfo.index + arrowInfo.match.length).trim();

  // Calculate column offset for the left node
  const leftOffset = statementOffset;

  // Parse the left node
  const fromNode = parseNode(leftPart, ctx, leftOffset);
  if (!fromNode) {
    // Error already recorded, skip this edge
    return [];
  }
  if (!nodes.has(fromNode.id)) {
    nodes.set(fromNode.id, fromNode);
  }

  // Parse edge style, checking for label after arrow (e.g., -->|label|)
  const { style, label, consumedAfter } = parseEdgeStyle(arrowInfo.match, rightPart);

  // If we consumed part of rightPart for the label, update rightPart
  if (consumedAfter !== undefined) {
    rightPart = consumedAfter;
  }

  // Check for missing node after arrow
  const rightPartTrimmed = rightPart.trim();
  if (!rightPartTrimmed) {
    const arrowEnd = statementOffset + arrowInfo.index + arrowInfo.match.length;
    ctx.errors.push(
      new ParseError(
        "Missing node identifier after arrow",
        ctx.currentLine,
        arrowEnd + 1,
        ctx.currentLineText,
        "Add a target node like 'A --> B'",
      ),
    );
    return [];
  }

  // Extract the first node from the right part (may be chained)
  const { nodeText: toNodeText, remaining } = extractFirstNode(rightPart);

  // Calculate column offset for the right node
  const rightOffset = statementOffset + arrowInfo.index + arrowInfo.match.length +
    (rightPart.length - rightPart.trimStart().length);

  const toNode = parseNode(toNodeText, ctx, rightOffset);
  if (!toNode) {
    // Error already recorded
    return [];
  }
  if (!nodes.has(toNode.id)) {
    nodes.set(toNode.id, toNode);
  }

  // Create the edge
  edges.push({
    from: fromNode.id,
    to: toNode.id,
    style,
    label,
  });

  // If there's remaining content, it's a chain - recursively parse it
  if (remaining) {
    // Calculate new offset for the remaining part
    const remainingOffset = statementOffset + statement.indexOf(remaining);
    const chainEdges = parseStatement(`${toNode.id} ${remaining}`, nodes, ctx, remainingOffset);
    edges.push(...chainEdges);
  }

  return edges;
}

/** Parse options */
export interface ParseOptions {
  /** Whether to throw on first error (default: false, collect all errors) */
  failFast?: boolean;
  /** Filename for error messages */
  filename?: string;
}

/** Parse a Mermaid flowchart */
export function parse(input: string, options?: ParseOptions): Flowchart {
  const lines = input.split("\n");
  let direction: Direction = "TB";
  const nodes = new Map<string, Node>();
  const edges: Edge[] = [];

  const ctx: ParseContext = {
    errors: [],
    currentLine: 0,
    currentLineText: "",
    lines,
  };

  let inFlowchart = false;

  for (let lineIndex = 0; lineIndex < lines.length; lineIndex++) {
    const line = lines[lineIndex];
    const trimmed = line.trim();

    // Update context for current line
    ctx.currentLine = lineIndex + 1;
    ctx.currentLineText = line;

    // Skip empty lines and comments
    if (!trimmed || trimmed.startsWith("%%")) {
      continue;
    }

    // Check for flowchart/graph declaration
    const flowchartMatch = trimmed.match(
      /^(flowchart|graph)\s*(TB|TD|BT|LR|RL)?/i,
    );

    let remainingContent = trimmed;
    let contentOffset = line.indexOf(trimmed);

    if (flowchartMatch) {
      inFlowchart = true;
      if (flowchartMatch[2]) {
        direction = parseDirection(flowchartMatch[2]);
      }
      // Remove the flowchart/graph declaration from the line
      remainingContent = trimmed.substring(flowchartMatch[0].length).trim();
      contentOffset = line.indexOf(remainingContent, contentOffset + flowchartMatch[0].length);
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

    let stmtOffset = contentOffset;
    for (const stmt of statements) {
      const stmtEdges = parseStatement(stmt, nodes, ctx, stmtOffset);
      edges.push(...stmtEdges);
      stmtOffset += stmt.length + 1; // +1 for semicolon

      // Fail fast if requested
      if (options?.failFast && ctx.errors.length > 0) {
        throw new ParseErrors(ctx.errors);
      }
    }
  }

  // Throw if there are any errors
  if (ctx.errors.length > 0) {
    throw new ParseErrors(ctx.errors);
  }

  return {
    direction,
    nodes,
    edges,
  };
}
