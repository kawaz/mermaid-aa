// Renderer module - ASCII art rendering for flowcharts

/// Node layout information
priv struct NodeLayout {
  node : @types.Node
  x : Int
  y : Int
  width : Int
  height : Int
}

/// Calculate node dimensions based on shape and label
fn calc_node_size(
  node : @types.Node,
  ambiguous_width : Int
) -> (Int, Int) {
  let label_width = @width.string_width(node.label(), ambiguous_width)
  // Add padding around label
  let padding = 2  // 1 char each side
  match node.shape() {
    @types.Rectangle | @types.Rounded | @types.Stadium => {
      let w = label_width + padding + 2  // +2 for borders
      let h = 3  // top border + label + bottom border
      (w, h)
    }
    @types.Diamond => {
      // Diamond needs more space
      let w = label_width + 4
      let h = 3
      (w, h)
    }
    @types.Hexagon => {
      let w = label_width + 6  // {{text}}
      let h = 3
      (w, h)
    }
    @types.Circle => {
      // Circle notation ((text))
      let w = label_width + 6
      let h = 3
      (w, h)
    }
    @types.Parallelogram | @types.Trapezoid => {
      let w = label_width + 4
      let h = 3
      (w, h)
    }
  }
}

/// Draw a node on the canvas
fn draw_node(
  canvas : Canvas,
  layout : NodeLayout,
  style : CharStyle,
  ambiguous_width : Int
) -> Unit {
  let x = layout.x
  let y = layout.y
  let w = layout.width
  let h = layout.height
  let label = layout.node.label()
  let label_x = x + (w - @width.string_width(label, ambiguous_width)) / 2
  let label_y = y + h / 2
  match layout.node.shape() {
    @types.Rectangle => {
      canvas.draw_rect(x, y, w, h, style)
      canvas.draw_text(label_x, label_y, label)
    }
    @types.Rounded => {
      // Use ( ) for sides
      canvas.set(x, y, style.paren_left)
      canvas.draw_h_line(x + 1, y, w - 2, style.h_line)
      canvas.set(x + w - 1, y, style.paren_right)
      canvas.set(x, y + h - 1, style.paren_left)
      canvas.draw_h_line(x + 1, y + h - 1, w - 2, style.h_line)
      canvas.set(x + w - 1, y + h - 1, style.paren_right)
      canvas.set(x, label_y, style.paren_left)
      canvas.set(x + w - 1, label_y, style.paren_right)
      canvas.draw_text(label_x, label_y, label)
    }
    @types.Stadium => {
      // ([text])
      canvas.set(x, y, style.paren_left)
      canvas.set(x + 1, y, '[')
      canvas.draw_h_line(x + 2, y, w - 4, style.h_line)
      canvas.set(x + w - 2, y, ']')
      canvas.set(x + w - 1, y, style.paren_right)
      canvas.set(x, y + h - 1, style.paren_left)
      canvas.set(x + 1, y + h - 1, '[')
      canvas.draw_h_line(x + 2, y + h - 1, w - 4, style.h_line)
      canvas.set(x + w - 2, y + h - 1, ']')
      canvas.set(x + w - 1, y + h - 1, style.paren_right)
      canvas.set(x, label_y, style.paren_left)
      canvas.set(x + w - 1, label_y, style.paren_right)
      canvas.draw_text(label_x, label_y, label)
    }
    @types.Diamond => {
      // <text>
      canvas.set(x, label_y, style.diamond_left)
      canvas.draw_h_line(x + 1, label_y, w - 2, style.h_line)
      canvas.set(x + w - 1, label_y, style.diamond_right)
      canvas.draw_text(label_x, label_y, label)
    }
    @types.Hexagon => {
      // {{text}}
      canvas.set(x, y, '{')
      canvas.set(x + 1, y, '{')
      canvas.draw_h_line(x + 2, y, w - 4, style.h_line)
      canvas.set(x + w - 2, y, '}')
      canvas.set(x + w - 1, y, '}')
      canvas.set(x, y + h - 1, '{')
      canvas.set(x + 1, y + h - 1, '{')
      canvas.draw_h_line(x + 2, y + h - 1, w - 4, style.h_line)
      canvas.set(x + w - 2, y + h - 1, '}')
      canvas.set(x + w - 1, y + h - 1, '}')
      canvas.set(x, label_y, '{')
      canvas.set(x + w - 1, label_y, '}')
      canvas.draw_text(label_x, label_y, label)
    }
    @types.Circle => {
      // ((text))
      canvas.set(x, y, '(')
      canvas.set(x + 1, y, '(')
      canvas.draw_h_line(x + 2, y, w - 4, style.h_line)
      canvas.set(x + w - 2, y, ')')
      canvas.set(x + w - 1, y, ')')
      canvas.set(x, y + h - 1, '(')
      canvas.set(x + 1, y + h - 1, '(')
      canvas.draw_h_line(x + 2, y + h - 1, w - 4, style.h_line)
      canvas.set(x + w - 2, y + h - 1, ')')
      canvas.set(x + w - 1, y + h - 1, ')')
      canvas.set(x, label_y, '(')
      canvas.set(x + w - 1, label_y, ')')
      canvas.draw_text(label_x, label_y, label)
    }
    @types.Parallelogram => {
      // /text/
      canvas.set(x, y, '/')
      canvas.draw_h_line(x + 1, y, w - 2, style.h_line)
      canvas.set(x + w - 1, y, '/')
      canvas.set(x, y + h - 1, '/')
      canvas.draw_h_line(x + 1, y + h - 1, w - 2, style.h_line)
      canvas.set(x + w - 1, y + h - 1, '/')
      canvas.set(x, label_y, '/')
      canvas.set(x + w - 1, label_y, '/')
      canvas.draw_text(label_x, label_y, label)
    }
    @types.Trapezoid => {
      // /text\
      canvas.set(x, y, '/')
      canvas.draw_h_line(x + 1, y, w - 2, style.h_line)
      canvas.set(x + w - 1, y, '\\')
      canvas.set(x, y + h - 1, '\\')
      canvas.draw_h_line(x + 1, y + h - 1, w - 2, style.h_line)
      canvas.set(x + w - 1, y + h - 1, '/')
      canvas.set(x, label_y, '/')
      canvas.set(x + w - 1, label_y, '\\')
      canvas.draw_text(label_x, label_y, label)
    }
  }
}

