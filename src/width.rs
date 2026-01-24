//! Character width calculation with ambiguous width support.

use std::str::FromStr;
use unicode_width::UnicodeWidthChar;

/// Ambiguous width setting for East Asian ambiguous characters.
///
/// Different terminals treat East Asian Ambiguous width characters differently:
/// - Western terminals: typically 1 cell (half-width)
/// - CJK terminals: typically 2 cells (full-width)
/// - Box Drawing characters (U+2500-U+257F) may be treated specially
#[derive(Debug, Clone, Copy, PartialEq, Eq, Default)]
pub enum AmbiguousWidth {
    /// Half-width (1 cell) for all ambiguous chars including Box Drawing.
    /// Suitable for Western terminals.
    #[default]
    Half,
    /// Full-width (2 cells) for ambiguous chars, but Box Drawing stays 1 cell.
    /// Recommended for CJK terminals.
    Full,
    /// Legacy mode: Full-width (2 cells) for all ambiguous chars including Box Drawing.
    /// For terminals that treat Box Drawing as full-width.
    Legacy,
}

impl AmbiguousWidth {
    /// Get the width for general ambiguous characters.
    pub fn ambiguous_char_width(self) -> usize {
        match self {
            Self::Half => 1,
            Self::Full | Self::Legacy => 2,
        }
    }

    /// Get the width for Box Drawing characters (U+2500-U+257F).
    pub fn box_drawing_width(self) -> usize {
        match self {
            Self::Half | Self::Full => 1,
            Self::Legacy => 2,
        }
    }
}

impl From<u8> for AmbiguousWidth {
    fn from(n: u8) -> Self {
        match n {
            2 => Self::Full,
            _ => Self::Half,
        }
    }
}

impl FromStr for AmbiguousWidth {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s.to_lowercase().as_str() {
            "1" | "half" | "console" => Ok(Self::Half),
            "2" | "full" => Ok(Self::Full),
            "legacy" => Ok(Self::Legacy),
            _ => Err(format!(
                "invalid ambiguous width '{}'. Valid values: 1, half, 2, full, console, legacy",
                s
            )),
        }
    }
}

/// Check if a character is a Box Drawing character (U+2500-U+257F).
fn is_box_drawing(c: char) -> bool {
    let cp = c as u32;
    (0x2500..=0x257F).contains(&cp)
}

/// Check if a character is in the East Asian Ambiguous category (excluding Box Drawing).
fn is_ambiguous(c: char) -> bool {
    let cp = c as u32;
    matches!(
        cp,
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
    if is_box_drawing(c) {
        return ambiguous_width.box_drawing_width();
    }
    if is_ambiguous(c) {
        return ambiguous_width.ambiguous_char_width();
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
        // Box Drawing - always 1 in Half mode
        assert_eq!(char_width('─', AmbiguousWidth::Half), 1);
        assert_eq!(char_width('│', AmbiguousWidth::Half), 1);
        // Other ambiguous chars
        assert_eq!(char_width('→', AmbiguousWidth::Half), 1);
    }

    #[test]
    fn test_ambiguous_full() {
        // Box Drawing - stays 1 even in Full mode
        assert_eq!(char_width('─', AmbiguousWidth::Full), 1);
        assert_eq!(char_width('│', AmbiguousWidth::Full), 1);
        // Other ambiguous chars - become 2 in Full mode
        assert_eq!(char_width('→', AmbiguousWidth::Full), 2);
    }

    #[test]
    fn test_ambiguous_legacy() {
        // Box Drawing - 2 in Legacy mode
        assert_eq!(char_width('─', AmbiguousWidth::Legacy), 2);
        assert_eq!(char_width('│', AmbiguousWidth::Legacy), 2);
        // Other ambiguous chars - also 2 in Legacy mode
        assert_eq!(char_width('→', AmbiguousWidth::Legacy), 2);
    }

    #[test]
    fn test_str_width() {
        assert_eq!(str_width("hello", AmbiguousWidth::Half), 5);
        assert_eq!(str_width("日本語", AmbiguousWidth::Half), 6);
        assert_eq!(str_width("A→B", AmbiguousWidth::Half), 3);
        assert_eq!(str_width("A→B", AmbiguousWidth::Full), 4);
    }

    #[test]
    fn test_box_drawing_vs_ambiguous() {
        // Box Drawing stays 1 in Full mode, but 2 in Legacy
        assert_eq!(str_width("┌─┐", AmbiguousWidth::Half), 3);
        assert_eq!(str_width("┌─┐", AmbiguousWidth::Full), 3);
        assert_eq!(str_width("┌─┐", AmbiguousWidth::Legacy), 6);

        // Mixed: arrow (ambiguous) + box drawing
        assert_eq!(str_width("→─", AmbiguousWidth::Half), 2);
        assert_eq!(str_width("→─", AmbiguousWidth::Full), 3); // arrow=2, box=1
        assert_eq!(str_width("→─", AmbiguousWidth::Legacy), 4); // both=2
    }

    #[test]
    fn test_from_str() {
        assert_eq!("1".parse::<AmbiguousWidth>().unwrap(), AmbiguousWidth::Half);
        assert_eq!(
            "half".parse::<AmbiguousWidth>().unwrap(),
            AmbiguousWidth::Half
        );
        assert_eq!(
            "console".parse::<AmbiguousWidth>().unwrap(),
            AmbiguousWidth::Half
        );
        assert_eq!("2".parse::<AmbiguousWidth>().unwrap(), AmbiguousWidth::Full);
        assert_eq!(
            "full".parse::<AmbiguousWidth>().unwrap(),
            AmbiguousWidth::Full
        );
        assert_eq!(
            "legacy".parse::<AmbiguousWidth>().unwrap(),
            AmbiguousWidth::Legacy
        );
        assert!("invalid".parse::<AmbiguousWidth>().is_err());
    }
}
