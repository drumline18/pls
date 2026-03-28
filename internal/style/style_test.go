package style

import (
	"strings"
	"testing"

	"pls/internal/types"
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
