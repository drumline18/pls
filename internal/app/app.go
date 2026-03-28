package app

import (
	"context"

	"pls/internal/policy"
	"pls/internal/prompt"
	"pls/internal/providers"
	"pls/internal/style"
	"pls/internal/types"
)

func GenerateSuggestion(ctx context.Context, request string, runtimeContext types.RuntimeContext, cfg types.Config) (types.Suggestion, error) {
	messages := prompt.Build(request, runtimeContext)
	raw, err := providers.Generate(ctx, cfg, messages)
	if err != nil {
		return types.Suggestion{}, err
	}

	validated, err := types.ValidateSuggestion(raw)
	if err != nil {
		return types.Suggestion{}, err
	}

	normalized := style.Normalize(request, runtimeContext, validated)
	return policy.Apply(normalized), nil
}
