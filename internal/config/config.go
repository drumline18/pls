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
	YoloMode     *bool  `json:"yoloMode"`
}

func Load(flags types.Flags) (types.Config, error) {
	globalPath, err := ResolvePath(flags.ConfigPath)
	if err != nil {
		return types.Config{}, err
	}

	globalCfg, err := loadFileConfig(globalPath, flags.ConfigPath != "" || os.Getenv("PLS_CONFIG") != "")
	if err != nil {
		return types.Config{}, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return types.Config{}, err
	}

	localPath, err := FindLocalConfigPath(cwd)
	if err != nil {
		return types.Config{}, err
	}

	localCfg, err := loadOptionalFileConfig(localPath)
	if err != nil {
		return types.Config{}, err
	}

	openAIAPIKey := firstNonEmpty(
		os.Getenv("PLS_OPENAI_API_KEY"),
		os.Getenv("OPENAI_API_KEY"),
		localCfg.OpenAIAPIKey,
		globalCfg.OpenAIAPIKey,
	)

	provider := firstNonEmpty(flags.Provider, os.Getenv("PLS_PROVIDER"), localCfg.Provider, globalCfg.Provider, detectDefaultProvider(openAIAPIKey, localCfg.Provider, globalCfg.Provider))
	model, err := defaultModelFor(provider)
	if err != nil {
		return types.Config{}, err
	}
	model = firstNonEmpty(flags.Model, os.Getenv("PLS_MODEL"), localCfg.Model, globalCfg.Model, model)

	host, err := defaultHostFor(provider)
	if err != nil {
		return types.Config{}, err
	}
	host = firstNonEmpty(flags.Host, os.Getenv("PLS_HOST"), providerHostEnv(provider), localCfg.Host, globalCfg.Host, host)

	yoloMode, yoloSource, err := resolveYoloMode(localCfg, globalCfg)
	if err != nil {
		return types.Config{}, err
	}

	return types.Config{
		Provider:       provider,
		Model:          model,
		Host:           host,
		ConfigPath:     globalPath,
		LocalConfigPath: localPath,
		YoloMode:       yoloMode,
		YoloSource:     yoloSource,
		OutputJSON:     flags.JSON,
		OpenAIAPIKey:   openAIAPIKey,
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

func FindLocalConfigPath(startDir string) (string, error) {
	dir := startDir
	for {
		candidate := filepath.Join(dir, "pls.json")
		info, err := os.Stat(candidate)
		if err == nil && !info.IsDir() {
			return candidate, nil
		}
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", nil
		}
		dir = parent
	}
}

func detectDefaultProvider(openAIAPIKey, localProvider, globalProvider string) string {
	if openAIAPIKey != "" {
		return "openai"
	}
	if localProvider != "" {
		return localProvider
	}
	if globalProvider != "" {
		return globalProvider
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

func resolveYoloMode(localCfg, globalCfg fileConfig) (bool, string, error) {
	if value, ok, err := envBool("PLS_YOLO_MODE"); err != nil {
		return false, "", err
	} else if ok {
		return value, "environment", nil
	}

	if localCfg.YoloMode != nil {
		return *localCfg.YoloMode, "local", nil
	}
	if globalCfg.YoloMode != nil {
		return *globalCfg.YoloMode, "global", nil
	}
	return false, "default", nil
}

func envBool(name string) (bool, bool, error) {
	value, ok := os.LookupEnv(name)
	if !ok || strings.TrimSpace(value) == "" {
		return false, false, nil
	}

	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "on":
		return true, true, nil
	case "0", "false", "no", "off":
		return false, true, nil
	default:
		return false, false, fmt.Errorf("invalid boolean value for %s: %s", name, value)
	}
}

func loadOptionalFileConfig(path string) (fileConfig, error) {
	if path == "" {
		return fileConfig{}, nil
	}
	return loadFileConfig(path, false)
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
