#!/usr/bin/env node

import * as fs from 'fs';
import * as readline from 'readline';
import { parse } from './parser';
import { render } from './renderer';

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
  -h, --help     Show this help message
  -f, --file     Read mermaid text from file
  -v, --version  Show version

Supported diagrams:
  - flowchart / graph (TB, TD, LR, RL, BT)
  - sequenceDiagram

Examples:
  mermaid-aa "graph LR; A-->B-->C"
  echo "graph TD; A-->B" | mermaid-aa
  mermaid-aa -f diagram.mmd
`);
}

function printVersion(): void {
  console.log('mermaid-aa v1.0.0');
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

  let input: string;

  // Check for file flag
  const fileIndex = args.findIndex((a) => a === '-f' || a === '--file');
  if (fileIndex !== -1 && args[fileIndex + 1]) {
    const filePath = args[fileIndex + 1];
    try {
      input = fs.readFileSync(filePath, 'utf-8');
    } catch (err) {
      console.error(`Error reading file: ${filePath}`);
      process.exit(1);
    }
  } else if (args.length > 0 && !args[0].startsWith('-')) {
    // Read from argument
    input = args.join(' ');
  } else if (!process.stdin.isTTY) {
    // Read from pipe/stdin
    input = await readStdin();
  } else {
    // Interactive mode - show help
    printHelp();
    return;
  }

  if (!input.trim()) {
    console.error('No input provided');
    process.exit(1);
  }

  try {
    const diagram = parse(input);
    const output = render(diagram);
    console.log(output);
  } catch (err) {
    console.error('Error parsing diagram:', err);
    process.exit(1);
  }
}

main();
