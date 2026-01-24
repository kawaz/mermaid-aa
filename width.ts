/**
 * Character width calculation with East Asian Width support
 */

/**
 * Ambiguous width mode configuration
 * - ambiguousWidth: Width for general ambiguous characters
 * - boxDrawingWidth: Width for Box Drawing characters (U+2500-U+257F)
 */
export interface AmbiguousWidthConfig {
  ambiguousWidth: 1 | 2;
  boxDrawingWidth: 1 | 2;
}

/**
 * Ambiguous width mode values:
 * - 1 or "half": All ambiguous = 1 (Western terminals)
 * - 2 or "full": Ambiguous = 2, Box Drawing = 1 (CJK terminals, recommended)
 * - "console": Same as "half" (explicit name)
 * - "legacy": All ambiguous including Box Drawing = 2 (legacy compatibility)
 */
export type AmbiguousWidthMode = 1 | 2 | "half" | "full" | "console" | "legacy";

/**
 * Parse ambiguous width mode to config
 */
export function parseAmbiguousWidthMode(
  mode: AmbiguousWidthMode,
): AmbiguousWidthConfig {
  switch (mode) {
    case 1:
    case "half":
    case "console":
      return { ambiguousWidth: 1, boxDrawingWidth: 1 };
    case 2:
    case "full":
      return { ambiguousWidth: 2, boxDrawingWidth: 1 };
    case "legacy":
      return { ambiguousWidth: 2, boxDrawingWidth: 2 };
  }
}

// Box Drawing characters range (U+2500-U+257F)
const BOX_DRAWING_RANGE: [number, number] = [0x2500, 0x257f];

/**
 * Check if a code point is a Box Drawing character
 */
export function isBoxDrawing(codePoint: number): boolean {
  return codePoint >= BOX_DRAWING_RANGE[0] && codePoint <= BOX_DRAWING_RANGE[1];
}

// Unicode East Asian Width ranges
// Full-width characters
const FULLWIDTH_RANGES: [number, number][] = [
  [0x1100, 0x115f], // Hangul Jamo
  [0x2e80, 0x2eff], // CJK Radicals Supplement
  [0x2f00, 0x2fdf], // Kangxi Radicals
  [0x2ff0, 0x2fff], // Ideographic Description Characters
  [0x3000, 0x303e], // CJK Symbols and Punctuation
  [0x3041, 0x3096], // Hiragana
  [0x30a1, 0x30ff], // Katakana
  [0x3105, 0x312d], // Bopomofo
  [0x3131, 0x318e], // Hangul Compatibility Jamo
  [0x3190, 0x31ba], // Kanbun
  [0x31c0, 0x31e3], // CJK Strokes
  [0x31f0, 0x31ff], // Katakana Phonetic Extensions
  [0x3200, 0x321e], // Enclosed CJK Letters and Months
  [0x3220, 0x3247], // Enclosed CJK Letters and Months (continued)
  [0x3250, 0x32fe], // Enclosed CJK Letters and Months (continued)
  [0x3300, 0x4dbf], // CJK Compatibility + Extension A
  [0x4e00, 0x9fff], // CJK Unified Ideographs
  [0xa000, 0xa4cf], // Yi Syllables + Radicals
  [0xa960, 0xa97f], // Hangul Jamo Extended-A
  [0xac00, 0xd7a3], // Hangul Syllables
  [0xf900, 0xfaff], // CJK Compatibility Ideographs
  [0xfe10, 0xfe1f], // Vertical Forms
  [0xfe30, 0xfe6f], // CJK Compatibility Forms
  [0xff00, 0xff60], // Fullwidth Forms
  [0xffe0, 0xffe6], // Fullwidth Forms
  [0x20000, 0x2fffd], // CJK Unified Ideographs Extension B+
  [0x30000, 0x3fffd], // CJK Unified Ideographs Extension G+
];

