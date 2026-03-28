import { extractJsonObject } from '../util.js';

export async function generateWithOpenAI({ config, messages }) {
  if (!config.openaiApiKey) {
    throw new Error('OPENAI_API_KEY or PLS_OPENAI_API_KEY is required for provider=openai');
  }

  const response = await fetch(`${config.host.replace(/\/$/, '')}/v1/chat/completions`, {
    method: 'POST',
    headers: {
      'content-type': 'application/json',
      authorization: `Bearer ${config.openaiApiKey}`,
    },
    body: JSON.stringify({
      model: config.model,
      temperature: 0.1,
      response_format: { type: 'json_object' },
      messages: [
        { role: 'system', content: messages.system },
        { role: 'user', content: JSON.stringify(messages.user) },
      ],
    }),
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(`OpenAI request failed (${response.status}): ${text}`);
  }

  const payload = await response.json();
  const content = payload?.choices?.[0]?.message?.content;
  if (typeof content !== 'string' || content.trim() === '') {
    throw new Error('OpenAI response did not contain message content');
  }

  return extractJsonObject(content);
}
