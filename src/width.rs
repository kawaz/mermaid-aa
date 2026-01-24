//! Character width calculation with ambiguous width support.

use unicode_width::UnicodeWidthChar;

/// Ambiguous width setting for East Asian ambiguous characters.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Default)]
pub enum AmbiguousWidth {
    /// Half-width (1): for Western terminals
    #[default]
    Half = 1,
    /// Full-width (2): for CJK terminals
    Full = 2,
}

impl From<u8> for AmbiguousWidth {
    fn from(n: u8) -> Self {
        match n {
            2 => Self::Full,
            _ => Self::Half,
        }
    }
}

/// Check if a character is in the East Asian Ambiguous category.
/// Includes box drawing characters (U+2500-U+257F) and some symbols.
fn is_ambiguous(c: char) -> bool {
    let cp = c as u32;
    matches!(
        cp,
        // Box drawing characters
        0x2500..=0x257F |
        // Block elements
        0x2580..=0x259F |
        // Arrows
        0x2190..=0x21FF |
        // Mathematical operators (partial)
        0x2200..=0x22FF |
        // Miscellaneous symbols (partial)
        0x2600..=0x26FF |
        // Dingbats (partial)
        0x2700..=0x27BF
    )
}

/// Calculate the display width of a character.
pub fn char_width(c: char, ambiguous_width: AmbiguousWidth) -> usize {
    if is_ambiguous(c) {
        return ambiguous_width as usize;
    }
    c.width().unwrap_or(0)
}

/// Calculate the display width of a string.
pub fn str_width(s: &str, ambiguous_width: AmbiguousWidth) -> usize {
    s.chars().map(|c| char_width(c, ambiguous_width)).sum()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_ascii() {
        assert_eq!(char_width('a', AmbiguousWidth::Half), 1);
        assert_eq!(char_width('Z', AmbiguousWidth::Half), 1);
    }

    #[test]
    fn test_fullwidth() {
        assert_eq!(char_width('日', AmbiguousWidth::Half), 2);
        assert_eq!(char_width('本', AmbiguousWidth::Half), 2);
    }

    #[test]
    fn test_ambiguous_half() {
        assert_eq!(char_width('─', AmbiguousWidth::Half), 1);
        assert_eq!(char_width('│', AmbiguousWidth::Half), 1);
        assert_eq!(char_width('→', AmbiguousWidth::Half), 1);
    }

    #[test]
    fn test_ambiguous_full() {
        assert_eq!(char_width('─', AmbiguousWidth::Full), 2);
        assert_eq!(char_width('│', AmbiguousWidth::Full), 2);
        assert_eq!(char_width('→', AmbiguousWidth::Full), 2);
    }

    #[test]
    fn test_str_width() {
        assert_eq!(str_width("hello", AmbiguousWidth::Half), 5);
        assert_eq!(str_width("日本語", AmbiguousWidth::Half), 6);
        assert_eq!(str_width("A→B", AmbiguousWidth::Half), 3);
        assert_eq!(str_width("A→B", AmbiguousWidth::Full), 4);
    }
}
