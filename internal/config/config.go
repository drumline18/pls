package config

import (
	"fmt"
	"os"

	"pls/internal/types"
)

func Load(flags types.Flags) (types.Config, error) {
	provider := firstNonEmpty(flags.Provider, os.Getenv("PLS_PROVIDER"), detectDefaultProvider())
	model, err := defaultModelFor(provider)
	if err != nil {
		return types.Config{}, err
	}
	model = firstNonEmpty(flags.Model, os.Getenv("PLS_MODEL"), model)

	host, err := defaultHostFor(provider)
	if err != nil {
		return types.Config{}, err
	}
	host = firstNonEmpty(flags.Host, host)

	return types.Config{
		Provider:     provider,
		Model:        model,
		Host:         host,
		OutputJSON:   flags.JSON,
		OpenAIAPIKey: firstNonEmpty(os.Getenv("PLS_OPENAI_API_KEY"), os.Getenv("OPENAI_API_KEY")),
	}, nil
}

func detectDefaultProvider() string {
	if os.Getenv("PLS_OPENAI_API_KEY") != "" || os.Getenv("OPENAI_API_KEY") != "" {
		return "openai"
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
		return firstNonEmpty(os.Getenv("PLS_OLLAMA_HOST"), os.Getenv("OLLAMA_HOST"), "http://127.0.0.1:11434"), nil
	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
