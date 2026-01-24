// mermaid-aa: Convert Mermaid diagrams to ASCII art
// Main library module

/// Parse and render Mermaid diagram to ASCII art
pub fn convert(
  input : String,
  charset : @types.Charset,
  ambiguous_width : Int
) -> String raise @parser.ParseError {
  let flowchart = @parser.parse(input)
  let options = @types.RenderOptions::new(charset, ambiguous_width)
  @renderer.render(flowchart, options)
}

/// Parse and render with default options
pub fn convert_default(input : String) -> String raise @parser.ParseError {
  convert(input, @types.Unicode, 1)
}

// Re-export parse function
pub fn parse(input : String) -> @types.Flowchart raise @parser.ParseError {
  @parser.parse(input)
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
