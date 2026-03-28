package prompt

import "pls/internal/types"

func Build(request string, runtimeContext types.RuntimeContext) types.Messages {
	system := `You generate safe shell command suggestions.
Return JSON only. Do not wrap in markdown fences.
Target the supplied operating system and shell exactly; do not emit Linux-only commands on macOS or PowerShell unless the user explicitly asked for them.
Prefer simple, readable, boring commands over clever shell tricks.
Prefer one direct command over pipelines when a direct command exists.
Never parse ls output with grep, awk, or sed when ls flags or find predicates can answer the request directly.
Never parse ps output with grep when a direct process or socket inspection command is more appropriate.
If the user already specified a unit, format, or scope, do not ask a clarification question about that same choice.
For batch file operations, it is acceptable to return a short shell loop when that is the clearest single command.
If the user explicitly named the destination shell config file or rc file and the alias text they want, do not ask a clarification question about that same alias setup.
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
}`

	user := map[string]any{
		"request":        request,
		"runtimeContext": runtimeContext,
		"examples":       examplesFor(runtimeContext),
		"instructions":   instructionsFor(runtimeContext),
	}

	return types.Messages{System: system, User: user}
}

func examplesFor(runtimeContext types.RuntimeContext) []map[string]any {
	switch runtimeContext.OS {
	case "macos":
		return macOSExamples(runtimeContext)
	case "windows":
		if runtimeContext.Shell == "powershell" {
			return powerShellExamples()
		}
		return windowsExamples()
	default:
		return linuxExamples(runtimeContext)
	}
}

func instructionsFor(runtimeContext types.RuntimeContext) []string {
	base := []string{
		"Target a single primary command where possible.",
		"Assume the command is for the current working directory unless the user says otherwise.",
		"Prefer direct predicates and flags over text-parsing pipelines.",
		"If multiple platform variants matter, put the main one in command and mention differences in notes.",
	}

	switch runtimeContext.OS {
	case "macos":
		return append(base,
			"Emit macOS-compatible commands and BSD-compatible flags.",
			"Prefer lsof for port/listener checks on macOS.",
			"Prefer launchctl, brew services, or pgrep for service/process checks on macOS depending on what the request implies.",
			"For shell config on macOS, prefer ~/.zshrc for zsh and ~/.bashrc or ~/.bash_profile for bash.",
			"For batch rename or move commands, use safe quoting and concise shell loops when needed.",
		)
	case "windows":
		if runtimeContext.Shell == "powershell" {
			return append(base,
				"Emit PowerShell-native commands and syntax.",
				"Prefer Get-ChildItem, Get-Service, Get-NetTCPConnection, Get-Process, and Rename-Item over Unix tools.",
				"For shell config, prefer the PowerShell profile in $PROFILE.",
				"Do not emit bash, grep, systemctl, or GNU find unless the user explicitly asked for WSL, MSYS2, or Git Bash.",
			)
		}
		return append(base,
			"Target native Windows commands for the detected shell.",
			"If PowerShell-specific features would help, mention PowerShell in notes rather than silently switching shells.",
		)
	default:
		return append(base,
			"Prefer inspection commands for this MVP. Do not add execution wrappers, aliases, or shell functions.",
			"If the user wants only files, only directories, or hidden entries, prefer direct find predicates or native flags instead of grepping ls output.",
			"For Linux service checks phrased like 'check if jellyfin is running', prefer systemctl is-active <service> when the target looks like a service name, and mention container/manual-process alternatives in notes if needed.",
			"For batch renames or metadata-based filename changes, a concise for-loop using purpose-built tools is acceptable.",
			"For batch rename or move commands, use safe quoting, guard globs with [ -e \"$f\" ] || continue, and prefer basename over brittle parameter expansion when prefixing filenames.",
		)
	}
}

func linuxExamples(runtimeContext types.RuntimeContext) []map[string]any {
	shellConfig := "~/.bashrc"
	shellName := "bashrc"
	if runtimeContext.Shell == "zsh" {
		shellConfig = "~/.zshrc"
		shellName = "zshrc"
	}

	return []map[string]any{
		example(
			"show me all dotfiles in this directory",
			"find . -maxdepth 1 -mindepth 1 -name '.*' -print",
			"Lists hidden files and directories in the current directory.",
			"low", false, false, "", "This directly lists dotfiles without parsing ls output.", "linux",
		),
		example(
			"check if jellyfin is running",
			"systemctl is-active jellyfin",
			"Checks whether the Jellyfin systemd service is active.",
			"low", false, false, "", "If Jellyfin was started with Docker or another supervisor, inspect that runtime instead.", "linux",
		),
		example(
			"prefix all the mp3s with their lengths in seconds",
			"for f in ./*.mp3; do [ -e \"$f\" ] || continue; secs=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 \"$f\" | awk '{printf \"%d\", $1}'); mv -- \"$f\" \"./${secs}s - $(basename \"$f\")\"; done",
			"Prefixes each MP3 filename with its duration in whole seconds.",
			"high", true, false, "", "This renames files in place and requires ffprobe from ffmpeg.", "linux",
		),
		example(
			"prefix all jpgs with vacation-",
			"for f in ./*.jpg; do [ -e \"$f\" ] || continue; base=$(basename \"$f\"); mv -- \"$f\" \"./vacation-$base\"; done",
			"Prefixes each JPG filename with 'vacation-'.",
			"high", true, false, "", "This renames files in place.", "linux",
		),
		example(
			"replace spaces in all filenames here with underscores",
			"for f in ./*\\ *; do [ -e \"$f\" ] || continue; mv -- \"$f\" \"${f// /_}\"; done",
			"Replaces spaces with underscores in filenames that contain spaces in the current directory.",
			"high", true, false, "", "This renames files in place.", "linux",
		),
		example(
			"move all srt files into a subtitles folder",
			"mkdir -p \"./subtitles\" && for f in ./*.srt; do [ -e \"$f\" ] || continue; mv -- \"$f\" \"./subtitles/\"; done",
			"Moves all .srt files from the current directory into the 'subtitles' folder.",
			"high", true, false, "", "This creates the destination folder if needed and moves files into it.", "linux",
		),
		example(
			"add cd.. as an alias to cd .. in "+shellName,
			"printf '\nalias cd..=\"cd ..\"\n' >> "+shellConfig,
			"Appends an alias named cd.. to "+shellConfig+" so it expands to 'cd ..'.",
			"high", true, false, "", "Run 'source "+shellConfig+"' or open a new shell afterward to use the alias in new sessions.", "linux",
		),
	}
}

