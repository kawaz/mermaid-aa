/**
 * Drawing canvas for ASCII art rendering
 */

import { type AmbiguousWidthMode, stringWidth } from "./width.ts";

/** 2D canvas for drawing ASCII art */
export class Canvas {
  private grid: string[][];
  public readonly width: number;
  public readonly height: number;
  private ambiguousWidth: AmbiguousWidthMode;

  constructor(
    width: number,
    height: number,
    ambiguousWidth: AmbiguousWidthMode = 1,
  ) {
    this.width = width;
    this.height = height;
    this.ambiguousWidth = ambiguousWidth;
    this.grid = [];

    for (let y = 0; y < height; y++) {
      this.grid.push(new Array(width).fill(" "));
    }
  }

  /** Set a character at the specified position */
  set(x: number, y: number, char: string): void {
    if (x >= 0 && x < this.width && y >= 0 && y < this.height) {
      this.grid[y][x] = char;
    }
  }

  /** Get a character at the specified position */
  get(x: number, y: number): string {
    if (x >= 0 && x < this.width && y >= 0 && y < this.height) {
      return this.grid[y][x];
    }
    return " ";
  }

  /** Draw a string at the specified position */
  drawString(x: number, y: number, str: string): void {
    if (y < 0 || y >= this.height) return;

    let currentX = x;
    for (const char of str) {
      if (currentX >= 0 && currentX < this.width) {
        this.grid[y][currentX] = char;
        // Handle wide characters
        const charW = stringWidth(char, this.ambiguousWidth);
        if (charW === 2 && currentX + 1 < this.width) {
          // Mark next cell as continuation (empty for wide char)
          this.grid[y][currentX + 1] = "";
        }
        currentX += charW;
      } else {
        currentX++;
      }
    }
  }

  /** Draw a horizontal line */
  drawHorizontalLine(x: number, y: number, length: number, char: string): void {
    for (let i = 0; i < length; i++) {
      this.set(x + i, y, char);
    }
  }

  /** Draw a vertical line */
  drawVerticalLine(x: number, y: number, length: number, char: string): void {
    for (let i = 0; i < length; i++) {
      this.set(x, y + i, char);
    }
  }

  /** Draw a box */
  drawBox(
    x: number,
    y: number,
    width: number,
    height: number,
    topLeft: string,
    topRight: string,
    bottomLeft: string,
    bottomRight: string,
    horizontal: string,
    vertical: string,
  ): void {
    // Corners
    this.set(x, y, topLeft);
    this.set(x + width - 1, y, topRight);
    this.set(x, y + height - 1, bottomLeft);
    this.set(x + width - 1, y + height - 1, bottomRight);

    // Horizontal lines
    for (let i = 1; i < width - 1; i++) {
      this.set(x + i, y, horizontal);
      this.set(x + i, y + height - 1, horizontal);
    }

    // Vertical lines
    for (let i = 1; i < height - 1; i++) {
      this.set(x, y + i, vertical);
      this.set(x + width - 1, y + i, vertical);
    }
  }

  /** Render the canvas to a string */
  render(): string {
    const lines: string[] = [];

    for (let y = 0; y < this.height; y++) {
      let line = "";
      for (let x = 0; x < this.width; x++) {
        const char = this.grid[y][x];
        // Skip empty cells (continuation of wide chars)
        if (char !== "") {
          line += char;
        }
      }
      // Trim trailing spaces
      lines.push(line.trimEnd());
    }

    // Remove trailing empty lines
    while (lines.length > 0 && lines[lines.length - 1] === "") {
      lines.pop();
    }

    return lines.join("\n");
  }
}
