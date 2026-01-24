//! Common types for mermaid-aa.

/// Diagram direction.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Default)]
pub enum Direction {
    #[default]
    TB, // Top to Bottom
    BT, // Bottom to Top
    LR, // Left to Right
    RL, // Right to Left
}

impl std::str::FromStr for Direction {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s.to_uppercase().as_str() {
            "TB" | "TD" => Ok(Direction::TB),
            "BT" => Ok(Direction::BT),
            "LR" => Ok(Direction::LR),
            "RL" => Ok(Direction::RL),
            _ => Err(format!("Unknown direction: {}", s)),
        }
    }
}

/// Node shape types.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Default)]
pub enum NodeShape {
    #[default]
    Rectangle, // A[text]
    Rounded,       // A(text)
    Stadium,       // A([text])
    Diamond,       // A{text}
    Hexagon,       // A{{text}}
    Circle,        // A((text))
    Parallelogram, // A[/text/]
    Trapezoid,     // A[/text\]
}

/// Edge style types.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Default)]
pub enum EdgeStyle {
    #[default]
    SolidArrow, // -->
    DottedArrow, // -.->
    ThickArrow,  // ==>
    Open,        // ---
}

/// A node in the flowchart.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Node {
    pub id: String,
    pub label: String,
    pub shape: NodeShape,
}

impl Node {
    pub fn new(id: impl Into<String>, label: impl Into<String>, shape: NodeShape) -> Self {
        Self {
            id: id.into(),
            label: label.into(),
            shape,
        }
    }
}

/// An edge between nodes.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Edge {
    pub from: String,
    pub to: String,
    pub style: EdgeStyle,
    pub label: Option<String>,
}

impl Edge {
    pub fn new(from: impl Into<String>, to: impl Into<String>, style: EdgeStyle) -> Self {
        Self {
            from: from.into(),
            to: to.into(),
            style,
            label: None,
        }
    }

    pub fn with_label(mut self, label: impl Into<String>) -> Self {
        self.label = Some(label.into());
        self
    }
}

/// A parsed flowchart diagram.
#[derive(Debug, Clone, Default)]
pub struct Flowchart {
    pub direction: Direction,
    pub nodes: Vec<Node>,
    pub edges: Vec<Edge>,
}

impl Flowchart {
    pub fn new(direction: Direction) -> Self {
        Self {
            direction,
            nodes: Vec::new(),
            edges: Vec::new(),
        }
    }

    pub fn add_node(&mut self, node: Node) {
        if !self.nodes.iter().any(|n| n.id == node.id) {
            self.nodes.push(node);
        }
    }

    pub fn add_edge(&mut self, edge: Edge) {
        self.edges.push(edge);
    }

    #[cfg(test)]
    pub fn get_node(&self, id: &str) -> Option<&Node> {
        self.nodes.iter().find(|n| n.id == id)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_direction_parse() {
        assert_eq!("TB".parse::<Direction>().unwrap(), Direction::TB);
        assert_eq!("TD".parse::<Direction>().unwrap(), Direction::TB);
        assert_eq!("LR".parse::<Direction>().unwrap(), Direction::LR);
        assert_eq!("lr".parse::<Direction>().unwrap(), Direction::LR);
    }

    #[test]
    fn test_flowchart() {
        let mut fc = Flowchart::new(Direction::LR);
        fc.add_node(Node::new("A", "Start", NodeShape::Rounded));
        fc.add_node(Node::new("B", "End", NodeShape::Rounded));
        fc.add_edge(Edge::new("A", "B", EdgeStyle::SolidArrow));

        assert_eq!(fc.nodes.len(), 2);
        assert_eq!(fc.edges.len(), 1);
        assert!(fc.get_node("A").is_some());
    }
}
