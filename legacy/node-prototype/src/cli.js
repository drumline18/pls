import process from 'node:process';
import { getRuntimeContext } from './context.js';
import { loadConfig, parseArgs } from './config.js';
import { generateSuggestion } from './generate.js';
import { formatHumanOutput } from './render.js';

const HELP_TEXT = `pls — natural-language shell command suggester

Usage:
  pls <request>
  pls --provider openai --model gpt-4.1-mini show hidden files here
  pls --provider ollama --model qwen2.5-coder:7b list large files in this directory
  pls --json find files bigger than 500mb

Flags:
  --provider <openai|ollama>   LLM provider
  --model <name>               Model name to use
  --json                       Emit JSON only
  --shell <name>               Override shell detection
  --os <name>                  Override OS detection
  --host <url>                 Override provider host/base URL
  --help, -h                   Show help

Environment:
  OPENAI_API_KEY               OpenAI API key
  PLS_OPENAI_API_KEY           Alternate OpenAI API key env var
  OLLAMA_HOST                  Ollama base URL (default: http://127.0.0.1:11434)
  PLS_OLLAMA_HOST              Alternate Ollama host env var
  PLS_PROVIDER                 Default provider when flag is omitted
  PLS_MODEL                    Default model when flag is omitted
`;

export async function main(argv = process.argv.slice(2)) {
  const parsed = parseArgs(argv);

  if (parsed.help || parsed.requestParts.length === 0) {
    console.log(HELP_TEXT);
    if (parsed.requestParts.length === 0) {
      process.exitCode = 1;
    }
    return;
  }

  const request = parsed.requestParts.join(' ').trim();
  const runtimeContext = getRuntimeContext({
    shellOverride: parsed.flags.shell,
    osOverride: parsed.flags.os,
  });
  const config = loadConfig(parsed.flags);
  const result = await generateSuggestion({
    request,
    runtimeContext,
    config,
  });

  if (config.outputJson) {
    console.log(JSON.stringify(result, null, 2));
    return;
  }

  console.log(formatHumanOutput(result));
}
