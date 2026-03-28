const VALID_RISKS = new Set(['low', 'medium', 'high', 'critical']);

export function validateSuggestion(payload) {
  if (!payload || typeof payload !== 'object' || Array.isArray(payload)) {
    throw new Error('model response was not a JSON object');
  }

  const normalized = {
    command: requiredString(payload.command, 'command'),
    explanation: requiredString(payload.explanation, 'explanation'),
    risk: requiredString(payload.risk, 'risk').toLowerCase(),
    requiresConfirmation: Boolean(payload.requiresConfirmation),
    needsClarification: Boolean(payload.needsClarification),
    clarificationQuestion: optionalString(payload.clarificationQuestion),
    notes: optionalString(payload.notes),
    platform: optionalString(payload.platform),
    refused: Boolean(payload.refused),
  };

  if (!VALID_RISKS.has(normalized.risk)) {
    throw new Error(`invalid risk level: ${normalized.risk}`);
  }

  if (normalized.needsClarification && !normalized.clarificationQuestion) {
    throw new Error('clarificationQuestion is required when needsClarification is true');
  }

  return normalized;
}

function requiredString(value, field) {
  if (typeof value !== 'string' || value.trim() === '') {
    throw new Error(`${field} must be a non-empty string`);
  }
  return value.trim();
}

function optionalString(value) {
  if (value == null) return '';
  if (typeof value !== 'string') {
    throw new Error('optional string field had unexpected type');
  }
  return value.trim();
}
