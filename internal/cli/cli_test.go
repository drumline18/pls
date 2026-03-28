package cli

import "testing"

func TestParseArgsCapturesFreeformRequestWithoutQuotes(t *testing.T) {
	parsed, err := ParseArgs([]string{"show", "me", "dotfiles", "here"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if got, want := len(parsed.RequestParts), 4; got != want {
		t.Fatalf("unexpected request part count: got %d want %d", got, want)
	}
}

func TestParseArgsSupportsKnownFlagsAndFreeformTail(t *testing.T) {
	parsed, err := ParseArgs([]string{"--provider", "ollama", "--model", "qwen2.5-coder:7b-instruct-q4_K_M", "show", "hidden", "files"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if parsed.Flags.Provider != "ollama" {
		t.Fatalf("unexpected provider: %s", parsed.Flags.Provider)
	}
	if parsed.Flags.Model != "qwen2.5-coder:7b-instruct-q4_K_M" {
		t.Fatalf("unexpected model: %s", parsed.Flags.Model)
	}
}
