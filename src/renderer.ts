import {
  Diagram,
  FlowchartDiagram,
  FlowchartNode,
  FlowchartEdge,
  SequenceDiagram,
} from './types';

export function render(diagram: Diagram): string {
  switch (diagram.type) {
    case 'flowchart':
      return renderFlowchart(diagram);
    case 'sequence':
      return renderSequence(diagram);
    default:
      return 'Unknown diagram type';
  }
}

// Canvas for drawing ASCII art
class Canvas {
  private grid: string[][];
  width: number;
  height: number;

  constructor(width: number, height: number) {
    this.width = width;
    this.height = height;
    this.grid = Array(height).fill(null).map(() => Array(width).fill(' '));
  }

  set(x: number, y: number, char: string): void {
    if (x >= 0 && x < this.width && y >= 0 && y < this.height) {
      this.grid[y][x] = char;
    }
  }

  get(x: number, y: number): string {
    if (x >= 0 && x < this.width && y >= 0 && y < this.height) {
      return this.grid[y][x];
    }
    return ' ';
  }

  drawText(x: number, y: number, text: string): void {
    for (let i = 0; i < text.length; i++) {
      this.set(x + i, y, text[i]);
    }
  }

  drawHLine(x: number, y: number, length: number, char = '-'): void {
    for (let i = 0; i < length; i++) {
      this.set(x + i, y, char);
    }
  }

  drawVLine(x: number, y: number, length: number, char = '|'): void {
    for (let i = 0; i < length; i++) {
      this.set(x, y + i, char);
    }
  }

  toString(): string {
    return this.grid.map(row => row.join('').trimEnd()).join('\n').trimEnd();
  }
}

// Flowchart renderer
function renderFlowchart(diagram: FlowchartDiagram): string {
  const { nodes, edges, direction } = diagram;

  if (nodes.length === 0) {
    return 'Empty flowchart';
  }

  const isHorizontal = direction === 'LR' || direction === 'RL';

  // Calculate node positions using topological ordering
  const nodePositions = layoutFlowchart(nodes, edges, isHorizontal);

  // Calculate canvas size
  const nodeWidth = 16;
  const nodeHeight = 5;
  const hSpacing = isHorizontal ? 8 : 4;
  const vSpacing = isHorizontal ? 4 : 4;

  let maxCol = 0;
  let maxRow = 0;
  for (const pos of nodePositions.values()) {
    maxCol = Math.max(maxCol, pos.col);
    maxRow = Math.max(maxRow, pos.row);
  }

  const canvasWidth = (maxCol + 1) * (nodeWidth + hSpacing) + 10;
  const canvasHeight = (maxRow + 1) * (nodeHeight + vSpacing) + 5;

  const canvas = new Canvas(canvasWidth, canvasHeight);

  // Draw nodes
  const nodeCoords: Map<string, { x: number; y: number; w: number; h: number }> = new Map();

  for (const node of nodes) {
    const pos = nodePositions.get(node.id);
    if (!pos) continue;

    const x = pos.col * (nodeWidth + hSpacing) + 2;
    const y = pos.row * (nodeHeight + vSpacing) + 1;
    const w = Math.max(nodeWidth, node.label.length + 4);
    const h = nodeHeight;

    nodeCoords.set(node.id, { x, y, w, h });
    drawNode(canvas, x, y, w, h, node);
  }

  // Draw edges
  for (const edge of edges) {
    const fromCoord = nodeCoords.get(edge.from);
    const toCoord = nodeCoords.get(edge.to);
    if (!fromCoord || !toCoord) continue;

    drawEdge(canvas, fromCoord, toCoord, edge, isHorizontal);
  }

  return canvas.toString();
}

function layoutFlowchart(
  nodes: FlowchartNode[],
  edges: FlowchartEdge[],
  isHorizontal: boolean
): Map<string, { row: number; col: number }> {
  const positions: Map<string, { row: number; col: number }> = new Map();

  // Build adjacency list
  const outgoing: Map<string, string[]> = new Map();
  const incoming: Map<string, string[]> = new Map();

  for (const node of nodes) {
    outgoing.set(node.id, []);
    incoming.set(node.id, []);
  }

  for (const edge of edges) {
    outgoing.get(edge.from)?.push(edge.to);
    incoming.get(edge.to)?.push(edge.from);
  }

  // Find root nodes (no incoming edges)
  const roots = nodes.filter(n => incoming.get(n.id)?.length === 0);
  if (roots.length === 0 && nodes.length > 0) {
    roots.push(nodes[0]);
  }

  // BFS to assign levels
  const levels: Map<string, number> = new Map();
  const queue: string[] = [];

  for (const root of roots) {
    levels.set(root.id, 0);
    queue.push(root.id);
  }

  while (queue.length > 0) {
    const current = queue.shift()!;
    const currentLevel = levels.get(current)!;

    for (const next of outgoing.get(current) || []) {
      const existingLevel = levels.get(next);
      if (existingLevel === undefined || existingLevel < currentLevel + 1) {
        levels.set(next, currentLevel + 1);
        queue.push(next);
      }
    }
  }

  // Group nodes by level
  const levelGroups: Map<number, string[]> = new Map();
  for (const node of nodes) {
    const level = levels.get(node.id) ?? 0;
    if (!levelGroups.has(level)) {
      levelGroups.set(level, []);
    }
    levelGroups.get(level)!.push(node.id);
  }

  // Assign positions
  for (const [level, nodeIds] of levelGroups) {
    for (let i = 0; i < nodeIds.length; i++) {
      if (isHorizontal) {
        positions.set(nodeIds[i], { row: i, col: level });
      } else {
        positions.set(nodeIds[i], { row: level, col: i });
      }
    }
  }

  return positions;
}