func macOSExamples(runtimeContext types.RuntimeContext) []map[string]any {
	shellConfig := "~/.zshrc"
	shellName := "zshrc"
	if runtimeContext.Shell == "bash" {
		shellConfig = "~/.bashrc"
		shellName = "bashrc"
	}

	return []map[string]any{
		example(
			"check if jellyfin is running",
			"pgrep -fl jellyfin",
			"Checks whether a Jellyfin process is running on macOS.",
			"low", false, false, "", "If Jellyfin was installed as a launchd or Homebrew service, you can also inspect launchctl or brew services.", "macos",
		),
		example(
			"show hidden files here",
			"ls -A",
			"Lists files in the current directory, including hidden dotfiles.",
			"low", false, false, "", "This uses BSD/macOS ls flags.", "macos",
		),
		example(
			"why is port 3000 busy",
			"lsof -nP -iTCP:3000 -sTCP:LISTEN",
			"Shows processes listening on TCP port 3000 on macOS.",
			"low", false, false, "", "lsof is a good default listener check on macOS.", "macos",
		),
		example(
			"prefix all jpgs with vacation-",
			"for f in ./*.jpg; do [ -e \"$f\" ] || continue; base=$(basename \"$f\"); mv \"$f\" \"./vacation-$base\"; done",
			"Prefixes each JPG filename with 'vacation-'.",
			"high", true, false, "", "This renames files in place on macOS.", "macos",
		),
		example(
			"add cd.. as an alias to cd .. in "+shellName,
			"printf '\nalias cd..=\"cd ..\"\n' >> "+shellConfig,
			"Appends an alias named cd.. to "+shellConfig+" so it expands to 'cd ..'.",
			"high", true, false, "", "Run 'source "+shellConfig+"' or open a new terminal afterward to use the alias in new sessions.", "macos",
		),
	}
}

func powerShellExamples() []map[string]any {
	return []map[string]any{
		example(
			"check if jellyfin is running",
			"Get-Service -Name jellyfin",
			"Checks the Jellyfin service status in PowerShell.",
			"low", false, false, "", "If Jellyfin is not registered as a Windows service, inspect running processes instead.", "windows-powershell",
		),
		example(
			"show hidden files here",
			"Get-ChildItem -Force",
			"Lists files in the current directory, including hidden items.",
			"low", false, false, "", "This is the PowerShell-native equivalent of showing hidden files.", "windows-powershell",
		),
		example(
			"why is port 3000 busy",
			"Get-NetTCPConnection -LocalPort 3000 | Select-Object LocalAddress,LocalPort,State,OwningProcess",
			"Shows TCP listeners or connections using local port 3000 in PowerShell.",
			"low", false, false, "", "You can pipe the owning process ID into Get-Process if you need the process name.", "windows-powershell",
		),
		example(
			"prefix all jpgs with vacation-",
			"Get-ChildItem -Filter *.jpg | ForEach-Object { Rename-Item -LiteralPath $_.FullName -NewName ('vacation-' + $_.Name) }",
			"Prefixes each JPG filename with 'vacation-' in PowerShell.",
			"high", true, false, "", "This renames files in place.", "windows-powershell",
		),
		example(
			"add cd.. as an alias to cd .. in powershell profile",
			"Add-Content -Path $PROFILE -Value \"`nfunction cd.. { Set-Location .. }\"",
			"Appends a cd.. helper function to the current PowerShell profile.",
			"high", true, false, "", "Reload the profile with '. $PROFILE' or open a new PowerShell session afterward.", "windows-powershell",
		),
	}
}

func windowsExamples() []map[string]any {
	return []map[string]any{
		example(
			"show hidden files here",
			"dir /a",
			"Lists files in the current directory, including hidden items.",
			"low", false, false, "", "PowerShell gives better cross-platform coverage than cmd for more advanced tasks.", "windows",
		),
		example(
			"check if jellyfin is running",
			"sc query jellyfin",
			"Checks whether a Windows service named jellyfin is installed and running.",
			"low", false, false, "", "If you use PowerShell, Get-Service is usually nicer for this.", "windows",
		),
	}
}

func example(request, command, explanation, risk string, requiresConfirmation, needsClarification bool, clarificationQuestion, notes, platform string) map[string]any {
	return map[string]any{
		"request": request,
		"response": map[string]any{
			"command":               command,
			"explanation":           explanation,
			"risk":                  risk,
			"requiresConfirmation":  requiresConfirmation,
			"needsClarification":    needsClarification,
			"clarificationQuestion": clarificationQuestion,
			"notes":                 notes,
			"platform":              platform,
			"refused":               false,
		},
	}
}
