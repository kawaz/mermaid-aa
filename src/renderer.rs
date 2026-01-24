//! ASCII Art renderer for flowcharts.

use crate::types::{Direction, EdgeStyle, Flowchart, Node, NodeShape};
use crate::width::{str_width, AmbiguousWidth};

/// Character set for rendering.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Default)]
pub enum Charset {
    Ascii,
    #[default]
    Unicode,
    UnicodeRound,
    UnicodeBold,
    UnicodeDouble,
}

impl std::str::FromStr for Charset {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s.to_lowercase().as_str() {
            "ascii" => Ok(Charset::Ascii),
            "unicode" => Ok(Charset::Unicode),
            "unicode-round" => Ok(Charset::UnicodeRound),
            "unicode-bold" => Ok(Charset::UnicodeBold),
            "unicode-double" => Ok(Charset::UnicodeDouble),
            _ => Err(format!("Unknown charset: {}", s)),
        }
    }
}

/// Drawing characters for a charset.
#[derive(Debug, Clone)]
pub struct DrawChars {
    pub h_line: char,
    pub v_line: char,
    pub corner_tl: char,
    pub corner_tr: char,
    pub corner_bl: char,
    pub corner_br: char,
    #[allow(dead_code)]
    pub cross: char,
    pub arrow_right: char,
    pub arrow_down: char,
    pub arrow_left: char,
    pub arrow_up: char,
    pub round_tl: char,
    pub round_tr: char,
    pub round_bl: char,
    pub round_br: char,
}

impl DrawChars {
    pub fn for_charset(charset: Charset) -> Self {
        match charset {
            Charset::Ascii => Self {
                h_line: '-',
                v_line: '|',
                corner_tl: '+',
                corner_tr: '+',
                corner_bl: '+',
                corner_br: '+',
                cross: '+',
                arrow_right: '>',
                arrow_down: 'v',
                arrow_left: '<',
                arrow_up: '^',
                round_tl: '(',
                round_tr: ')',
                round_bl: '(',
                round_br: ')',
            },
            Charset::Unicode => Self {
                h_line: '─',
                v_line: '│',
                corner_tl: '┌',
                corner_tr: '┐',
                corner_bl: '└',
                corner_br: '┘',
                cross: '┼',
                arrow_right: '→',
                arrow_down: '↓',
                arrow_left: '←',
                arrow_up: '↑',
                round_tl: '╭',
                round_tr: '╮',
                round_bl: '╰',
                round_br: '╯',
            },
            Charset::UnicodeRound => Self {
                h_line: '─',
                v_line: '│',
                corner_tl: '╭',
                corner_tr: '╮',
                corner_bl: '╰',
                corner_br: '╯',
                cross: '┼',
                arrow_right: '→',
                arrow_down: '↓',
                arrow_left: '←',
                arrow_up: '↑',
                round_tl: '╭',
                round_tr: '╮',
                round_bl: '╰',
                round_br: '╯',
            },
            Charset::UnicodeBold => Self {
                h_line: '━',
                v_line: '┃',
                corner_tl: '┏',
                corner_tr: '┓',
                corner_bl: '┗',
                corner_br: '┛',
                cross: '╋',
                arrow_right: '➤',
                arrow_down: '▼',
                arrow_left: '◄',
                arrow_up: '▲',
                round_tl: '┏',
                round_tr: '┓',
                round_bl: '┗',
                round_br: '┛',
            },
            Charset::UnicodeDouble => Self {
                h_line: '═',
                v_line: '║',
                corner_tl: '╔',
                corner_tr: '╗',
                corner_bl: '╚',
                corner_br: '╝',
                cross: '╬',
                arrow_right: '⇒',
                arrow_down: '⇓',
                arrow_left: '⇐',
                arrow_up: '⇑',
                round_tl: '╔',
                round_tr: '╗',
                round_bl: '╚',
                round_br: '╝',
            },
        }
    }
}

/// Render options.
#[derive(Debug, Clone)]
pub struct RenderOptions {
    pub charset: Charset,
    pub ambiguous_width: AmbiguousWidth,
    pub direction_override: Option<Direction>,
}

impl Default for RenderOptions {
    fn default() -> Self {
        Self {
            charset: Charset::Unicode,
            ambiguous_width: AmbiguousWidth::Half,
            direction_override: None,
        }
    }
}

/// Sentinel character for wide char continuation (not rendered).
const WIDE_CHAR_CONT: char = '\x00';

