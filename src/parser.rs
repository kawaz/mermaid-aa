//! Mermaid flowchart parser.

use crate::types::{Direction, Edge, EdgeStyle, Flowchart, Node, NodeShape};

/// Parse error with location information.
#[derive(Debug, Clone)]
pub struct ParseError {
    pub message: String,
    pub line: usize,
    pub column: usize,
}

impl std::fmt::Display for ParseError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(
            f,
            "Parse error at line {}, column {}: {}",
            self.line, self.column, self.message
        )
    }
}

impl std::error::Error for ParseError {}

/// Parse a Mermaid flowchart diagram.
pub fn parse(input: &str) -> Result<Flowchart, ParseError> {
    let lines: Vec<&str> = input.lines().collect();

    // Find the flowchart/graph declaration
    let (direction, start_idx) = find_flowchart_declaration(&lines)?;

    let mut flowchart = Flowchart::new(direction);

    // Parse remaining lines
    for (line_num, line) in lines.iter().enumerate().skip(start_idx) {
        let line = line.trim();
        if line.is_empty() || line.starts_with("%%") {
            continue;
        }

        // Handle subgraph, end, etc. (skip for now)
        if line.starts_with("subgraph") || line == "end" {
            continue;
        }

        // Parse statements (nodes and edges)
        parse_statement(line, &mut flowchart, line_num + 1)?;
    }

    Ok(flowchart)
}

fn find_flowchart_declaration(lines: &[&str]) -> Result<(Direction, usize), ParseError> {
    for (i, line) in lines.iter().enumerate() {
        let trimmed = line.trim();
        if trimmed.starts_with("flowchart") || trimmed.starts_with("graph") {
            let parts: Vec<&str> = trimmed.split_whitespace().collect();
            let direction = if parts.len() > 1 {
                parts[1].parse().unwrap_or(Direction::TB)
            } else {
                Direction::TB
            };
            return Ok((direction, i + 1));
        }
    }
    Err(ParseError {
        message: "No flowchart or graph declaration found".to_string(),
        line: 1,
        column: 1,
    })
}

fn parse_statement(
    line: &str,
    flowchart: &mut Flowchart,
    _line_num: usize,
) -> Result<(), ParseError> {
    // Remove trailing semicolon
    let line = line.trim_end_matches(';').trim();

    // Try to parse as edge (may contain node definitions like A[Start] --> B[End])
    if let Some((from_node, to_node, style, label)) = try_parse_edge(line) {
        // Extract IDs before moving nodes
        let from_id = from_node.id.clone();
        let to_id = to_node.id.clone();

        // Add nodes if they don't exist (or update if they have labels)
        add_or_update_node(flowchart, from_node);
        add_or_update_node(flowchart, to_node);

        let mut edge = Edge::new(from_id, to_id, style);
        if let Some(lbl) = label {
            edge = edge.with_label(lbl);
        }
        flowchart.add_edge(edge);
        return Ok(());
    }

    // Try to parse as node definition
    if let Some(node) = try_parse_node(line) {
        flowchart.add_node(node);
        return Ok(());
    }

    // Ignore unparseable lines (could be style, class, etc.)
    Ok(())
}

fn add_or_update_node(flowchart: &mut Flowchart, node: Node) {
    // Only add if not exists, or update if new one has a real label (not just ID)
    if let Some(existing) = flowchart.nodes.iter_mut().find(|n| n.id == node.id) {
        // Update only if the new node has a meaningful label
        if node.label != node.id {
            existing.label = node.label;
            existing.shape = node.shape;
        }
    } else {
        flowchart.nodes.push(node);
    }
}

fn try_parse_edge(line: &str) -> Option<(Node, Node, EdgeStyle, Option<String>)> {
    // Edge patterns: A --> B, A -.-> B, A ==> B, A --- B
    // With labels: A -->|label| B, A -- label --> B

    let edge_patterns = [
        ("==>", EdgeStyle::ThickArrow),
        ("-.->", EdgeStyle::DottedArrow),
        ("-->", EdgeStyle::SolidArrow),
        ("---", EdgeStyle::Open),
    ];

    for (pattern, style) in edge_patterns {
        if let Some(pos) = line.find(pattern) {
            let left = &line[..pos];
            let right = &line[pos + pattern.len()..];

            // Parse left side (node with optional label/shape)
            let from_node = parse_node_from_edge_part(left.trim())?;

            // Check for label after arrow: -->|label| or --|label|-->
            let (to_part, label) = extract_label(right.trim());

            // Parse right side
            let to_node = parse_node_from_edge_part(to_part.trim())?;

            return Some((from_node, to_node, style, label));
        }
    }

    None
}

fn parse_node_from_edge_part(s: &str) -> Option<Node> {
    let s = s.trim();
    if s.is_empty() {
        return None;
    }

    // Find where node ID ends
    let id_end = s
        .find(|c: char| !c.is_alphanumeric() && c != '_' && c != '-')
        .unwrap_or(s.len());

    if id_end == 0 {
        return None;
    }

    let id = s[..id_end].to_string();
    let rest = &s[id_end..];

    if rest.is_empty() {
        return Some(Node::new(&id, &id, NodeShape::Rectangle));
    }

    // Try to parse shape/label
    if let Some((label, shape)) = parse_node_shape(s, &id) {
        Some(Node::new(id, label, shape))
    } else {
        Some(Node::new(&id, &id, NodeShape::Rectangle))
    }
}

fn extract_label(s: &str) -> (&str, Option<String>) {
    // Pattern: |label| NodeB
    if let Some(stripped) = s.strip_prefix('|') {
        if let Some(end) = stripped.find('|') {
            let label = stripped[..end].to_string();
            let rest = stripped[end + 1..].trim();
            return (rest, Some(label));
        }
    }
    (s, None)
}

