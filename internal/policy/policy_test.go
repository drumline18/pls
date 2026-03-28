package policy

import (
	"strings"
	"testing"

	"pls/internal/types"
)

func TestApplyEscalatesFileMutationCommands(t *testing.T) {
	result := Apply(types.Suggestion{
		Command:              "mkdir -p \"./subtitles\" && for f in ./*.srt; do [ -e \"$f\" ] || continue; mv -- \"$f\" \"./subtitles/\"; done",
		Explanation:          "Moves subtitle files.",
		Risk:                 "low",
		RequiresConfirmation: false,
	})

	if result.Risk != "high" {
		t.Fatalf("unexpected risk: %s", result.Risk)
	}
	if !result.RequiresConfirmation {
		t.Fatalf("expected confirmation to be required")
	}
	if !strings.Contains(result.Notes, "renames or moves files") {
		t.Fatalf("expected file-mutation note, got %q", result.Notes)
	}
}