/// A 2D canvas for drawing.
#[derive(Debug, Clone)]
pub struct Canvas {
    cells: Vec<Vec<char>>,
    width: usize,
    height: usize,
    ambiguous_width: AmbiguousWidth,
}

impl Canvas {
    pub fn new(width: usize, height: usize, ambiguous_width: AmbiguousWidth) -> Self {
        Self {
            cells: vec![vec![' '; width]; height],
            width,
            height,
            ambiguous_width,
        }
    }

    pub fn set(&mut self, x: usize, y: usize, c: char) {
        if x < self.width && y < self.height {
            self.cells[y][x] = c;
            // Mark continuation cells for wide characters
            let w = crate::width::char_width(c, self.ambiguous_width);
            for i in 1..w {
                if x + i < self.width {
                    self.cells[y][x + i] = WIDE_CHAR_CONT;
                }
            }
        }
    }

    pub fn put_str(&mut self, x: usize, y: usize, s: &str, aw: AmbiguousWidth) {
        let mut col = x;
        for c in s.chars() {
            if col < self.width && y < self.height {
                self.cells[y][col] = c;
                let w = crate::width::char_width(c, aw);
                // Mark continuation cells for wide characters
                for i in 1..w {
                    if col + i < self.width {
                        self.cells[y][col + i] = WIDE_CHAR_CONT;
                    }
                }
                col += w;
            }
        }
    }

    pub fn h_line(&mut self, x: usize, y: usize, len: usize, c: char) {
        for i in 0..len {
            self.set(x + i, y, c);
        }
    }

    pub fn v_line(&mut self, x: usize, y: usize, len: usize, c: char) {
        for i in 0..len {
            self.set(x, y + i, c);
        }
    }

    pub fn render(&self) -> String {
        self.cells
            .iter()
            .map(|row| {
                row.iter()
                    .filter(|&&c| c != WIDE_CHAR_CONT)
                    .collect::<String>()
                    .trim_end()
                    .to_string()
            })
            .collect::<Vec<_>>()
            .join("\n")
    }
}

impl std::fmt::Display for Canvas {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.render())
    }
}

/// Rendered node with position.
#[derive(Debug, Clone)]
struct RenderedNode {
    id: String,
    x: usize,
    y: usize,
    width: usize,
    height: usize,
}

/// Render a flowchart to ASCII art.
pub fn render(flowchart: &Flowchart, options: &RenderOptions) -> String {
    let direction = options.direction_override.unwrap_or(flowchart.direction);
    let chars = DrawChars::for_charset(options.charset);
    let aw = options.ambiguous_width;

    // Calculate node sizes
    let node_sizes: Vec<(usize, usize)> = flowchart
        .nodes
        .iter()
        .map(|n| calculate_node_size(n, aw))
        .collect();

    // Layout nodes
    let rendered_nodes = layout_nodes(flowchart, &node_sizes, direction);

    // Calculate canvas size
    let (canvas_width, canvas_height) = calculate_canvas_size(&rendered_nodes);
    let mut canvas = Canvas::new(canvas_width + 2, canvas_height + 2, aw);

    // Draw edges first (so nodes draw over them)
    draw_edges(
        &mut canvas,
        flowchart,
        &rendered_nodes,
        &chars,
        direction,
        aw,
    );

    // Draw nodes
    for (i, node) in flowchart.nodes.iter().enumerate() {
        let rn = &rendered_nodes[i];
        draw_node(&mut canvas, node, rn.x, rn.y, &chars, aw);
    }

    canvas.render()
}

fn calculate_node_size(node: &Node, aw: AmbiguousWidth) -> (usize, usize) {
    let label_width = str_width(&node.label, aw);
    let padding = get_shape_padding(node.shape);
    let width = label_width + padding * 2;
    let height = 3; // All nodes are 3 lines tall
    (width.max(5), height)
}

fn get_shape_padding(shape: NodeShape) -> usize {
    match shape {
        NodeShape::Circle | NodeShape::Stadium | NodeShape::Hexagon => 3,
        _ => 2,
    }
}