// Ambiguous width characters (treated as width 1 or 2 depending on context)
const AMBIGUOUS_RANGES: [number, number][] = [
  [0x00a1, 0x00a1], // Inverted Exclamation Mark
  [0x00a4, 0x00a4], // Currency Sign
  [0x00a7, 0x00a8], // Section Sign, Diaeresis
  [0x00aa, 0x00aa], // Feminine Ordinal Indicator
  [0x00ad, 0x00ae], // Soft Hyphen, Registered Sign
  [0x00b0, 0x00b4], // Degree Sign through Acute Accent
  [0x00b6, 0x00ba], // Pilcrow Sign through Masculine Ordinal Indicator
  [0x00bc, 0x00bf], // Vulgar Fractions, Inverted Question Mark
  [0x00c6, 0x00c6], // Latin Capital Letter Ae
  [0x00d0, 0x00d0], // Latin Capital Letter Eth
  [0x00d7, 0x00d8], // Multiplication Sign, Latin Capital Letter O with Stroke
  [0x00de, 0x00e1], // Latin Capital Letter Thorn through Small Letter A with Acute
  [0x00e6, 0x00e6], // Latin Small Letter Ae
  [0x00e8, 0x00ea], // Latin Small Letter E with Grave through E with Circumflex
  [0x00ec, 0x00ed], // Latin Small Letter I with Grave, I with Acute
  [0x00f0, 0x00f0], // Latin Small Letter Eth
  [0x00f2, 0x00f3], // Latin Small Letter O with Grave, O with Acute
  [0x00f7, 0x00fa], // Division Sign through U with Acute
  [0x00fc, 0x00fc], // Latin Small Letter U with Diaeresis
  [0x00fe, 0x00fe], // Latin Small Letter Thorn
  [0x0101, 0x0101], // Latin Small Letter A with Macron
  [0x0111, 0x0111], // Latin Small Letter D with Stroke
  [0x0113, 0x0113], // Latin Small Letter E with Macron
  [0x011b, 0x011b], // Latin Small Letter E with Caron
  [0x0126, 0x0127], // Latin Capital/Small Letter H with Stroke
  [0x012b, 0x012b], // Latin Small Letter I with Macron
  [0x0131, 0x0133], // Latin Small Letter Dotless I through IJ
  [0x0138, 0x0138], // Latin Small Letter Kra
  [0x013f, 0x0142], // Latin Capital/Small Letter L with Middle Dot, Stroke
  [0x0144, 0x0144], // Latin Small Letter N with Acute
  [0x0148, 0x014b], // Latin Small Letter N with Caron through Eng
  [0x014d, 0x014d], // Latin Small Letter O with Macron
  [0x0152, 0x0153], // Latin Capital/Small Ligature OE
  [0x0166, 0x0167], // Latin Capital/Small Letter T with Stroke
  [0x016b, 0x016b], // Latin Small Letter U with Macron
  [0x01ce, 0x01ce], // Latin Small Letter A with Caron
  [0x01d0, 0x01d0], // Latin Small Letter I with Caron
  [0x01d2, 0x01d2], // Latin Small Letter O with Caron
  [0x01d4, 0x01d4], // Latin Small Letter U with Caron
  [0x01d6, 0x01d6], // Latin Small Letter U with Diaeresis and Macron
  [0x01d8, 0x01d8], // Latin Small Letter U with Diaeresis and Acute
  [0x01da, 0x01da], // Latin Small Letter U with Diaeresis and Caron
  [0x01dc, 0x01dc], // Latin Small Letter U with Diaeresis and Grave
  [0x0251, 0x0251], // Latin Small Letter Alpha
  [0x0261, 0x0261], // Latin Small Letter Script G
  [0x02c4, 0x02c4], // Modifier Letter Up Arrowhead
  [0x02c7, 0x02c7], // Caron
  [0x02c9, 0x02cb], // Modifier Letter Macron through Grave Accent
  [0x02cd, 0x02cd], // Modifier Letter Low Macron
  [0x02d0, 0x02d0], // Modifier Letter Triangular Colon
  [0x02d8, 0x02db], // Breve through Ogonek
  [0x02dd, 0x02dd], // Double Acute Accent
  [0x02df, 0x02df], // Modifier Letter Cross Accent
  [0x0300, 0x036f], // Combining Diacritical Marks (zero width, but treated as ambiguous in some contexts)
  [0x0391, 0x03a1], // Greek Capital Letters
  [0x03a3, 0x03a9], // Greek Capital Letters (Sigma through Omega)
  [0x03b1, 0x03c1], // Greek Small Letters
  [0x03c3, 0x03c9], // Greek Small Letters (Sigma through Omega)
  [0x0401, 0x0401], // Cyrillic Capital Letter Io
  [0x0410, 0x044f], // Cyrillic Letters
  [0x0451, 0x0451], // Cyrillic Small Letter Io
  [0x2010, 0x2010], // Hyphen
  [0x2013, 0x2016], // En Dash through Double Vertical Line
  [0x2018, 0x2019], // Single Quotation Marks
  [0x201c, 0x201d], // Double Quotation Marks
  [0x2020, 0x2022], // Dagger, Double Dagger, Bullet
  [0x2024, 0x2027], // One Dot Leader through Hyphenation Point
  [0x2030, 0x2030], // Per Mille Sign
  [0x2032, 0x2033], // Prime, Double Prime
  [0x2035, 0x2035], // Reversed Prime
  [0x203b, 0x203b], // Reference Mark
  [0x203e, 0x203e], // Overline
  [0x2074, 0x2074], // Superscript Four
  [0x207f, 0x207f], // Superscript Latin Small Letter N
  [0x2081, 0x2084], // Subscript One through Four
  [0x20ac, 0x20ac], // Euro Sign
  [0x2103, 0x2103], // Degree Celsius
  [0x2105, 0x2105], // Care Of
  [0x2109, 0x2109], // Degree Fahrenheit
  [0x2113, 0x2113], // Script Small L
  [0x2116, 0x2116], // Numero Sign
  [0x2121, 0x2122], // Telephone Sign, Trade Mark Sign
  [0x2126, 0x2126], // Ohm Sign
  [0x212b, 0x212b], // Angstrom Sign
  [0x2153, 0x2154], // Vulgar Fraction One Third, Two Thirds
  [0x215b, 0x215e], // Vulgar Fractions
  [0x2160, 0x216b], // Roman Numerals
  [0x2170, 0x2179], // Small Roman Numerals
  [0x2189, 0x2189], // Vulgar Fraction Zero Thirds
  [0x2190, 0x2199], // Arrows
  [0x21b8, 0x21b9], // North West Arrow to Long Bar, Leftwards Arrow to Bar over Rightwards Arrow to Bar
  [0x21d2, 0x21d2], // Rightwards Double Arrow
  [0x21d4, 0x21d4], // Left Right Double Arrow
  [0x21e7, 0x21e7], // Upwards White Arrow
  [0x2200, 0x2200], // For All
  [0x2202, 0x2203], // Partial Differential, There Exists
  [0x2207, 0x2208], // Nabla, Element Of
  [0x220b, 0x220b], // Contains as Member
  [0x220f, 0x220f], // N-Ary Product
  [0x2211, 0x2211], // N-Ary Summation
  [0x2215, 0x2215], // Division Slash
  [0x221a, 0x221a], // Square Root
  [0x221d, 0x2220], // Proportional To through Angle
  [0x2223, 0x2223], // Divides
  [0x2225, 0x2225], // Parallel To
  [0x2227, 0x222c], // Logical And through Double Integral
  [0x222e, 0x222e], // Contour Integral
  [0x2234, 0x2237], // Therefore through Proportion
  [0x223c, 0x223d], // Tilde Operator, Reversed Tilde
  [0x2248, 0x2248], // Almost Equal To
  [0x224c, 0x224c], // All Equal To
  [0x2252, 0x2252], // Approximately Equal To or the Image Of
  [0x2260, 0x2261], // Not Equal To, Identical To
  [0x2264, 0x2267], // Less-Than or Equal To through Greater-Than or Equal To
  [0x226a, 0x226b], // Much Less-Than, Much Greater-Than
  [0x226e, 0x226f], // Not Less-Than, Not Greater-Than
  [0x2282, 0x2283], // Subset Of, Superset Of
  [0x2286, 0x2287], // Subset Of or Equal To, Superset Of or Equal To
  [0x2295, 0x2295], // Circled Plus
  [0x2299, 0x2299], // Circled Dot Operator
  [0x22a5, 0x22a5], // Up Tack
  [0x22bf, 0x22bf], // Right Triangle
  [0x2312, 0x2312], // Arc
  [0x2500, 0x257f], // Box Drawing
  [0x2580, 0x258f], // Block Elements
  [0x2592, 0x2595], // Medium Shade through Right One Eighth Block
  [0x25a0, 0x25a1], // Black Square, White Square
  [0x25a3, 0x25a9], // White Square with Rounded Corners through Square with Diagonal Crosshatch Fill
  [0x25b2, 0x25b3], // Black Up-Pointing Triangle, White Up-Pointing Triangle
  [0x25b6, 0x25b7], // Black Right-Pointing Triangle, White Right-Pointing Triangle
  [0x25bc, 0x25bd], // Black Down-Pointing Triangle, White Down-Pointing Triangle
  [0x25c0, 0x25c1], // Black Left-Pointing Triangle, White Left-Pointing Triangle
  [0x25c6, 0x25c8], // Black Diamond through White Diamond Containing Black Small Diamond
  [0x25cb, 0x25cb], // White Circle
  [0x25ce, 0x25d1], // Bullseye through Circle with Right Half Black
  [0x25e2, 0x25e5], // Black Lower Right Triangle through Black Upper Right Triangle
  [0x25ef, 0x25ef], // Large Circle
  [0x2605, 0x2606], // Black Star, White Star
  [0x2609, 0x2609], // Sun
  [0x260e, 0x260f], // Black Telephone, White Telephone
  [0x2614, 0x2615], // Umbrella with Rain Drops, Hot Beverage
  [0x261c, 0x261c], // White Left Pointing Index
  [0x261e, 0x261e], // White Right Pointing Index
  [0x2640, 0x2640], // Female Sign
  [0x2642, 0x2642], // Male Sign
  [0x2660, 0x2661], // Black Spade Suit, White Heart Suit
  [0x2663, 0x2665], // Black Club Suit through Black Heart Suit
  [0x2667, 0x266a], // White Club Suit through Eighth Note
  [0x266c, 0x266d], // Beamed Sixteenth Notes, Music Flat Sign
  [0x266f, 0x266f], // Music Sharp Sign
  [0x269e, 0x269f], // Three Lines Converging Right, Three Lines Converging Left
  [0x26be, 0x26bf], // Baseball, Squared Key
  [0x26c4, 0x26cd], // Snowman without Snow through Disabled Car
  [0x26cf, 0x26e1], // Pick through Restricted Left Entry-2
  [0x26e3, 0x26e3], // Heavy Circle with Stroke and Two Dots Above
  [0x26e8, 0x26ff], // Black Cross on Shield through White Flag with Horizontal Middle Black Stripe
  [0x273d, 0x273d], // Heavy Teardrop-Spoked Asterisk
  [0x2757, 0x2757], // Heavy Exclamation Mark Symbol
  [0x2776, 0x277f], // Dingbat Negative Circled Digit One through Dingbat Negative Circled Number Ten
  [0x2b55, 0x2b59], // Heavy Large Circle through Heavy Circled Saltire
  [0xe000, 0xf8ff], // Private Use Area
  [0xfe00, 0xfe0f], // Variation Selectors
  [0xfffd, 0xfffd], // Replacement Character
  [0x1f100, 0x1f10a], // Enclosed Alphanumeric Supplement
  [0x1f110, 0x1f12d], // Enclosed Alphanumeric Supplement
  [0x1f130, 0x1f169], // Enclosed Alphanumeric Supplement
  [0x1f170, 0x1f19a], // Enclosed Alphanumeric Supplement
  [0xe0100, 0xe01ef], // Variation Selectors Supplement
  [0xf0000, 0xffffd], // Supplementary Private Use Area-A
  [0x100000, 0x10fffd], // Supplementary Private Use Area-B
];

