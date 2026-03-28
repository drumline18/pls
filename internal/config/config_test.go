package config

import (
	"testing"

	"pls/internal/types"
)

func TestDefaultProviderFallsBackToOllama(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")

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
}
