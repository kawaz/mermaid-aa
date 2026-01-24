///|
fn main {
  // Example: convert a simple flowchart
  let input = "graph TD; A[Start] --> B[End]"
  // Note: main cannot handle errors, so we use try-catch pattern
  try {
    let result = @lib.convert_default(input)
    println(result)
  } catch {
    @parser.ParseError(e) => println("Error: \{e}")
  }
}
