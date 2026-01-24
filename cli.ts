/**
 * CLI entry point for mermaid-aa
 */

import { parseArgs } from "@std/cli/parse-args";
import type { Charset, Direction } from "./types.ts";
import { parse } from "./parser.ts";
import { render } from "./renderer.ts";

const VERSION = "0.1.0";

const HELP = `mermaid-aa - Convert Mermaid diagrams to ASCII art

Usage:
  mermaid-aa [OPTIONS] [INPUT]

Arguments:
  [INPUT]  Mermaid text or file path

Options:
  -f, --file <FILE>           Read from file
  -c, --charset <CHARSET>     Character set [default: unicode]
                              [ascii|unicode|unicode-round|unicode-bold|unicode-double]
  -a, --ambiguous-width <N>   Width for ambiguous chars [default: 1]
                              [1: half-width, 2: full-width]
  -d, --direction <DIR>       Override direction [TB|TD|BT|LR|RL]
  -h, --help                  Print help
  -V, --version               Print version

Input methods (priority order):
  1. -f/--file option
  2. [INPUT] argument
  3. stdin (pipe or interactive)

Examples:
  mermaid-aa 'graph TD; A-->B'
  mermaid-aa -f diagram.mmd
  echo 'graph LR; A-->B' | mermaid-aa
  mermaid-aa -c ascii -a 2 'graph TD; A[日本語]-->B'

Environment Variables:
  MERMAID_AA_CHARSET          Default character set
  MERMAID_AA_AMBIGUOUS_WIDTH  Default ambiguous width (1 or 2)
`;

function isValidCharset(value: string): value is Charset {
  return ["ascii", "unicode", "unicode-round", "unicode-bold", "unicode-double"]
    .includes(value);
}

function isValidDirection(value: string): value is Direction {
  return ["TB", "TD", "BT", "LR", "RL"].includes(value.toUpperCase());
}

async function readStdin(): Promise<string> {
  const decoder = new TextDecoder();
  const chunks: Uint8Array[] = [];

  for await (const chunk of Deno.stdin.readable) {
    chunks.push(chunk);
  }

  const totalLength = chunks.reduce((sum, chunk) => sum + chunk.length, 0);
  const result = new Uint8Array(totalLength);
  let offset = 0;
  for (const chunk of chunks) {
    result.set(chunk, offset);
    offset += chunk.length;
  }

  return decoder.decode(result);
}

async function getInput(args: ReturnType<typeof parseArgs>): Promise<string> {
  // Priority 1: --file option
  if (args.file) {
    try {
      return await Deno.readTextFile(args.file as string);
    } catch (error) {
      if (error instanceof Deno.errors.NotFound) {
        console.error(`Error: File not found: ${args.file}`);
        Deno.exit(2);
      }
      throw error;
    }
  }

  // Priority 2: Positional argument
  if (args._.length > 0) {
    const input = String(args._[0]);
    // Check if it's a file path
    try {
      const stat = await Deno.stat(input);
      if (stat.isFile) {
        return await Deno.readTextFile(input);
      }
    } catch {
      // Not a file, treat as mermaid text
    }
    return input;
  }

  // Priority 3: stdin
  if (!Deno.stdin.isTerminal()) {
    return await readStdin();
  }

  // No input provided
  console.error("Error: No input provided");
  console.error("Use --help for usage information");
  Deno.exit(1);
}

export async function main(): Promise<void> {
  const args = parseArgs(Deno.args, {
    string: ["file", "charset", "direction", "ambiguous-width"],
    boolean: ["help", "version"],
    alias: {
      f: "file",
      c: "charset",
      d: "direction",
      a: "ambiguous-width",
      h: "help",
      V: "version",
    },
    default: {
      charset: Deno.env.get("MERMAID_AA_CHARSET") || "unicode",
      "ambiguous-width": Deno.env.get("MERMAID_AA_AMBIGUOUS_WIDTH") || "1",
    },
  });

  // Handle --help
  if (args.help) {
    console.log(HELP);
    Deno.exit(0);
  }

  // Handle --version
  if (args.version) {
    console.log(`mermaid-aa ${VERSION}`);
    Deno.exit(0);
  }

  // Validate charset
  const charset = args.charset as string;
  if (!isValidCharset(charset)) {
    console.error(`Error: Invalid charset: ${charset}`);
    console.error(
      "Valid values: ascii, unicode, unicode-round, unicode-bold, unicode-double",
    );
    Deno.exit(3);
  }

  // Validate ambiguous-width
  const ambiguousWidthStr = args["ambiguous-width"] as string;
  const ambiguousWidth = parseInt(ambiguousWidthStr, 10);
  if (ambiguousWidth !== 1 && ambiguousWidth !== 2) {
    console.error(`Error: Invalid ambiguous-width: ${ambiguousWidthStr}`);
    console.error("Valid values: 1 (half-width), 2 (full-width)");
    Deno.exit(3);
  }

  // Validate direction if provided
  let direction: Direction | undefined;
  if (args.direction) {
    const dirStr = args.direction as string;
    if (!isValidDirection(dirStr)) {
      console.error(`Error: Invalid direction: ${dirStr}`);
      console.error("Valid values: TB, TD, BT, LR, RL");
      Deno.exit(3);
    }
    direction = dirStr.toUpperCase() as Direction;
  }

  // Get input
  const input = await getInput(args);

  // Parse and render
  try {
    const flowchart = parse(input);

    // Override direction if specified
    if (direction) {
      flowchart.direction = direction;
    }

    const output = render(flowchart, charset, ambiguousWidth as 1 | 2);
    console.log(output);
  } catch (error) {
    if (error instanceof Error) {
      console.error(`Error: ${error.message}`);
    } else {
      console.error("Error: Unknown error occurred");
    }
    Deno.exit(1);
  }
}

// Run if this is the main module
if (import.meta.main) {
  await main();
}
