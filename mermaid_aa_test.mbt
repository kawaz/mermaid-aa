// Test module for mermaid-aa

///|
test "parse simple flowchart" {
  let input = "graph TD; A --> B"
  let flowchart = @parser.parse(input)
  inspect!(flowchart.direction(), content="TB")
  inspect!(flowchart.nodes().length(), content="2")
  inspect!(flowchart.edges().length(), content="1")
}

///|
test "parse flowchart keyword" {
  let input = "flowchart LR; A[Hello] --> B[World]"
  let flowchart = @parser.parse(input)
  inspect!(flowchart.direction(), content="LR")
  inspect!(flowchart.nodes().length(), content="2")
  inspect!(flowchart.nodes()[0].label(), content="Hello")
  inspect!(flowchart.nodes()[1].label(), content="World")
}

///|
test "parse node shapes" {
  let input = "graph TD; A[rect] --> B(round) --> C{diamond}"
  let flowchart = @parser.parse(input)
  inspect!(flowchart.nodes().length(), content="3")
  inspect!(flowchart.nodes()[0].shape(), content="Rectangle")
  inspect!(flowchart.nodes()[1].shape(), content="Rounded")
  inspect!(flowchart.nodes()[2].shape(), content="Diamond")
}

///|
test "string width - ASCII" {
  inspect!(@width.string_width("hello", 1), content="5")
  inspect!(@width.string_width("hello", 2), content="5")
}

///|
test "string width - CJK" {
  inspect!(@width.string_width("日本語", 1), content="6")
  inspect!(@width.string_width("日本語", 2), content="6")
}

///|
test "string width - ambiguous" {
  // Box drawing characters are ambiguous
  inspect!(@width.string_width("───", 1), content="3")
  inspect!(@width.string_width("───", 2), content="6")
}

///|
test "render simple flowchart" {
  let input = "graph LR; A --> B"
  let result = convert_default(input)
  // Just check it doesn't crash and produces output
  assert_true!(result.length() > 0)
}

///|
test "direction parsing" {
  inspect!(@types.Direction::from_string("TB"), content="Some(TB)")
  inspect!(@types.Direction::from_string("td"), content="Some(TD)")
  inspect!(@types.Direction::from_string("LR"), content="Some(LR)")
  inspect!(@types.Direction::from_string("RL"), content="Some(RL)")
  inspect!(@types.Direction::from_string("BT"), content="Some(BT)")
  inspect!(@types.Direction::from_string("invalid"), content="None")
}

///|
test "charset parsing" {
  inspect!(@types.Charset::from_string("ascii"), content="Some(Ascii)")
  inspect!(@types.Charset::from_string("unicode"), content="Some(Unicode)")
  inspect!(
    @types.Charset::from_string("unicode-round"),
    content="Some(UnicodeRound)",
  )
  inspect!(
    @types.Charset::from_string("unicode-bold"),
    content="Some(UnicodeBold)",
  )
  inspect!(
    @types.Charset::from_string("unicode-double"),
    content="Some(UnicodeDouble)",
  )
  inspect!(@types.Charset::from_string("invalid"), content="None")
}
