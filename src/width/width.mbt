// Width module - Unicode character width calculation
// Implements East Asian Width algorithm for proper text alignment

/// Helper to get char from string at index
fn char_at(s : String, i : Int) -> Char {
  s[i].to_int().unsafe_to_char()
}

/// Check if a character is fullwidth (always 2 cells)
fn is_fullwidth(cp : Int) -> Bool {
  // CJK Unified Ideographs and related
  (cp >= 0x4E00 && cp <= 0x9FFF) ||    // CJK Unified Ideographs
  (cp >= 0x3400 && cp <= 0x4DBF) ||    // CJK Extension A
  (cp >= 0x20000 && cp <= 0x2A6DF) ||  // CJK Extension B
  (cp >= 0x2A700 && cp <= 0x2B73F) ||  // CJK Extension C
  (cp >= 0x2B740 && cp <= 0x2B81F) ||  // CJK Extension D
  (cp >= 0x2B820 && cp <= 0x2CEAF) ||  // CJK Extension E
  (cp >= 0x2CEB0 && cp <= 0x2EBEF) ||  // CJK Extension F
  (cp >= 0x30000 && cp <= 0x3134F) ||  // CJK Extension G
  (cp >= 0xF900 && cp <= 0xFAFF) ||    // CJK Compatibility Ideographs
  (cp >= 0x2F800 && cp <= 0x2FA1F) ||  // CJK Compatibility Ideographs Supplement
  // Hangul
  (cp >= 0xAC00 && cp <= 0xD7AF) ||    // Hangul Syllables
  (cp >= 0x1100 && cp <= 0x11FF) ||    // Hangul Jamo
  (cp >= 0x3130 && cp <= 0x318F) ||    // Hangul Compatibility Jamo
  (cp >= 0xA960 && cp <= 0xA97F) ||    // Hangul Jamo Extended-A
  (cp >= 0xD7B0 && cp <= 0xD7FF) ||    // Hangul Jamo Extended-B
  // Japanese Kana
  (cp >= 0x3040 && cp <= 0x309F) ||    // Hiragana
  (cp >= 0x30A0 && cp <= 0x30FF) ||    // Katakana
  (cp >= 0x31F0 && cp <= 0x31FF) ||    // Katakana Phonetic Extensions
  // Fullwidth Forms
  (cp >= 0xFF00 && cp <= 0xFF60) ||    // Fullwidth ASCII variants
  (cp >= 0xFFE0 && cp <= 0xFFE6) ||    // Fullwidth symbol variants
  // Other Wide
  (cp >= 0x3000 && cp <= 0x303F) ||    // CJK Symbols and Punctuation
  (cp >= 0x3200 && cp <= 0x32FF) ||    // Enclosed CJK Letters and Months
  (cp >= 0x3300 && cp <= 0x33FF) ||    // CJK Compatibility
  (cp >= 0xFE30 && cp <= 0xFE4F)       // CJK Compatibility Forms
}

/// Check if a character is a Box Drawing character (U+2500-U+257F)
/// These are treated specially in some width modes
fn is_box_drawing(cp : Int) -> Bool {
  cp >= 0x2500 && cp <= 0x257F
}

/// Check if a character is ambiguous width
/// These characters are width 1 or 2 depending on context
fn is_ambiguous(cp : Int) -> Bool {
  // Box Drawing (also ambiguous, but handled separately)
  (cp >= 0x2500 && cp <= 0x257F) ||
  // Block Elements
  (cp >= 0x2580 && cp <= 0x259F) ||
  // Geometric Shapes (partial)
  (cp >= 0x25A0 && cp <= 0x25FF) ||
  // Arrows (many are ambiguous)
  (cp >= 0x2190 && cp <= 0x21FF) ||
  // Mathematical Operators (some)
  (cp >= 0x2200 && cp <= 0x22FF) ||
  // Miscellaneous Symbols (some)
  (cp >= 0x2600 && cp <= 0x26FF) ||
  // Dingbats (some)
  (cp >= 0x2700 && cp <= 0x27BF) ||
  // Specific characters
  cp == 0x203B ||  // Reference mark
  cp == 0x2015 ||  // Horizontal bar
  cp == 0x2014 ||  // Em dash
  cp == 0x00A7 ||  // Section sign
  cp == 0x00A8 ||  // Diaeresis
  cp == 0x00B0 ||  // Degree
  cp == 0x00B1 ||  // Plus-minus
  cp == 0x00B4 ||  // Acute accent
  cp == 0x00D7 ||  // Multiplication
  cp == 0x00F7 ||  // Division
  // Greek letters (common in CJK contexts)
  (cp >= 0x0391 && cp <= 0x03C9) ||
  // Cyrillic letters (some)
  (cp >= 0x0410 && cp <= 0x044F)
}

