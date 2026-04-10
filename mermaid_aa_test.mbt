// Test module for mermaid-aa

///|
test "parse simple flowchart" {
  let input = "graph TD; A --> B"
  let flowchart = @parser.parse(input)
  inspect(flowchart.direction(), content="TB")
  inspect(flowchart.nodes().length(), content="2")
  inspect(flowchart.edges().length(), content="1")
}

///|
test "parse flowchart keyword" {
  let input = "flowchart LR; A[Hello] --> B[World]"
  let flowchart = @parser.parse(input)
  inspect(flowchart.direction(), content="LR")
  inspect(flowchart.nodes().length(), content="2")
  inspect(flowchart.nodes()[0].label(), content="Hello")
  inspect(flowchart.nodes()[1].label(), content="World")
}

///|
test "parse node shapes" {
  let input = "graph TD; A[rect] --> B(round) --> C{diamond}"
  let flowchart = @parser.parse(input)
  inspect(flowchart.nodes().length(), content="3")
  inspect(flowchart.nodes()[0].shape(), content="Rectangle")
  inspect(flowchart.nodes()[1].shape(), content="Rounded")
  inspect(flowchart.nodes()[2].shape(), content="Diamond")
}

///|
test "string width - ASCII" {
  inspect(@width.string_width("hello", 1), content="5")
  inspect(@width.string_width("hello", 2), content="5")
}

///|
test "string width - CJK" {
  inspect(@width.string_width("日本語", 1), content="6")
  inspect(@width.string_width("日本語", 2), content="6")
}

///|
test "string width - ambiguous (legacy API)" {
  // Box drawing characters are ambiguous
  inspect(@width.string_width("───", 1), content="3")
  inspect(@width.string_width("───", 2), content="6")
}

///|
test "WidthMode parsing" {
  // half / 1
  inspect(@types.WidthMode::from_string("1"), content="Some(Half)")
  inspect(@types.WidthMode::from_string("half"), content="Some(Half)")
  inspect(@types.WidthMode::from_string("Half"), content="Some(Half)")
  // full / 2
  inspect(@types.WidthMode::from_string("2"), content="Some(Full)")
  inspect(@types.WidthMode::from_string("full"), content="Some(Full)")
  inspect(@types.WidthMode::from_string("FULL"), content="Some(Full)")
  // console
  inspect(@types.WidthMode::from_string("console"), content="Some(Console)")
  // legacy
  inspect(@types.WidthMode::from_string("legacy"), content="Some(Legacy)")
  // invalid
  inspect(@types.WidthMode::from_string("invalid"), content="None")
}

///|
test "WidthMode - ambiguous_width" {
  inspect(@types.Half.ambiguous_width(), content="1")
  inspect(@types.Full.ambiguous_width(), content="2")
  inspect(@types.Console.ambiguous_width(), content="1")
  inspect(@types.Legacy.ambiguous_width(), content="2")
}

///|
test "WidthMode - box_drawing_width" {
  inspect(@types.Half.box_drawing_width(), content="1")
  inspect(@types.Full.box_drawing_width(), content="1")
  inspect(@types.Console.box_drawing_width(), content="1")
  inspect(@types.Legacy.box_drawing_width(), content="2")
}

///|
test "string_width_mode - box drawing" {
  // Box drawing: U+2500 (─)
  // Half mode: box drawing = 1
  inspect(@width.string_width_mode("───", @types.Half), content="3")
  // Full mode: box drawing = 1 (not 2!)
  inspect(@width.string_width_mode("───", @types.Full), content="3")
  // Console mode: box drawing = 1
  inspect(@width.string_width_mode("───", @types.Console), content="3")
  // Legacy mode: box drawing = 2
  inspect(@width.string_width_mode("───", @types.Legacy), content="6")
}

///|
test "string_width_mode - other ambiguous" {
  // Greek letter alpha (U+03B1): ambiguous, not box drawing
  // Half mode: ambiguous = 1
  inspect(@width.string_width_mode("αβγ", @types.Half), content="3")
  // Full mode: ambiguous = 2
  inspect(@width.string_width_mode("αβγ", @types.Full), content="6")
  // Console mode: ambiguous = 1
  inspect(@width.string_width_mode("αβγ", @types.Console), content="3")
  // Legacy mode: ambiguous = 2
  inspect(@width.string_width_mode("αβγ", @types.Legacy), content="6")
}

///|
test "render simple flowchart" {
  let input = "graph LR; A --> B"
  let result = convert_default(input)
  // Just check it doesn't crash and produces output
  assert_true(result.length() > 0)
}

///|
test "direction parsing" {
  inspect(@types.Direction::from_string("TB"), content="Some(TB)")
  inspect(@types.Direction::from_string("td"), content="Some(TD)")
  inspect(@types.Direction::from_string("LR"), content="Some(LR)")
  inspect(@types.Direction::from_string("RL"), content="Some(RL)")
  inspect(@types.Direction::from_string("BT"), content="Some(BT)")
  inspect(@types.Direction::from_string("invalid"), content="None")
}

///|
test "charset parsing" {
  inspect(@types.Charset::from_string("ascii"), content="Some(Ascii)")
  inspect(@types.Charset::from_string("unicode"), content="Some(Unicode)")
  inspect(
    @types.Charset::from_string("unicode-round"),
    content="Some(UnicodeRound)",
  )
  inspect(
    @types.Charset::from_string("unicode-bold"),
    content="Some(UnicodeBold)",
  )
  inspect(
    @types.Charset::from_string("unicode-double"),
    content="Some(UnicodeDouble)",
  )
  inspect(@types.Charset::from_string("invalid"), content="None")
}

///|
test "error format - missing keyword" {
  try {
    let _ = @parser.parse("invalid input")
    assert_true(false, msg="should have raised error")
  } catch {
    @parser.ParseError(msg) => {
      assert_true(msg.contains("flowchart") || msg.contains("graph"))
    }
    @parser.ParseErrorDetailed(info) => {
      let formatted = info.format()
      assert_true(formatted.contains("Error:"))
      assert_true(formatted.contains("-->"))
    }
    _ => assert_true(true) // Other errors are acceptable
  }
}

///|
test "error format - missing node after arrow" {
  try {
    let _ = @parser.parse_with_filename(
      "graph TD\n    A -->",
      "test.mmd",
    )
    assert_true(false, msg="should have raised error")
  } catch {
    @parser.ParseError(msg) => {
      assert_true(msg.contains("Missing") || msg.contains("error"))
    }
    @parser.ParseErrorDetailed(info) => {
      let formatted = info.format()
      assert_true(formatted.contains("Error: Missing node identifier after arrow"))
      assert_true(formatted.contains("--> test.mmd:2:"))
      assert_true(formatted.contains("= help:"))
    }
    _ => assert_true(true) // Other errors are acceptable
  }
}

///|
test "ParseErrorInfo format" {
  let err_info = @parser.ParseErrorInfo::new(
    "Missing node identifier after arrow",
    3,
    7,
    9,
    "    A -->",
    "input.mmd",
    Some("add a node identifier like 'B' after the arrow"),
  )
  let formatted = err_info.format()
  // Verify the format structure
  assert_true(formatted.contains("Error: Missing node identifier after arrow"))
  assert_true(formatted.contains("--> input.mmd:3:7"))
  assert_true(formatted.contains("3 |     A -->"))
  assert_true(formatted.contains("^^^"))
  assert_true(formatted.contains("= help: add a node identifier"))
}
