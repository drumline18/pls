package render

import (
	"fmt"
	"strings"

	"pls/internal/types"
)

func Human(result types.Suggestion) string {
	if result.NeedsClarification {
		return strings.Join([]string{
			"Need clarification:",
			fmt.Sprintf("  %s", result.ClarificationQuestion),
		}, "\n")
	}

	if result.Refused {
		lines := []string{
			"Refused:",
			fmt.Sprintf("  %s", result.Explanation),
		}
		if result.Notes != "" {
			lines = append(lines, "", "Notes:", fmt.Sprintf("  %s", result.Notes))
		}
		return strings.Join(lines, "\n")
	}

	lines := []string{
		"Command:",
		fmt.Sprintf("  %s", result.Command),
		"",
		"Why:",
		fmt.Sprintf("  %s", result.Explanation),
		"",
		"Risk:",
		fmt.Sprintf("  %s", result.Risk),
	}

	if result.Platform != "" {
		lines = append(lines, "", "Platform:", fmt.Sprintf("  %s", result.Platform))
	}
	if result.Notes != "" {
		lines = append(lines, "", "Notes:", fmt.Sprintf("  %s", result.Notes))
	}
	if result.RequiresConfirmation {
		lines = append(lines, "", "Execution:", "  This command will ask for confirmation before execution.")
	}

	return strings.Join(lines, "\n")
}
