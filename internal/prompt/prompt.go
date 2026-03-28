package prompt

import "pls/internal/types"

func Build(request string, runtimeContext types.RuntimeContext) types.Messages {
	system := `You generate safe shell command suggestions.
Return JSON only. Do not wrap in markdown fences.
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
		"examples": []map[string]any{
			{
				"request": "check if jellyfin is running",
				"response": map[string]any{
					"command":              "systemctl is-active jellyfin",
					"explanation":          "Checks whether the Jellyfin systemd service is active.",
					"risk":                 "low",
					"requiresConfirmation": false,
					"needsClarification":   false,
					"clarificationQuestion": "",
					"notes":                "If Jellyfin was started with Docker or another supervisor, inspect that runtime instead.",
					"platform":             "linux",
					"refused":              false,
				},
			},
			{
				"request": "prefix all the mp3s with their lengths in seconds",
				"response": map[string]any{
					"command":              "for f in ./*.mp3; do [ -e \"$f\" ] || continue; secs=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 \"$f\" | awk '{printf \"%d\", $1}'); mv -- \"$f\" \"./${secs}s - $(basename \"$f\")\"; done",
					"explanation":          "Prefixes each MP3 filename with its duration in whole seconds.",
					"risk":                 "high",
					"requiresConfirmation": true,
					"needsClarification":   false,
					"clarificationQuestion": "",
					"notes":                "This renames files in place and requires ffprobe from ffmpeg.",
					"platform":             "linux",
					"refused":              false,
				},
			},
			{
				"request": "prefix all jpgs with vacation-",
				"response": map[string]any{
					"command":              "for f in ./*.jpg; do [ -e \"$f\" ] || continue; base=$(basename \"$f\"); mv -- \"$f\" \"./vacation-$base\"; done",
					"explanation":          "Prefixes each JPG filename with 'vacation-'.",
					"risk":                 "high",
					"requiresConfirmation": true,
					"needsClarification":   false,
					"clarificationQuestion": "",
					"notes":                "This renames files in place.",
					"platform":             "linux",
					"refused":              false,
				},
			},
			{
				"request": "replace spaces in all filenames here with underscores",
				"response": map[string]any{
					"command":              "for f in ./*\\ *; do [ -e \"$f\" ] || continue; mv -- \"$f\" \"${f// /_}\"; done",
					"explanation":          "Replaces spaces with underscores in filenames that contain spaces in the current directory.",
					"risk":                 "high",
					"requiresConfirmation": true,
					"needsClarification":   false,
					"clarificationQuestion": "",
					"notes":                "This renames files in place.",
					"platform":             "linux",
					"refused":              false,
				},
			},
			{
				"request": "move all srt files into a subtitles folder",
				"response": map[string]any{
					"command":              "mkdir -p \"./subtitles\" && for f in ./*.srt; do [ -e \"$f\" ] || continue; mv -- \"$f\" \"./subtitles/\"; done",
					"explanation":          "Moves all .srt files from the current directory into the 'subtitles' folder.",
					"risk":                 "high",
					"requiresConfirmation": true,
					"needsClarification":   false,
					"clarificationQuestion": "",
					"notes":                "This creates the destination folder if needed and moves files into it.",
					"platform":             "linux",
					"refused":              false,
				},
			},
			{
				"request": "add cd.. as an alias to cd .. in bashrc",
				"response": map[string]any{
					"command":              "printf '\nalias cd..=\"cd ..\"\n' >> ~/.bashrc",
					"explanation":          "Appends an alias named cd.. to ~/.bashrc so it expands to 'cd ..'.",
					"risk":                 "high",
					"requiresConfirmation": true,
					"needsClarification":   false,
					"clarificationQuestion": "",
					"notes":                "Run 'source ~/.bashrc' or open a new shell afterward to use the alias in new sessions.",
					"platform":             "linux",
					"refused":              false,
				},
			},
		},
		"instructions": []string{
			"Target a single primary command where possible.",
			"Assume the command is for the current working directory unless the user says otherwise.",
			"Prefer inspection commands for this MVP. Do not add execution wrappers, aliases, or shell functions.",
			"Prefer direct predicates and flags over text-parsing pipelines.",
			"If the user wants only files, only directories, or hidden entries, prefer direct find predicates or native flags instead of grepping ls output.",
			"For Linux service checks phrased like 'check if jellyfin is running', prefer systemctl is-active <service> when the target looks like a service name, and mention container/manual-process alternatives in notes if needed.",
			"For batch renames or metadata-based filename changes, a concise for-loop using purpose-built tools is acceptable.",
			"For batch rename or move commands, use safe quoting, guard globs with [ -e \"$f\" ] || continue, and prefer basename over brittle parameter expansion when prefixing filenames.",
			"If multiple platform variants matter, put the main one in command and mention differences in notes.",
		},
	}

	return types.Messages{System: system, User: user}
}
