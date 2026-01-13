use clap::Parser;
use regex::Regex;
use std::collections::HashMap;
use std::io::{self, BufRead};

#[derive(Parser)]
#[command(name = "mermaid-aa")]
#[command(about = "Convert Mermaid diagrams to ASCII Art", long_about = None)]
struct Cli {
    /// Mermaid text (if not provided, reads from stdin)
    #[arg(short, long)]
    input: Option<String>,

    /// Input file path
    #[arg(short, long)]
    file: Option<String>,
}

#[derive(Debug, Clone)]
enum NodeShape {
    Rectangle,   // [text]
    Round,       // (text)
    Diamond,     // {text}
    Hexagon,     // {{text}}
    Circle,      // ((text))
    Stadium,     // ([text])
}

#[derive(Debug, Clone)]
struct Node {
    id: String,
    label: String,
    shape: NodeShape,
}

#[derive(Debug, Clone)]
struct Edge {
    from: String,
    to: String,
    label: Option<String>,
    arrow_type: ArrowType,
}

#[derive(Debug, Clone)]
enum ArrowType {
    Normal,      // -->
    Dotted,      // -.->
    Thick,       // ==>
    Open,        // ---
}

#[derive(Debug, Clone)]
enum Direction {
    TopDown,    // TD, TB
    BottomUp,   // BT
    LeftRight,  // LR
    RightLeft,  // RL
}

#[derive(Debug)]
struct Flowchart {
    direction: Direction,
    nodes: HashMap<String, Node>,
    edges: Vec<Edge>,
}

fn parse_node_shape(text: &str) -> (String, NodeShape) {
    let text = text.trim();

    // {{text}} - Hexagon
    if text.starts_with("{{") && text.ends_with("}}") {
        return (text[2..text.len()-2].to_string(), NodeShape::Hexagon);
    }
    // ((text)) - Circle
    if text.starts_with("((") && text.ends_with("))") {
        return (text[2..text.len()-2].to_string(), NodeShape::Circle);
    }
    // ([text]) - Stadium
    if text.starts_with("([") && text.ends_with("])") {
        return (text[2..text.len()-2].to_string(), NodeShape::Stadium);
    }
    // {text} - Diamond
    if text.starts_with('{') && text.ends_with('}') {
        return (text[1..text.len()-1].to_string(), NodeShape::Diamond);
    }
    // (text) - Round
    if text.starts_with('(') && text.ends_with(')') {
        return (text[1..text.len()-1].to_string(), NodeShape::Round);
    }
    // [text] - Rectangle
    if text.starts_with('[') && text.ends_with(']') {
        return (text[1..text.len()-1].to_string(), NodeShape::Rectangle);
    }

    (text.to_string(), NodeShape::Rectangle)
}

fn parse_arrow(arrow: &str) -> ArrowType {
    if arrow.contains("-.->") || arrow.contains("-.-") {
        ArrowType::Dotted
    } else if arrow.contains("==>") || arrow.contains("===") {
        ArrowType::Thick
    } else if arrow.contains("-->") {
        ArrowType::Normal
    } else {
        ArrowType::Open
    }
}

