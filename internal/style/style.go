package style

import (
	"regexp"
	"strings"

	"pls/internal/types"
)

var hiddenListPattern = regexp.MustCompile(`(?i)^ls\b.*\|\s*grep\s+['"]\^\.[^'"]*['"]$`)
var directoryOnlyPattern = regexp.MustCompile(`(?i)^ls\b.*\|\s*grep\s+['"]\^d['"]$`)
var fileOnlyPattern = regexp.MustCompile(`(?i)^ls\b.*\|\s*grep\s+-v\s+['"]\^d['"]$`)
var portPipePattern = regexp.MustCompile(`(?i)^(netstat|ss)\b.*\|\s*grep\b.*$`)
var portNumberPattern = regexp.MustCompile(`(?i)\bport\s+(\d{1,5})\b`)
var psGrepPattern = regexp.MustCompile(`(?i)^ps\b.*\|\s*grep\b.*$`)
var serviceRequestPattern = regexp.MustCompile(`(?i)^(?:check\s+if|is|whether)\s+([a-z0-9_.@-]+)\s+(?:service\s+)?(?:is\s+)?running\b`)
var prefixRequestPattern = regexp.MustCompile(`(?i)^prefix all (?:the )?([a-z0-9*_.-]+?)s?\s+with\s+(.+)$`)
var moveIntoFolderRequestPattern = regexp.MustCompile(`(?i)^move all (?:the )?([a-z0-9*_.-]+?)s?\s+files?\s+into\s+(?:a |an |the )?(.+?)\s+folder$`)
var safeExtensionPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

func Normalize(request string, runtimeContext types.RuntimeContext, suggestion types.Suggestion) types.Suggestion {
	requestTrimmed := strings.TrimSpace(request)
	requestLower := strings.ToLower(requestTrimmed)
	command := normalizeWhitespace(suggestion.Command)

	if replacement, explanation, matched := normalizeHiddenListing(requestLower, runtimeContext.OS, command); matched {
		suggestion.Command = replacement
		suggestion.Explanation = explanation
		suggestion.Notes = joinNotes(suggestion.Notes, "Style policy normalized this to a simpler command and avoided parsing ls output.")
		return suggestion
	}

	if replacement, explanation, matched := normalizeDirectoryListing(requestLower, runtimeContext.OS, command); matched {
		suggestion.Command = replacement
		suggestion.Explanation = explanation
		suggestion.Notes = joinNotes(suggestion.Notes, "Style policy normalized this to a direct directory filter instead of parsing ls output.")
		return suggestion
	}

	if replacement, explanation, matched := normalizeFileListing(requestLower, runtimeContext.OS, command); matched {
		suggestion.Command = replacement
		suggestion.Explanation = explanation
		suggestion.Notes = joinNotes(suggestion.Notes, "Style policy normalized this to a direct file filter instead of parsing ls output.")
		return suggestion
	}

	if replacement, explanation, matched := normalizePortInspection(requestLower, runtimeContext.OS, command); matched {
		suggestion.Command = replacement
		suggestion.Explanation = explanation
		suggestion.Notes = joinNotes(suggestion.Notes, "Style policy normalized this to a direct socket inspection command instead of a grep pipeline.")
		return suggestion
	}

	if replacement, explanation, matched := normalizeServiceInspection(requestLower, runtimeContext.OS, command); matched {
		suggestion.Command = replacement
		suggestion.Explanation = explanation
		suggestion.Notes = joinNotes(suggestion.Notes, "Style policy normalized this to a direct service-status check instead of a process-name search.")
		return suggestion
	}

	if replacement, explanation, matched := normalizePrefixRename(requestTrimmed, runtimeContext.OS, command); matched {
		suggestion.Command = replacement
		suggestion.Explanation = explanation
		suggestion.Notes = joinNotes(suggestion.Notes, "Style policy normalized this to a safer quoted rename loop using basename.")
		return suggestion
	}

	if replacement, explanation, matched := normalizeReplaceSpaces(requestLower, runtimeContext.OS, command); matched {
		suggestion.Command = replacement
		suggestion.Explanation = explanation
		suggestion.Notes = joinNotes(suggestion.Notes, "Style policy normalized this to only touch filenames that actually contain spaces.")
		return suggestion
	}

	if replacement, explanation, matched := normalizeMoveIntoFolder(requestTrimmed, runtimeContext.OS, command); matched {
		suggestion.Command = replacement
		suggestion.Explanation = explanation
		suggestion.Notes = joinNotes(suggestion.Notes, "Style policy normalized this to a guarded move loop with an explicit destination directory.")
		return suggestion
	}

	return suggestion
}

func normalizeHiddenListing(request, osName, command string) (string, string, bool) {
	if osName != "linux" || !isHiddenListingRequest(request) {
		return "", "", false
	}

	if hiddenListPattern.MatchString(command) || isAnyCommand(command, "ls -a", "ls -a .", "ls -la", "ls -la .", "ls -A", "ls -A .") {
		return "find . -maxdepth 1 -mindepth 1 -name '.*' -print", "Lists hidden files and directories in the current directory.", true
	}

	return "", "", false
}

func normalizeDirectoryListing(request, osName, command string) (string, string, bool) {
	if osName != "linux" || !isDirectoryOnlyRequest(request) {
		return "", "", false
	}

	if directoryOnlyPattern.MatchString(command) {
		return "find . -maxdepth 1 -mindepth 1 -type d -print", "Lists only directories in the current directory.", true
	}

	return "", "", false
}

