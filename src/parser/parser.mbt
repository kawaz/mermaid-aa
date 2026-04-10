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
  filename : String
}

/// Create a new parser with filename
fn Parser::new_with_filename(input : String, filename : String) -> Parser {
  { input, pos: 0, line: 1, column: 1, filename }
}

/// Get the source line at the given line number (1-indexed)
fn Parser::get_source_line(self : Parser, line_num : Int) -> String {
  let lines = split_lines(self.input)
  if line_num >= 1 && line_num <= lines.length() {
    lines[line_num - 1]
  } else {
    ""
  }
}

/// Split input into lines
fn split_lines(input : String) -> Array[String] {
  let lines : Array[String] = []
  let buf = StringBuilder::new()
  for i = 0; i < input.length(); i = i + 1 {
    let c = char_at(input, i)
    if c == '\n' {
      lines.push(buf.to_string())
      buf.reset()
    } else if c != '\r' {
      buf.write_char(c)
    }
  }
  // Don't forget the last line
  if buf.to_string().length() > 0 {
    lines.push(buf.to_string())
  }
  lines
}

/// Create a parse error with detailed information
fn Parser::make_error(
  self : Parser,
  message : String,
  help : String?
) -> ParseErrorInfo {
  ParseErrorInfo::new(
    message,
    self.line,
    self.column,
    self.column,
    self.get_source_line(self.line),
    self.filename,
    help,
  )
}

