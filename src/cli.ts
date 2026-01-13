#!/usr/bin/env node

import * as fs from 'fs';
import * as readline from 'readline';
import { parse } from './parser';
import { render } from './renderer';
import { CharsetName } from './types';

const VALID_CHARSETS: CharsetName[] = ['ascii', 'unicode', 'unicode-round', 'unicode-bold', 'unicode-double'];

async function readStdin(): Promise<string> {
  return new Promise((resolve) => {
    let input = '';
    const rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout,
      terminal: false,
    });

    rl.on('line', (line) => {
      input += line + '\n';
    });

    rl.on('close', () => {
      resolve(input);
    });
  });
}

function printHelp(): void {
  console.log(`
mermaid-aa - Render Mermaid diagrams as ASCII art

Usage:
  mermaid-aa <mermaid-text>     Render from argument
  mermaid-aa -f <file>          Render from file
  echo "..." | mermaid-aa       Render from stdin (pipe)
  mermaid-aa                    Interactive mode (stdin)

Options:
  -h, --help           Show this help message
  -f, --file <path>    Read mermaid text from file
  -c, --charset <name> Character set for rendering (default: ascii)
  -v, --version        Show version

Available charsets:
  ascii          ASCII characters (+, -, |, >, <, v, ^)
  unicode        Unicode box drawing (┌, ─, │, →, ↓)
  unicode-round  Unicode with rounded corners (╭, ╮, ╰, ╯)
  unicode-bold   Unicode bold lines (┏, ━, ┃, ▶, ▼)
  unicode-double Unicode double lines (╔, ═, ║, ▷, ▽)

Supported diagrams:
  - flowchart / graph (TB, TD, LR, RL, BT)
  - sequenceDiagram

Examples:
  mermaid-aa "graph LR; A-->B-->C"
  mermaid-aa -c unicode "graph TD; A-->B"
  echo "graph TD; A-->B" | mermaid-aa --charset unicode-round
  mermaid-aa -f diagram.mmd -c unicode-bold
`);
}

function printVersion(): void {
  console.log('mermaid-aa v1.0.0');
}

function getArgValue(args: string[], flags: string[]): string | undefined {
  for (const flag of flags) {
    const index = args.indexOf(flag);
    if (index !== -1 && args[index + 1]) {
      return args[index + 1];
    }
  }
  return undefined;
}

function removeArgWithValue(args: string[], flags: string[]): string[] {
  const result = [...args];
  for (const flag of flags) {
    const index = result.indexOf(flag);
    if (index !== -1) {
      result.splice(index, 2);
    }
  }
  return result;
}

async function main(): Promise<void> {
  const args = process.argv.slice(2);

  // Handle flags
  if (args.includes('-h') || args.includes('--help')) {
    printHelp();
    return;
  }

  if (args.includes('-v') || args.includes('--version')) {
    printVersion();
    return;
  }

  // Parse charset option
  const charsetArg = getArgValue(args, ['-c', '--charset']);
  let charset: CharsetName = 'ascii';
  if (charsetArg) {
    if (!VALID_CHARSETS.includes(charsetArg as CharsetName)) {
      console.error(`Invalid charset: ${charsetArg}`);
      console.error(`Valid charsets: ${VALID_CHARSETS.join(', ')}`);
      process.exit(1);
    }
    charset = charsetArg as CharsetName;
  }

  let input: string;

  // Check for file flag
  const filePath = getArgValue(args, ['-f', '--file']);
  if (filePath) {
    try {
      input = fs.readFileSync(filePath, 'utf-8');
    } catch (err) {
      console.error(`Error reading file: ${filePath}`);
      process.exit(1);
    }
  } else {
    // Remove option flags to get remaining arguments
    let remainingArgs = removeArgWithValue(args, ['-f', '--file', '-c', '--charset']);

    if (remainingArgs.length > 0 && !remainingArgs[0].startsWith('-')) {
      // Read from argument
      input = remainingArgs.join(' ');
    } else if (!process.stdin.isTTY) {
      // Read from pipe/stdin
      input = await readStdin();
    } else {
      // Interactive mode - show help
      printHelp();
      return;
    }
  }

  if (!input.trim()) {
    console.error('No input provided');
    process.exit(1);
  }

  try {
    const diagram = parse(input);
    const output = render(diagram, { charset });
    console.log(output);
  } catch (err) {
    console.error('Error parsing diagram:', err);
    process.exit(1);
  }
}

main();