fn try_parse_node(line: &str) -> Option<Node> {
    let (id, rest) = parse_node_part(line)?;

    // If there's no shape info, it's just an ID reference
    if rest.is_empty() && !line.contains('[') && !line.contains('(') && !line.contains('{') {
        return None;
    }

    let (label, shape) = parse_node_shape(line, &id)?;
    Some(Node::new(id, label, shape))
}

fn parse_node_part(s: &str) -> Option<(String, &str)> {
    let s = s.trim();
    if s.is_empty() {
        return None;
    }

    // Find where node ID ends and shape/label begins
    let id_end = s
        .find(|c: char| !c.is_alphanumeric() && c != '_' && c != '-')
        .unwrap_or(s.len());

    if id_end == 0 {
        return None;
    }

    let id = s[..id_end].to_string();
    let rest = &s[id_end..];

    Some((id, rest))
}

fn parse_node_shape(line: &str, id: &str) -> Option<(String, NodeShape)> {
    let after_id = &line[line.find(id)? + id.len()..];
    let after_id = after_id.trim();

    if after_id.is_empty() {
        return Some((id.to_string(), NodeShape::Rectangle));
    }

    // Try each shape pattern
    // Order matters: more specific patterns first

    // Stadium: ([text])
    if after_id.starts_with("([") {
        let label = extract_between(after_id, "([", "])")?;
        return Some((label, NodeShape::Stadium));
    }

    // Circle: ((text))
    if after_id.starts_with("((") {
        let label = extract_between(after_id, "((", "))")?;
        return Some((label, NodeShape::Circle));
    }

    // Hexagon: {{text}}
    if after_id.starts_with("{{") {
        let label = extract_between(after_id, "{{", "}}")?;
        return Some((label, NodeShape::Hexagon));
    }

    // Trapezoid: [/text\]
    if after_id.starts_with("[/") && after_id.contains("\\]") {
        let label = extract_between(after_id, "[/", "\\]")?;
        return Some((label, NodeShape::Trapezoid));
    }

    // Parallelogram: [/text/]
    if after_id.starts_with("[/") {
        let label = extract_between(after_id, "[/", "/]")?;
        return Some((label, NodeShape::Parallelogram));
    }

    // Diamond: {text}
    if after_id.starts_with('{') {
        let label = extract_between(after_id, "{", "}")?;
        return Some((label, NodeShape::Diamond));
    }

    // Rounded: (text)
    if after_id.starts_with('(') {
        let label = extract_between(after_id, "(", ")")?;
        return Some((label, NodeShape::Rounded));
    }

    // Rectangle: [text]
    if after_id.starts_with('[') {
        let label = extract_between(after_id, "[", "]")?;
        return Some((label, NodeShape::Rectangle));
    }

    Some((id.to_string(), NodeShape::Rectangle))
}

fn extract_between(s: &str, start: &str, end: &str) -> Option<String> {
    let start_pos = s.find(start)? + start.len();
    let end_pos = s.rfind(end)?;
    if start_pos <= end_pos {
        Some(s[start_pos..end_pos].to_string())
    } else {
        None
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_simple() {
        let input = "graph TD\n  A --> B";
        let fc = parse(input).unwrap();
        assert_eq!(fc.direction, Direction::TB);
        assert_eq!(fc.nodes.len(), 2);
        assert_eq!(fc.edges.len(), 1);
    }

    #[test]
    fn test_parse_with_labels() {
        let input = "flowchart LR\n  A[Start] --> B[End]";
        let fc = parse(input).unwrap();
        assert_eq!(fc.direction, Direction::LR);
        assert_eq!(fc.get_node("A").unwrap().label, "Start");
        assert_eq!(fc.get_node("B").unwrap().label, "End");
    }

    #[test]
    fn test_parse_shapes() {
        let input = r#"
flowchart TD
    A[Rectangle]
    B(Rounded)
    C([Stadium])
    D{Diamond}
    E{{Hexagon}}
    F((Circle))
    G[/Parallelogram/]
    H[/Trapezoid\]
"#;
        let fc = parse(input).unwrap();
        assert_eq!(fc.get_node("A").unwrap().shape, NodeShape::Rectangle);
        assert_eq!(fc.get_node("B").unwrap().shape, NodeShape::Rounded);
        assert_eq!(fc.get_node("C").unwrap().shape, NodeShape::Stadium);
        assert_eq!(fc.get_node("D").unwrap().shape, NodeShape::Diamond);
        assert_eq!(fc.get_node("E").unwrap().shape, NodeShape::Hexagon);
        assert_eq!(fc.get_node("F").unwrap().shape, NodeShape::Circle);
        assert_eq!(fc.get_node("G").unwrap().shape, NodeShape::Parallelogram);
        assert_eq!(fc.get_node("H").unwrap().shape, NodeShape::Trapezoid);
    }

    #[test]
    fn test_parse_edge_styles() {
        let input = r#"
flowchart TD
    A --> B
    B -.-> C
    C ==> D
    D --- E
"#;
        let fc = parse(input).unwrap();
        assert_eq!(fc.edges[0].style, EdgeStyle::SolidArrow);
        assert_eq!(fc.edges[1].style, EdgeStyle::DottedArrow);
        assert_eq!(fc.edges[2].style, EdgeStyle::ThickArrow);
        assert_eq!(fc.edges[3].style, EdgeStyle::Open);
    }

    #[test]
    fn test_parse_edge_with_label() {
        let input = "flowchart TD\n  A -->|yes| B";
        let fc = parse(input).unwrap();
        assert_eq!(fc.edges[0].label.as_deref(), Some("yes"));
    }
}