function drawNode(
  canvas: Canvas,
  x: number,
  y: number,
  w: number,
  h: number,
  node: FlowchartNode
): void {
  const label = node.label.substring(0, w - 4);
  const labelX = x + Math.floor((w - label.length) / 2);
  const labelY = y + Math.floor(h / 2);

  switch (node.shape) {
    case 'rectangle':
      // +--------+
      // |  text  |
      // +--------+
      canvas.drawHLine(x, y, w, '-');
      canvas.drawHLine(x, y + h - 1, w, '-');
      canvas.set(x, y, '+');
      canvas.set(x + w - 1, y, '+');
      canvas.set(x, y + h - 1, '+');
      canvas.set(x + w - 1, y + h - 1, '+');
      for (let i = 1; i < h - 1; i++) {
        canvas.set(x, y + i, '|');
        canvas.set(x + w - 1, y + i, '|');
      }
      canvas.drawText(labelX, labelY, label);
      break;

    case 'rounded':
      // /--------\
      // |  text  |
      // \--------/
      canvas.drawHLine(x + 1, y, w - 2, '-');
      canvas.drawHLine(x + 1, y + h - 1, w - 2, '-');
      canvas.set(x, y, '/');
      canvas.set(x + w - 1, y, '\\');
      canvas.set(x, y + h - 1, '\\');
      canvas.set(x + w - 1, y + h - 1, '/');
      for (let i = 1; i < h - 1; i++) {
        canvas.set(x, y + i, '|');
        canvas.set(x + w - 1, y + i, '|');
      }
      canvas.drawText(labelX, labelY, label);
      break;

    case 'diamond':
      //     /\
      //    /  \
      //   < tx >
      //    \  /
      //     \/
      const midX = x + Math.floor(w / 2);
      const midY = y + Math.floor(h / 2);
      canvas.set(midX, y, '/');
      canvas.set(midX + 1, y, '\\');
      canvas.set(midX - 1, y + 1, '/');
      canvas.set(midX + 2, y + 1, '\\');
      canvas.set(x + 2, midY, '<');
      canvas.set(x + w - 3, midY, '>');
      canvas.set(midX - 1, y + h - 2, '\\');
      canvas.set(midX + 2, y + h - 2, '/');
      canvas.set(midX, y + h - 1, '\\');
      canvas.set(midX + 1, y + h - 1, '/');
      canvas.drawText(labelX, labelY, label);
      break;

    case 'circle':
      //  .--.
      // ( tx )
      //  '--'
      canvas.drawText(x + 1, y, '.--.');
      canvas.set(x, y + 1, '(');
      canvas.set(x + w - 4, y + 1, ')');
      canvas.drawText(x + 1, y + 2, "'--'");
      canvas.drawText(x + 2, y + 1, label.substring(0, w - 6));
      break;

    case 'stadium':
      // (-------)
      // |  text |
      // (-------)
      canvas.set(x, y, '(');
      canvas.drawHLine(x + 1, y, w - 2, '-');
      canvas.set(x + w - 1, y, ')');
      for (let i = 1; i < h - 1; i++) {
        canvas.set(x, y + i, '|');
        canvas.set(x + w - 1, y + i, '|');
      }
      canvas.set(x, y + h - 1, '(');
      canvas.drawHLine(x + 1, y + h - 1, w - 2, '-');
      canvas.set(x + w - 1, y + h - 1, ')');
      canvas.drawText(labelX, labelY, label);
      break;
  }
}

