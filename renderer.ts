/**
 * ASCII art renderer for flowcharts
 */

import type {
  Charset,
  Direction,
  Edge,
  Flowchart,
  RenderedNode,
} from "./types.ts";
import { Canvas } from "./canvas.ts";
import { calculateLayout } from "./layout.ts";
import { type CharacterSet, getCharset } from "./style.ts";
import { centerString } from "./width.ts";

/** Render a flowchart to ASCII art */
export function render(
  flowchart: Flowchart,
  charset: Charset = "unicode",
  ambiguousWidth: 1 | 2 = 1,
): string {
  const chars = getCharset(charset);
  const { nodes, width, height } = calculateLayout(flowchart, ambiguousWidth);

  if (nodes.size === 0) {
    return "";
  }

  // Create canvas with padding
  const canvas = new Canvas(width + 2, height + 2, ambiguousWidth);

  // Draw nodes
  for (const node of nodes.values()) {
    drawNode(canvas, node, chars, ambiguousWidth);
  }

  // Draw edges
  for (const edge of flowchart.edges) {
    const fromNode = nodes.get(edge.from);
    const toNode = nodes.get(edge.to);
    if (fromNode && toNode) {
      drawEdge(
        canvas,
        fromNode,
        toNode,
        edge,
        chars,
        flowchart.direction,
        ambiguousWidth,
      );
    }
  }

  return canvas.render();
}

/** Draw a node on the canvas */
function drawNode(
  canvas: Canvas,
  node: RenderedNode,
  chars: CharacterSet,
  ambiguousWidth: 1 | 2,
): void {
  const { x, y, width, height, label, shape } = node;

  switch (shape) {
    case "rectangle":
      drawRectangle(canvas, x, y, width, height, label, chars, ambiguousWidth);
      break;
    case "rounded":
      drawRounded(canvas, x, y, width, height, label, chars, ambiguousWidth);
      break;
    case "stadium":
      drawStadium(canvas, x, y, width, height, label, chars, ambiguousWidth);
      break;
    case "diamond":
      drawDiamond(canvas, x, y, width, height, label, chars, ambiguousWidth);
      break;
    case "hexagon":
      drawHexagon(canvas, x, y, width, height, label, chars, ambiguousWidth);
      break;
    case "circle":
      drawCircle(canvas, x, y, width, height, label, chars, ambiguousWidth);
      break;
    case "parallelogram":
      drawParallelogram(
        canvas,
        x,
        y,
        width,
        height,
        label,
        chars,
        ambiguousWidth,
      );
      break;
    case "trapezoid":
      drawTrapezoid(canvas, x, y, width, height, label, chars, ambiguousWidth);
      break;
    default:
      drawRectangle(canvas, x, y, width, height, label, chars, ambiguousWidth);
  }
}

/** Draw a rectangle node */
function drawRectangle(
  canvas: Canvas,
  x: number,
  y: number,
  width: number,
  height: number,
  label: string,
  chars: CharacterSet,
  ambiguousWidth: 1 | 2,
): void {
  canvas.drawBox(
    x,
    y,
    width,
    height,
    chars.topLeft,
    chars.topRight,
    chars.bottomLeft,
    chars.bottomRight,
    chars.horizontal,
    chars.vertical,
  );

  // Draw label centered
  const labelY = y + Math.floor(height / 2);
  const centeredLabel = centerString(label, width - 2, " ", ambiguousWidth);
  canvas.drawString(x + 1, labelY, centeredLabel);
}

/** Draw a rounded node */
function drawRounded(
  canvas: Canvas,
  x: number,
  y: number,
  width: number,
  height: number,
  label: string,
  chars: CharacterSet,
  ambiguousWidth: 1 | 2,
): void {
  // Top and bottom with curved brackets
  canvas.set(x, y, chars.roundedTopLeft);
  canvas.set(x + width - 1, y, chars.roundedTopRight);
  canvas.set(x, y + height - 1, chars.roundedBottomLeft);
  canvas.set(x + width - 1, y + height - 1, chars.roundedBottomRight);

  // Horizontal lines
  for (let i = 1; i < width - 1; i++) {
    canvas.set(x + i, y, chars.horizontal);
    canvas.set(x + i, y + height - 1, chars.horizontal);
  }

  // Vertical lines (using vertical char for middle rows)
  for (let i = 1; i < height - 1; i++) {
    canvas.set(x, y + i, chars.vertical);
    canvas.set(x + width - 1, y + i, chars.vertical);
  }

  // Draw label centered
  const labelY = y + Math.floor(height / 2);
  const centeredLabel = centerString(label, width - 2, " ", ambiguousWidth);
  canvas.drawString(x + 1, labelY, centeredLabel);
}