/// Draw an edge between nodes
fn draw_edge(
  canvas : Canvas,
  from_layout : NodeLayout,
  to_layout : NodeLayout,
  edge : @types.Edge,
  style : CharStyle,
  direction : @types.Direction
) -> Unit {
  match direction {
    @types.TB | @types.TD => {
      // Vertical flow (top to bottom)
      let from_x = from_layout.x + from_layout.width / 2
      let from_y = from_layout.y + from_layout.height
      let to_x = to_layout.x + to_layout.width / 2
      let to_y = to_layout.y
      // Draw vertical line
      let line_char = match edge.style() {
        @types.Solid | @types.Open => style.v_line
        @types.Dotted => ':'
        @types.Thick => style.v_line
      }
      for y = from_y; y < to_y; y = y + 1 {
        canvas.set(from_x, y, line_char)
      }
      // Draw arrow if not open style
      if edge.style() != @types.Open {
        canvas.set(to_x, to_y - 1, style.arrow_down)
      }
      // Draw label if present
      match edge.label() {
        Some(label) => {
          let mid_y = (from_y + to_y) / 2
          let label_x = from_x + 1
          canvas.draw_text(label_x, mid_y, label)
        }
        None => ()
      }
    }
    @types.BT => {
      // Vertical flow (bottom to top)
      let from_x = from_layout.x + from_layout.width / 2
      let from_y = from_layout.y
      let to_x = to_layout.x + to_layout.width / 2
      let to_y = to_layout.y + to_layout.height
      // Draw vertical line
      let line_char = match edge.style() {
        @types.Solid | @types.Open => style.v_line
        @types.Dotted => ':'
        @types.Thick => style.v_line
      }
      for y = to_y; y < from_y; y = y + 1 {
        canvas.set(from_x, y, line_char)
      }
      // Draw arrow
      if edge.style() != @types.Open {
        canvas.set(to_x, to_y + 1, style.arrow_up)
      }
    }
    @types.LR => {
      // Horizontal flow (left to right)
      let from_x = from_layout.x + from_layout.width
      let from_y = from_layout.y + from_layout.height / 2
      let to_x = to_layout.x
      let to_y = to_layout.y + to_layout.height / 2
      // Draw horizontal line
      let line_char = match edge.style() {
        @types.Solid | @types.Open => style.h_line
        @types.Dotted => '-'
        @types.Thick => style.h_line
      }
      for x = from_x; x < to_x; x = x + 1 {
        canvas.set(x, from_y, line_char)
      }
      // Draw arrow
      if edge.style() != @types.Open {
        canvas.set(to_x - 1, to_y, style.arrow_right)
      }
      // Draw label if present
      match edge.label() {
        Some(label) => {
          let mid_x = (from_x + to_x) / 2 -
            @width.string_width(label, canvas.ambiguous_width) / 2
          canvas.draw_text(mid_x, from_y - 1, label)
        }
        None => ()
      }
    }
    @types.RL => {
      // Horizontal flow (right to left)
      let from_x = from_layout.x
      let from_y = from_layout.y + from_layout.height / 2
      let to_x = to_layout.x + to_layout.width
      let to_y = to_layout.y + to_layout.height / 2
      // Draw horizontal line
      let line_char = match edge.style() {
        @types.Solid | @types.Open => style.h_line
        @types.Dotted => '-'
        @types.Thick => style.h_line
      }
      for x = to_x; x < from_x; x = x + 1 {
        canvas.set(x, from_y, line_char)
      }
      // Draw arrow
      if edge.style() != @types.Open {
        canvas.set(to_x + 1, to_y, style.arrow_left)
      }
    }
  }
}

