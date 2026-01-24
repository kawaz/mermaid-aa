// Parser module - Mermaid flowchart syntax parser

/// Helper to get char from string at index
fn char_at(s : String, i : Int) -> Char {
  s[i].to_int().unsafe_to_char()
}

/// Parser state
priv struct Parser {
  input : String
  mut pos : Int
  mut line : Int
  mut column : Int
}

/// Create a new parser
fn Parser::new(input : String) -> Parser {
  { input, pos: 0, line: 1, column: 1 }
}

/// Check if at end of input
fn Parser::is_eof(self : Parser) -> Bool {
  self.pos >= self.input.length()
}

/// Get current character
fn Parser::current(self : Parser) -> Char? {
  if self.is_eof() {
    None
  } else {
    Some(char_at(self.input, self.pos))
  }
}

/// Advance one character
fn Parser::advance(self : Parser) -> Unit {
  if not(self.is_eof()) {
    if char_at(self.input, self.pos) == '\n' {
      self.line = self.line + 1
      self.column = 1
    } else {
      self.column = self.column + 1
    }
    self.pos = self.pos + 1
  }
}

/// Skip whitespace (but not newlines for some contexts)
fn Parser::skip_whitespace(self : Parser) -> Unit {
  while not(self.is_eof()) {
    match self.current() {
      Some(' ') | Some('\t') | Some('\r') | Some('\n') => self.advance()
      _ => break
    }
  }
}

/// Skip whitespace except newlines
fn Parser::skip_spaces(self : Parser) -> Unit {
  while not(self.is_eof()) {
    match self.current() {
      Some(' ') | Some('\t') => self.advance()
      _ => break
    }
  }
}

/// Check if current position starts with given string
fn Parser::starts_with(self : Parser, s : String) -> Bool {
  if self.pos + s.length() > self.input.length() {
    return false
  }
  for i = 0; i < s.length(); i = i + 1 {
    if self.input[self.pos + i] != s[i] {
      return false
    }
  }
  true
}

/// Consume expected string
fn Parser::consume(self : Parser, s : String) -> Bool {
  if self.starts_with(s) {
    for i = 0; i < s.length(); i = i + 1 {
      self.advance()
    }
    true
  } else {
    false
  }
}

