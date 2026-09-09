package translation

import (
	"context"
	"fmt"

	"github.com/leijux/go-clean-template/internal/entity"
	"github.com/leijux/go-clean-template/internal/repo"
	"github.com/leijux/go-clean-template/internal/usecase"
)

// UseCase -.
type UseCase struct {
	repo   repo.TranslationRepo
	webAPI repo.TranslationWebAPI
	cache  repo.TranslationCache
}

// New returns a Translation usecase instrumented with OpenTelemetry tracing spans.
func New(r repo.TranslationRepo, w repo.TranslationWebAPI, c repo.TranslationCache) usecase.Translation {
	return newTraced(&UseCase{
		repo:   r,
		webAPI: w,
		cache:  c,
	})
}

// History - getting translate history from store.
func (uc *UseCase) History(ctx context.Context, userID string) (entity.TranslationHistory, error) {
	translations, err := uc.repo.GetHistory(ctx, userID)
	if err != nil {
		return entity.TranslationHistory{}, fmt.Errorf("TranslationUseCase - History - s.repo.GetHistory: %w", err)
	}

	return entity.TranslationHistory{History: translations}, nil
}

// Translate -. Serves the result from cache when the same original text was
// already translated for the same language pair, otherwise calls the external
// API and stores the result both in Postgres and in the cache.
func (uc *UseCase) Translate(ctx context.Context, userID string, t entity.Translation) (entity.Translation, error) {
	cached, found, err := uc.cache.Get(ctx, userID, t)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - uc.cache.Get: %w", err)
	}

	if found {
		return cached, nil
	}

	translation, err := uc.webAPI.Translate(ctx, t)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - s.webAPI.Translate: %w", err)
	}

	err = uc.repo.Store(ctx, userID, translation)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - s.repo.Store: %w", err)
	}

	err = uc.cache.Set(ctx, userID, translation)
	if err != nil {
		return entity.Translation{}, fmt.Errorf("TranslationUseCase - Translate - uc.cache.Set: %w", err)
	}

	return translation, nil
}