/** Draw a stadium node */
function drawStadium(
  canvas: Canvas,
  x: number,
  y: number,
  width: number,
  height: number,
  label: string,
  chars: CharacterSet,
  ambiguousWidth: 1 | 2,
): void {
  // Stadium shape: ([text])
  canvas.set(x, y, "(");
  canvas.set(x + 1, y, "[");
  canvas.set(x + width - 2, y, "]");
  canvas.set(x + width - 1, y, ")");

  canvas.set(x, y + height - 1, "(");
  canvas.set(x + 1, y + height - 1, "[");
  canvas.set(x + width - 2, y + height - 1, "]");
  canvas.set(x + width - 1, y + height - 1, ")");

  // Horizontal lines
  for (let i = 2; i < width - 2; i++) {
    canvas.set(x + i, y, chars.horizontal);
    canvas.set(x + i, y + height - 1, chars.horizontal);
  }

  // Vertical lines
  for (let i = 1; i < height - 1; i++) {
    canvas.set(x, y + i, chars.vertical);
    canvas.set(x + width - 1, y + i, chars.vertical);
  }

  // Draw label centered
  const labelY = y + Math.floor(height / 2);
  const centeredLabel = centerString(label, width - 4, " ", ambiguousWidth);
  canvas.drawString(x + 2, labelY, centeredLabel);
}

/** Draw a diamond node */
function drawDiamond(
  canvas: Canvas,
  x: number,
  y: number,
  width: number,
  height: number,
  label: string,
  chars: CharacterSet,
  ambiguousWidth: 1 | 2,
): void {
  // Diamond: <text>
  const midY = y + Math.floor(height / 2);

  canvas.set(x, midY, chars.diamondLeft);
  canvas.set(x + width - 1, midY, chars.diamondRight);

  // Top slant
  for (let i = 1; i < width - 1; i++) {
    canvas.set(x + i, y, chars.horizontal);
  }

  // Bottom slant
  for (let i = 1; i < width - 1; i++) {
    canvas.set(x + i, y + height - 1, chars.horizontal);
  }

  // Draw label centered
  const centeredLabel = centerString(label, width - 2, " ", ambiguousWidth);
  canvas.drawString(x + 1, midY, centeredLabel);
}

/** Draw a hexagon node */
function drawHexagon(
  canvas: Canvas,
  x: number,
  y: number,
  width: number,
  height: number,
  label: string,
  chars: CharacterSet,
  ambiguousWidth: 1 | 2,
): void {
  // Hexagon: {{text}}
  canvas.set(x, y, "{");
  canvas.set(x + 1, y, "{");
  canvas.set(x + width - 2, y, "}");
  canvas.set(x + width - 1, y, "}");

  canvas.set(x, y + height - 1, "{");
  canvas.set(x + 1, y + height - 1, "{");
  canvas.set(x + width - 2, y + height - 1, "}");
  canvas.set(x + width - 1, y + height - 1, "}");

  // Horizontal lines
  for (let i = 2; i < width - 2; i++) {
    canvas.set(x + i, y, chars.horizontal);
    canvas.set(x + i, y + height - 1, chars.horizontal);
  }

  // Vertical lines
  for (let i = 1; i < height - 1; i++) {
    canvas.set(x, y + i, chars.vertical);
    canvas.set(x + width - 1, y + i, chars.vertical);
  }

  // Draw label centered
  const labelY = y + Math.floor(height / 2);
  const centeredLabel = centerString(label, width - 4, " ", ambiguousWidth);
  canvas.drawString(x + 2, labelY, centeredLabel);
}

/** Draw a circle node */
function drawCircle(
  canvas: Canvas,
  x: number,
  y: number,
  width: number,
  height: number,
  label: string,
  chars: CharacterSet,
  ambiguousWidth: 1 | 2,
): void {
  // Circle: ((text))
  canvas.set(x, y, "(");
  canvas.set(x + 1, y, "(");
  canvas.set(x + width - 2, y, ")");
  canvas.set(x + width - 1, y, ")");

  canvas.set(x, y + height - 1, "(");
  canvas.set(x + 1, y + height - 1, "(");
  canvas.set(x + width - 2, y + height - 1, ")");
  canvas.set(x + width - 1, y + height - 1, ")");

  // Horizontal lines
  for (let i = 2; i < width - 2; i++) {
    canvas.set(x + i, y, chars.horizontal);
    canvas.set(x + i, y + height - 1, chars.horizontal);
  }

  // Vertical lines
  for (let i = 1; i < height - 1; i++) {
    canvas.set(x, y + i, chars.vertical);
    canvas.set(x + width - 1, y + i, chars.vertical);
  }

  // Draw label centered
  const labelY = y + Math.floor(height / 2);
  const centeredLabel = centerString(label, width - 4, " ", ambiguousWidth);
  canvas.drawString(x + 2, labelY, centeredLabel);
}

