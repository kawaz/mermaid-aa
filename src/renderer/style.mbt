// Style module - Character set definitions for ASCII art rendering

/// Character set for rendering
pub struct CharStyle {
  // Box drawing
  h_line : Char       // Horizontal line
  v_line : Char       // Vertical line
  corner_tl : Char    // Top-left corner
  corner_tr : Char    // Top-right corner
  corner_bl : Char    // Bottom-left corner
  corner_br : Char    // Bottom-right corner
  cross : Char        // Cross intersection
  t_down : Char       // T pointing down
  t_up : Char         // T pointing up
  t_right : Char      // T pointing right
  t_left : Char       // T pointing left
  // Arrows
  arrow_right : Char  // Right arrow
  arrow_left : Char   // Left arrow
  arrow_down : Char   // Down arrow
  arrow_up : Char     // Up arrow
  // Special shapes
  round_tl : Char     // Rounded top-left
  round_tr : Char     // Rounded top-right
  round_bl : Char     // Rounded bottom-left
  round_br : Char     // Rounded bottom-right
  diamond_left : Char
  diamond_right : Char
  paren_left : Char
  paren_right : Char
}

/// Get ASCII character set
pub fn ascii_style() -> CharStyle {
  {
    h_line: '-',
    v_line: '|',
    corner_tl: '+',
    corner_tr: '+',
    corner_bl: '+',
    corner_br: '+',
    cross: '+',
    t_down: '+',
    t_up: '+',
    t_right: '+',
    t_left: '+',
    arrow_right: '>',
    arrow_left: '<',
    arrow_down: 'v',
    arrow_up: '^',
    round_tl: '(',
    round_tr: ')',
    round_bl: '(',
    round_br: ')',
    diamond_left: '<',
    diamond_right: '>',
    paren_left: '(',
    paren_right: ')',
  }
}

/// Get Unicode character set
pub fn unicode_style() -> CharStyle {
  {
    h_line: '\u2500',         // horizontal line
    v_line: '\u2502',         // vertical line
    corner_tl: '\u250C',      // top-left corner
    corner_tr: '\u2510',      // top-right corner
    corner_bl: '\u2514',      // bottom-left corner
    corner_br: '\u2518',      // bottom-right corner
    cross: '\u253C',          // cross
    t_down: '\u252C',         // t pointing down
    t_up: '\u2534',           // t pointing up
    t_right: '\u251C',        // t pointing right
    t_left: '\u2524',         // t pointing left
    arrow_right: '\u2192',    // right arrow
    arrow_left: '\u2190',     // left arrow
    arrow_down: '\u2193',     // down arrow
    arrow_up: '\u2191',       // up arrow
    round_tl: '(',
    round_tr: ')',
    round_bl: '(',
    round_br: ')',
    diamond_left: '<',
    diamond_right: '>',
    paren_left: '(',
    paren_right: ')',
  }
}

/// Get Unicode rounded character set
pub fn unicode_round_style() -> CharStyle {
  {
    h_line: '\u2500',         // horizontal line
    v_line: '\u2502',         // vertical line
    corner_tl: '\u256D',      // rounded top-left
    corner_tr: '\u256E',      // rounded top-right
    corner_bl: '\u2570',      // rounded bottom-left
    corner_br: '\u256F',      // rounded bottom-right
    cross: '\u253C',          // cross
    t_down: '\u252C',         // t pointing down
    t_up: '\u2534',           // t pointing up
    t_right: '\u251C',        // t pointing right
    t_left: '\u2524',         // t pointing left
    arrow_right: '\u2192',    // right arrow
    arrow_left: '\u2190',     // left arrow
    arrow_down: '\u2193',     // down arrow
    arrow_up: '\u2191',       // up arrow
    round_tl: '\u256D',       // rounded top-left
    round_tr: '\u256E',       // rounded top-right
    round_bl: '\u2570',       // rounded bottom-left
    round_br: '\u256F',       // rounded bottom-right
    diamond_left: '<',
    diamond_right: '>',
    paren_left: '(',
    paren_right: ')',
  }
}

/// Get Unicode bold character set
pub fn unicode_bold_style() -> CharStyle {
  {
    h_line: '\u2501',         // bold horizontal line
    v_line: '\u2503',         // bold vertical line
    corner_tl: '\u250F',      // bold top-left
    corner_tr: '\u2513',      // bold top-right
    corner_bl: '\u2517',      // bold bottom-left
    corner_br: '\u251B',      // bold bottom-right
    cross: '\u254B',          // bold cross
    t_down: '\u2533',         // bold t pointing down
    t_up: '\u253B',           // bold t pointing up
    t_right: '\u2523',        // bold t pointing right
    t_left: '\u252B',         // bold t pointing left
    arrow_right: '\u279C',    // heavy right arrow
    arrow_left: '\u2190',     // left arrow
    arrow_down: '\u2193',     // down arrow
    arrow_up: '\u2191',       // up arrow
    round_tl: '(',
    round_tr: ')',
    round_bl: '(',
    round_br: ')',
    diamond_left: '<',
    diamond_right: '>',
    paren_left: '(',
    paren_right: ')',
  }
}

/// Get Unicode double character set
pub fn unicode_double_style() -> CharStyle {
  {
    h_line: '\u2550',         // double horizontal line
    v_line: '\u2551',         // double vertical line
    corner_tl: '\u2554',      // double top-left
    corner_tr: '\u2557',      // double top-right
    corner_bl: '\u255A',      // double bottom-left
    corner_br: '\u255D',      // double bottom-right
    cross: '\u256C',          // double cross
    t_down: '\u2566',         // double t pointing down
    t_up: '\u2569',           // double t pointing up
    t_right: '\u2560',        // double t pointing right
    t_left: '\u2563',         // double t pointing left
    arrow_right: '\u21D2',    // double right arrow
    arrow_left: '\u21D0',     // double left arrow
    arrow_down: '\u21D3',     // double down arrow
    arrow_up: '\u21D1',       // double up arrow
    round_tl: '(',
    round_tr: ')',
    round_bl: '(',
    round_br: ')',
    diamond_left: '<',
    diamond_right: '>',
    paren_left: '(',
    paren_right: ')',
  }
}

/// Get character style for charset
pub fn get_style(charset : @types.Charset) -> CharStyle {
  match charset {
    @types.Ascii => ascii_style()
    @types.Unicode => unicode_style()
    @types.UnicodeRound => unicode_round_style()
    @types.UnicodeBold => unicode_bold_style()
    @types.UnicodeDouble => unicode_double_style()
  }
}