func normalizeFileListing(request, osName, command string) (string, string, bool) {
	if osName != "linux" || !isFileOnlyRequest(request) {
		return "", "", false
	}

	if fileOnlyPattern.MatchString(command) {
		return "find . -maxdepth 1 -mindepth 1 -type f -print", "Lists only regular files in the current directory.", true
	}

	return "", "", false
}

func normalizePortInspection(request, osName, command string) (string, string, bool) {
	if osName != "linux" || !isPortInspectionRequest(request) || !portPipePattern.MatchString(command) {
		return "", "", false
	}

	match := portNumberPattern.FindStringSubmatch(request)
	if len(match) != 2 {
		return "", "", false
	}

	port := match[1]
	return "ss -ltnp 'sport = :" + port + "'", "Shows listening sockets bound to port " + port + ", including the owning process when available.", true
}

func normalizeServiceInspection(request, osName, command string) (string, string, bool) {
	if osName != "linux" || !(strings.HasPrefix(command, "pgrep ") || psGrepPattern.MatchString(command)) {
		return "", "", false
	}

	match := serviceRequestPattern.FindStringSubmatch(request)
	if len(match) != 2 {
		return "", "", false
	}

	service := match[1]
	return "systemctl is-active " + service, "Checks whether the " + service + " systemd service is active.", true
}

func normalizePrefixRename(request, osName, command string) (string, string, bool) {
	if osName != "linux" || !strings.Contains(command, "mv --") {
		return "", "", false
	}

	match := prefixRequestPattern.FindStringSubmatch(request)
	if len(match) != 3 {
		return "", "", false
	}

	prefix := strings.TrimSpace(match[2])
	if prefix == "" || strings.Contains(strings.ToLower(prefix), "their ") {
		return "", "", false
	}

	extension, ok := extensionFromToken(match[1])
	if !ok {
		return "", "", false
	}

	escapedPrefix := escapeDoubleQuoted(prefix)
	return "for f in ./*." + extension + "; do [ -e \"$f\" ] || continue; base=$(basename \"$f\"); mv -- \"$f\" \"./" + escapedPrefix + "$base\"; done", "Prefixes each ." + extension + " filename with '" + prefix + "'.", true
}

func normalizeReplaceSpaces(request, osName, command string) (string, string, bool) {
	if osName != "linux" || !isReplaceSpacesRequest(request) || !strings.Contains(command, "${f// /_}") {
		return "", "", false
	}

	return "for f in ./*\\ *; do [ -e \"$f\" ] || continue; mv -- \"$f\" \"${f// /_}\"; done", "Replaces spaces with underscores in filenames that contain spaces in the current directory.", true
}

func normalizeMoveIntoFolder(request, osName, command string) (string, string, bool) {
	if osName != "linux" || !strings.Contains(command, "mv") {
		return "", "", false
	}

	match := moveIntoFolderRequestPattern.FindStringSubmatch(request)
	if len(match) != 3 {
		return "", "", false
	}

	extension, ok := extensionFromToken(match[1])
	if !ok {
		return "", "", false
	}

	folder := strings.TrimSpace(match[2])
	if folder == "" {
		return "", "", false
	}

	escapedFolder := escapeDoubleQuoted(folder)
	return "mkdir -p \"./" + escapedFolder + "\" && for f in ./*." + extension + "; do [ -e \"$f\" ] || continue; mv -- \"$f\" \"./" + escapedFolder + "/\"; done", "Moves all ." + extension + " files from the current directory into the '" + folder + "' folder.", true
}

func isHiddenListingRequest(request string) bool {
	return containsAny(request,
		"dotfile",
		"dotfiles",
		"hidden file",
		"hidden files",
		"hidden directory",
		"hidden directories",
		"show hidden",
		"list hidden",
	)
}

func isDirectoryOnlyRequest(request string) bool {
	return containsAny(request,
		"only directories",
		"just directories",
		"directories only",
		"show directories",
		"list directories",
	)
}

func isFileOnlyRequest(request string) bool {
	return containsAny(request,
		"only files",
		"just files",
		"files only",
		"show files",
		"list files",
	)
}

func isPortInspectionRequest(request string) bool {
	if !portNumberPattern.MatchString(request) {
		return false
	}
	return containsAny(request,
		"busy",
		"in use",
		"using",
		"listening",
		"bound",
	)
}

func isReplaceSpacesRequest(request string) bool {
	return strings.Contains(request, "replace spaces") && strings.Contains(request, "underscores")
}

func extensionFromToken(token string) (string, bool) {
	value := strings.ToLower(strings.TrimSpace(token))
	switch {
	case strings.HasPrefix(value, "*."):
		value = value[2:]
	case strings.HasPrefix(value, "."):
		value = value[1:]
	case strings.HasSuffix(value, "s"):
		value = strings.TrimSuffix(value, "s")
	}

	if !safeExtensionPattern.MatchString(value) {
		return "", false
	}

	return value, true
}

func escapeDoubleQuoted(value string) string {
	replacer := strings.NewReplacer(`\\`, `\\\\`, `"`, `\\"`, `$`, `\\$`, "`", "\\`")
	return replacer.Replace(value)
}

func containsAny(value string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(value, needle) {
			return true
		}
	}
	return false
}

func isAnyCommand(command string, variants ...string) bool {
	for _, variant := range variants {
		if command == variant {
			return true
		}
	}
	return false
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func joinNotes(existing, extra string) string {
	if existing == "" {
		return extra
	}
	if extra == "" {
		return existing
	}
	return existing + " " + extra
}