/// Calculate layout for all nodes
fn calc_layout(
  flowchart : @types.Flowchart,
  ambiguous_width : Int
) -> (Array[NodeLayout], Int, Int) {
  let layouts : Array[NodeLayout] = []
  let direction = flowchart.direction()
  let nodes = flowchart.nodes()
  let is_horizontal = direction.is_horizontal()
  // Calculate node sizes first
  let node_sizes : Array[(Int, Int)] = []
  for node in nodes {
    node_sizes.push(calc_node_size(node, ambiguous_width))
  }
  // Simple layout algorithm
  let spacing = 4  // Space between nodes
  let mut current_pos = 0
  let mut max_cross = 0  // Max size in cross direction
  for i = 0; i < nodes.length(); i = i + 1 {
    let node = nodes[i]
    let (w, h) = node_sizes[i]
    let (x, y) = if is_horizontal {
      max_cross = if h > max_cross { h } else { max_cross }
      let pos = current_pos
      current_pos = current_pos + w + spacing
      (pos, 0)
    } else {
      max_cross = if w > max_cross { w } else { max_cross }
      let pos = current_pos
      current_pos = current_pos + h + spacing
      (0, pos)
    }
    layouts.push({ node, x, y, width: w, height: h })
  }
  // Center nodes on cross axis
  for i = 0; i < layouts.length(); i = i + 1 {
    let layout = layouts[i]
    if is_horizontal {
      let centered_y = (max_cross - layout.height) / 2
      layouts[i] = { ..layout, y: centered_y }
    } else {
      let centered_x = (max_cross - layout.width) / 2
      layouts[i] = { ..layout, x: centered_x }
    }
  }
  // Calculate canvas size
  let (canvas_width, canvas_height) = if is_horizontal {
    (current_pos, max_cross)
  } else {
    (max_cross, current_pos)
  }
  (layouts, canvas_width, canvas_height)
}

/// Find layout by node ID
fn find_layout(layouts : Array[NodeLayout], id : String) -> NodeLayout? {
  for layout in layouts {
    if layout.node.id() == id {
      return Some(layout)
    }
  }
  None
}

/// Render flowchart to ASCII art
pub fn render(
  flowchart : @types.Flowchart,
  options : @types.RenderOptions
) -> String {
  let nodes = flowchart.nodes()
  if nodes.is_empty() {
    return ""
  }
  let (layouts, canvas_width, canvas_height) = calc_layout(
    flowchart,
    options.ambiguous_width(),
  )
  let canvas = Canvas::new(canvas_width, canvas_height, options.ambiguous_width())
  let style = get_style(options.charset())
  // Draw edges first (so nodes draw over them)
  for edge in flowchart.edges() {
    match (find_layout(layouts, edge.from()), find_layout(layouts, edge.to())) {
      (Some(from_l), Some(to_l)) =>
        draw_edge(canvas, from_l, to_l, edge, style, flowchart.direction())
      _ => ()
    }
  }
  // Draw nodes
  for layout in layouts {
    draw_node(canvas, layout, style, options.ambiguous_width())
  }
  canvas.to_string()
}