fn parse_mermaid(input: &str) -> Option<Flowchart> {
    let lines: Vec<&str> = input.lines().collect();

    if lines.is_empty() {
        return None;
    }

    let first_line = lines[0].trim().to_lowercase();
    let direction = if first_line.contains("lr") {
        Direction::LeftRight
    } else if first_line.contains("rl") {
        Direction::RightLeft
    } else if first_line.contains("bt") {
        Direction::BottomUp
    } else {
        Direction::TopDown
    };

    let mut nodes: HashMap<String, Node> = HashMap::new();
    let mut edges: Vec<Edge> = Vec::new();

    // Regex for parsing connections
    let edge_re = Regex::new(r"(\w+)(\[.*?\]|\(.*?\)|\{.*?\})?(\s*)(-->|-.->|==>|---)(\|([^|]*)\|)?(\s*)(\w+)(\[.*?\]|\(.*?\)|\{.*?\})?").unwrap();
    let node_re = Regex::new(r"^\s*(\w+)(\[.*?\]|\(.*?\)|\{.*?\}|\(\[.*?\]\)|\(\(.*?\)\)|\{\{.*?\}\})?\s*$").unwrap();

    for line in lines.iter().skip(1) {
        let line = line.trim();
        if line.is_empty() || line.starts_with("%%") {
            continue;
        }

        // Parse edge with connection
        if let Some(caps) = edge_re.captures(line) {
            let from_id = caps.get(1).unwrap().as_str().to_string();
            let from_shape = caps.get(2).map(|m| m.as_str());
            let arrow = caps.get(4).unwrap().as_str();
            let edge_label = caps.get(6).map(|m| m.as_str().to_string());
            let to_id = caps.get(8).unwrap().as_str().to_string();
            let to_shape = caps.get(9).map(|m| m.as_str());

            // Register from node
            if !nodes.contains_key(&from_id) {
                let (label, shape) = if let Some(s) = from_shape {
                    parse_node_shape(s)
                } else {
                    (from_id.clone(), NodeShape::Rectangle)
                };
                nodes.insert(from_id.clone(), Node {
                    id: from_id.clone(),
                    label,
                    shape,
                });
            }

            // Register to node
            if !nodes.contains_key(&to_id) {
                let (label, shape) = if let Some(s) = to_shape {
                    parse_node_shape(s)
                } else {
                    (to_id.clone(), NodeShape::Rectangle)
                };
                nodes.insert(to_id.clone(), Node {
                    id: to_id.clone(),
                    label,
                    shape,
                });
            }

            edges.push(Edge {
                from: from_id,
                to: to_id,
                label: edge_label,
                arrow_type: parse_arrow(arrow),
            });
        } else if let Some(caps) = node_re.captures(line) {
            // Standalone node definition
            let node_id = caps.get(1).unwrap().as_str().to_string();
            let shape_text = caps.get(2).map(|m| m.as_str());

            if !nodes.contains_key(&node_id) {
                let (label, shape) = if let Some(s) = shape_text {
                    parse_node_shape(s)
                } else {
                    (node_id.clone(), NodeShape::Rectangle)
                };
                nodes.insert(node_id.clone(), Node {
                    id: node_id.clone(),
                    label,
                    shape,
                });
            }
        }
    }

    Some(Flowchart {
        direction,
        nodes,
        edges,
    })
}

struct Canvas {
    width: usize,
    height: usize,
    data: Vec<Vec<char>>,
}

impl Canvas {
    fn new(width: usize, height: usize) -> Self {
        Canvas {
            width,
            height,
            data: vec![vec![' '; width]; height],
        }
    }

    fn set(&mut self, x: usize, y: usize, c: char) {
        if x < self.width && y < self.height {
            self.data[y][x] = c;
        }
    }

    fn draw_text(&mut self, x: usize, y: usize, text: &str) {
        for (i, c) in text.chars().enumerate() {
            self.set(x + i, y, c);
        }
    }

