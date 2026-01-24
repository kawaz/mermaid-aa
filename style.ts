/**
 * Character set definitions for rendering
 */

import type { Charset } from "./types.ts";

/** Character set for box drawing */
export interface CharacterSet {
  // Box corners
  topLeft: string;
  topRight: string;
  bottomLeft: string;
  bottomRight: string;

  // Box lines
  horizontal: string;
  vertical: string;

  // Arrows
  arrowRight: string;
  arrowLeft: string;
  arrowUp: string;
  arrowDown: string;

  // Diamond corners
  diamondLeft: string;
  diamondRight: string;

  // Rounded corners
  roundedTopLeft: string;
  roundedTopRight: string;
  roundedBottomLeft: string;
  roundedBottomRight: string;

  // Hexagon
  hexLeft: string;
  hexRight: string;

  // Circle
  circleLeft: string;
  circleRight: string;

  // Stadium
  stadiumLeft: string;
  stadiumRight: string;

  // Parallelogram / Trapezoid
  slashForward: string;
  slashBack: string;

  // Edge styles
  dotted: string;
  thick: string;
}

const ASCII: CharacterSet = {
  topLeft: "+",
  topRight: "+",
  bottomLeft: "+",
  bottomRight: "+",
  horizontal: "-",
  vertical: "|",
  arrowRight: ">",
  arrowLeft: "<",
  arrowUp: "^",
  arrowDown: "v",
  diamondLeft: "<",
  diamondRight: ">",
  roundedTopLeft: "(",
  roundedTopRight: ")",
  roundedBottomLeft: "(",
  roundedBottomRight: ")",
  hexLeft: "{",
  hexRight: "}",
  circleLeft: "(",
  circleRight: ")",
  stadiumLeft: "(",
  stadiumRight: ")",
  slashForward: "/",
  slashBack: "\\",
  dotted: "-",
  thick: "=",
};

const UNICODE: CharacterSet = {
  topLeft: "┌",
  topRight: "┐",
  bottomLeft: "└",
  bottomRight: "┘",
  horizontal: "─",
  vertical: "│",
  arrowRight: "→",
  arrowLeft: "←",
  arrowUp: "↑",
  arrowDown: "↓",
  diamondLeft: "◇",
  diamondRight: "◇",
  roundedTopLeft: "(",
  roundedTopRight: ")",
  roundedBottomLeft: "(",
  roundedBottomRight: ")",
  hexLeft: "⬡",
  hexRight: "⬡",
  circleLeft: "(",
  circleRight: ")",
  stadiumLeft: "(",
  stadiumRight: ")",
  slashForward: "/",
  slashBack: "\\",
  dotted: "╌",
  thick: "━",
};

const UNICODE_ROUND: CharacterSet = {
  ...UNICODE,
  topLeft: "╭",
  topRight: "╮",
  bottomLeft: "╰",
  bottomRight: "╯",
};

const UNICODE_BOLD: CharacterSet = {
  topLeft: "┏",
  topRight: "┓",
  bottomLeft: "┗",
  bottomRight: "┛",
  horizontal: "━",
  vertical: "┃",
  arrowRight: "➤",
  arrowLeft: "◄",
  arrowUp: "▲",
  arrowDown: "▼",
  diamondLeft: "◆",
  diamondRight: "◆",
  roundedTopLeft: "(",
  roundedTopRight: ")",
  roundedBottomLeft: "(",
  roundedBottomRight: ")",
  hexLeft: "⬢",
  hexRight: "⬢",
  circleLeft: "(",
  circleRight: ")",
  stadiumLeft: "(",
  stadiumRight: ")",
  slashForward: "/",
  slashBack: "\\",
  dotted: "┄",
  thick: "━",
};

const UNICODE_DOUBLE: CharacterSet = {
  topLeft: "╔",
  topRight: "╗",
  bottomLeft: "╚",
  bottomRight: "╝",
  horizontal: "═",
  vertical: "║",
  arrowRight: "⇒",
  arrowLeft: "⇐",
  arrowUp: "⇑",
  arrowDown: "⇓",
  diamondLeft: "◇",
  diamondRight: "◇",
  roundedTopLeft: "(",
  roundedTopRight: ")",
  roundedBottomLeft: "(",
  roundedBottomRight: ")",
  hexLeft: "⬡",
  hexRight: "⬡",
  circleLeft: "(",
  circleRight: ")",
  stadiumLeft: "(",
  stadiumRight: ")",
  slashForward: "/",
  slashBack: "\\",
  dotted: "╌",
  thick: "═",
};

/** Get the character set for a given charset name */
export function getCharset(charset: Charset): CharacterSet {
  switch (charset) {
    case "ascii":
      return ASCII;
    case "unicode":
      return UNICODE;
    case "unicode-round":
      return UNICODE_ROUND;
    case "unicode-bold":
      return UNICODE_BOLD;
    case "unicode-double":
      return UNICODE_DOUBLE;
    default:
      return UNICODE;
  }
}
