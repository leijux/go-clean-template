package translation

import (
	"context"

	"github.com/leijux/go-clean-template/internal/entity"
	"github.com/leijux/go-clean-template/internal/repo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/leijux/go-clean-template/internal/repo/cache/translation"

// tracedRepo wraps a TranslationCache with OpenTelemetry spans, giving a
// semantic "TranslationCache.<Method>" view on top of the redis client spans.
type tracedRepo struct {
	next repo.TranslationCache
}

func newTraced(next repo.TranslationCache) repo.TranslationCache {
	return &tracedRepo{next: next}
}

func startSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return otel.Tracer(_tracerName).Start(ctx, name, trace.WithAttributes(attrs...))
}

func endSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()
}

func (r *tracedRepo) Get(ctx context.Context, userID string, t entity.Translation) (entity.Translation, bool, error) {
	ctx, span := startSpan(
		ctx, "TranslationCache.Get",
		attribute.String("user.id", userID),
		attribute.String("translation.source", t.Source),
		attribute.String("translation.destination", t.Destination),
	)

	cached, found, err := r.next.Get(ctx, userID, t)
	endSpan(span, err)

	return cached, found, err
}

func (r *tracedRepo) Set(ctx context.Context, userID string, t entity.Translation) error {
	ctx, span := startSpan(
		ctx, "TranslationCache.Set",
		attribute.String("user.id", userID),
		attribute.String("translation.source", t.Source),
		attribute.String("translation.destination", t.Destination),
	)

	err := r.next.Set(ctx, userID, t)
	endSpan(span, err)

	return err
}
