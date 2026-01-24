// mermaid-aa: Convert Mermaid diagrams to ASCII art
// Main library module

/// Parse and render Mermaid diagram to ASCII art (legacy API)
pub fn convert(
  input : String,
  charset : @types.Charset,
  ambiguous_width : Int
) -> String raise @parser.ParseError {
  let flowchart = @parser.parse(input)
  let options = @types.RenderOptions::new(charset, ambiguous_width)
  @renderer.render(flowchart, options)
}

/// Parse and render Mermaid diagram to ASCII art with WidthMode
pub fn convert_with_mode(
  input : String,
  charset : @types.Charset,
  mode : @types.WidthMode
) -> String raise @parser.ParseError {
  let flowchart = @parser.parse(input)
  let options = @types.RenderOptions::with_mode(charset, mode)
  @renderer.render(flowchart, options)
}

/// Parse and render with default options
pub fn convert_default(input : String) -> String raise @parser.ParseError {
  convert_with_mode(input, @types.Unicode, @types.Half)
}

// Re-export parse function
pub fn parse(input : String) -> @types.Flowchart raise @parser.ParseError {
  @parser.parse(input)
}

// Parse with custom filename for error messages
pub fn parse_with_filename(
  input : String,
  filename : String
) -> @types.Flowchart raise @parser.ParseError {
  @parser.parse_with_filename(input, filename)
}

// Format a ParseError for display
pub fn format_error(err : @parser.ParseError) -> String {
  @parser.format_error(err)
}

pub fn render(
  flowchart : @types.Flowchart,
  options : @types.RenderOptions
) -> String {
  @renderer.render(flowchart, options)
}

// Width utilities
pub fn string_width(s : String, ambiguous_width : Int) -> Int {
  @width.string_width(s, ambiguous_width)
}

pub fn string_width_mode(s : String, mode : @types.WidthMode) -> Int {
  @width.string_width_mode(s, mode)
}

pub fn parse_width_mode(s : String) -> @types.WidthMode? {
  @types.WidthMode::from_string(s)
}
