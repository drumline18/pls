package app

import (
	"context"

	"github.com/drumline18/pls/internal/policy"
	"github.com/drumline18/pls/internal/prompt"
	"github.com/drumline18/pls/internal/providers"
	"github.com/drumline18/pls/internal/style"
	"github.com/drumline18/pls/internal/types"
)

func GenerateSuggestion(ctx context.Context, request string, runtimeContext types.RuntimeContext, cfg types.Config) (types.Suggestion, error) {
	if direct, ok := style.DirectSuggestion(request, runtimeContext); ok {
		return policy.Apply(direct), nil
	}

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