function isInRanges(codePoint: number, ranges: [number, number][]): boolean {
  for (const [start, end] of ranges) {
    if (codePoint >= start && codePoint <= end) {
      return true;
    }
  }
  return false;
}

/** Check if a character is fullwidth */
export function isFullwidth(codePoint: number): boolean {
  return isInRanges(codePoint, FULLWIDTH_RANGES);
}

/** Check if a character is ambiguous width */
export function isAmbiguous(codePoint: number): boolean {
  return isInRanges(codePoint, AMBIGUOUS_RANGES);
}

/** Get the display width of a character */
export function charWidth(
  char: string,
  ambiguousWidthMode: AmbiguousWidthMode = 1,
): number {
  const codePoint = char.codePointAt(0);
  if (codePoint === undefined) return 0;

  // Control characters have zero width
  if (codePoint < 0x20 || (codePoint >= 0x7f && codePoint < 0xa0)) {
    return 0;
  }

  // Combining characters have zero width
  if (codePoint >= 0x0300 && codePoint <= 0x036f) {
    return 0;
  }

  if (isFullwidth(codePoint)) {
    return 2;
  }

  if (isAmbiguous(codePoint)) {
    const config = parseAmbiguousWidthMode(ambiguousWidthMode);
    // Box Drawing characters have separate width setting
    if (isBoxDrawing(codePoint)) {
      return config.boxDrawingWidth;
    }
    return config.ambiguousWidth;
  }

  return 1;
}

