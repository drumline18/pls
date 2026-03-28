package doctor

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"pls/internal/types"
)

func TestHumanIncludesJokeAndConfigSections(t *testing.T) {
	output := Human(Report{
		Joke:                 "Why did the CLI go to the doctor? It had too many terminal conditions.",
		OverallStatus:        "ok",
		ConfigPath:           "/tmp/global.json",
		ConfigExists:         true,
		LocalConfigPath:      "/tmp/project/pls.json",
		LocalConfigExists:    true,
		YoloMode:             true,
		YoloSource:           "local",
		Provider:             "ollama",
		Model:                "qwen2.5-coder:7b-instruct-q4_K_M",
		Host:                 "http://127.0.0.1:11434",
		Shell:                "bash",
		OS:                   "linux",
		CWD:                  "/tmp",
		ExecutablePath:       "/tmp/pls",
		InstalledCommandPath: "/usr/local/bin/pls",
		InPath:               true,
		ProviderStatus:       "ok",
		ProviderMessage:      "connected to Ollama successfully",
	})

	for _, needle := range []string{"terminal conditions", "global path", "local override", "yolo mode: yes", "connected to Ollama successfully"} {
		if !strings.Contains(output, needle) {
			t.Fatalf("expected output to contain %q", needle)
		}
	}
}

func TestRunHandlesMissingOptionalConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("PLS_CONFIG", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")
	t.Setenv("OLLAMA_HOST", "")
	t.Setenv("PLS_OLLAMA_HOST", "")
	t.Setenv("PLS_YOLO_MODE", "")

	report, err := Run(context.Background(), types.Flags{})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if report.ConfigPath == "" {
		t.Fatalf("expected config path to be set")
	}
	if report.Provider != "ollama" {
		t.Fatalf("expected default provider ollama, got %s", report.Provider)
	}
	if report.YoloMode {
		t.Fatalf("did not expect yolo mode to be enabled")
	}
}
