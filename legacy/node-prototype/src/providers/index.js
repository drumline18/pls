import { generateWithOpenAI } from './openai.js';
import { generateWithOllama } from './ollama.js';

export async function generateWithProvider({ config, messages }) {
  switch (config.provider) {
    case 'openai':
      return generateWithOpenAI({ config, messages });
    case 'ollama':
      return generateWithOllama({ config, messages });
    default:
      throw new Error(`unsupported provider: ${config.provider}`);
  }
}
