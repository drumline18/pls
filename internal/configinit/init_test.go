package configinit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"pls/internal/types"
)

func TestRunWritesOpenAIConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")
	t.Setenv("OLLAMA_HOST", "")
	t.Setenv("PLS_OLLAMA_HOST", "")

	input := strings.NewReader("openai\n\n\nmy-test-key\ny\ny\n")
	var output strings.Builder

	if err := run(types.Flags{}, input, &output); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	path := filepath.Join(home, ".config", "pls", "config.json")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	text := string(content)
	for _, needle := range []string{"\"provider\": \"openai\"", "\"model\": \"gpt-4.1-mini\"", "\"openaiApiKey\": \"my-test-key\"", "\"yoloMode\": true"} {
		if !strings.Contains(text, needle) {
			t.Fatalf("expected config file to contain %q, got: %s", needle, text)
		}
	}
}

func TestRunCanCancelWithoutWriting(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")
	t.Setenv("OLLAMA_HOST", "")
	t.Setenv("PLS_OLLAMA_HOST", "")

	input := strings.NewReader("openai\n\n\n\n\nn\n")
	var output strings.Builder

	if err := run(types.Flags{}, input, &output); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	path := filepath.Join(home, ".config", "pls", "config.json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected config file not to exist after cancellation")
	}
}

func TestRunWritesAnthropicConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "anthropic-test-key")
	t.Setenv("OLLAMA_HOST", "")
	t.Setenv("PLS_OLLAMA_HOST", "")

	input := strings.NewReader("anthropic\n\ny\ny\n")
	var output strings.Builder

	if err := run(types.Flags{}, input, &output); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	path := filepath.Join(home, ".config", "pls", "config.json")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	text := string(content)
	for _, needle := range []string{"\"provider\": \"anthropic\"", "\"model\": \"claude-3-5-haiku-latest\"", "\"yoloMode\": true"} {
		if !strings.Contains(text, needle) {
			t.Fatalf("expected config file to contain %q, got: %s", needle, text)
		}
	}
	if strings.Contains(text, "openaiApiKey") {
		t.Fatalf("did not expect anthropic config to contain an OpenAI API key: %s", text)
	}
}

func TestRunWritesLlamaCppConfigWithAlias(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")
	t.Setenv("OLLAMA_HOST", "")
	t.Setenv("PLS_OLLAMA_HOST", "")

	input := strings.NewReader("llama.cpp\n\nmy-local-model\nn\ny\n")
	var output strings.Builder

	if err := run(types.Flags{}, input, &output); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	path := filepath.Join(home, ".config", "pls", "config.json")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	text := string(content)
	for _, needle := range []string{"\"provider\": \"llamacpp\"", "\"host\": \"http://127.0.0.1:8080/v1\"", "\"model\": \"my-local-model\"", "\"yoloMode\": false"} {
		if !strings.Contains(text, needle) {
			t.Fatalf("expected config file to contain %q, got: %s", needle, text)
		}
	}
}

func TestRunSwitchingProviderDoesNotReuseOtherProviderDefaults(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "anthropic-test-key")
	t.Setenv("OLLAMA_HOST", "")
	t.Setenv("PLS_OLLAMA_HOST", "")

	path := filepath.Join(home, ".config", "pls", "config.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte(`{"provider":"openai","model":"gpt-4.1-mini","host":"https://api.openai.com","openaiApiKey":"old-key","yoloMode":false}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	input := strings.NewReader("anthropic\n\ny\ny\n")
	var output strings.Builder

	if err := run(types.Flags{}, input, &output); err != nil {
		t.Fatalf("run returned error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	text := string(content)
	if strings.Contains(text, `"host": "https://api.openai.com"`) {
		t.Fatalf("did not expect OpenAI host to leak into Anthropi‍c config: %s", text)
	}
	if strings.Contains(text, `"model": "gpt-4.1-mini"`) {
		t.Fatalf("did not expect OpenAI model to leak into Anthropi‍c config: %s", text)
	}
	for _, needle := range []string{"\"provider\": \"anthropic\"", "\"model\": \"claude-3-5-haiku-latest\""} {
		if !strings.Contains(text, needle) {
			t.Fatalf("expected config file to contain %q, got: %s", needle, text)
		}
	}
}
