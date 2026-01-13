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

    /// Use Unicode box drawing characters
    #[arg(short, long)]
    unicode: bool,
}

#[derive(Debug, Clone, Copy)]
struct BoxChars {
    h_line: char,
    v_line: char,
    corner_tl: char,
    corner_tr: char,
    corner_bl: char,
    corner_br: char,
    arrow_down: char,
    arrow_up: char,
    arrow_right: char,
    arrow_left: char,
    thick_h: char,
    thick_v: char,
    round_tl: char,
    round_tr: char,
    round_bl: char,
    round_br: char,
}

impl BoxChars {
    fn ascii() -> Self {
        BoxChars {
            h_line: '-',
            v_line: '|',
            corner_tl: '+',
            corner_tr: '+',
            corner_bl: '+',
            corner_br: '+',
            arrow_down: 'v',
            arrow_up: '^',
            arrow_right: '>',
            arrow_left: '<',
            thick_h: '=',
            thick_v: '#',
            round_tl: '/',
            round_tr: '\\',
            round_bl: '\\',
            round_br: '/',
        }
    }

    fn unicode() -> Self {
        BoxChars {
            h_line: '─',
            v_line: '│',
            corner_tl: '┌',
            corner_tr: '┐',
            corner_bl: '└',
            corner_br: '┘',
            arrow_down: '▼',
            arrow_up: '▲',
            arrow_right: '▶',
            arrow_left: '◀',
            thick_h: '━',
            thick_v: '┃',
            round_tl: '╭',
            round_tr: '╮',
            round_bl: '╰',
            round_br: '╯',
        }
    }
}

#[derive(Debug, Clone)]
enum NodeShape {
    Rectangle,   // [text]
    Round,       // (text)
    Diamond,     // {text}
    Hexagon,     // {{text}}
    Circle,      // ((text))
    Stadium,     // ([text])
    Parallelogram, // [/text/]
    Trapezoid,   // [/text\]
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

#[derive(Debug, Clone, Copy)]
enum ArrowType {
    Normal,      // -->
    Dotted,      // -.->
    Thick,       // ==>
    Open,        // ---
}

#[derive(Debug, Clone, Copy)]
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
    // [/text/] - Parallelogram
    if text.starts_with("[/") && text.ends_with("/]") {
        return (text[2..text.len()-2].to_string(), NodeShape::Parallelogram);
    }
    // [/text\] or [\text/] - Trapezoid
    if (text.starts_with("[/") && text.ends_with("\\]")) ||
       (text.starts_with("[\\") && text.ends_with("/]")) {
        return (text[2..text.len()-2].to_string(), NodeShape::Trapezoid);
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
    } else if arrow.contains("-->") || arrow.contains("->") {
        ArrowType::Normal
    } else {
        ArrowType::Open
    }
}

