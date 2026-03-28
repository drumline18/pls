package types

import "testing"

func TestValidateSuggestionAcceptsSchemaCompliantPayloads(t *testing.T) {
	result, err := ValidateSuggestion(Suggestion{
		Command:              "ls -la",
		Explanation:          "Lists files including dotfiles.",
		Risk:                 "low",
		RequiresConfirmation: false,
		NeedsClarification:   false,
		Platform:             "linux",
	})
	if err != nil {
		t.Fatalf("ValidateSuggestion returned error: %v", err)
	}

	if result.Command != "ls -la" {
		t.Fatalf("unexpected command: %s", result.Command)
	}
}