fn layout_nodes(
    flowchart: &Flowchart,
    sizes: &[(usize, usize)],
    direction: Direction,
) -> Vec<RenderedNode> {
    let mut result = Vec::new();
    let spacing = 4;

    match direction {
        Direction::TB | Direction::BT => {
            // Vertical layout
            let mut y = 1;
            let max_width = sizes.iter().map(|(w, _)| *w).max().unwrap_or(10);

            for (i, node) in flowchart.nodes.iter().enumerate() {
                let (w, h) = sizes[i];
                let x = (max_width - w) / 2 + 1;
                result.push(RenderedNode {
                    id: node.id.clone(),
                    x,
                    y,
                    width: w,
                    height: h,
                });
                y += h + spacing;
            }

            if direction == Direction::BT {
                // Reverse Y positions for bottom-to-top
                let total_height = y - spacing;
                for rn in &mut result {
                    rn.y = total_height - rn.y - rn.height + 2;
                }
            }
        }
        Direction::LR | Direction::RL => {
            // Horizontal layout
            let mut x = 1;
            let max_height = sizes.iter().map(|(_, h)| *h).max().unwrap_or(3);

            for (i, node) in flowchart.nodes.iter().enumerate() {
                let (w, h) = sizes[i];
                let y = (max_height - h) / 2 + 1;
                result.push(RenderedNode {
                    id: node.id.clone(),
                    x,
                    y,
                    width: w,
                    height: h,
                });
                x += w + spacing + 2;
            }

            if direction == Direction::RL {
                // Reverse X positions for right-to-left
                let total_width = x - spacing - 2;
                for rn in &mut result {
                    rn.x = total_width - rn.x - rn.width + 2;
                }
            }
        }
    }

    result
}

fn calculate_canvas_size(nodes: &[RenderedNode]) -> (usize, usize) {
    let max_x = nodes.iter().map(|n| n.x + n.width).max().unwrap_or(10);
    let max_y = nodes.iter().map(|n| n.y + n.height).max().unwrap_or(5);
    (max_x + 1, max_y + 1)
}