fn parse_mermaid(input: &str) -> Option<Flowchart> {
    // Normalize input: replace semicolons with newlines
    let normalized = input.replace(';', "\n");

    let lines: Vec<&str> = normalized.lines().collect();

    if lines.is_empty() {
        return None;
    }

    // Find direction from the first line(s)
    let mut direction = Direction::TopDown;
    let mut start_line = 0;

    for (i, line) in lines.iter().enumerate() {
        let line_lower = line.trim().to_lowercase();
        if line_lower.starts_with("graph") || line_lower.starts_with("flowchart") {
            // Extract direction from this line
            if line_lower.contains(" lr") || line_lower.ends_with("lr") {
                direction = Direction::LeftRight;
            } else if line_lower.contains(" rl") || line_lower.ends_with("rl") {
                direction = Direction::RightLeft;
            } else if line_lower.contains(" bt") || line_lower.ends_with("bt") {
                direction = Direction::BottomUp;
            }
            start_line = i + 1;
            break;
        }
    }

    let mut nodes: HashMap<String, Node> = HashMap::new();
    let mut edges: Vec<Edge> = Vec::new();

    // Simpler regex: match node with optional shape
    // Order matters: check double-brackets first, then single brackets
    // Use .*? for non-greedy matching inside nested brackets
    let node_re = Regex::new(r"(\w+)(\{\{.*?\}\}|\(\(.*?\)\)|\(\[.*?\]\)|\[[^\]]*\]|\([^)]*\)|\{[^}]*\})?").unwrap();
    // Match arrows with optional labels
    let arrow_re = Regex::new(r"\s*(-->|-.->|==>|---|->)(\|([^|]*)\|)?\s*").unwrap();

    for line in lines.iter().skip(start_line) {
        let line = line.trim();
        // Skip direction-only lines, comments, subgraph markers
        if line.is_empty() || line.starts_with("%%") || line.starts_with("subgraph") || line == "end"
            || line.eq_ignore_ascii_case("td") || line.eq_ignore_ascii_case("tb")
            || line.eq_ignore_ascii_case("lr") || line.eq_ignore_ascii_case("rl")
            || line.eq_ignore_ascii_case("bt") {
            continue;
        }

        // Parse line: find all nodes and arrows
        let mut pos = 0;
        let mut prev_node: Option<String> = None;
        let mut prev_arrow: Option<(ArrowType, Option<String>)> = None;

        while pos < line.len() {
            let remaining = &line[pos..];

            // Try to match a node first
            if let Some(node_caps) = node_re.captures(remaining) {
                if node_caps.get(0).unwrap().start() == 0 {
                    let node_id = node_caps.get(1).unwrap().as_str().to_string();
                    let node_shape_str = node_caps.get(2).map(|m| m.as_str());

                    // Register node if not exists
                    if !nodes.contains_key(&node_id) {
                        let (label, shape) = if let Some(s) = node_shape_str {
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

                    // Create edge from previous node if we had an arrow
                    if let (Some(from_id), Some((arrow_type, label))) = (&prev_node, prev_arrow.take()) {
                        edges.push(Edge {
                            from: from_id.clone(),
                            to: node_id.clone(),
                            label,
                            arrow_type,
                        });
                    }

                    prev_node = Some(node_id);
                    pos += node_caps.get(0).unwrap().end();
                    continue;
                }
            }

            // Try to match an arrow
            if let Some(arrow_caps) = arrow_re.captures(remaining) {
                if arrow_caps.get(0).unwrap().start() == 0 {
                    let arrow_str = arrow_caps.get(1).unwrap().as_str();
                    let label = arrow_caps.get(3).map(|m| m.as_str().to_string());
                    prev_arrow = Some((parse_arrow(arrow_str), label));
                    pos += arrow_caps.get(0).unwrap().end();
                    continue;
                }
            }

            // Skip whitespace
            if remaining.starts_with(char::is_whitespace) {
                pos += 1;
                continue;
            }

            // Unknown character, skip
            pos += 1;
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
    chars: BoxChars,
}

impl Canvas {
    fn new(width: usize, height: usize, unicode: bool) -> Self {
        Canvas {
            width,
            height,
            data: vec![vec![' '; width]; height],
            chars: if unicode { BoxChars::unicode() } else { BoxChars::ascii() },
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
        let c = self.chars; // Copy the struct
        match shape {
            NodeShape::Rectangle => {
                self.set(x, y, c.corner_tl);
                for i in 1..width-1 {
                    self.set(x + i, y, c.h_line);
                }
                self.set(x + width - 1, y, c.corner_tr);

                for j in 1..height-1 {
                    self.set(x, y + j, c.v_line);
                    self.set(x + width - 1, y + j, c.v_line);
                }

                self.set(x, y + height - 1, c.corner_bl);
                for i in 1..width-1 {
                    self.set(x + i, y + height - 1, c.h_line);
                }
                self.set(x + width - 1, y + height - 1, c.corner_br);
            }
            NodeShape::Round => {
                self.set(x, y, c.round_tl);
                for i in 1..width-1 {
                    self.set(x + i, y, c.h_line);
                }
                self.set(x + width - 1, y, c.round_tr);

                for j in 1..height-1 {
                    self.set(x, y + j, c.v_line);
                    self.set(x + width - 1, y + j, c.v_line);
                }

                self.set(x, y + height - 1, c.round_bl);
                for i in 1..width-1 {
                    self.set(x + i, y + height - 1, c.h_line);
                }
                self.set(x + width - 1, y + height - 1, c.round_br);
            }
            NodeShape::Diamond => {
                let mid_x = width / 2;
                self.set(x + mid_x, y, c.arrow_up);
                for j in 1..height-1 {
                    self.set(x + mid_x - j.min(mid_x), y + j, '/');
                    self.set(x + mid_x + j.min(mid_x), y + j, '\\');
                }
                self.set(x + mid_x, y + height - 1, c.arrow_down);
            }
            NodeShape::Hexagon => {
                self.set(x + 1, y, '/');
                for i in 2..width-2 {
                    self.set(x + i, y, c.h_line);
                }
                self.set(x + width - 2, y, '\\');

                for j in 1..height-1 {
                    self.set(x, y + j, c.v_line);
                    self.set(x + width - 1, y + j, c.v_line);
                }

                self.set(x + 1, y + height - 1, '\\');
                for i in 2..width-2 {
                    self.set(x + i, y + height - 1, c.h_line);
                }
                self.set(x + width - 2, y + height - 1, '/');
            }
            NodeShape::Circle | NodeShape::Stadium => {
                self.set(x, y, '(');
                for i in 1..width-1 {
                    self.set(x + i, y, c.h_line);
                }
                self.set(x + width - 1, y, ')');

                for j in 1..height-1 {
                    self.set(x, y + j, '(');
                    self.set(x + width - 1, y + j, ')');
                }

                self.set(x, y + height - 1, '(');
                for i in 1..width-1 {
                    self.set(x + i, y + height - 1, c.h_line);
                }
                self.set(x + width - 1, y + height - 1, ')');
            }
            NodeShape::Parallelogram => {
                self.set(x + 1, y, '/');
                for i in 2..width-1 {
                    self.set(x + i, y, c.h_line);
                }
                self.set(x + width - 1, y, '/');

                for j in 1..height-1 {
                    self.set(x, y + j, '/');
                    self.set(x + width - 1, y + j, '/');
                }

                self.set(x, y + height - 1, '/');
                for i in 1..width-2 {
                    self.set(x + i, y + height - 1, c.h_line);
                }
                self.set(x + width - 2, y + height - 1, '/');
            }
            NodeShape::Trapezoid => {
                self.set(x + 1, y, '/');
                for i in 2..width-2 {
                    self.set(x + i, y, c.h_line);
                }
                self.set(x + width - 2, y, '\\');

                for j in 1..height-1 {
                    self.set(x, y + j, c.v_line);
                    self.set(x + width - 1, y + j, c.v_line);
                }

                self.set(x, y + height - 1, c.corner_bl);
                for i in 1..width-1 {
                    self.set(x + i, y + height - 1, c.h_line);
                }
                self.set(x + width - 1, y + height - 1, c.corner_br);
            }
        }
    }

    fn draw_line_vertical(&mut self, x: usize, y1: usize, y2: usize, dotted: bool, thick: bool) {
        let (start, end) = if y1 < y2 { (y1, y2) } else { (y2, y1) };
        let ch = if thick { self.chars.thick_v } else { self.chars.v_line };
        for y in start..=end {
            if dotted && (y - start) % 2 == 1 {
                continue;
            }
            self.set(x, y, ch);
        }
    }

    fn draw_line_horizontal(&mut self, x1: usize, x2: usize, y: usize, dotted: bool, thick: bool) {
        let (start, end) = if x1 < x2 { (x1, x2) } else { (x2, x1) };
        let ch = if thick { self.chars.thick_h } else { self.chars.h_line };
        for x in start..=end {
            if dotted && (x - start) % 2 == 1 {
                continue;
            }
            self.set(x, y, ch);
        }
    }

    fn draw_arrow_down(&mut self, x: usize, y: usize) {
        self.set(x, y, self.chars.arrow_down);
    }

    fn draw_arrow_up(&mut self, x: usize, y: usize) {
        self.set(x, y, self.chars.arrow_up);
    }

    fn draw_arrow_right(&mut self, x: usize, y: usize) {
        self.set(x, y, self.chars.arrow_right);
    }

    fn draw_arrow_left(&mut self, x: usize, y: usize) {
        self.set(x, y, self.chars.arrow_left);
    }

    fn render(&self) -> String {
        let lines: Vec<String> = self.data
            .iter()
            .map(|row| row.iter().collect::<String>().trim_end().to_string())
            .collect();

        let start = lines.iter().position(|l| !l.is_empty()).unwrap_or(0);
        let end = lines.iter().rposition(|l| !l.is_empty()).unwrap_or(lines.len());

        lines[start..=end].join("\n")
    }
}

fn render_flowchart(flowchart: &Flowchart, unicode: bool) -> String {
    if flowchart.nodes.is_empty() {
        return String::from("(empty diagram)");
    }

    let node_width = 16;
    let node_height = 3;
    let h_spacing = 20;
    let v_spacing = 6;

    let mut levels: Vec<Vec<String>> = Vec::new();
    let mut placed: HashMap<String, (usize, usize)> = HashMap::new();

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

    let mut queue: Vec<(String, usize)> = roots.iter().map(|r| (r.clone(), 0)).collect();
    if queue.is_empty() && !flowchart.nodes.is_empty() {
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

        for edge in &flowchart.edges {
            if edge.from == node_id && !placed.contains_key(&edge.to) {
                queue.push((edge.to.clone(), level + 1));
            }
        }
    }

    for node_id in flowchart.nodes.keys() {
        if !placed.contains_key(node_id) {
            let level = levels.len();
            levels.push(vec![node_id.clone()]);
            placed.insert(node_id.clone(), (level, 0));
        }
    }

    let max_nodes_per_level = levels.iter().map(|l| l.len()).max().unwrap_or(1);
    let num_levels = levels.len();

    let is_horizontal = matches!(flowchart.direction, Direction::LeftRight | Direction::RightLeft);

    let (canvas_width, canvas_height) = if is_horizontal {
        (num_levels * h_spacing + node_width + 4, max_nodes_per_level * v_spacing + node_height + 4)
    } else {
        (max_nodes_per_level * h_spacing + node_width + 4, num_levels * v_spacing + node_height + 4)
    };

    let mut canvas = Canvas::new(canvas_width, canvas_height, unicode);

    let mut node_positions: HashMap<String, (usize, usize, usize, usize)> = HashMap::new();

    for (level, nodes_in_level) in levels.iter().enumerate() {
        for (pos, node_id) in nodes_in_level.iter().enumerate() {
            let node = flowchart.nodes.get(node_id).unwrap();
            let label_width = node.label.chars().count().max(4);
            let box_width = label_width + 4;
            let box_height = node_height;

            let (x, y) = if is_horizontal {
                (2 + level * h_spacing, 2 + pos * v_spacing)
            } else {
                let level_width = nodes_in_level.len();
                let start_x = (canvas_width - level_width * h_spacing) / 2;
                (start_x + pos * h_spacing, 2 + level * v_spacing)
            };

            node_positions.insert(node_id.clone(), (x, y, box_width, box_height));

            canvas.draw_box(x, y, box_width, box_height, &node.shape);
            canvas.draw_text(x + 2, y + 1, &node.label);
        }
    }

    for edge in &flowchart.edges {
        if let (Some(&(from_x, from_y, from_w, from_h)), Some(&(to_x, to_y, to_w, to_h))) =
            (node_positions.get(&edge.from), node_positions.get(&edge.to)) {

            let (dotted, thick) = match edge.arrow_type {
                ArrowType::Dotted => (true, false),
                ArrowType::Thick => (false, true),
                _ => (false, false),
            };

            if is_horizontal {
                let start_x = from_x + from_w;
                let end_x = to_x;
                let y = from_y + from_h / 2;
                let to_y_mid = to_y + to_h / 2;

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
                let start_x = from_x + from_w / 2;
                let start_y = from_y + from_h;
                let end_x = to_x + to_w / 2;
                let end_y = to_y;

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

                if let Some(label) = &edge.label {
                    let label_x = (start_x.min(end_x) + start_x.max(end_x)) / 2;
                    let label_y = (start_y + end_y) / 2;
                    canvas.draw_text(label_x.saturating_sub(label.len() / 2), label_y, label);
                }
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
            let ascii_art = render_flowchart(&flowchart, cli.unicode);
            println!("{}", ascii_art);
        }
        None => {
            eprintln!("Error: Could not parse mermaid diagram");
            std::process::exit(1);
        }
    }
}
