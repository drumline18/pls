package policy

import (
	"regexp"
	"strings"

	"pls/internal/types"
)

var highRiskPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bsudo\b`),
	regexp.MustCompile(`(?i)\brm\b\s+(-[a-zA-Z]*[rf][a-zA-Z]*|--recursive|--force)`),
	regexp.MustCompile(`(?i)\bdd\b`),
	regexp.MustCompile(`(?i)\bmkfs\b`),
	regexp.MustCompile(`(?i)\bshutdown\b`),
	regexp.MustCompile(`(?i)\breboot\b`),
	regexp.MustCompile(`(?i)\bpoweroff\b`),
	regexp.MustCompile(`(?i)curl\s+[^|\n]+\|\s*(sh|bash|zsh)`),
	regexp.MustCompile(`(?i)wget\s+[^|\n]+\|\s*(sh|bash|zsh)`),
	regexp.MustCompile(`>\s*/dev/`),
}

func Apply(s types.Suggestion) types.Suggestion {
	for _, pattern := range highRiskPatterns {
		if pattern.MatchString(s.Command) {
			s = escalateHigh(s, "Safety policy flagged this command for manual review before any future execution support.")
			break
		}
	}

	if isFileMutationCommand(s.Command) {
		s = escalateHigh(s, "This command renames or moves files and should be reviewed before execution.")
	}

	return s
}

func isFileMutationCommand(command string) bool {
	normalized := strings.TrimSpace(command)
	return strings.Contains(normalized, " mv ") || strings.HasPrefix(normalized, "mv ") || strings.Contains(normalized, "rename ") || strings.Contains(normalized, "mmv ")
}

func escalateHigh(s types.Suggestion, note string) types.Suggestion {
	if riskRank(s.Risk) < riskRank("high") {
		s.Risk = "high"
	}
	s.RequiresConfirmation = true
	s.Notes = joinNotes(s.Notes, note)
	return s
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

func riskRank(risk string) int {
	switch risk {
	case "low":
		return 1
	case "medium":
		return 2
	case "high":
		return 3
	case "critical":
		return 4
	default:
		return 0
	}
}
