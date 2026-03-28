export function formatHumanOutput(result) {
  if (result.needsClarification) {
    return [
      'Need clarification:',
      `  ${result.clarificationQuestion}`,
    ].join('\n');
  }

  if (result.refused) {
    return [
      'Refused:',
      `  ${result.explanation}`,
      result.notes ? `\nNotes:\n  ${result.notes}` : '',
    ].filter(Boolean).join('\n');
  }

  return [
    'Command:',
    `  ${result.command}`,
    '',
    'Why:',
    `  ${result.explanation}`,
    '',
    'Risk:',
    `  ${result.risk}`,
    result.platform ? `\nPlatform:\n  ${result.platform}` : '',
    result.notes ? `\nNotes:\n  ${result.notes}` : '',
    result.requiresConfirmation ? '\nExecution:\n  This command should require confirmation in a future execution-enabled version.' : '',
  ].filter(Boolean).join('\n');
}
