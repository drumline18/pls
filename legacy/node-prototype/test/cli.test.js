import test from 'node:test';
import assert from 'node:assert/strict';
import { parseArgs, loadConfig } from '../src/config.js';
import { validateSuggestion } from '../src/schema.js';

test('parseArgs captures freeform request without quotes', () => {
  const parsed = parseArgs(['show', 'me', 'dotfiles', 'here']);
  assert.deepEqual(parsed.requestParts, ['show', 'me', 'dotfiles', 'here']);
});

test('parseArgs supports known flags and freeform tail', () => {
  const parsed = parseArgs(['--provider', 'ollama', '--model', 'qwen2.5-coder:7b', 'show', 'hidden', 'files']);
  assert.equal(parsed.flags.provider, 'ollama');
  assert.equal(parsed.flags.model, 'qwen2.5-coder:7b');
  assert.deepEqual(parsed.requestParts, ['show', 'hidden', 'files']);
});

test('loadConfig defaults to ollama when no OpenAI key exists', () => {
  const originalOpenAI = process.env.OPENAI_API_KEY;
  const originalAlt = process.env.PLS_OPENAI_API_KEY;
  delete process.env.OPENAI_API_KEY;
  delete process.env.PLS_OPENAI_API_KEY;

  const config = loadConfig({});
  assert.equal(config.provider, 'ollama');
  assert.equal(config.model, 'qwen2.5-coder:7b-instruct-q4_K_M');

  process.env.OPENAI_API_KEY = originalOpenAI;
  process.env.PLS_OPENAI_API_KEY = originalAlt;
});

test('validateSuggestion accepts schema-compliant payloads', () => {
  const result = validateSuggestion({
    command: 'ls -la',
    explanation: 'Lists files including dotfiles.',
    risk: 'low',
    requiresConfirmation: false,
    needsClarification: false,
    clarificationQuestion: '',
    notes: '',
    platform: 'linux',
    refused: false,
  });

  assert.equal(result.command, 'ls -la');
  assert.equal(result.risk, 'low');
});
