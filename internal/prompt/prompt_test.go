package prompt

import (
	"testing"

	"pls/internal/types"
)

func TestBuildIncludesLinuxExamples(t *testing.T) {
	messages := Build("check if jellyfin is running", types.RuntimeContext{OS: "linux", Shell: "bash"})
	examples, ok := messages.User["examples"].([]map[string]any)
	if !ok || len(examples) == 0 {
		t.Fatalf("expected examples in prompt payload")
	}

	foundJellyfin := false
	foundDotfiles := false
	for _, item := range examples {
		response, _ := item["response"].(map[string]any)
		if response["command"] == "systemctl is-active jellyfin" {
			foundJellyfin = true
		}
		if response["command"] == "find . -maxdepth 1 -mindepth 1 -name '.*' -print" {
			foundDotfiles = true
		}
	}
	if !foundJellyfin {
		t.Fatalf("expected linux jellyfin example in prompt payload")
	}
	if !foundDotfiles {
		t.Fatalf("expected linux dotfiles example in prompt payload")
	}
}

func TestBuildIncludesMacOSExamples(t *testing.T) {
	messages := Build("why is port 3000 busy", types.RuntimeContext{OS: "macos", Shell: "zsh"})
	examples, ok := messages.User["examples"].([]map[string]any)
	if !ok || len(examples) == 0 {
		t.Fatalf("expected examples in prompt payload")
	}

	found := false
	for _, item := range examples {
		response, _ := item["response"].(map[string]any)
		if response["command"] == "lsof -nP -iTCP:3000 -sTCP:LISTEN" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected macOS port example in prompt payload")
	}
}

func TestBuildIncludesPowerShellExamples(t *testing.T) {
	messages := Build("show hidden files here", types.RuntimeContext{OS: "windows", Shell: "powershell"})
	examples, ok := messages.User["examples"].([]map[string]any)
	if !ok || len(examples) == 0 {
		t.Fatalf("expected examples in prompt payload")
	}

	found := false
	for _, item := range examples {
		response, _ := item["response"].(map[string]any)
		if response["command"] == "Get-ChildItem -Force" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected PowerShell hidden-files example in prompt payload")
	}
}
