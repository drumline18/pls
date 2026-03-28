package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"pls/internal/types"
	"pls/internal/util"
)

func Generate(ctx context.Context, cfg types.Config, messages types.Messages) (types.Suggestion, error) {
	switch cfg.Provider {
	case "openai":
		return generateWithOpenAI(ctx, cfg, messages)
	case "ollama":
		return generateWithOllama(ctx, cfg, messages)
	default:
		return types.Suggestion{}, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}

func generateWithOpenAI(ctx context.Context, cfg types.Config, messages types.Messages) (types.Suggestion, error) {
	if cfg.OpenAIAPIKey == "" {
		return types.Suggestion{}, fmt.Errorf("OPENAI_API_KEY or PLS_OPENAI_API_KEY is required for provider=openai")
	}

	payload := map[string]any{
		"model":       cfg.Model,
		"temperature": 0.1,
		"response_format": map[string]any{
			"type": "json_object",
		},
		"messages": []map[string]string{
			{"role": "system", "content": messages.System},
			{"role": "user", "content": util.MustJSON(messages.User)},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return types.Suggestion{}, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(cfg.Host, "/")+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return types.Suggestion{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+cfg.OpenAIAPIKey)

	response, err := httpClient().Do(request)
	if err != nil {
		return types.Suggestion{}, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return types.Suggestion{}, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return types.Suggestion{}, fmt.Errorf("OpenAI request failed (%d): %s", response.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return types.Suggestion{}, err
	}
	if len(parsed.Choices) == 0 || strings.TrimSpace(parsed.Choices[0].Message.Content) == "" {
		return types.Suggestion{}, fmt.Errorf("OpenAI response did not contain message content")
	}

	jsonBody, err := util.ExtractJSONObject(parsed.Choices[0].Message.Content)
	if err != nil {
		return types.Suggestion{}, err
	}

	var suggestion types.Suggestion
	if err := json.Unmarshal([]byte(jsonBody), &suggestion); err != nil {
		return types.Suggestion{}, err
	}

	return suggestion, nil
}

func generateWithOllama(ctx context.Context, cfg types.Config, messages types.Messages) (types.Suggestion, error) {
	payload := map[string]any{
		"model":  cfg.Model,
		"format": "json",
		"stream": false,
		"options": map[string]any{
			"temperature": 0.1,
		},
		"messages": []map[string]string{
			{"role": "system", "content": messages.System},
			{"role": "user", "content": util.MustJSON(messages.User)},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return types.Suggestion{}, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(cfg.Host, "/")+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return types.Suggestion{}, err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := httpClient().Do(request)
	if err != nil {
		return types.Suggestion{}, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return types.Suggestion{}, err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return types.Suggestion{}, fmt.Errorf("Ollama request failed (%d): %s", response.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var parsed struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return types.Suggestion{}, err
	}
	if strings.TrimSpace(parsed.Message.Content) == "" {
		return types.Suggestion{}, fmt.Errorf("Ollama response did not contain message content")
	}

	var suggestion types.Suggestion
	if err := json.Unmarshal([]byte(parsed.Message.Content), &suggestion); err != nil {
		return types.Suggestion{}, err
	}

	return suggestion, nil
}

func httpClient() *http.Client {
	return &http.Client{Timeout: 45 * time.Second}
}
