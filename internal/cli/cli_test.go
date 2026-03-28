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

func TestParseArgsSupportsConfigFlags(t *testing.T) {
	parsed, err := ParseArgs([]string{"--config", "~/test.json", "--print-config-path"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if parsed.Flags.ConfigPath != "~/test.json" {
		t.Fatalf("unexpected config path: %s", parsed.Flags.ConfigPath)
	}
	if !parsed.Flags.PrintConfigPath {
		t.Fatalf("expected PrintConfigPath to be true")
	}
}

func TestParseArgsLeavesDoctorAsRequest(t *testing.T) {
	parsed, err := ParseArgs([]string{"doctor"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if len(parsed.RequestParts) != 1 || parsed.RequestParts[0] != "doctor" {
		t.Fatalf("unexpected request parts: %#v", parsed.RequestParts)
	}
}
