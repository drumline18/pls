package cli

import (
	"reflect"
	"testing"
)

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
	parsed, err := ParseArgs([]string{"--provider", "ollama", "--model", "qwen3.5:4b", "show", "hidden", "files"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if parsed.Flags.Provider != "ollama" {
		t.Fatalf("unexpected provider: %s", parsed.Flags.Provider)
	}
	if parsed.Flags.Model != "qwen3.5:4b" {
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

func TestParseArgsLeavesConfigInitAsRequestSequence(t *testing.T) {
	parsed, err := ParseArgs([]string{"config", "init"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if !reflect.DeepEqual(parsed.RequestParts, []string{"config", "init"}) {
		t.Fatalf("unexpected request parts: %#v", parsed.RequestParts)
	}
}

func TestParseArgsLeavesConfigLocalInitAsRequestSequence(t *testing.T) {
	parsed, err := ParseArgs([]string{"config", "local", "init"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if !reflect.DeepEqual(parsed.RequestParts, []string{"config", "local", "init"}) {
		t.Fatalf("unexpected request parts: %#v", parsed.RequestParts)
	}
}

func TestParseArgsLeavesSetupAsRequest(t *testing.T) {
	parsed, err := ParseArgs([]string{"setup"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if !reflect.DeepEqual(parsed.RequestParts, []string{"setup"}) {
		t.Fatalf("unexpected request parts: %#v", parsed.RequestParts)
	}
}

func TestParseArgsSupportsExecutionFlagsBeforeRequest(t *testing.T) {
	parsed, err := ParseArgs([]string{"--yes", "--no-exec", "show", "hidden", "files"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if !parsed.Flags.Yes || !parsed.Flags.NoExec {
		t.Fatalf("expected yes and no-exec flags to be set: %#v", parsed.Flags)
	}
	if !reflect.DeepEqual(parsed.RequestParts, []string{"show", "hidden", "files"}) {
		t.Fatalf("unexpected request parts: %#v", parsed.RequestParts)
	}
}

func TestParseArgsStopsParsingFlagsAfterRequestStarts(t *testing.T) {
	parsed, err := ParseArgs([]string{"show", "files", "--json", "now"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if parsed.Flags.JSON {
		t.Fatalf("did not expect --json to be parsed as a flag after request started")
	}
	if !reflect.DeepEqual(parsed.RequestParts, []string{"show", "files", "--json", "now"}) {
		t.Fatalf("unexpected request parts: %#v", parsed.RequestParts)
	}
}

func TestParseArgsSupportsDoubleDashSeparator(t *testing.T) {
	parsed, err := ParseArgs([]string{"--yes", "--", "show", "me", "--json"})
	if err != nil {
		t.Fatalf("ParseArgs returned error: %v", err)
	}

	if !parsed.Flags.Yes {
		t.Fatalf("expected --yes flag to be set")
	}
	if !reflect.DeepEqual(parsed.RequestParts, []string{"show", "me", "--json"}) {
		t.Fatalf("unexpected request parts: %#v", parsed.RequestParts)
	}
}