fn draw_node(
    canvas: &mut Canvas,
    node: &Node,
    x: usize,
    y: usize,
    chars: &DrawChars,
    aw: AmbiguousWidth,
) {
    let label_width = str_width(&node.label, aw);
    let padding = get_shape_padding(node.shape);
    let width = label_width + padding * 2;

    match node.shape {
        NodeShape::Rectangle => {
            // Top border
            canvas.set(x, y, chars.corner_tl);
            canvas.h_line(x + 1, y, width - 2, chars.h_line);
            canvas.set(x + width - 1, y, chars.corner_tr);
            // Middle with label
            canvas.set(x, y + 1, chars.v_line);
            let label_x = x + (width - label_width) / 2;
            canvas.put_str(label_x, y + 1, &node.label, aw);
            canvas.set(x + width - 1, y + 1, chars.v_line);
            // Bottom border
            canvas.set(x, y + 2, chars.corner_bl);
            canvas.h_line(x + 1, y + 2, width - 2, chars.h_line);
            canvas.set(x + width - 1, y + 2, chars.corner_br);
        }
        NodeShape::Rounded => {
            canvas.set(x, y, chars.round_tl);
            canvas.h_line(x + 1, y, width - 2, chars.h_line);
            canvas.set(x + width - 1, y, chars.round_tr);
            canvas.set(x, y + 1, chars.v_line);
            let label_x = x + (width - label_width) / 2;
            canvas.put_str(label_x, y + 1, &node.label, aw);
            canvas.set(x + width - 1, y + 1, chars.v_line);
            canvas.set(x, y + 2, chars.round_bl);
            canvas.h_line(x + 1, y + 2, width - 2, chars.h_line);
            canvas.set(x + width - 1, y + 2, chars.round_br);
        }
        NodeShape::Stadium => {
            // ([text])
            canvas.set(x, y, '(');
            canvas.set(x + 1, y, '[');
            canvas.h_line(x + 2, y, width - 4, chars.h_line);
            canvas.set(x + width - 2, y, ']');
            canvas.set(x + width - 1, y, ')');
            // Middle
            canvas.set(x, y + 1, '(');
            canvas.set(x + 1, y + 1, ' ');
            let label_x = x + (width - label_width) / 2;
            canvas.put_str(label_x, y + 1, &node.label, aw);
            canvas.set(x + width - 2, y + 1, ' ');
            canvas.set(x + width - 1, y + 1, ')');
            // Bottom
            canvas.set(x, y + 2, '(');
            canvas.set(x + 1, y + 2, '[');
            canvas.h_line(x + 2, y + 2, width - 4, chars.h_line);
            canvas.set(x + width - 2, y + 2, ']');
            canvas.set(x + width - 1, y + 2, ')');
        }
        NodeShape::Diamond => {
            // Diamond shape: < text >
            let mid = width / 2;
            canvas.set(x + mid, y, '/');
            canvas.set(x + mid + 1, y, '\\');
            canvas.set(x, y + 1, '<');
            let label_x = x + (width - label_width) / 2;
            canvas.put_str(label_x, y + 1, &node.label, aw);
            canvas.set(x + width - 1, y + 1, '>');
            canvas.set(x + mid, y + 2, '\\');
            canvas.set(x + mid + 1, y + 2, '/');
        }
        NodeShape::Hexagon => {
            // {{text}}
            canvas.set(x, y, '{');
            canvas.set(x + 1, y, '{');
            canvas.h_line(x + 2, y, width - 4, chars.h_line);
            canvas.set(x + width - 2, y, '}');
            canvas.set(x + width - 1, y, '}');
            canvas.set(x, y + 1, '{');
            canvas.set(x + 1, y + 1, ' ');
            let label_x = x + (width - label_width) / 2;
            canvas.put_str(label_x, y + 1, &node.label, aw);
            canvas.set(x + width - 2, y + 1, ' ');
            canvas.set(x + width - 1, y + 1, '}');
            canvas.set(x, y + 2, '{');
            canvas.set(x + 1, y + 2, '{');
            canvas.h_line(x + 2, y + 2, width - 4, chars.h_line);
            canvas.set(x + width - 2, y + 2, '}');
            canvas.set(x + width - 1, y + 2, '}');
        }
        NodeShape::Circle => {
            // ((text))
            canvas.set(x, y, '(');
            canvas.set(x + 1, y, '(');
            canvas.h_line(x + 2, y, width - 4, chars.h_line);
            canvas.set(x + width - 2, y, ')');
            canvas.set(x + width - 1, y, ')');
            canvas.set(x, y + 1, '(');
            canvas.set(x + 1, y + 1, '(');
            let label_x = x + (width - label_width) / 2;
            canvas.put_str(label_x, y + 1, &node.label, aw);
            canvas.set(x + width - 2, y + 1, ')');
            canvas.set(x + width - 1, y + 1, ')');
            canvas.set(x, y + 2, '(');
            canvas.set(x + 1, y + 2, '(');
            canvas.h_line(x + 2, y + 2, width - 4, chars.h_line);
            canvas.set(x + width - 2, y + 2, ')');
            canvas.set(x + width - 1, y + 2, ')');
        }
        NodeShape::Parallelogram => {
            // /text/
            canvas.set(x, y, '/');
            canvas.h_line(x + 1, y, width - 2, chars.h_line);
            canvas.set(x + width - 1, y, '/');
            canvas.set(x, y + 1, '/');
            let label_x = x + (width - label_width) / 2;
            canvas.put_str(label_x, y + 1, &node.label, aw);
            canvas.set(x + width - 1, y + 1, '/');
            canvas.set(x, y + 2, '/');
            canvas.h_line(x + 1, y + 2, width - 2, chars.h_line);
            canvas.set(x + width - 1, y + 2, '/');
        }
        NodeShape::Trapezoid => {
            // /text\
            canvas.set(x, y, '/');
            canvas.h_line(x + 1, y, width - 2, chars.h_line);
            canvas.set(x + width - 1, y, '\\');
            canvas.set(x, y + 1, '/');
            let label_x = x + (width - label_width) / 2;
            canvas.put_str(label_x, y + 1, &node.label, aw);
            canvas.set(x + width - 1, y + 1, '\\');
            canvas.set(x, y + 2, '/');
            canvas.h_line(x + 1, y + 2, width - 2, chars.h_line);
            canvas.set(x + width - 1, y + 2, '\\');
        }
    }
}

fn draw_edges(
    canvas: &mut Canvas,
    flowchart: &Flowchart,
    nodes: &[RenderedNode],
    chars: &DrawChars,
    direction: Direction,
    aw: AmbiguousWidth,
) {
    for edge in &flowchart.edges {
        let from_node = nodes.iter().find(|n| n.id == edge.from);
        let to_node = nodes.iter().find(|n| n.id == edge.to);

        if let (Some(from), Some(to)) = (from_node, to_node) {
            draw_edge(
                canvas,
                from,
                to,
                &edge.style,
                &edge.label,
                chars,
                direction,
                aw,
            );
        }
    }
}

