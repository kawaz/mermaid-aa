#!/usr/bin/env node
// Node.js CLI wrapper for mermaid-aa WASM

const fs = require('fs');
const path = require('path');

async function main() {
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
    }
  } catch (e) {
    console.error('Error running WASM:', e.message);
    process.exit(1);
  }
}

main();
