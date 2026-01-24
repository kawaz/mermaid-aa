// Types module - Common type definitions for mermaid-aa

/// Direction for flowchart layout
pub(all) enum Direction {
  TB  // Top to Bottom
  TD  // Top to Down (same as TB)
  BT  // Bottom to Top
  LR  // Left to Right
  RL  // Right to Left
} derive(Eq, Show)

/// Parse direction from string
pub fn Direction::from_string(s : String) -> Direction? {
  match s.to_upper() {
    "TB" => Some(TB)
    "TD" => Some(TD)
    "BT" => Some(BT)
    "LR" => Some(LR)
    "RL" => Some(RL)
    _ => None
  }
}

/// Check if direction is horizontal
pub fn Direction::is_horizontal(self : Direction) -> Bool {
  match self {
    LR | RL => true
    _ => false
  }
}

/// Node shape types (8 kinds)
pub(all) enum NodeShape {
  Rectangle       // A[text]
  Rounded         // A(text)
  Stadium         // A([text])
  Diamond         // A{text}
  Hexagon         // A{{text}}
  Circle          // A((text))
  Parallelogram   // A[/text/]
  Trapezoid       // A[/text\]
} derive(Eq, Show)

/// Edge/Arrow style
pub(all) enum EdgeStyle {
  Solid     // -->
  Dotted    // -.->
  Thick     // ==>
  Open      // ---
} derive(Eq, Show)

/// Node in the flowchart
pub(all) struct Node {
  id : String
  label : String
  shape : NodeShape
} derive(Eq, Show)

/// Create a new Node
pub fn Node::new(id : String, label : String, shape : NodeShape) -> Node {
  { id, label, shape }
}

/// Get node id
pub fn Node::id(self : Node) -> String {
  self.id
}

/// Get node label
pub fn Node::label(self : Node) -> String {
  self.label
}

/// Get node shape
pub fn Node::shape(self : Node) -> NodeShape {
  self.shape
}

/// Edge connecting two nodes
pub(all) struct Edge {
  from : String
  to : String
  style : EdgeStyle
  label : String?
} derive(Eq, Show)

/// Create a new Edge
pub fn Edge::new(
  from : String,
  to : String,
  style : EdgeStyle,
  label : String?
) -> Edge {
  { from, to, style, label }
}

/// Get edge source
pub fn Edge::from(self : Edge) -> String {
  self.from
}

/// Get edge target
pub fn Edge::to(self : Edge) -> String {
  self.to
}

/// Get edge style
pub fn Edge::style(self : Edge) -> EdgeStyle {
  self.style
}

/// Get edge label
pub fn Edge::label(self : Edge) -> String? {
  self.label
}

/// Flowchart AST
pub(all) struct Flowchart {
  direction : Direction
  nodes : Array[Node]
  edges : Array[Edge]
} derive(Eq, Show)

/// Create a new Flowchart
pub fn Flowchart::new(direction : Direction) -> Flowchart {
  { direction, nodes: [], edges: [] }
}

/// Get flowchart direction
pub fn Flowchart::direction(self : Flowchart) -> Direction {
  self.direction
}

/// Get flowchart nodes
pub fn Flowchart::nodes(self : Flowchart) -> Array[Node] {
  self.nodes
}

/// Get flowchart edges
pub fn Flowchart::edges(self : Flowchart) -> Array[Edge] {
  self.edges
}

/// Add a node to the flowchart
pub fn Flowchart::add_node(self : Flowchart, node : Node) -> Unit {
  // Check if node already exists
  for n in self.nodes {
    if n.id == node.id {
      return
    }
  }
  self.nodes.push(node)
}

/// Add an edge to the flowchart
pub fn Flowchart::add_edge(self : Flowchart, edge : Edge) -> Unit {
  self.edges.push(edge)
}

/// Character set for rendering
pub(all) enum Charset {
  Ascii         // +--+ | - + >
  Unicode       // Box drawing
  UnicodeRound  // Rounded box drawing
  UnicodeBold   // Bold box drawing
  UnicodeDouble // Double box drawing
} derive(Eq, Show)

/// Parse charset from string
pub fn Charset::from_string(s : String) -> Charset? {
  match s.to_lower() {
    "ascii" => Some(Ascii)
    "unicode" => Some(Unicode)
    "unicode-round" => Some(UnicodeRound)
    "unicode-bold" => Some(UnicodeBold)
    "unicode-double" => Some(UnicodeDouble)
    _ => None
  }
}

/// Render options
pub(all) struct RenderOptions {
  charset : Charset
  ambiguous_width : Int  // 1 or 2
} derive(Eq, Show)

/// Create default render options
pub fn RenderOptions::default() -> RenderOptions {
  { charset: Unicode, ambiguous_width: 1 }
}

/// Create render options with specified values
pub fn RenderOptions::new(charset : Charset, ambiguous_width : Int) -> RenderOptions {
  { charset, ambiguous_width }
}

/// Get charset
pub fn RenderOptions::charset(self : RenderOptions) -> Charset {
  self.charset
}

/// Get ambiguous width
pub fn RenderOptions::ambiguous_width(self : RenderOptions) -> Int {
  self.ambiguous_width
}
