import { buildMessages } from './prompt.js';
import { applySafetyPolicy } from './policy.js';
import { generateWithProvider } from './providers/index.js';
import { validateSuggestion } from './schema.js';

export async function generateSuggestion({ request, runtimeContext, config }) {
  const messages = buildMessages({ request, runtimeContext });
  const raw = await generateWithProvider({ config, messages });
  const validated = validateSuggestion(raw);
  return applySafetyPolicy(validated);
}
