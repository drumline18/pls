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
