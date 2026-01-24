/**
 * Common type definitions for mermaid-aa
 */

/** Layout direction */
export type Direction = "TB" | "TD" | "BT" | "LR" | "RL";

/** Character set for rendering */
export type Charset =
  | "ascii"
  | "unicode"
  | "unicode-round"
  | "unicode-bold"
  | "unicode-double";

/** Node shape types */
export type NodeShape =
  | "rectangle" // A[text]
  | "rounded" // A(text)
  | "stadium" // A([text])
  | "diamond" // A{text}
  | "hexagon" // A{{text}}
  | "circle" // A((text))
  | "parallelogram" // A[/text/]
  | "trapezoid"; // A[/text\]

/** Edge style types */
export type EdgeStyle =
  | "solid" // -->
  | "dotted" // -.->
  | "thick" // ==>
  | "open"; // ---

/** AST Node */
export interface Node {
  id: string;
  label: string;
  shape: NodeShape;
}

/** AST Edge */
export interface Edge {
  from: string;
  to: string;
  style: EdgeStyle;
  label?: string;
}

/** Parsed flowchart */
export interface Flowchart {
  direction: Direction;
  nodes: Map<string, Node>;
  edges: Edge[];
}

/** Rendered node with position */
export interface RenderedNode extends Node {
  x: number;
  y: number;
  width: number;
  height: number;
}

/** CLI options */
export interface Options {
  file?: string;
  charset: Charset;
  ambiguousWidth: 1 | 2;
  direction?: Direction;
}
