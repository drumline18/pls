#!/usr/bin/env node
import { main } from '../src/cli.js';

main().catch((error) => {
  const message = error instanceof Error ? error.message : String(error);
  console.error(`pls: ${message}`);
  process.exit(1);
});