/// Check if character is valid for identifier start
fn is_id_start(c : Char) -> Bool {
  (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

/// Check if character is valid for identifier continuation
fn is_id_char(c : Char) -> Bool {
  is_id_start(c) || (c >= '0' && c <= '9')
}

/// Get substring by building manually
fn get_substring(s : String, start : Int, end : Int) -> String {
  let buf = StringBuilder::new()
  for i = start; i < end; i = i + 1 {
    buf.write_char(char_at(s, i))
  }
  buf.to_string()
}

/// Parse an identifier
fn Parser::parse_identifier(self : Parser) -> String? {
  self.skip_spaces()
  let start = self.pos
  // First character must be letter or underscore or digit
  match self.current() {
    Some(c) => {
      if not(is_id_start(c)) && not(c >= '0' && c <= '9') {
        return None
      }
    }
    None => return None
  }
  // Parse remaining characters
  while not(self.is_eof()) {
    match self.current() {
      Some(c) => {
        if is_id_char(c) || (c >= '0' && c <= '9') {
          self.advance()
        } else {
          break
        }
      }
      None => break
    }
  }
  if self.pos > start {
    Some(get_substring(self.input, start, self.pos))
  } else {
    None
  }
}

/// Parse text until closing delimiter
fn Parser::parse_until(self : Parser, close : String) -> String? {
  let start = self.pos
  while not(self.is_eof()) && not(self.starts_with(close)) {
    self.advance()
  }
  if self.starts_with(close) {
    Some(get_substring(self.input, start, self.pos))
  } else {
    None
  }
}

/// Parse node shape and label: [text], (text), ((text)), {text}, etc.
fn Parser::parse_node_shape(self : Parser) -> (@types.NodeShape, String)? {
  self.skip_spaces()
  // Circle: ((text))
  if self.consume("((") {
    match self.parse_until("))") {
      Some(label) => {
        let _ = self.consume("))")
        return Some((@types.Circle, label))
      }
      None => return None
    }
  }
  // Stadium: ([text])
  if self.consume("([") {
    match self.parse_until("])") {
      Some(label) => {
        let _ = self.consume("])")
        return Some((@types.Stadium, label))
      }
      None => return None
    }
  }
  // Hexagon: {{text}}
  if self.consume("{{") {
    match self.parse_until("}}") {
      Some(label) => {
        let _ = self.consume("}}")
        return Some((@types.Hexagon, label))
      }
      None => return None
    }
  }
  // Parallelogram: [/text/]
  if self.consume("[/") {
    // Check for trapezoid [/text\]
    match self.parse_until_any(["/]", "\\]"]) {
      Some((label, ending)) => {
        if ending == "/]" {
          let _ = self.consume("/]")
          return Some((@types.Parallelogram, label))
        } else {
          let _ = self.consume("\\]")
          return Some((@types.Trapezoid, label))
        }
      }
      None => return None
    }
  }
  // Diamond: {text}
  if self.consume("{") {
    match self.parse_until("}") {
      Some(label) => {
        let _ = self.consume("}")
        return Some((@types.Diamond, label))
      }
      None => return None
    }
  }
  // Rounded: (text)
  if self.consume("(") {
    match self.parse_until(")") {
      Some(label) => {
        let _ = self.consume(")")
        return Some((@types.Rounded, label))
      }
      None => return None
    }
  }
  // Rectangle: [text]
  if self.consume("[") {
    match self.parse_until("]") {
      Some(label) => {
        let _ = self.consume("]")
        return Some((@types.Rectangle, label))
      }
      None => return None
    }
  }
  None
}

/// Parse until any of the given endings
fn Parser::parse_until_any(
  self : Parser,
  endings : Array[String]
) -> (String, String)? {
  let start = self.pos
  while not(self.is_eof()) {
    for ending in endings {
      if self.starts_with(ending) {
        return Some((get_substring(self.input, start, self.pos), ending))
      }
    }
    self.advance()
  }
  None
}

/// Parse edge style: -->, -.->, ==>, ---
fn Parser::parse_edge(self : Parser) -> (@types.EdgeStyle, String?)? {
  self.skip_spaces()
  // With label: --|text|-->
  if self.consume("--") {
    // Check for label
    if self.consume("|") {
      match self.parse_until("|") {
        Some(label) => {
          let _ = self.consume("|")
          // Now expect arrow ending
          if self.consume("-->") {
            return Some((@types.Solid, Some(label)))
          } else if self.consume("->") {
            return Some((@types.Solid, Some(label)))
          } else if self.consume("-") {
            return Some((@types.Open, Some(label)))
          }
          return None
        }
        None => return None
      }
    }
    // Solid arrow: -->
    if self.consume(">") {
      return Some((@types.Solid, None))
    }
    // Open: ---
    if self.consume("-") {
      return Some((@types.Open, None))
    }
    return None
  }
  // Dotted: -.->
  if self.consume("-.") {
    if self.consume("->") {
      return Some((@types.Dotted, None))
    }
    // Dotted with label: -.-|text|-.-
    if self.consume("-") {
      return Some((@types.Dotted, None))
    }
    return None
  }
  // Thick: ==>
  if self.consume("==>") {
    return Some((@types.Thick, None))
  }
  None
}

/// Parse a single node with optional shape
fn Parser::parse_node(self : Parser) -> @types.Node? {
  self.skip_spaces()
  match self.parse_identifier() {
    Some(id) => {
      // Check for shape/label
      self.skip_spaces()
      match self.parse_node_shape() {
        Some((shape, label)) => Some(@types.Node::new(id, label, shape))
        None => Some(@types.Node::new(id, id, @types.Rectangle))
      }
    }
    None => None
  }
}

/// Skip comment (starts with %%)
fn Parser::skip_comment(self : Parser) -> Unit {
  if self.starts_with("%%") {
    while not(self.is_eof()) {
      match self.current() {
        Some('\n') => {
          self.advance()
          break
        }
        _ => self.advance()
      }
    }
  }
}

/// Parse a line (node definition or edge)
fn Parser::parse_line(self : Parser, flowchart : @types.Flowchart) -> Bool {
  self.skip_whitespace()
  self.skip_comment()
  if self.is_eof() {
    return true
  }
  // Try to parse a node
  match self.parse_node() {
    Some(node) => {
      flowchart.add_node(node)
      // Check for edges in a while loop
      while true {
        self.skip_spaces()
        match self.parse_edge() {
          Some((style, label)) => {
            // Parse target node
            match self.parse_node() {
              Some(target) => {
                flowchart.add_node(target)
                flowchart.add_edge(
                  @types.Edge::new(node.id(), target.id(), style, label),
                )
                // Allow chaining: A --> B --> C
              }
              None => return false
            }
          }
          None => break
        }
      }
      // Skip to next statement (semicolon or newline)
      self.skip_spaces()
      let _ = self.consume(";")
      true
    }
    None => {
      // Skip semicolons
      if self.consume(";") {
        return true
      }
      // Skip newlines
      if self.consume("\n") {
        return true
      }
      // Unknown content - skip to next line
      while not(self.is_eof()) {
        match self.current() {
          Some('\n') | Some(';') => break
          _ => self.advance()
        }
      }
      true
    }
  }
}

/// Parse error type
pub suberror ParseError {
  ParseError(String)
}

/// Main parse function
pub fn parse(input : String) -> @types.Flowchart raise ParseError {
  let parser = Parser::new(input)
  parser.skip_whitespace()
  // Parse "graph" or "flowchart" keyword
  let is_flowchart = parser.consume("flowchart")
  if not(is_flowchart) {
    if not(parser.consume("graph")) {
      raise ParseError("Expected 'flowchart' or 'graph' keyword")
    }
  }
  parser.skip_spaces()
  // Parse direction
  let direction : @types.Direction = if parser.consume("TD") ||
    parser.consume("TB") {
    @types.TB
  } else if parser.consume("BT") {
    @types.BT
  } else if parser.consume("LR") {
    @types.LR
  } else if parser.consume("RL") {
    @types.RL
  } else {
    @types.TB // Default
  }
  let flowchart = @types.Flowchart::new(direction)
  // Skip semicolon or newline after direction
  parser.skip_spaces()
  let _ = parser.consume(";")
  // Parse nodes and edges
  while not(parser.is_eof()) {
    if not(parser.parse_line(flowchart)) {
      raise ParseError(
        "Parse error at line \{parser.line}, column \{parser.column}",
      )
    }
  }
  flowchart
}
