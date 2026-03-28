export function extractJsonObject(text) {
  const trimmed = text.trim();
  try {
    return JSON.parse(trimmed);
  } catch {
    const start = trimmed.indexOf('{');
    const end = trimmed.lastIndexOf('}');
    if (start === -1 || end === -1 || end <= start) {
      throw new Error('response did not contain a JSON object');
    }
    return JSON.parse(trimmed.slice(start, end + 1));
  }
}
