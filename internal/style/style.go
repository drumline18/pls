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

func Normalize(request string, runtimeContext types.RuntimeContext, suggestion types.Suggestion) types.Suggestion {
	requestLower := strings.ToLower(strings.TrimSpace(request))
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