/** Draw a parallelogram node */
function drawParallelogram(
  canvas: Canvas,
  x: number,
  y: number,
  width: number,
  height: number,
  label: string,
  chars: CharacterSet,
  ambiguousWidth: 1 | 2,
): void {
  // Parallelogram: /text/
  canvas.set(x, y, chars.slashForward);
  canvas.set(x + width - 1, y, chars.slashForward);
  canvas.set(x, y + height - 1, chars.slashForward);
  canvas.set(x + width - 1, y + height - 1, chars.slashForward);

  // Horizontal lines
  for (let i = 1; i < width - 1; i++) {
    canvas.set(x + i, y, chars.horizontal);
    canvas.set(x + i, y + height - 1, chars.horizontal);
  }

  // Vertical lines
  for (let i = 1; i < height - 1; i++) {
    canvas.set(x, y + i, chars.vertical);
    canvas.set(x + width - 1, y + i, chars.vertical);
  }

  // Draw label centered
  const labelY = y + Math.floor(height / 2);
  const centeredLabel = centerString(label, width - 2, " ", ambiguousWidth);
  canvas.drawString(x + 1, labelY, centeredLabel);
}

/** Draw a trapezoid node */
function drawTrapezoid(
  canvas: Canvas,
  x: number,
  y: number,
  width: number,
  height: number,
  label: string,
  chars: CharacterSet,
  ambiguousWidth: 1 | 2,
): void {
  // Trapezoid: /text\
  canvas.set(x, y, chars.slashForward);
  canvas.set(x + width - 1, y, chars.slashBack);
  canvas.set(x, y + height - 1, chars.slashBack);
  canvas.set(x + width - 1, y + height - 1, chars.slashForward);

  // Horizontal lines
  for (let i = 1; i < width - 1; i++) {
    canvas.set(x + i, y, chars.horizontal);
    canvas.set(x + i, y + height - 1, chars.horizontal);
  }

  // Vertical lines
  for (let i = 1; i < height - 1; i++) {
    canvas.set(x, y + i, chars.vertical);
    canvas.set(x + width - 1, y + i, chars.vertical);
  }

  // Draw label centered
  const labelY = y + Math.floor(height / 2);
  const centeredLabel = centerString(label, width - 2, " ", ambiguousWidth);
  canvas.drawString(x + 1, labelY, centeredLabel);
}

/** Draw an edge between two nodes */
function drawEdge(
  canvas: Canvas,
  from: RenderedNode,
  to: RenderedNode,
  edge: Edge,
  chars: CharacterSet,
  direction: Direction,
  _ambiguousWidth: 1 | 2,
): void {
  const isVertical = direction === "TB" || direction === "TD" ||
    direction === "BT";

  // Calculate connection points
  let fromX: number, fromY: number, toX: number, toY: number;

  if (isVertical) {
    // Connect from bottom of from node to top of to node
    fromX = from.x + Math.floor(from.width / 2);
    fromY = from.y + from.height;
    toX = to.x + Math.floor(to.width / 2);
    toY = to.y - 1;
  } else {
    // Connect from right of from node to left of to node
    fromX = from.x + from.width;
    fromY = from.y + Math.floor(from.height / 2);
    toX = to.x - 1;
    toY = to.y + Math.floor(to.height / 2);
  }

  // Choose line character based on edge style
  let lineChar = chars.horizontal;
  if (edge.style === "dotted") {
    lineChar = chars.dotted;
  } else if (edge.style === "thick") {
    lineChar = chars.thick;
  }

  // Draw simple straight line
  if (isVertical) {
    // Draw vertical line
    const vertChar = edge.style === "dotted" ? ":" : chars.vertical;
    for (let y = fromY; y <= toY; y++) {
      canvas.set(fromX, y, vertChar);
    }

    // Draw arrow
    if (edge.style !== "open") {
      canvas.set(toX, toY, chars.arrowDown);
    }
  } else {
    // Draw horizontal line
    const startX = Math.min(fromX, toX);
    const endX = Math.max(fromX, toX);
    for (let x = startX; x <= endX; x++) {
      canvas.set(x, fromY, lineChar);
    }

    // Draw arrow
    if (edge.style !== "open") {
      if (toX > fromX) {
        canvas.set(toX, toY, chars.arrowRight);
      } else {
        canvas.set(toX, toY, chars.arrowLeft);
      }
    }
  }

  // Draw edge label if present
  if (edge.label) {
    if (isVertical) {
      const labelY = Math.floor((fromY + toY) / 2);
      canvas.drawString(fromX + 1, labelY, edge.label);
    } else {
      const labelX = Math.floor((fromX + toX) / 2) -
        Math.floor(edge.label.length / 2);
      canvas.drawString(labelX, fromY - 1, edge.label);
    }
  }
}
