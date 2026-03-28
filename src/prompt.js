export function buildMessages({ request, runtimeContext }) {
  const system = `You generate safe shell command suggestions.
Return JSON only. Do not wrap in markdown fences.
Prefer simple, readable, boring commands over clever shell tricks.
Do not assume GNU-only flags unless the supplied OS strongly suggests Linux.
If the request is ambiguous, set needsClarification to true and ask exactly one short question.
If the request is dangerous, still respond with JSON but set refused=true when you should not provide a direct command.

JSON schema:
{
  "command": "string",
  "explanation": "string",
  "risk": "low|medium|high|critical",
  "requiresConfirmation": true,
  "needsClarification": false,
  "clarificationQuestion": "string",
  "notes": "string",
  "platform": "string",
  "refused": false
}`;

  const user = {
    request,
    runtimeContext,
    instructions: [
      'Target a single primary command where possible.',
      'Assume the command is for the current working directory unless the user says otherwise.',
      'Prefer inspection commands for this MVP. Do not add execution wrappers, aliases, or shell functions.',
      'If multiple platform variants matter, put the main one in command and mention differences in notes.',
    ],
  };

  return { system, user };
}
