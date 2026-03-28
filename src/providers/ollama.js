export async function generateWithOllama({ config, messages }) {
  const response = await fetch(`${config.host.replace(/\/$/, '')}/api/chat`, {
    method: 'POST',
    headers: {
      'content-type': 'application/json',
    },
    body: JSON.stringify({
      model: config.model,
      format: 'json',
      stream: false,
      options: {
        temperature: 0.1,
      },
      messages: [
        { role: 'system', content: messages.system },
        { role: 'user', content: JSON.stringify(messages.user) },
      ],
    }),
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(`Ollama request failed (${response.status}): ${text}`);
  }

  const payload = await response.json();
  const content = payload?.message?.content;
  if (typeof content !== 'string' || content.trim() === '') {
    throw new Error('Ollama response did not contain message content');
  }

  return JSON.parse(content);
}
