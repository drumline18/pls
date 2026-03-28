package configshow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"pls/internal/types"
)

func TestRunReportsEffectiveConfigAndPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("PLS_CONFIG", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")
	t.Setenv("PLS_YOLO_MODE", "")

	globalPath := filepath.Join(home, ".config", "pls", "config.json")
	if err := os.MkdirAll(filepath.Dir(globalPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(globalPath, []byte(`{"provider":"ollama","model":"qwen2.5-coder:7b-instruct-q4_K_M","host":"http://192.168.2.166:11434","yoloMode":false}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	projectRoot := filepath.Join(home, "project")
	if err := os.MkdirAll(projectRoot, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectRoot, "pls.json"), []byte(`{"yoloMode":true}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	defer func() { _ = os.Chdir(previousWD) }()
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}

	report, err := Run(types.Flags{})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if report.GlobalConfigPath != globalPath {
		t.Fatalf("unexpected global path: %s", report.GlobalConfigPath)
	}
	if !report.GlobalConfigExists {
		t.Fatalf("expected global config to exist")
	}
	if report.LocalConfigPath != filepath.Join(projectRoot, "pls.json") {
		t.Fatalf("unexpected local config path: %s", report.LocalConfigPath)
	}
	if !report.YoloMode || report.YoloSource != "local" {
		t.Fatalf("unexpected yolo state: %#v", report)
	}
}

func TestHumanIncludesKeyFields(t *testing.T) {
	output := Human(Report{
		GlobalConfigPath:       "/tmp/global.json",
		GlobalConfigExists:     true,
		LocalConfigPath:        "/tmp/project/pls.json",
		LocalConfigExists:      true,
		EffectiveProvider:      "ollama",
		EffectiveModel:         "qwen2.5-coder:7b-instruct-q4_K_M",
		EffectiveHost:          "http://127.0.0.1:11434",
		ConfigStoredAPIKeyConfigured: false,
		YoloMode:               true,
		YoloSource:             "local",
	})

	for _, needle := range []string{"pls config show", "global path", "local override", "provider: ollama", "yolo mode: yes", "config-stored api key configured: no"} {
		if !strings.Contains(output, needle) {
			t.Fatalf("expected output to contain %q", needle)
		}
	}
}