/// Create a parse error with span (multi-character highlight)
fn Parser::make_error_span(
  self : Parser,
  message : String,
  start_column : Int,
  end_column : Int,
  help : String?
) -> ParseErrorInfo {
  ParseErrorInfo::new(
    message,
    self.line,
    start_column,
    end_column,
    self.get_source_line(self.line),
    self.filename,
    help,
  )
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
  if !(self.is_eof()) {
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
  while !(self.is_eof()) {
    match self.current() {
      Some(' ') | Some('\t') | Some('\r') | Some('\n') => self.advance()
      _ => break
    }
  }
}

/// Skip whitespace except newlines
fn Parser::skip_spaces(self : Parser) -> Unit {
  while !(self.is_eof()) {
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
      if !(is_id_start(c)) && !(c >= '0' && c <= '9') {
        return None
      }
    }
    None => return None
  }
  // Parse remaining characters
  while !(self.is_eof()) {
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
  while !(self.is_eof()) && !(self.starts_with(close)) {
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
  while !(self.is_eof()) {
    for ending in endings {
      if self.starts_with(ending) {
        return Some((get_substring(self.input, start, self.pos), ending))
      }
    }
    self.advance()
  }
  None
}

/// Parse edge label: |text|
fn Parser::parse_edge_label(self : Parser) -> String? {
  if self.consume("|") {
    match self.parse_until("|") {
      Some(label) => {
        let _ = self.consume("|")
        Some(label)
      }
      None => None
    }
  } else {
    None
  }
}

/// Parse edge style: -->, -.->, ==>, ---
/// Supports both pre-label (--|text|-->) and post-label (-->|text|) syntax
fn Parser::parse_edge(self : Parser) -> (@types.EdgeStyle, String?)? {
  self.skip_spaces()
  // With label: --|text|-->
  if self.consume("--") {
    // Check for pre-label: --|text|-->
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
      // Check for post-label: -->|text|
      let label = self.parse_edge_label()
      return Some((@types.Solid, label))
    }
    // Open: ---
    if self.consume("-") {
      // Check for post-label: ---|text|
      let label = self.parse_edge_label()
      return Some((@types.Open, label))
    }
    return None
  }
  // Dotted: -.->
  if self.consume("-.") {
    if self.consume("->") {
      // Check for post-label: -.->|text|
      let label = self.parse_edge_label()
      return Some((@types.Dotted, label))
    }
    // Dotted open: -.-
    if self.consume("-") {
      // Check for post-label: -.-|text|
      let label = self.parse_edge_label()
      return Some((@types.Dotted, label))
    }
    return None
  }
  // Thick: ==>
  if self.consume("==>") {
    // Check for post-label: ==>|text|
    let label = self.parse_edge_label()
    return Some((@types.Thick, label))
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
    while !(self.is_eof()) {
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

/// Result type for parse_line with detailed error
priv enum ParseLineResult {
  Ok
  Error(ParseErrorInfo)
}

/// Parse a line (node definition or edge)
fn Parser::parse_line(self : Parser, flowchart : @types.Flowchart) -> ParseLineResult {
  self.skip_whitespace()
  self.skip_comment()
  if self.is_eof() {
    return Ok
  }
  // Try to parse a node
  match self.parse_node() {
    Some(node) => {
      flowchart.add_node(node)
      // Track the last source node for chained edges
      let mut last_node = node
      // Check for edges in a while loop
      while true {
        self.skip_spaces()
        let edge_start_col = self.column
        match self.parse_edge() {
          Some((style, label)) => {
            // Record position after edge for error reporting
            let after_edge_col = self.column
            self.skip_spaces()
            // Parse target node
            match self.parse_node() {
              Some(target) => {
                flowchart.add_node(target)
                flowchart.add_edge(
                  @types.Edge::new(last_node.id(), target.id(), style, label),
                )
                // Allow chaining: A --> B --> C
                last_node = target
              }
              None => {
                // Missing target node after arrow
                return Error(
                  self.make_error_span(
                    "Missing node identifier after arrow",
                    edge_start_col,
                    after_edge_col - 1,
                    Some("add a node identifier like 'B' or 'B[Label]' after the arrow"),
                  ),
                )
              }
            }
          }
          None => break
        }
      }
      // Skip to next statement (semicolon or newline)
      self.skip_spaces()
      let _ = self.consume(";")
      Ok
    }
    None => {
      // Skip semicolons
      if self.consume(";") {
        return Ok
      }
      // Skip newlines
      if self.consume("\n") {
        return Ok
      }
      // Unknown content - report error
      let char_str = match self.current() {
        Some(c) => c.to_string()
        None => "EOF"
      }
      return Error(
        self.make_error(
          "Unexpected character '\{char_str}'",
          Some("expected node identifier (e.g., 'A', 'Start') or comment (%%)"),
        ),
      )
    }
  }
}

/// Parse error with detailed location information
pub(all) struct ParseErrorInfo {
  message : String
  line : Int
  column : Int
  end_column : Int  // For multi-character highlights
  source_line : String
  filename : String
  help : String?
}

/// Create a new ParseErrorInfo
pub fn ParseErrorInfo::new(
  message : String,
  line : Int,
  column : Int,
  end_column : Int,
  source_line : String,
  filename : String,
  help : String?
) -> ParseErrorInfo {
  { message, line, column, end_column, source_line, filename, help }
}

/// Write n spaces to buffer
fn write_spaces(buf : StringBuilder, n : Int) -> Unit {
  for i = 0; i < n; i = i + 1 {
    buf.write_char(' ')
  }
}

/// Write n copies of a character to buffer
fn write_chars(buf : StringBuilder, c : Char, n : Int) -> Unit {
  for i = 0; i < n; i = i + 1 {
    buf.write_char(c)
  }
}

/// Format error for display (visual format)
pub fn ParseErrorInfo::format(self : ParseErrorInfo) -> String {
  let buf = StringBuilder::new()
  // Error header: "Error: message"
  buf.write_string("Error: ")
  buf.write_string(self.message)
  buf.write_char('\n')
  // Location: " --> filename:line:column"
  buf.write_string(" --> ")
  buf.write_string(self.filename)
  buf.write_char(':')
  buf.write_string(self.line.to_string())
  buf.write_char(':')
  buf.write_string(self.column.to_string())
  buf.write_char('\n')
  // Separator line
  let line_num_str = self.line.to_string()
  let padding = line_num_str.length()
  write_spaces(buf, padding + 1)
  buf.write_string("|\n")
  // Source line with line number
  buf.write_string(line_num_str)
  buf.write_string(" | ")
  buf.write_string(self.source_line)
  buf.write_char('\n')
  // Indicator line with ^
  write_spaces(buf, padding + 1)
  buf.write_string("| ")
  // Spaces before the indicator
  write_spaces(buf, self.column - 1)
  // The ^ indicators
  let indicator_len = if self.end_column > self.column {
    self.end_column - self.column + 1
  } else {
    1
  }
  write_chars(buf, '^', indicator_len)
  buf.write_char('\n')
  // Help message if present
  match self.help {
    Some(help_msg) => {
      write_spaces(buf, padding + 1)
      buf.write_string("= help: ")
      buf.write_string(help_msg)
      buf.write_char('\n')
    }
    None => ()
  }
  buf.to_string()
}

/// Parse error type
pub suberror ParseError {
  ParseError(String)
  ParseErrorDetailed(ParseErrorInfo)
}

/// Main parse function
pub fn parse(input : String) -> @types.Flowchart raise ParseError {
  parse_with_filename(input, "input.mmd")
}

/// Parse with custom filename for error messages
pub fn parse_with_filename(
  input : String,
  filename : String
) -> @types.Flowchart raise ParseError {
  let parser = Parser::new_with_filename(input, filename)
  parser.skip_whitespace()
  // Parse "graph" or "flowchart" keyword
  let keyword_col = parser.column
  let is_flowchart = parser.consume("flowchart")
  if !(is_flowchart) {
    if !(parser.consume("graph")) {
      raise ParseErrorDetailed(
        parser.make_error_span(
          "Expected 'flowchart' or 'graph' keyword",
          keyword_col,
          keyword_col + 4,
          Some("Mermaid diagrams must start with 'flowchart' or 'graph'"),
        ),
      )
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
  while !(parser.is_eof()) {
    match parser.parse_line(flowchart) {
      Ok => ()
      Error(err) => raise ParseErrorDetailed(err)
    }
  }
  flowchart
}

/// Format a ParseError for display
pub fn format_error(err : ParseError) -> String {
  match err {
    ParseError(msg) => "Error: \{msg}\n"
    ParseErrorDetailed(info) => info.format()
  }
}