/// Check if a character is zero-width (combining marks, etc.)
fn is_zero_width(cp : Int) -> Bool {
  // Combining Diacritical Marks
  (cp >= 0x0300 && cp <= 0x036F) ||
  // Combining Diacritical Marks Extended
  (cp >= 0x1AB0 && cp <= 0x1AFF) ||
  // Combining Diacritical Marks Supplement
  (cp >= 0x1DC0 && cp <= 0x1DFF) ||
  // Combining Diacritical Marks for Symbols
  (cp >= 0x20D0 && cp <= 0x20FF) ||
  // Combining Half Marks
  (cp >= 0xFE20 && cp <= 0xFE2F) ||
  // Zero Width characters
  cp == 0x200B ||  // Zero Width Space
  cp == 0x200C ||  // Zero Width Non-Joiner
  cp == 0x200D ||  // Zero Width Joiner
  cp == 0xFEFF     // Zero Width No-Break Space / BOM
}

/// Calculate display width of a single character
pub fn char_width(ch : Char, ambiguous_width : Int) -> Int {
  let cp = ch.to_int()
  if is_zero_width(cp) {
    0
  } else if is_fullwidth(cp) {
    2
  } else if is_ambiguous(cp) {
    ambiguous_width
  } else {
    1
  }
}

/// Calculate display width of a single character using WidthMode
pub fn char_width_mode(ch : Char, mode : @types.WidthMode) -> Int {
  let cp = ch.to_int()
  if is_zero_width(cp) {
    0
  } else if is_fullwidth(cp) {
    2
  } else if is_box_drawing(cp) {
    mode.box_drawing_width()
  } else if is_ambiguous(cp) {
    mode.ambiguous_width()
  } else {
    1
  }
}

/// Calculate display width of a string
pub fn string_width(s : String, ambiguous_width : Int) -> Int {
  let mut width = 0
  for i = 0; i < s.length(); i = i + 1 {
    width = width + char_width(char_at(s, i), ambiguous_width)
  }
  width
}

/// Calculate display width of a string using WidthMode
pub fn string_width_mode(s : String, mode : @types.WidthMode) -> Int {
  let mut width = 0
  for i = 0; i < s.length(); i = i + 1 {
    width = width + char_width_mode(char_at(s, i), mode)
  }
  width
}

/// Pad a string to a given display width (left-aligned)
pub fn pad_right(s : String, target_width : Int, ambiguous_width : Int) -> String {
  let current_width = string_width(s, ambiguous_width)
  if current_width >= target_width {
    s
  } else {
    let padding = target_width - current_width
    let result = StringBuilder::new()
    result.write_string(s)
    for i = 0; i < padding; i = i + 1 {
      result.write_char(' ')
    }
    result.to_string()
  }
}

/// Pad a string to a given display width (left-aligned) using WidthMode
pub fn pad_right_mode(s : String, target_width : Int, mode : @types.WidthMode) -> String {
  let current_width = string_width_mode(s, mode)
  if current_width >= target_width {
    s
  } else {
    let padding = target_width - current_width
    let result = StringBuilder::new()
    result.write_string(s)
    for i = 0; i < padding; i = i + 1 {
      result.write_char(' ')
    }
    result.to_string()
  }
}

/// Pad a string to a given display width (centered)
pub fn pad_center(s : String, target_width : Int, ambiguous_width : Int) -> String {
  let current_width = string_width(s, ambiguous_width)
  if current_width >= target_width {
    s
  } else {
    let padding = target_width - current_width
    let left_padding = padding / 2
    let right_padding = padding - left_padding
    let result = StringBuilder::new()
    for i = 0; i < left_padding; i = i + 1 {
      result.write_char(' ')
    }
    result.write_string(s)
    for i = 0; i < right_padding; i = i + 1 {
      result.write_char(' ')
    }
    result.to_string()
  }
}

/// Pad a string to a given display width (centered) using WidthMode
pub fn pad_center_mode(s : String, target_width : Int, mode : @types.WidthMode) -> String {
  let current_width = string_width_mode(s, mode)
  if current_width >= target_width {
    s
  } else {
    let padding = target_width - current_width
    let left_padding = padding / 2
    let right_padding = padding - left_padding
    let result = StringBuilder::new()
    for i = 0; i < left_padding; i = i + 1 {
      result.write_char(' ')
    }
    result.write_string(s)
    for i = 0; i < right_padding; i = i + 1 {
      result.write_char(' ')
    }
    result.to_string()
  }
}

/// Pad a string to a given display width (right-aligned)
pub fn pad_left(s : String, target_width : Int, ambiguous_width : Int) -> String {
  let current_width = string_width(s, ambiguous_width)
  if current_width >= target_width {
    s
  } else {
    let padding = target_width - current_width
    let result = StringBuilder::new()
    for i = 0; i < padding; i = i + 1 {
      result.write_char(' ')
    }
    result.write_string(s)
    result.to_string()
  }
}

/// Pad a string to a given display width (right-aligned) using WidthMode
pub fn pad_left_mode(s : String, target_width : Int, mode : @types.WidthMode) -> String {
  let current_width = string_width_mode(s, mode)
  if current_width >= target_width {
    s
  } else {
    let padding = target_width - current_width
    let result = StringBuilder::new()
    for i = 0; i < padding; i = i + 1 {
      result.write_char(' ')
    }
    result.write_string(s)
    result.to_string()
  }
}
