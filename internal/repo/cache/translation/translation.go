// Package translation implements the Redis-backed Translation cache.
package translation

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/leijux/go-clean-template/internal/entity"
	"github.com/leijux/go-clean-template/internal/repo"
	"github.com/redis/go-redis/v9"
)

const (
	_keyPrefix = "translation"
	_ttl       = 24 * time.Hour
)

// Repo -.
type Repo struct {
	client *redis.Client
}

// New returns a Translation cache instrumented with OpenTelemetry tracing spans.
func New(client *redis.Client) repo.TranslationCache {
	return newTraced(&Repo{client: client})
}

// Get -. Returns the cached translation and whether it was found.
func (r *Repo) Get(ctx context.Context, userID string, t entity.Translation) (entity.Translation, bool, error) {
	raw, err := r.client.Get(ctx, cacheKey(userID, t)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return entity.Translation{}, false, nil
		}

		return entity.Translation{}, false, fmt.Errorf("TranslationCache - Get - r.client.Get: %w", err)
	}

	var cached entity.Translation

	err = json.Unmarshal(raw, &cached)
	if err != nil {
		return entity.Translation{}, false, fmt.Errorf("TranslationCache - Get - json.Unmarshal: %w", err)
	}

	return cached, true, nil
}

// Set -. Stores the translation for the default TTL.
func (r *Repo) Set(ctx context.Context, userID string, t entity.Translation) error {
	raw, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("TranslationCache - Set - json.Marshal: %w", err)
	}

	err = r.client.Set(ctx, cacheKey(userID, t), raw, _ttl).Err()
	if err != nil {
		return fmt.Errorf("TranslationCache - Set - r.client.Set: %w", err)
	}

	return nil
}

// cacheKey is content-addressed: the same original text for the same language
// pair always maps to the same entry, so no invalidation is needed.
func cacheKey(userID string, t entity.Translation) string {
	sum := sha256.Sum256([]byte(t.Original))

	return _keyPrefix + ":" + userID + ":" + t.Source + ":" + t.Destination + ":" + hex.EncodeToString(sum[:])
}
