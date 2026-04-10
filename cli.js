#!/usr/bin/env node
// Node.js CLI wrapper for mermaid-aa WASM

const fs = require('fs');
const path = require('path');

// Parse command line arguments
function parseArgs() {
  const args = process.argv.slice(2);
  const options = {
    help: false,
    ambiguousWidth: 'half', // default
  };

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg === '-h' || arg === '--help') {
      options.help = true;
    } else if (arg === '-w' || arg === '--ambiguous-width') {
      options.ambiguousWidth = args[++i];
    } else if (arg.startsWith('--ambiguous-width=')) {
      options.ambiguousWidth = arg.split('=')[1];
    }
  }

  return options;
}

// Validate width mode
function validateWidthMode(mode) {
  const validModes = ['1', 'half', '2', 'full', 'console', 'legacy'];
  if (!validModes.includes(mode.toLowerCase())) {
    console.error(`Invalid ambiguous-width value: ${mode}`);
    console.error(`Valid values: ${validModes.join(', ')}`);
    console.error('');
    console.error('  1 | half    : Ambiguous=1, BoxDrawing=1 (Western terminals, default)');
    console.error('  2 | full    : Ambiguous=2, BoxDrawing=1 (CJK terminals, recommended)');
    console.error('  console     : Ambiguous=1, BoxDrawing=1 (same as half)');
    console.error('  legacy      : Ambiguous=2, BoxDrawing=2 (legacy compatibility)');
    process.exit(1);
  }
  return mode.toLowerCase();
}

function printHelp() {
  console.log('Usage: mermaid-aa [options]');
  console.log('');
  console.log('Options:');
  console.log('  -h, --help                  Show this help message');
  console.log('  -w, --ambiguous-width MODE  Set width mode for ambiguous characters');
  console.log('');
  console.log('Width Modes:');
  console.log('  1 | half    : Ambiguous=1, BoxDrawing=1 (Western terminals, default)');
  console.log('  2 | full    : Ambiguous=2, BoxDrawing=1 (CJK terminals, recommended)');
  console.log('  console     : Ambiguous=1, BoxDrawing=1 (same as half)');
  console.log('  legacy      : Ambiguous=2, BoxDrawing=2 (legacy compatibility)');
}

async function main() {
  const options = parseArgs();

  if (options.help) {
    printHelp();
    process.exit(0);
  }

  const widthMode = validateWidthMode(options.ambiguousWidth);

  const wasmPath = path.join(__dirname, '_build/wasm-gc/release/build/cmd/main/main.wasm');

  if (!fs.existsSync(wasmPath)) {
    console.error('WASM file not found. Run "moon build" first.');
    process.exit(1);
  }

  const wasmBuffer = fs.readFileSync(wasmPath);
  const wasmModule = await WebAssembly.compile(wasmBuffer);

  // Create memory and other imports
  const memory = new WebAssembly.Memory({ initial: 1, maximum: 100 });

  const imports = {
    spectest: {
      print_i32: (x) => console.log(x),
      print_char: (x) => process.stdout.write(String.fromCodePoint(x)),
    },
    env: {
      memory: memory,
      width_mode: widthMode, // Pass width mode to WASM
    },
  };

  try {
    const instance = await WebAssembly.instantiate(wasmModule, imports);

    // Call _start if it exists (WASI style)
    if (instance.exports._start) {
      instance.exports._start();
    } else if (instance.exports.main) {
      instance.exports.main();
    } else {
      console.log('WASM module loaded successfully.');
      console.log('Exports:', Object.keys(instance.exports));
      console.log('Width mode:', widthMode);
    }
  } catch (e) {
    console.error('Error running WASM:', e.message);
    process.exit(1);
  }
}

main();