/** Get the display width of a string */
export function stringWidth(
  str: string,
  ambiguousWidthMode: AmbiguousWidthMode = 1,
): number {
  let width = 0;
  for (const char of str) {
    width += charWidth(char, ambiguousWidthMode);
  }
  return width;
}

/** Pad a string to the specified display width */
export function padString(
  str: string,
  targetWidth: number,
  padChar: string = " ",
  ambiguousWidthMode: AmbiguousWidthMode = 1,
): string {
  const currentWidth = stringWidth(str, ambiguousWidthMode);
  if (currentWidth >= targetWidth) {
    return str;
  }
  const padWidth = targetWidth - currentWidth;
  const padCharWidth = charWidth(padChar, ambiguousWidthMode);
  const padCount = Math.floor(padWidth / padCharWidth);
  return str + padChar.repeat(padCount);
}

/** Center a string within the specified display width */
export function centerString(
  str: string,
  targetWidth: number,
  padChar: string = " ",
  ambiguousWidthMode: AmbiguousWidthMode = 1,
): string {
  const currentWidth = stringWidth(str, ambiguousWidthMode);
  if (currentWidth >= targetWidth) {
    return str;
  }
  const totalPad = targetWidth - currentWidth;
  const leftPad = Math.floor(totalPad / 2);
  const rightPad = totalPad - leftPad;
  return padChar.repeat(leftPad) + str + padChar.repeat(rightPad);
}