function drawEdge(
  canvas: Canvas,
  from: { x: number; y: number; w: number; h: number },
  to: { x: number; y: number; w: number; h: number },
  edge: FlowchartEdge,
  isHorizontal: boolean
): void {
  const lineChar = edge.style === 'dotted' ? ':' : edge.style === 'thick' ? '=' : '-';
  const vLineChar = edge.style === 'dotted' ? ':' : '|';

  if (isHorizontal) {
    // Draw horizontal arrow
    const startX = from.x + from.w;
    const startY = from.y + Math.floor(from.h / 2);
    const endX = to.x - 1;
    const endY = to.y + Math.floor(to.h / 2);

    if (startY === endY) {
      // Straight horizontal line
      canvas.drawHLine(startX, startY, endX - startX, lineChar);
      if (edge.arrow === 'arrow') {
        canvas.set(endX, endY, '>');
      }
    } else {
      // Need to route around
      const midX = startX + Math.floor((endX - startX) / 2);
      canvas.drawHLine(startX, startY, midX - startX, lineChar);
      canvas.drawVLine(midX, Math.min(startY, endY), Math.abs(endY - startY) + 1, vLineChar);
      canvas.drawHLine(midX, endY, endX - midX, lineChar);
      if (edge.arrow === 'arrow') {
        canvas.set(endX, endY, '>');
      }
    }
  } else {
    // Draw vertical arrow
    const startX = from.x + Math.floor(from.w / 2);
    const startY = from.y + from.h;
    const endX = to.x + Math.floor(to.w / 2);
    const endY = to.y - 1;

    if (startX === endX) {
      // Straight vertical line
      canvas.drawVLine(startX, startY, endY - startY, vLineChar);
      if (edge.arrow === 'arrow') {
        canvas.set(endX, endY, 'v');
      }
    } else {
      // Need to route around
      const midY = startY + Math.floor((endY - startY) / 2);
      canvas.drawVLine(startX, startY, midY - startY, vLineChar);
      canvas.drawHLine(Math.min(startX, endX), midY, Math.abs(endX - startX) + 1, lineChar);
      canvas.drawVLine(endX, midY, endY - midY, vLineChar);
      if (edge.arrow === 'arrow') {
        canvas.set(endX, endY, 'v');
      }
    }
  }

  // Draw edge label if present
  if (edge.label) {
    const labelX = isHorizontal
      ? from.x + from.w + 2
      : Math.min(from.x, to.x) + Math.floor(Math.abs(to.x - from.x) / 2);
    const labelY = isHorizontal
      ? from.y + Math.floor(from.h / 2) - 1
      : from.y + from.h + 1;
    canvas.drawText(labelX, labelY, edge.label);
  }
}

// Sequence diagram renderer
function renderSequence(diagram: SequenceDiagram): string {
  const { participants, messages } = diagram;

  if (participants.length === 0) {
    return 'Empty sequence diagram';
  }

  const colWidth = 20;
  const rowHeight = 2;
  const headerHeight = 4;

  const width = participants.length * colWidth + 10;
  const height = headerHeight + messages.length * rowHeight + 4;

  const canvas = new Canvas(width, height);

  // Calculate participant positions
  const participantX: Map<string, number> = new Map();
  participants.forEach((p, i) => {
    const x = i * colWidth + 10;
    participantX.set(p.id, x);
  });

  // Draw participant boxes at top
  for (const participant of participants) {
    const x = participantX.get(participant.id)!;
    const label = participant.label.substring(0, colWidth - 4);
    const boxWidth = Math.max(label.length + 4, 10);
    const boxX = x - Math.floor(boxWidth / 2);

    // Draw box
    canvas.drawHLine(boxX, 0, boxWidth, '-');
    canvas.set(boxX, 0, '+');
    canvas.set(boxX + boxWidth - 1, 0, '+');
    canvas.set(boxX, 1, '|');
    canvas.set(boxX + boxWidth - 1, 1, '|');
    canvas.drawText(boxX + 2, 1, label);
    canvas.drawHLine(boxX, 2, boxWidth, '-');
    canvas.set(boxX, 2, '+');
    canvas.set(boxX + boxWidth - 1, 2, '+');
  }

  // Draw lifelines
  for (const participant of participants) {
    const x = participantX.get(participant.id)!;
    canvas.drawVLine(x, 3, height - 4, '|');
  }

  // Draw messages
  for (let i = 0; i < messages.length; i++) {
    const msg = messages[i];
    const y = headerHeight + i * rowHeight;
    const fromX = participantX.get(msg.from)!;
    const toX = participantX.get(msg.to)!;

    const lineChar = msg.style === 'dotted' ? '-' : '-';
    const arrowChar = msg.arrow === 'arrow' ? '>' : '>';

    if (fromX < toX) {
      // Arrow going right
      canvas.drawHLine(fromX + 1, y, toX - fromX - 2, lineChar);
      canvas.set(toX - 1, y, arrowChar);
      // Label above arrow
      const label = msg.label.substring(0, toX - fromX - 4);
      canvas.drawText(fromX + 2, y - 1, label);
    } else if (fromX > toX) {
      // Arrow going left
      canvas.drawHLine(toX + 2, y, fromX - toX - 2, lineChar);
      canvas.set(toX + 1, y, '<');
      // Label above arrow
      const label = msg.label.substring(0, fromX - toX - 4);
      canvas.drawText(toX + 3, y - 1, label);
    } else {
      // Self-message
      canvas.drawText(fromX + 1, y, '--+');
      canvas.set(fromX + 3, y + 1, '|');
      canvas.drawText(fromX + 1, y + 2, '<-+');
      canvas.drawText(fromX + 4, y, msg.label.substring(0, 10));
    }
  }

  return canvas.toString();
}
