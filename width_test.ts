/**
 * Tests for width calculation
 */

import { assertEquals } from "@std/assert";
import {
  centerString,
  charWidth,
  isAmbiguous,
  isBoxDrawing,
  isFullwidth,
  padString,
  parseAmbiguousWidthMode,
  stringWidth,
} from "./width.ts";

Deno.test("charWidth - ASCII characters have width 1", () => {
  assertEquals(charWidth("a"), 1);
  assertEquals(charWidth("Z"), 1);
  assertEquals(charWidth("0"), 1);
  assertEquals(charWidth(" "), 1);
  assertEquals(charWidth("@"), 1);
});

Deno.test("charWidth - CJK characters have width 2", () => {
  assertEquals(charWidth("日"), 2);
  assertEquals(charWidth("本"), 2);
  assertEquals(charWidth("語"), 2);
  assertEquals(charWidth("中"), 2);
  assertEquals(charWidth("国"), 2);
});

Deno.test("charWidth - Box drawing characters are ambiguous", () => {
  // With ambiguous width = 1 (half)
  assertEquals(charWidth("─", 1), 1);
  assertEquals(charWidth("│", 1), 1);
  assertEquals(charWidth("┌", 1), 1);
  assertEquals(charWidth("└", 1), 1);

  // With ambiguous width = 2 (full) - Box Drawing stays at 1
  assertEquals(charWidth("─", 2), 1);
  assertEquals(charWidth("│", 2), 1);
  assertEquals(charWidth("┌", 2), 1);
  assertEquals(charWidth("└", 2), 1);

  // With legacy mode - Box Drawing becomes 2
  assertEquals(charWidth("─", "legacy"), 2);
  assertEquals(charWidth("│", "legacy"), 2);
  assertEquals(charWidth("┌", "legacy"), 2);
  assertEquals(charWidth("└", "legacy"), 2);
});

Deno.test("charWidth - Arrows are ambiguous", () => {
  assertEquals(charWidth("→", 1), 1);
  assertEquals(charWidth("←", 1), 1);
  assertEquals(charWidth("↑", 1), 1);
  assertEquals(charWidth("↓", 1), 1);

  assertEquals(charWidth("→", 2), 2);
  assertEquals(charWidth("←", 2), 2);
});

Deno.test("isFullwidth - CJK characters", () => {
  assertEquals(isFullwidth("日".codePointAt(0)!), true);
  assertEquals(isFullwidth("あ".codePointAt(0)!), true);
  assertEquals(isFullwidth("ア".codePointAt(0)!), true);
  assertEquals(isFullwidth("가".codePointAt(0)!), true);
});

Deno.test("isFullwidth - ASCII characters", () => {
  assertEquals(isFullwidth("a".codePointAt(0)!), false);
  assertEquals(isFullwidth("A".codePointAt(0)!), false);
  assertEquals(isFullwidth("1".codePointAt(0)!), false);
});

Deno.test("isAmbiguous - Box drawing characters", () => {
  assertEquals(isAmbiguous("─".codePointAt(0)!), true);
  assertEquals(isAmbiguous("│".codePointAt(0)!), true);
  assertEquals(isAmbiguous("┌".codePointAt(0)!), true);
  assertEquals(isAmbiguous("┐".codePointAt(0)!), true);
  assertEquals(isAmbiguous("└".codePointAt(0)!), true);
  assertEquals(isAmbiguous("┘".codePointAt(0)!), true);
});

Deno.test("stringWidth - ASCII string", () => {
  assertEquals(stringWidth("hello"), 5);
  assertEquals(stringWidth("world"), 5);
  assertEquals(stringWidth("hello world"), 11);
});

Deno.test("stringWidth - CJK string", () => {
  assertEquals(stringWidth("日本語"), 6);
  assertEquals(stringWidth("中国語"), 6);
});

Deno.test("stringWidth - Mixed string", () => {
  assertEquals(stringWidth("Hello日本"), 9); // 5 + 2 + 2
  assertEquals(stringWidth("abc日def"), 8); // 3 + 2 + 3
});

