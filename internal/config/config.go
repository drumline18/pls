package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pls/internal/types"
)

type fileConfig struct {
	Provider     string `json:"provider"`
	Model        string `json:"model"`
	Host         string `json:"host"`
	OpenAIAPIKey string `json:"openaiApiKey"`
}

func Load(flags types.Flags) (types.Config, error) {
	resolvedPath, err := ResolvePath(flags.ConfigPath)
	if err != nil {
		return types.Config{}, err
	}

	fileCfg, err := loadFileConfig(resolvedPath, flags.ConfigPath != "" || os.Getenv("PLS_CONFIG") != "")
	if err != nil {
		return types.Config{}, err
	}

	provider := firstNonEmpty(flags.Provider, os.Getenv("PLS_PROVIDER"), fileCfg.Provider, detectDefaultProvider(fileCfg))
	model, err := defaultModelFor(provider)
	if err != nil {
		return types.Config{}, err
	}
	model = firstNonEmpty(flags.Model, os.Getenv("PLS_MODEL"), fileCfg.Model, model)

	host, err := defaultHostFor(provider)
	if err != nil {
		return types.Config{}, err
	}
	host = firstNonEmpty(flags.Host, os.Getenv("PLS_HOST"), providerHostEnv(provider), fileCfg.Host, host)

	return types.Config{
		Provider:     provider,
		Model:        model,
		Host:         host,
		ConfigPath:   resolvedPath,
		OutputJSON:   flags.JSON,
		OpenAIAPIKey: firstNonEmpty(os.Getenv("PLS_OPENAI_API_KEY"), os.Getenv("OPENAI_API_KEY"), fileCfg.OpenAIAPIKey),
	}, nil
}

func ResolvePath(override string) (string, error) {
	if override != "" {
		return expandHome(override)
	}
	if envPath := os.Getenv("PLS_CONFIG"); envPath != "" {
		return expandHome(envPath)
	}

	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configHome = filepath.Join(home, ".config")
	}

	return filepath.Join(configHome, "pls", "config.json"), nil
}

func detectDefaultProvider(fileCfg fileConfig) string {
	if os.Getenv("PLS_OPENAI_API_KEY") != "" || os.Getenv("OPENAI_API_KEY") != "" || fileCfg.OpenAIAPIKey != "" {
		return "openai"
	}
	if fileCfg.Provider != "" {
		return fileCfg.Provider
	}
	return "ollama"
}

func defaultModelFor(provider string) (string, error) {
	switch provider {
	case "openai":
		return "gpt-4.1-mini", nil
	case "ollama":
		return "qwen2.5-coder:7b-instruct-q4_K_M", nil
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}

func defaultHostFor(provider string) (string, error) {
	switch provider {
	case "openai":
		return "https://api.openai.com", nil
	case "ollama":
		return "http://127.0.0.1:11434", nil
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}

func providerHostEnv(provider string) string {
	switch provider {
	case "ollama":
		return firstNonEmpty(os.Getenv("PLS_OLLAMA_HOST"), os.Getenv("OLLAMA_HOST"))
	default:
		return ""
	}
}

func loadFileConfig(path string, required bool) (fileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) && !required {
			return fileConfig{}, nil
		}
		return fileConfig{}, err
	}

	var cfg fileConfig
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return fileConfig{}, fmt.Errorf("invalid config file %s: %w", path, err)
	}

	return cfg, nil
}

func expandHome(path string) (string, error) {
	if path == "~" {
		return os.UserHomeDir()
	}
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
