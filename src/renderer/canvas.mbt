// Canvas module - Drawing canvas for ASCII art

/// Helper to get char from string at index
fn char_at(s : String, i : Int) -> Char {
  s[i].to_int().unsafe_to_char()
}

/// 2D canvas for ASCII art rendering
pub struct Canvas {
  width : Int
  height : Int
  cells : Array[Array[Char]]
  ambiguous_width : Int
}

/// Create a new canvas
pub fn Canvas::new(width : Int, height : Int, ambiguous_width : Int) -> Canvas {
  let cells : Array[Array[Char]] = []
  for y = 0; y < height; y = y + 1 {
    let row : Array[Char] = []
    for x = 0; x < width; x = x + 1 {
      row.push(' ')
    }
    cells.push(row)
  }
  { width, height, cells, ambiguous_width }
}

/// Get character at position
pub fn Canvas::get(self : Canvas, x : Int, y : Int) -> Char {
  if x >= 0 && x < self.width && y >= 0 && y < self.height {
    self.cells[y][x]
  } else {
    ' '
  }
}

/// Set character at position
pub fn Canvas::set(self : Canvas, x : Int, y : Int, ch : Char) -> Unit {
  if x >= 0 && x < self.width && y >= 0 && y < self.height {
    self.cells[y][x] = ch
  }
}

/// Draw horizontal line
pub fn Canvas::draw_h_line(
  self : Canvas,
  x : Int,
  y : Int,
  length : Int,
  ch : Char
) -> Unit {
  for i = 0; i < length; i = i + 1 {
    self.set(x + i, y, ch)
  }
}

/// Draw vertical line
pub fn Canvas::draw_v_line(
  self : Canvas,
  x : Int,
  y : Int,
  length : Int,
  ch : Char
) -> Unit {
  for i = 0; i < length; i = i + 1 {
    self.set(x, y + i, ch)
  }
}

/// Draw text at position
pub fn Canvas::draw_text(self : Canvas, x : Int, y : Int, text : String) -> Unit {
  let mut col = x
  for i = 0; i < text.length(); i = i + 1 {
    let ch = char_at(text, i)
    self.set(col, y, ch)
    col = col + @width.char_width(ch, self.ambiguous_width)
  }
}

/// Draw rectangle with corners and edges
pub fn Canvas::draw_rect(
  self : Canvas,
  x : Int,
  y : Int,
  w : Int,
  h : Int,
  style : CharStyle
) -> Unit {
  // Top edge
  self.set(x, y, style.corner_tl)
  self.draw_h_line(x + 1, y, w - 2, style.h_line)
  self.set(x + w - 1, y, style.corner_tr)
  // Bottom edge
  self.set(x, y + h - 1, style.corner_bl)
  self.draw_h_line(x + 1, y + h - 1, w - 2, style.h_line)
  self.set(x + w - 1, y + h - 1, style.corner_br)
  // Side edges
  self.draw_v_line(x, y + 1, h - 2, style.v_line)
  self.draw_v_line(x + w - 1, y + 1, h - 2, style.v_line)
}

/// Draw rounded rectangle
pub fn Canvas::draw_rounded_rect(
  self : Canvas,
  x : Int,
  y : Int,
  w : Int,
  h : Int,
  style : CharStyle
) -> Unit {
  // Top edge
  self.set(x, y, style.round_tl)
  self.draw_h_line(x + 1, y, w - 2, style.h_line)
  self.set(x + w - 1, y, style.round_tr)
  // Bottom edge
  self.set(x, y + h - 1, style.round_bl)
  self.draw_h_line(x + 1, y + h - 1, w - 2, style.h_line)
  self.set(x + w - 1, y + h - 1, style.round_br)
  // Side edges
  self.draw_v_line(x, y + 1, h - 2, style.v_line)
  self.draw_v_line(x + w - 1, y + 1, h - 2, style.v_line)
}

/// Convert canvas to string
pub fn Canvas::to_string(self : Canvas) -> String {
  let result = StringBuilder::new()
  for y = 0; y < self.height; y = y + 1 {
    // Find last non-space character to trim trailing spaces
    let mut last_non_space = -1
    for x = self.width - 1; x >= 0; x = x - 1 {
      if self.cells[y][x] != ' ' {
        last_non_space = x
        break
      }
    }
    for x = 0; x <= last_non_space; x = x + 1 {
      result.write_char(self.cells[y][x])
    }
    if y < self.height - 1 {
      result.write_char('\n')
    }
  }
  result.to_string()
}
