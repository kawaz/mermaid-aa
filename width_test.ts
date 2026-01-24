/**
 * Tests for width calculation
 */

import { assertEquals } from "@std/assert";
import {
  centerString,
  charWidth,
  isAmbiguous,
  isFullwidth,
  padString,
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
  // With ambiguous width = 1
  assertEquals(charWidth("─", 1), 1);
  assertEquals(charWidth("│", 1), 1);
  assertEquals(charWidth("┌", 1), 1);
  assertEquals(charWidth("└", 1), 1);

  // With ambiguous width = 2
  assertEquals(charWidth("─", 2), 2);
  assertEquals(charWidth("│", 2), 2);
  assertEquals(charWidth("┌", 2), 2);
  assertEquals(charWidth("└", 2), 2);
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

Deno.test("stringWidth - Box drawing with ambiguous width", () => {
  assertEquals(stringWidth("┌──┐", 1), 4);
  assertEquals(stringWidth("┌──┐", 2), 8);
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
