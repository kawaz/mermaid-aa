/**
 * Layout calculation for flowchart nodes
 */

import type { AmbiguousWidthMode, Flowchart, RenderedNode } from "./types.ts";
import { stringWidth } from "./width.ts";

/** Minimum dimensions */
const MIN_NODE_WIDTH = 5;
const MIN_NODE_HEIGHT = 3;
const NODE_PADDING = 2;
const NODE_SPACING_X = 4;
const NODE_SPACING_Y = 2;

/** Calculate layout for all nodes */
export function calculateLayout(
  flowchart: Flowchart,
  ambiguousWidth: AmbiguousWidthMode = 1,
): {
  nodes: Map<string, RenderedNode>;
  width: number;
  height: number;
} {
  const { direction, nodes, edges } = flowchart;
  const renderedNodes = new Map<string, RenderedNode>();

  // Build adjacency lists
  const outgoing = new Map<string, string[]>();
  const incoming = new Map<string, string[]>();

  for (const node of nodes.values()) {
    outgoing.set(node.id, []);
    incoming.set(node.id, []);
  }

  for (const edge of edges) {
    outgoing.get(edge.from)?.push(edge.to);
    incoming.get(edge.to)?.push(edge.from);
  }

  // Calculate node dimensions
  const nodeDimensions = new Map<string, { width: number; height: number }>();
  for (const node of nodes.values()) {
    const labelWidth = stringWidth(node.label, ambiguousWidth);
    const width = Math.max(MIN_NODE_WIDTH, labelWidth + NODE_PADDING * 2);
    const height = MIN_NODE_HEIGHT;
    nodeDimensions.set(node.id, { width, height });
  }

  // Assign levels using topological sort
  const levels = new Map<string, number>();
  const visited = new Set<string>();

  function assignLevel(nodeId: string, level: number): void {
    if (visited.has(nodeId)) {
      // Update level if this path is longer
      const currentLevel = levels.get(nodeId) ?? 0;
      if (level > currentLevel) {
        levels.set(nodeId, level);
      }
      return;
    }
    visited.add(nodeId);
    levels.set(nodeId, level);

    const children = outgoing.get(nodeId) ?? [];
    for (const child of children) {
      assignLevel(child, level + 1);
    }
  }

  // Find root nodes (no incoming edges)
  const roots: string[] = [];
  for (const node of nodes.values()) {
    const incomingEdges = incoming.get(node.id) ?? [];
    if (incomingEdges.length === 0) {
      roots.push(node.id);
    }
  }

  // If no roots found, start from first node
  if (roots.length === 0 && nodes.size > 0) {
    roots.push(nodes.keys().next().value!);
  }

  // Assign levels starting from roots
  for (const root of roots) {
    assignLevel(root, 0);
  }

  // Handle any unvisited nodes (disconnected components)
  for (const node of nodes.values()) {
    if (!visited.has(node.id)) {
      assignLevel(node.id, 0);
    }
  }

  // Group nodes by level
  const nodesByLevel = new Map<number, string[]>();
  for (const [nodeId, level] of levels.entries()) {
    if (!nodesByLevel.has(level)) {
      nodesByLevel.set(level, []);
    }
    nodesByLevel.get(level)!.push(nodeId);
  }

  // Calculate positions based on direction
  const isVertical = direction === "TB" || direction === "TD" ||
    direction === "BT";
  const isReversed = direction === "BT" || direction === "RL";

  let totalWidth = 0;
  let totalHeight = 0;

  // Calculate level sizes
  const levelSizes: { primary: number; secondary: number }[] = [];
  const maxLevel = Math.max(...levels.values(), 0);

  for (let level = 0; level <= maxLevel; level++) {
    const nodesAtLevel = nodesByLevel.get(level) ?? [];
    let primarySize = 0;
    let secondarySize = 0;

    for (const nodeId of nodesAtLevel) {
      const dim = nodeDimensions.get(nodeId)!;
      if (isVertical) {
        secondarySize += dim.width + NODE_SPACING_X;
        primarySize = Math.max(primarySize, dim.height);
      } else {
        secondarySize += dim.height + NODE_SPACING_Y;
        primarySize = Math.max(primarySize, dim.width);
      }
    }

    // Remove trailing spacing
    if (nodesAtLevel.length > 0) {
      secondarySize -= isVertical ? NODE_SPACING_X : NODE_SPACING_Y;
    }

    levelSizes.push({ primary: primarySize, secondary: secondarySize });
  }

  // Calculate max secondary size for centering
  const maxSecondarySize = Math.max(...levelSizes.map((s) => s.secondary), 0);

  // Position nodes
  let primaryOffset = 0;
  const primarySpacing = isVertical ? NODE_SPACING_Y : NODE_SPACING_X;

  for (let level = 0; level <= maxLevel; level++) {
    const actualLevel = isReversed ? maxLevel - level : level;
    const nodesAtLevel = nodesByLevel.get(actualLevel) ?? [];
    const levelSize = levelSizes[actualLevel];

    // Center nodes in their level
    let secondaryOffset = Math.floor(
      (maxSecondarySize - levelSize.secondary) / 2,
    );
    const secondarySpacing = isVertical ? NODE_SPACING_X : NODE_SPACING_Y;

    for (const nodeId of nodesAtLevel) {
      const node = nodes.get(nodeId)!;
      const dim = nodeDimensions.get(nodeId)!;

      let x: number, y: number;
      if (isVertical) {
        x = secondaryOffset;
        y = primaryOffset;
      } else {
        x = primaryOffset;
        y = secondaryOffset;
      }

      renderedNodes.set(nodeId, {
        ...node,
        x,
        y,
        width: dim.width,
        height: dim.height,
      });

      secondaryOffset += (isVertical ? dim.width : dim.height) +
        secondarySpacing;
    }

    primaryOffset += levelSize.primary + primarySpacing;
  }

  // Calculate total dimensions
  for (const node of renderedNodes.values()) {
    totalWidth = Math.max(totalWidth, node.x + node.width);
    totalHeight = Math.max(totalHeight, node.y + node.height);
  }

  return {
    nodes: renderedNodes,
    width: totalWidth,
    height: totalHeight,
  };
}
