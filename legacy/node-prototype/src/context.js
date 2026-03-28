import os from 'node:os';
import process from 'node:process';

export function getRuntimeContext({ shellOverride, osOverride } = {}) {
  const detectedOs = osOverride || normalizeOs(process.platform);
  const shellPath = shellOverride || process.env.SHELL || process.env.ComSpec || 'unknown';
  const shell = normalizeShell(shellPath);

  return {
    cwd: process.cwd(),
    os: detectedOs,
    shell,
    homeDirectory: os.homedir(),
    isTTY: Boolean(process.stdout.isTTY),
  };
}

function normalizeOs(platform) {
  switch (platform) {
    case 'darwin':
      return 'macos';
    case 'linux':
      return 'linux';
    case 'win32':
      return 'windows';
    default:
      return platform;
  }
}

function normalizeShell(shellPath) {
  const lower = shellPath.toLowerCase();
  if (lower.includes('fish')) return 'fish';
  if (lower.includes('zsh')) return 'zsh';
  if (lower.includes('bash')) return 'bash';
  if (lower.includes('powershell') || lower.includes('pwsh')) return 'powershell';
  if (lower.includes('cmd.exe')) return 'cmd';
  return shellPath;
}
