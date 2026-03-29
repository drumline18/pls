package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/drumline18/pls/internal/types"
)

func TestExplicitMissingConfigPathReturnsError(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")
	t.Setenv("PLS_CONFIG", filepath.Join(t.TempDir(), "missing.json"))

	_, err := Load(types.Flags{})
	if err == nil {
		t.Fatalf("expected error for explicit missing config path")
	}
}

func TestDefaultPathWithoutConfigFallsBackToOllama(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("PLS_CONFIG", "")

	cfg, err := Load(types.Flags{})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Provider != "ollama" {
		t.Fatalf("unexpected provider: %s", cfg.Provider)
	}
	if cfg.Model != "qwen2.5-coder:7b-instruct-q4_K_M" {
		t.Fatalf("unexpected model: %s", cfg.Model)
	}
	if want := filepath.Join(home, ".config", "pls", "config.json"); cfg.ConfigPath != want {
		t.Fatalf("unexpected config path: got %s want %s", cfg.ConfigPath, want)
	}
	if cfg.LocalConfigPath != "" {
		t.Fatalf("did not expect a local config path: %s", cfg.LocalConfigPath)
	}
}

func TestLoadReadsGlobalConfigFileDefaults(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "pls.json")
	content := []byte(`{"provider":"ollama","model":"qwen3.5:9b","host":"http://192.168.2.166:11434","yoloMode":false}`)
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	cfg, err := Load(types.Flags{ConfigPath: configPath})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Provider != "ollama" {
		t.Fatalf("unexpected provider: %s", cfg.Provider)
	}
	if cfg.Model != "qwen3.5:9b" {
		t.Fatalf("unexpected model: %s", cfg.Model)
	}
	if cfg.Host != "http://192.168.2.166:11434" {
		t.Fatalf("unexpected host: %s", cfg.Host)
	}
	if cfg.YoloMode {
		t.Fatalf("did not expect yolo mode to be enabled")
	}
	if cfg.YoloSource != "global" {
		t.Fatalf("unexpected yolo source: %s", cfg.YoloSource)
	}
}

func TestFlagsOverrideConfigFile(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "pls.json")
	content := []byte(`{"provider":"ollama","model":"qwen3.5:9b","host":"http://192.168.2.166:11434"}`)
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	cfg, err := Load(types.Flags{
		ConfigPath: configPath,
		Model:      "qwen2.5-coder:7b-instruct-q4_K_M",
		Host:       "http://127.0.0.1:11434",
	})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Model != "qwen2.5-coder:7b-instruct-q4_K_M" {
		t.Fatalf("unexpected model: %s", cfg.Model)
	}
	if cfg.Host != "http://127.0.0.1:11434" {
		t.Fatalf("unexpected host: %s", cfg.Host)
	}
}

func TestLocalConfigOverridesGlobalConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("PLS_CONFIG", "")
	globalPath := filepath.Join(home, ".config", "pls", "config.json")
	if err := os.MkdirAll(filepath.Dir(globalPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(globalPath, []byte(`{"provider":"ollama","model":"qwen3.5:9b","host":"http://127.0.0.1:11434","yoloMode":false}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	projectRoot := filepath.Join(home, "projects", "demo")
	nested := filepath.Join(projectRoot, "subdir")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectRoot, "pls.json"), []byte(`{"model":"qwen2.5-coder:7b-instruct-q4_K_M","host":"http://192.168.2.166:11434","yoloMode":true}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	defer func() { _ = os.Chdir(previousWD) }()
	if err := os.Chdir(nested); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}

	cfg, err := Load(types.Flags{})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Model != "qwen2.5-coder:7b-instruct-q4_K_M" {
		t.Fatalf("unexpected model: %s", cfg.Model)
	}
	if cfg.Host != "http://192.168.2.166:11434" {
		t.Fatalf("unexpected host: %s", cfg.Host)
	}
	if !cfg.YoloMode {
		t.Fatalf("expected yolo mode to be enabled")
	}
	if cfg.YoloSource != "local" {
		t.Fatalf("unexpected yolo source: %s", cfg.YoloSource)
	}
	if cfg.LocalConfigPath != filepath.Join(projectRoot, "pls.json") {
		t.Fatalf("unexpected local config path: %s", cfg.LocalConfigPath)
	}
}

func TestEnvironmentOverridesLocalYoloMode(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("PLS_CONFIG", "")
	t.Setenv("PLS_YOLO_MODE", "false")

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

	cfg, err := Load(types.Flags{})
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.YoloMode {
		t.Fatalf("expected env override to disable yolo mode")
	}
	if cfg.YoloSource != "environment" {
		t.Fatalf("unexpected yolo source: %s", cfg.YoloSource)
	}
}
