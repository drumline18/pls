package style

import (
	"strings"
	"testing"

	"github.com/drumline18/pls/internal/types"
)

func TestNormalizeHiddenListingRewritesLSVariants(t *testing.T) {
	suggestion := Normalize("show me all dotfiles in this directory", types.RuntimeContext{OS: "linux"}, types.Suggestion{
		Command:     "ls -a . | grep '^.'",
		Explanation: "Lists dotfiles.",
		Risk:        "low",
	})

	if suggestion.Command != "find . -maxdepth 1 -mindepth 1 -name '.*' -print" {
		t.Fatalf("unexpected command: %s", suggestion.Command)
	}
	if !strings.Contains(suggestion.Notes, "Style policy normalized") {
		t.Fatalf("expected normalization note, got %q", suggestion.Notes)
	}
}

func TestNormalizeDirectoryListingRewritesLSPipeGrep(t *testing.T) {
	suggestion := Normalize("show only directories here", types.RuntimeContext{OS: "linux"}, types.Suggestion{
		Command:     "ls -l | grep '^d'",
		Explanation: "Lists only directories.",
		Risk:        "low",
	})

	if suggestion.Command != "find . -maxdepth 1 -mindepth 1 -type d -print" {
		t.Fatalf("unexpected command: %s", suggestion.Command)
	}
}

func TestNormalizePortInspectionRewritesNetstatGrep(t *testing.T) {
	suggestion := Normalize("why is port 3000 busy", types.RuntimeContext{OS: "linux"}, types.Suggestion{
		Command:     "netstat -tuln | grep :3000",
		Explanation: "Checks port 3000.",
		Risk:        "low",
	})

	if suggestion.Command != "ss -ltnp 'sport = :3000'" {
		t.Fatalf("unexpected command: %s", suggestion.Command)
	}
}

func TestNormalizeServiceInspectionRewritesPgrep(t *testing.T) {
	suggestion := Normalize("check if jellyfin is running", types.RuntimeContext{OS: "linux"}, types.Suggestion{
		Command:     "pgrep -f jellyfin",
		Explanation: "Checks for a Jellyfin process.",
		Risk:        "low",
	})

	if suggestion.Command != "systemctl is-active jellyfin" {
		t.Fatalf("unexpected command: %s", suggestion.Command)
	}
}

func TestNormalizePrefixRenameUsesBasenameLoop(t *testing.T) {
	suggestion := Normalize("prefix all jpgs with vacation-", types.RuntimeContext{OS: "linux"}, types.Suggestion{
		Command:     "for f in ./*.jpg; do [ -e \"$f\" ] || continue; mv -- \"$f\" \"./vacation-${f#*.jpg}.jpg\"; done",
		Explanation: "Prefixes JPG names.",
		Risk:        "high",
	})

	want := "for f in ./*.jpg; do [ -e \"$f\" ] || continue; base=$(basename \"$f\"); mv -- \"$f\" \"./vacation-$base\"; done"
	if suggestion.Command != want {
		t.Fatalf("unexpected command: %s", suggestion.Command)
	}
}

func TestNormalizeReplaceSpacesTouchesOnlyMatchingFiles(t *testing.T) {
	suggestion := Normalize("replace spaces in all filenames here with underscores", types.RuntimeContext{OS: "linux"}, types.Suggestion{
		Command:     "for f in *; do mv -- \"$f\" \"${f// /_}\"; done",
		Explanation: "Replaces spaces.",
		Risk:        "high",
	})

	want := "for f in ./*\\ *; do [ -e \"$f\" ] || continue; mv -- \"$f\" \"${f// /_}\"; done"
	if suggestion.Command != want {
		t.Fatalf("unexpected command: %s", suggestion.Command)
	}
}

func TestNormalizeMoveIntoFolderUsesGuardedLoop(t *testing.T) {
	suggestion := Normalize("move all srt files into a subtitles folder", types.RuntimeContext{OS: "linux"}, types.Suggestion{
		Command:     "mkdir -p subtitles && mv *.srt subtitles/",
		Explanation: "Moves subtitle files.",
		Risk:        "low",
	})

	want := "mkdir -p \"./subtitles\" && for f in ./*.srt; do [ -e \"$f\" ] || continue; mv -- \"$f\" \"./subtitles/\"; done"
	if suggestion.Command != want {
		t.Fatalf("unexpected command: %s", suggestion.Command)
	}
}

func TestDirectSuggestionHandlesFindFilesBiggerThan(t *testing.T) {
	suggestion, ok := DirectSuggestion("find files bigger than 500mb under the current directory", types.RuntimeContext{OS: "linux"})
	if !ok {
		t.Fatal("expected direct suggestion")
	}
	if suggestion.Command != "find . -type f -size +500M -print" {
		t.Fatalf("unexpected command: %s", suggestion.Command)
	}
	if suggestion.Risk != "low" {
		t.Fatalf("unexpected risk: %s", suggestion.Risk)
	}
}

func TestDirectSuggestionHandlesBiggestFiles(t *testing.T) {
	suggestion, ok := DirectSuggestion("show the 10 biggest files under the current directory", types.RuntimeContext{OS: "linux"})
	if !ok {
		t.Fatal("expected direct suggestion")
	}
	if suggestion.Command != "find . -type f -printf '%s %p\\n' | sort -rn | head -10" {
		t.Fatalf("unexpected command: %s", suggestion.Command)
	}
}

func TestNormalizeLeavesUnrelatedCommandsAlone(t *testing.T) {
	original := types.Suggestion{
		Command:     "ss -ltnp 'sport = :3000'",
		Explanation: "Shows listeners on port 3000.",
		Risk:        "low",
	}

	normalized := Normalize("why is port 3000 busy", types.RuntimeContext{OS: "linux"}, original)
	if normalized.Command != original.Command {
		t.Fatalf("unexpected command rewrite: %s", normalized.Command)
	}
}
