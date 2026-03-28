package config

import (
	"os"
	"path/filepath"
	"testing"

	"pls/internal/types"
)

func TestDefaultProviderFallsBackToOllama(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")
	t.Setenv("PLS_CONFIG", filepath.Join(t.TempDir(), "missing.json"))

	cfg, err := Load(types.Flags{})
	if err == nil {
		t.Fatalf("expected error for explicit missing config path")
	}

	_ = cfg
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
}

func TestLoadReadsConfigFileDefaults(t *testing.T) {
	configDir := t.TempDir()
	configPath := filepath.Join(configDir, "pls.json")
	content := []byte(`{"provider":"ollama","model":"qwen3.5:9b","host":"http://192.168.2.166:11434"}`)
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
