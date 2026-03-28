const HIGH_RISK_PATTERNS = [
  /\bsudo\b/i,
  /\brm\b\s+(-[a-zA-Z]*[rf][a-zA-Z]*|--recursive|--force)/i,
  /\bdd\b/i,
  /\bmkfs\b/i,
  /\bshutdown\b/i,
  /\breboot\b/i,
  /\bpoweroff\b/i,
  /curl\s+[^|\n]+\|\s*(sh|bash|zsh)/i,
  /wget\s+[^|\n]+\|\s*(sh|bash|zsh)/i,
  />\s*\/dev\//i,
];

export function applySafetyPolicy(suggestion) {
  const command = suggestion.command;
  const hasHighRiskPattern = HIGH_RISK_PATTERNS.some((pattern) => pattern.test(command));

  if (!hasHighRiskPattern) {
    return suggestion;
  }

  const escalatedRisk = riskRank(suggestion.risk) < riskRank('high') ? 'high' : suggestion.risk;
  return {
    ...suggestion,
    risk: escalatedRisk,
    requiresConfirmation: true,
    notes: joinNotes(
      suggestion.notes,
      'Safety policy flagged this command for manual review before any future execution support.'
    ),
  };
}

function joinNotes(existing, extra) {
  return [existing, extra].filter(Boolean).join(' ');
}

function riskRank(risk) {
  switch (risk) {
    case 'low':
      return 1;
    case 'medium':
      return 2;
    case 'high':
      return 3;
    case 'critical':
      return 4;
    default:
      return 0;
  }
}