    fn draw_box(&mut self, x: usize, y: usize, width: usize, height: usize, shape: &NodeShape) {
        match shape {
            NodeShape::Rectangle => {
                // Top border
                self.set(x, y, '+');
                for i in 1..width-1 {
                    self.set(x + i, y, '-');
                }
                self.set(x + width - 1, y, '+');

                // Sides
                for j in 1..height-1 {
                    self.set(x, y + j, '|');
                    self.set(x + width - 1, y + j, '|');
                }

                // Bottom border
                self.set(x, y + height - 1, '+');
                for i in 1..width-1 {
                    self.set(x + i, y + height - 1, '-');
                }
                self.set(x + width - 1, y + height - 1, '+');
            }
            NodeShape::Round => {
                // Top border
                self.set(x, y, '/');
                for i in 1..width-1 {
                    self.set(x + i, y, '-');
                }
                self.set(x + width - 1, y, '\\');

                // Sides
                for j in 1..height-1 {
                    self.set(x, y + j, '|');
                    self.set(x + width - 1, y + j, '|');
                }

                // Bottom border
                self.set(x, y + height - 1, '\\');
                for i in 1..width-1 {
                    self.set(x + i, y + height - 1, '-');
                }
                self.set(x + width - 1, y + height - 1, '/');
            }
            NodeShape::Diamond => {
                let mid_x = width / 2;
                let mid_y = height / 2;

                // Draw diamond shape
                for j in 0..height {
                    let dist = if j <= mid_y {
                        mid_y - j
                    } else {
                        j - mid_y
                    };

                    let left = mid_x.saturating_sub(mid_y.saturating_sub(dist));
                    let right = mid_x + mid_y.saturating_sub(dist);

                    if left < width {
                        if j == 0 || j == height - 1 {
                            self.set(x + mid_x, y + j, if j == 0 { '^' } else { 'v' });
                        } else {
                            self.set(x + left, y + j, '/');
                            if right < width {
                                self.set(x + right, y + j, '\\');
                            }
                        }
                    }
                }
            }
            NodeShape::Hexagon => {
                // Top
                self.set(x + 1, y, '/');
                for i in 2..width-2 {
                    self.set(x + i, y, '-');
                }
                self.set(x + width - 2, y, '\\');

                // Sides
                for j in 1..height-1 {
                    self.set(x, y + j, '|');
                    self.set(x + width - 1, y + j, '|');
                }

                // Bottom
                self.set(x + 1, y + height - 1, '\\');
                for i in 2..width-2 {
                    self.set(x + i, y + height - 1, '-');
                }
                self.set(x + width - 2, y + height - 1, '/');
            }
            NodeShape::Circle | NodeShape::Stadium => {
                // Top
                self.set(x, y, '(');
                for i in 1..width-1 {
                    self.set(x + i, y, '-');
                }
                self.set(x + width - 1, y, ')');

                // Sides
                for j in 1..height-1 {
                    self.set(x, y + j, '(');
                    self.set(x + width - 1, y + j, ')');
                }

                // Bottom
                self.set(x, y + height - 1, '(');
                for i in 1..width-1 {
                    self.set(x + i, y + height - 1, '-');
                }
                self.set(x + width - 1, y + height - 1, ')');
            }
        }
    }

    fn draw_line_vertical(&mut self, x: usize, y1: usize, y2: usize, dotted: bool, thick: bool) {
        let (start, end) = if y1 < y2 { (y1, y2) } else { (y2, y1) };
        for y in start..=end {
            let c = if thick {
                '#'
            } else if dotted && (y - start) % 2 == 1 {
                ' '
            } else {
                '|'
            };
            self.set(x, y, c);
        }
    }

    fn draw_line_horizontal(&mut self, x1: usize, x2: usize, y: usize, dotted: bool, thick: bool) {
        let (start, end) = if x1 < x2 { (x1, x2) } else { (x2, x1) };
        for x in start..=end {
            let c = if thick {
                '='
            } else if dotted && (x - start) % 2 == 1 {
                ' '
            } else {
                '-'
            };
            self.set(x, y, c);
        }
    }

    fn draw_arrow_down(&mut self, x: usize, y: usize) {
        self.set(x, y, 'v');
    }

    fn draw_arrow_up(&mut self, x: usize, y: usize) {
        self.set(x, y, '^');
    }

    fn draw_arrow_right(&mut self, x: usize, y: usize) {
        self.set(x, y, '>');
    }

    fn draw_arrow_left(&mut self, x: usize, y: usize) {
        self.set(x, y, '<');
    }

