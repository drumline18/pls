package configinit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/drumline18/pls/internal/types"
)

func TestRunLocalWritesProjectOverrides(t *testing.T) {
	home := t.TempDir()
	project := filepath.Join(home, "project")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("PLS_CONFIG", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")
	t.Setenv("OLLAMA_HOST", "")
	t.Setenv("PLS_OLLAMA_HOST", "")
	t.Setenv("PLS_YOLO_MODE", "")

	globalPath := filepath.Join(home, ".config", "pls", "config.json")
	if err := os.MkdirAll(filepath.Dir(globalPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(globalPath, []byte(`{"provider":"ollama","model":"qwen2.5-coder:7b-instruct-q4_K_M","host":"http://192.168.2.166:11434","yoloMode":false}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	defer func() { _ = os.Chdir(previousWD) }()
	if err := os.Chdir(project); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}

	input := strings.NewReader("n\nn\ny\nqwen3.5:14b\non\ny\n")
	var output strings.Builder
	if err := runLocal(types.Flags{}, input, &output); err != nil {
		t.Fatalf("runLocal returned error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(project, "pls.json"))
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	text := string(content)
	for _, needle := range []string{"\"model\": \"qwen3.5:14b\"", "\"yoloMode\": true"} {
		if !strings.Contains(text, needle) {
			t.Fatalf("expected local config to contain %q, got: %s", needle, text)
		}
	}
	if strings.Contains(text, "openaiApiKey") {
		t.Fatalf("did not expect local config to contain an API key: %s", text)
	}
}

func TestRunLocalCanRemoveExistingOverrides(t *testing.T) {
	home := t.TempDir()
	project := filepath.Join(home, "project")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("PLS_CONFIG", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")
	t.Setenv("OLLAMA_HOST", "")
	t.Setenv("PLS_OLLAMA_HOST", "")
	t.Setenv("PLS_YOLO_MODE", "")

	globalPath := filepath.Join(home, ".config", "pls", "config.json")
	if err := os.MkdirAll(filepath.Dir(globalPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(globalPath, []byte(`{"provider":"ollama","model":"qwen2.5-coder:7b-instruct-q4_K_M","host":"http://192.168.2.166:11434","yoloMode":false}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(project, "pls.json"), []byte(`{"model":"qwen3.5:14b","yoloMode":true}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	defer func() { _ = os.Chdir(previousWD) }()
	if err := os.Chdir(project); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}

	input := strings.NewReader("n\nn\nn\ninherit\ny\n")
	var output strings.Builder
	if err := runLocal(types.Flags{}, input, &output); err != nil {
		t.Fatalf("runLocal returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(project, "pls.json")); !os.IsNotExist(err) {
		t.Fatalf("expected local config file to be removed")
	}
}

func TestRunLocalCanOverrideProviderWithAlias(t *testing.T) {
	home := t.TempDir()
	project := filepath.Join(home, "project")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("PLS_CONFIG", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("PLS_OPENAI_API_KEY", "")
	t.Setenv("OLLAMA_HOST", "")
	t.Setenv("PLS_OLLAMA_HOST", "")
	t.Setenv("PLS_YOLO_MODE", "")

	globalPath := filepath.Join(home, ".config", "pls", "config.json")
	if err := os.MkdirAll(filepath.Dir(globalPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(globalPath, []byte(`{"provider":"openai","model":"gpt-4.1-mini","host":"https://api.openai.com","yoloMode":false}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	defer func() { _ = os.Chdir(previousWD) }()
	if err := os.Chdir(project); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}

	input := strings.NewReader("y\nz.ai\ny\nhttps://api.z.ai/api/paas/v4\ny\nglm-4.5-air\ninherit\ny\n")
	var output strings.Builder
	if err := runLocal(types.Flags{}, input, &output); err != nil {
		t.Fatalf("runLocal returned error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(project, "pls.json"))
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	text := string(content)
	for _, needle := range []string{"\"provider\": \"zai\"", "\"host\": \"https://api.z.ai/api/paas/v4\"", "\"model\": \"glm-4.5-air\""} {
		if !strings.Contains(text, needle) {
			t.Fatalf("expected local config to contain %q, got: %s", needle, text)
		}
	}
	if strings.Contains(text, "openaiApiKey") {
		t.Fatalf("did not expect local config to contain an API key: %s", text)
	}
}

func TestRunLocalSwitchingProviderDoesNotReuseCurrentProviderDefaults(t *testing.T) {
	home := t.TempDir()
	project := filepath.Join(home, "project")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
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
	if err := os.WriteFile(globalPath, []byte(`{"provider":"openai","model":"gpt-4.1-mini","host":"https://api.openai.com","yoloMode":false}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(project, "pls.json"), []byte(`{"host":"https://api.openai.com","model":"gpt-4.1-mini"}`), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}
	defer func() { _ = os.Chdir(previousWD) }()
	if err := os.Chdir(project); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}

	input := strings.NewReader("y\nanthropic\nn\ny\nclaude-3-5-haiku-latest\ninherit\ny\n")
	var output strings.Builder
	if err := runLocal(types.Flags{}, input, &output); err != nil {
		t.Fatalf("runLocal returned error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(project, "pls.json"))
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	text := string(content)
	if strings.Contains(text, `"host": "https://api.openai.com"`) {
		t.Fatalf("did not expect OpenAI host to leak into Anthropi‍c local config: %s", text)
	}
	if strings.Contains(text, `"model": "gpt-4.1-mini"`) {
		t.Fatalf("did not expect OpenAI model to leak into Anthropi‍c local config: %s", text)
	}
	for _, needle := range []string{"\"provider\": \"anthropic\"", "\"model\": \"claude-3-5-haiku-latest\""} {
		if !strings.Contains(text, needle) {
			t.Fatalf("expected local config to contain %q, got: %s", needle, text)
		}
	}
}
