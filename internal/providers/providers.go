package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	anyllm "github.com/mozilla-ai/any-llm-go"
	"github.com/mozilla-ai/any-llm-go/providers/anthropic"
	"github.com/mozilla-ai/any-llm-go/providers/deepseek"
	"github.com/mozilla-ai/any-llm-go/providers/gemini"
	"github.com/mozilla-ai/any-llm-go/providers/groq"
	"github.com/mozilla-ai/any-llm-go/providers/llamacpp"
	"github.com/mozilla-ai/any-llm-go/providers/llamafile"
	"github.com/mozilla-ai/any-llm-go/providers/mistral"
	"github.com/mozilla-ai/any-llm-go/providers/ollama"
	"github.com/mozilla-ai/any-llm-go/providers/openai"
	"github.com/mozilla-ai/any-llm-go/providers/zai"
	"github.com/drumline18/pls/internal/types"
	"github.com/drumline18/pls/internal/util"
)

func Generate(ctx context.Context, cfg types.Config, messages types.Messages) (types.Suggestion, error) {
	provider, err := buildProvider(cfg)
	if err != nil {
		return types.Suggestion{}, err
	}

	temperature := 0.1
	completion, err := provider.Completion(ctx, anyllm.CompletionParams{
		Model:       cfg.Model,
		Temperature: &temperature,
		Messages: []anyllm.Message{
			{Role: anyllm.RoleSystem, Content: messages.System},
			{Role: anyllm.RoleUser, Content: util.MustJSON(messages.User)},
		},
		ResponseFormat: suggestionResponseFormat(),
	})
	if err != nil {
		return types.Suggestion{}, normalizeProviderError(cfg.Provider, err)
	}

	content, err := completionContent(completion)
	if err != nil {
		return types.Suggestion{}, err
	}

	jsonBody, err := util.ExtractJSONObject(content)
	if err != nil {
		return types.Suggestion{}, err
	}

	var suggestion types.Suggestion
	if err := json.Unmarshal([]byte(jsonBody), &suggestion); err != nil {
		return types.Suggestion{}, err
	}

	return suggestion, nil
}

func buildProvider(cfg types.Config) (anyllm.Provider, error) {
	providerName := strings.ToLower(strings.TrimSpace(cfg.Provider))
	opts := make([]anyllm.Option, 0, 2)
	if strings.TrimSpace(cfg.Host) != "" {
		opts = append(opts, anyllm.WithBaseURL(cfg.Host))
	}
	if providerName == "openai" && strings.TrimSpace(cfg.OpenAIAPIKey) != "" {
		opts = append(opts, anyllm.WithAPIKey(cfg.OpenAIAPIKey))
	}

	switch providerName {
	case "openai":
		return openai.New(opts...)
	case "ollama":
		return ollama.New(opts...)
	case "anthropic":
		return anthropic.New(opts...)
	case "gemini":
		return gemini.New(opts...)
	case "groq":
		return groq.New(opts...)
	case "deepseek":
		return deepseek.New(opts...)
	case "mistral":
		return mistral.New(opts...)
	case "zai":
		return zai.New(opts...)
	case "llamacpp":
		return llamacpp.New(opts...)
	case "llamafile":
		return llamafile.New(opts...)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}

func suggestionResponseFormat() *anyllm.ResponseFormat {
	strict := false
	return &anyllm.ResponseFormat{
		Type: "json_schema",
		JSONSchema: &anyllm.JSONSchema{
			Name:        "pls_suggestion",
			Description: "Shell command suggestion payload for the pls CLI.",
			Strict:      &strict,
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"command": map[string]any{"type": "string"},
					"explanation": map[string]any{"type": "string"},
					"risk": map[string]any{"type": "string", "enum": []string{"low", "medium", "high", "critical"}},
					"requiresConfirmation": map[string]any{"type": "boolean"},
					"needsClarification": map[string]any{"type": "boolean"},
					"clarificationQuestion": map[string]any{"type": "string"},
					"notes": map[string]any{"type": "string"},
					"platform": map[string]any{"type": "string"},
					"refused": map[string]any{"type": "boolean"},
				},
				"required": []string{"risk", "requiresConfirmation", "needsClarification", "refused"},
				"additionalProperties": false,
			},
		},
	}
}

func completionContent(completion *anyllm.ChatCompletion) (string, error) {
	if completion == nil || len(completion.Choices) == 0 {
		return "", fmt.Errorf("provider response did not contain any choices")
	}

	content := completion.Choices[0].Message.Content
	switch value := content.(type) {
	case string:
		if strings.TrimSpace(value) == "" {
			return "", fmt.Errorf("provider response did not contain message content")
		}
		return value, nil
	case []anyllm.ContentPart:
		var builder strings.Builder
		for _, part := range value {
			if part.Type == "text" {
				builder.WriteString(part.Text)
			}
		}
		if strings.TrimSpace(builder.String()) == "" {
			return "", fmt.Errorf("provider response did not contain text content")
		}
		return builder.String(), nil
	default:
		bytes, err := json.Marshal(content)
		if err != nil {
			return "", fmt.Errorf("provider response had unsupported content type %T", content)
		}
		if strings.TrimSpace(string(bytes)) == "" || string(bytes) == "null" {
			return "", fmt.Errorf("provider response did not contain usable content")
		}
		return string(bytes), nil
	}
}

func normalizeProviderError(provider string, err error) error {
	if err == nil {
		return nil
	}
	message := err.Error()
	if provider == "openai" && strings.Contains(strings.ToLower(message), "api key") {
		return fmt.Errorf("OPENAI_API_KEY or PLS_OPENAI_API_KEY is required for provider=openai")
	}
	return err
}