#[allow(clippy::too_many_arguments)]
fn draw_edge(
    canvas: &mut Canvas,
    from: &RenderedNode,
    to: &RenderedNode,
    style: &EdgeStyle,
    label: &Option<String>,
    chars: &DrawChars,
    direction: Direction,
    aw: AmbiguousWidth,
) {
    let (line_char, arrow_h, arrow_v_down, arrow_v_up) = match style {
        EdgeStyle::SolidArrow => (
            chars.h_line,
            chars.arrow_right,
            chars.arrow_down,
            chars.arrow_up,
        ),
        EdgeStyle::DottedArrow => ('-', '>', 'v', '^'),
        EdgeStyle::ThickArrow => ('=', '>', 'V', '^'),
        EdgeStyle::Open => (chars.h_line, chars.h_line, chars.v_line, chars.v_line),
    };

    match direction {
        Direction::TB => {
            // Vertical: connect bottom of 'from' to top of 'to'
            let x = from.x + from.width / 2;
            let y1 = from.y + from.height;
            let y2 = to.y - 1;

            if y2 > y1 {
                canvas.v_line(x, y1, y2 - y1, chars.v_line);
            }
            canvas.set(x, y2, arrow_v_down);

            // Draw label if present
            if let Some(lbl) = label {
                let lbl_y = (y1 + y2) / 2;
                let lbl_x = x + 2;
                canvas.put_str(lbl_x, lbl_y, lbl, aw);
            }
        }
        Direction::BT => {
            // Bottom to top: connect top of 'from' to bottom of 'to'
            let x = from.x + from.width / 2;
            let y1 = from.y - 1;
            let y2 = to.y + to.height;

            if y1 > y2 {
                canvas.v_line(x, y2 + 1, y1 - y2, chars.v_line);
            }
            canvas.set(x, y2, arrow_v_up);
        }
        Direction::LR => {
            // Horizontal: connect right of 'from' to left of 'to'
            let y = from.y + from.height / 2;
            let x1 = from.x + from.width;
            let x2 = to.x - 1;

            for x in x1..=x2 {
                if x == x2 {
                    canvas.set(x, y, arrow_h);
                } else {
                    canvas.set(x, y, line_char);
                }
            }

            // Draw label if present
            if let Some(lbl) = label {
                let lbl_x = (x1 + x2) / 2 - str_width(lbl, aw) / 2;
                let lbl_y = y - 1;
                canvas.put_str(lbl_x, lbl_y, lbl, aw);
            }
        }
        Direction::RL => {
            // Right to left: connect left of 'from' to right of 'to'
            let y = from.y + from.height / 2;
            let x1 = to.x + to.width;
            let x2 = from.x - 1;

            for x in x1..=x2 {
                if x == x1 {
                    canvas.set(x, y, chars.arrow_left);
                } else {
                    canvas.set(x, y, line_char);
                }
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::parser::parse;

    #[test]
    fn test_render_simple() {
        let fc = parse("graph TD\n  A --> B").unwrap();
        let output = render(&fc, &RenderOptions::default());
        assert!(output.contains('A'));
        assert!(output.contains('B'));
    }

    #[test]
    fn test_render_lr() {
        let fc = parse("graph LR\n  A --> B").unwrap();
        let output = render(&fc, &RenderOptions::default());
        assert!(output.contains('→'));
    }

    #[test]
    fn test_render_ascii() {
        let fc = parse("graph TD\n  A --> B").unwrap();
        let options = RenderOptions {
            charset: Charset::Ascii,
            ..Default::default()
        };
        let output = render(&fc, &options);
        assert!(output.contains('+'));
        assert!(output.contains('-'));
    }

    #[test]
    fn test_canvas() {
        let mut canvas = Canvas::new(10, 5, AmbiguousWidth::Half);
        canvas.set(5, 0, 'X');
        canvas.h_line(0, 1, 5, '-');
        canvas.v_line(0, 2, 3, '|');
        let s = canvas.render();
        assert!(s.contains('X'));
        assert!(s.contains("-----"));
        assert!(s.contains('|'));
    }

    #[test]
    fn test_canvas_wide_chars() {
        let mut canvas = Canvas::new(10, 3, AmbiguousWidth::Half);
        canvas.put_str(0, 0, "日本語", AmbiguousWidth::Half);
        let s = canvas.render();
        assert!(s.contains("日本語"));
        // Each Japanese char is width 2, so total display width is 6
    }
}
