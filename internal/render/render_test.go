package render

import (
	"strings"
	"testing"

	"github.com/drumline18/pls/internal/types"
)

func TestHumanMentionsCurrentConfirmationBehavior(t *testing.T) {
	output := Human(types.Suggestion{
		Command:              "printf 'alias cd..=\"cd ..\"\n' >> ~/.bashrc",
		Explanation:          "Appends an alias to ~/.bashrc.",
		Risk:                 "high",
		RequiresConfirmation: true,
	})

	if !strings.Contains(output, "This command will ask for confirmation before execution.") {
		t.Fatalf("expected current execution wording, got: %s", output)
	}
	if strings.Contains(output, "future execution-enabled version") {
		t.Fatalf("found stale execution wording: %s", output)
	}
}
