import process from 'node:process';

const FLAG_NAMES_WITH_VALUES = new Set(['--provider', '--model', '--shell', '--os', '--host']);

export function parseArgs(argv) {
  const flags = {};
  const requestParts = [];

  for (let index = 0; index < argv.length; index += 1) {
    const arg = argv[index];

    if (arg === '--help' || arg === '-h') {
      flags.help = true;
      continue;
    }

    if (arg === '--json') {
      flags.json = true;
      continue;
    }

    if (FLAG_NAMES_WITH_VALUES.has(arg)) {
      const value = argv[index + 1];
      if (!value || value.startsWith('--')) {
        throw new Error(`missing value for ${arg}`);
      }
      flags[arg.slice(2)] = value;
      index += 1;
      continue;
    }

    requestParts.push(arg);
  }

  return {
    help: Boolean(flags.help),
    flags,
    requestParts,
  };
}

export function loadConfig(flags) {
  const provider = flags.provider || process.env.PLS_PROVIDER || detectDefaultProvider();
  const model = flags.model || process.env.PLS_MODEL || defaultModelFor(provider);
  const host = flags.host || defaultHostFor(provider);

  return {
    provider,
    model,
    host,
    outputJson: Boolean(flags.json),
    openaiApiKey: process.env.PLS_OPENAI_API_KEY || process.env.OPENAI_API_KEY || '',
  };
}

function detectDefaultProvider() {
  if (process.env.PLS_OPENAI_API_KEY || process.env.OPENAI_API_KEY) {
    return 'openai';
  }
  return 'ollama';
}

function defaultModelFor(provider) {
  if (provider === 'openai') {
    return 'gpt-4.1-mini';
  }
  if (provider === 'ollama') {
    return 'qwen2.5-coder:7b-instruct-q4_K_M';
  }
  throw new Error(`unsupported provider: ${provider}`);
}

function defaultHostFor(provider) {
  if (provider === 'openai') {
    return 'https://api.openai.com';
  }
  if (provider === 'ollama') {
    return process.env.PLS_OLLAMA_HOST || process.env.OLLAMA_HOST || 'http://127.0.0.1:11434';
  }
  throw new Error(`unsupported provider: ${provider}`);
}