Deno.test("stringWidth - Box drawing with ambiguous width modes", () => {
  // half/console: Box Drawing = 1
  assertEquals(stringWidth("┌──┐", 1), 4);
  assertEquals(stringWidth("┌──┐", "half"), 4);
  assertEquals(stringWidth("┌──┐", "console"), 4);

  // full: Box Drawing = 1 (special case)
  assertEquals(stringWidth("┌──┐", 2), 4);
  assertEquals(stringWidth("┌──┐", "full"), 4);

  // legacy: Box Drawing = 2
  assertEquals(stringWidth("┌──┐", "legacy"), 8);
});

Deno.test("padString - Pads ASCII string", () => {
  assertEquals(padString("abc", 6), "abc   ");
  assertEquals(padString("hello", 10), "hello     ");
});

Deno.test("padString - Does not truncate", () => {
  assertEquals(padString("hello", 3), "hello");
});

Deno.test("centerString - Centers ASCII string", () => {
  assertEquals(centerString("a", 5), "  a  ");
  assertEquals(centerString("ab", 6), "  ab  ");
});

Deno.test("centerString - Centers CJK string", () => {
  assertEquals(centerString("日", 6), "  日  ");
});

// New tests for extended ambiguous width modes
Deno.test("isBoxDrawing - identifies Box Drawing characters", () => {
  assertEquals(isBoxDrawing("─".codePointAt(0)!), true); // U+2500
  assertEquals(isBoxDrawing("│".codePointAt(0)!), true); // U+2502
  assertEquals(isBoxDrawing("┌".codePointAt(0)!), true); // U+250C
  assertEquals(isBoxDrawing("┐".codePointAt(0)!), true); // U+2510
  assertEquals(isBoxDrawing("└".codePointAt(0)!), true); // U+2514
  assertEquals(isBoxDrawing("┘".codePointAt(0)!), true); // U+2518
  assertEquals(isBoxDrawing("├".codePointAt(0)!), true); // U+251C
  assertEquals(isBoxDrawing("┤".codePointAt(0)!), true); // U+2524
  assertEquals(isBoxDrawing("┬".codePointAt(0)!), true); // U+252C
  assertEquals(isBoxDrawing("┴".codePointAt(0)!), true); // U+2534
  assertEquals(isBoxDrawing("┼".codePointAt(0)!), true); // U+253C

  // Not Box Drawing
  assertEquals(isBoxDrawing("a".codePointAt(0)!), false);
  assertEquals(isBoxDrawing("日".codePointAt(0)!), false);
  assertEquals(isBoxDrawing("→".codePointAt(0)!), false); // Arrow is ambiguous but not box drawing
});

Deno.test("parseAmbiguousWidthMode - numeric modes", () => {
  assertEquals(parseAmbiguousWidthMode(1), {
    ambiguousWidth: 1,
    boxDrawingWidth: 1,
  });
  assertEquals(parseAmbiguousWidthMode(2), {
    ambiguousWidth: 2,
    boxDrawingWidth: 1,
  });
});

Deno.test("parseAmbiguousWidthMode - string modes", () => {
  assertEquals(parseAmbiguousWidthMode("half"), {
    ambiguousWidth: 1,
    boxDrawingWidth: 1,
  });
  assertEquals(parseAmbiguousWidthMode("full"), {
    ambiguousWidth: 2,
    boxDrawingWidth: 1,
  });
  assertEquals(parseAmbiguousWidthMode("console"), {
    ambiguousWidth: 1,
    boxDrawingWidth: 1,
  });
  assertEquals(parseAmbiguousWidthMode("legacy"), {
    ambiguousWidth: 2,
    boxDrawingWidth: 2,
  });
});

Deno.test("charWidth - non-box-drawing ambiguous chars with full mode", () => {
  // Arrows should be width 2 with "full" mode
  assertEquals(charWidth("→", "full"), 2);
  assertEquals(charWidth("←", "full"), 2);
  assertEquals(charWidth("↑", "full"), 2);
  assertEquals(charWidth("↓", "full"), 2);

  // But Box Drawing should stay at 1
  assertEquals(charWidth("─", "full"), 1);
  assertEquals(charWidth("│", "full"), 1);
});

Deno.test("stringWidth - mixed content with full mode", () => {
  // Box drawing (4 chars, each width 1) + CJK (1 char, width 2)
  assertEquals(stringWidth("┌──┐日", "full"), 6);
  // Box drawing (4 chars, each width 2 in legacy) + CJK (1 char, width 2)
  assertEquals(stringWidth("┌──┐日", "legacy"), 10);
});