    fn render(&self) -> String {
        self.data
            .iter()
            .map(|row| row.iter().collect::<String>().trim_end().to_string())
            .collect::<Vec<_>>()
            .join("\n")
            .trim_end()
            .to_string()
    }
}

fn render_flowchart(flowchart: &Flowchart) -> String {
    if flowchart.nodes.is_empty() {
        return String::from("(empty diagram)");
    }

    // Calculate layout
    let node_width = 16;
    let node_height = 3;
    let h_spacing = 20;
    let v_spacing = 6;

    // Simple layout: arrange nodes in levels based on edges
    let mut levels: Vec<Vec<String>> = Vec::new();
    let mut placed: HashMap<String, (usize, usize)> = HashMap::new(); // node_id -> (level, position)

    // Find root nodes (nodes with no incoming edges)
    let mut has_incoming: HashMap<String, bool> = HashMap::new();
    for node_id in flowchart.nodes.keys() {
        has_incoming.insert(node_id.clone(), false);
    }
    for edge in &flowchart.edges {
        has_incoming.insert(edge.to.clone(), true);
    }

    let roots: Vec<String> = flowchart.nodes.keys()
        .filter(|id| !has_incoming.get(*id).unwrap_or(&false))
        .cloned()
        .collect();

    // BFS to assign levels
    let mut queue: Vec<(String, usize)> = roots.iter().map(|r| (r.clone(), 0)).collect();
    if queue.is_empty() && !flowchart.nodes.is_empty() {
        // If no roots found, start with first node
        queue.push((flowchart.nodes.keys().next().unwrap().clone(), 0));
    }

    while let Some((node_id, level)) = queue.pop() {
        if placed.contains_key(&node_id) {
            continue;
        }

        while levels.len() <= level {
            levels.push(Vec::new());
        }

        let pos = levels[level].len();
        levels[level].push(node_id.clone());
        placed.insert(node_id.clone(), (level, pos));

        // Add children
        for edge in &flowchart.edges {
            if edge.from == node_id && !placed.contains_key(&edge.to) {
                queue.push((edge.to.clone(), level + 1));
            }
        }
    }

    // Add any unplaced nodes
    for node_id in flowchart.nodes.keys() {
        if !placed.contains_key(node_id) {
            let level = levels.len();
            levels.push(vec![node_id.clone()]);
            placed.insert(node_id.clone(), (level, 0));
        }
    }

    // Calculate canvas size
    let max_nodes_per_level = levels.iter().map(|l| l.len()).max().unwrap_or(1);
    let num_levels = levels.len();

    let is_horizontal = matches!(flowchart.direction, Direction::LeftRight | Direction::RightLeft);

    let (canvas_width, canvas_height) = if is_horizontal {
        // Horizontal: levels go left-to-right, nodes in level go top-to-bottom
        (num_levels * h_spacing + node_width + 4, max_nodes_per_level * v_spacing + node_height + 4)
    } else {
        // Vertical: levels go top-to-bottom, nodes in level go left-to-right
        (max_nodes_per_level * h_spacing + node_width + 4, num_levels * v_spacing + node_height + 4)
    };

    let mut canvas = Canvas::new(canvas_width, canvas_height);

    // Calculate node positions
    let mut node_positions: HashMap<String, (usize, usize, usize, usize)> = HashMap::new(); // id -> (x, y, w, h)

    for (level, nodes_in_level) in levels.iter().enumerate() {
        for (pos, node_id) in nodes_in_level.iter().enumerate() {
            let node = flowchart.nodes.get(node_id).unwrap();
            let label_width = node.label.len().max(4);
            let box_width = label_width + 4;
            let box_height = node_height;

            let (x, y) = if is_horizontal {
                let x = 2 + level * h_spacing;
                let y = 2 + pos * v_spacing;
                (x, y)
            } else {
                let level_width = nodes_in_level.len();
                let start_x = (canvas_width - level_width * h_spacing) / 2;
                let x = start_x + pos * h_spacing;
                let y = 2 + level * v_spacing;
                (x, y)
            };

            node_positions.insert(node_id.clone(), (x, y, box_width, box_height));

            // Draw node
            canvas.draw_box(x, y, box_width, box_height, &node.shape);
            canvas.draw_text(x + 2, y + 1, &node.label);
        }
    }

    // Draw edges
    for edge in &flowchart.edges {
        let (from_x, from_y, from_w, from_h) = node_positions.get(&edge.from).unwrap();
        let (to_x, to_y, _to_w, _to_h) = node_positions.get(&edge.to).unwrap();

        let (dotted, thick) = match edge.arrow_type {
            ArrowType::Dotted => (true, false),
            ArrowType::Thick => (false, true),
            _ => (false, false),
        };

        if is_horizontal {
            // Horizontal layout: draw horizontal lines
            let start_x = from_x + from_w;
            let end_x = *to_x;
            let y = from_y + from_h / 2;
            let to_y_mid = to_y + _to_h / 2;

            if y == to_y_mid {
                canvas.draw_line_horizontal(start_x, end_x.saturating_sub(1), y, dotted, thick);
                canvas.draw_arrow_right(end_x.saturating_sub(1), y);
            } else {
                let mid_x = (start_x + end_x) / 2;
                canvas.draw_line_horizontal(start_x, mid_x, y, dotted, thick);
                canvas.draw_line_vertical(mid_x, y, to_y_mid, dotted, thick);
                canvas.draw_line_horizontal(mid_x, end_x.saturating_sub(1), to_y_mid, dotted, thick);
                canvas.draw_arrow_right(end_x.saturating_sub(1), to_y_mid);
            }
        } else {
            // Vertical layout: draw vertical lines
            let start_x = from_x + from_w / 2;
            let start_y = from_y + from_h;
            let end_x = to_x + _to_w / 2;
            let end_y = *to_y;

            if start_x == end_x {
                canvas.draw_line_vertical(start_x, start_y, end_y.saturating_sub(1), dotted, thick);
                canvas.draw_arrow_down(start_x, end_y.saturating_sub(1));
            } else {
                let mid_y = (start_y + end_y) / 2;
                canvas.draw_line_vertical(start_x, start_y, mid_y, dotted, thick);
                canvas.draw_line_horizontal(start_x, end_x, mid_y, dotted, thick);
                canvas.draw_line_vertical(end_x, mid_y, end_y.saturating_sub(1), dotted, thick);
                canvas.draw_arrow_down(end_x, end_y.saturating_sub(1));
            }

            // Draw edge label if present
            if let Some(label) = &edge.label {
                let label_x = (start_x.min(end_x) + start_x.max(end_x)) / 2;
                let label_y = (start_y + end_y) / 2;
                canvas.draw_text(label_x.saturating_sub(label.len() / 2), label_y, label);
            }
        }
    }

    canvas.render()
}

fn main() {
    let cli = Cli::parse();

    let input = if let Some(text) = cli.input {
        text
    } else if let Some(file_path) = cli.file {
        std::fs::read_to_string(&file_path)
            .unwrap_or_else(|e| {
                eprintln!("Error reading file '{}': {}", file_path, e);
                std::process::exit(1);
            })
    } else {
        // Read from stdin
        let stdin = io::stdin();
        let lines: Vec<String> = stdin.lock().lines()
            .map(|l| l.unwrap_or_default())
            .collect();
        lines.join("\n")
    };

    if input.trim().is_empty() {
        eprintln!("Error: No input provided");
        eprintln!("Usage: mermaid-aa -i 'graph TD; A-->B' or pipe input");
        std::process::exit(1);
    }

    match parse_mermaid(&input) {
        Some(flowchart) => {
            let ascii_art = render_flowchart(&flowchart);
            println!("{}", ascii_art);
        }
        None => {
            eprintln!("Error: Could not parse mermaid diagram");
            std::process::exit(1);
        }
    }
}
