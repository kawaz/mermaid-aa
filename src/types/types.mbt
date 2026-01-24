// Types module - Common type definitions for mermaid-aa

/// Width mode for ambiguous characters
/// Controls how ambiguous-width characters and box drawing characters are handled
pub(all) enum WidthMode {
  /// Ambiguous=1, BoxDrawing=1 (Western terminals, default)
  Half
  /// Ambiguous=2, BoxDrawing=1 (CJK terminals, recommended)
  Full
  /// Ambiguous=1, BoxDrawing=1 (same as Half, explicit name)
  Console
  /// Ambiguous=2, BoxDrawing=2 (legacy compatibility)
  Legacy
} derive(Eq, Show)

/// Parse width mode from string
pub fn WidthMode::from_string(s : String) -> WidthMode? {
  match s.to_lower() {
    "1" | "half" => Some(Half)
    "2" | "full" => Some(Full)
    "console" => Some(Console)
    "legacy" => Some(Legacy)
    _ => None
  }
}

/// Get ambiguous width for this mode
pub fn WidthMode::ambiguous_width(self : WidthMode) -> Int {
  match self {
    Half | Console => 1
    Full | Legacy => 2
  }
}

/// Get box drawing width for this mode
pub fn WidthMode::box_drawing_width(self : WidthMode) -> Int {
  match self {
    Half | Full | Console => 1
    Legacy => 2
  }
}

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
  width_mode : WidthMode
} derive(Eq, Show)

/// Create default render options
pub fn RenderOptions::default() -> RenderOptions {
  { charset: Unicode, width_mode: Half }
}

/// Create render options with specified values (legacy, for compatibility)
pub fn RenderOptions::new(charset : Charset, ambiguous_width : Int) -> RenderOptions {
  let mode : WidthMode = if ambiguous_width == 2 { Full } else { Half }
  { charset, width_mode: mode }
}

/// Create render options with WidthMode
pub fn RenderOptions::with_mode(charset : Charset, mode : WidthMode) -> RenderOptions {
  { charset, width_mode: mode }
}

/// Get charset
pub fn RenderOptions::charset(self : RenderOptions) -> Charset {
  self.charset
}

/// Get width mode
pub fn RenderOptions::width_mode(self : RenderOptions) -> WidthMode {
  self.width_mode
}

/// Get ambiguous width (for compatibility)
pub fn RenderOptions::ambiguous_width(self : RenderOptions) -> Int {
  self.width_mode.ambiguous_width()
}
