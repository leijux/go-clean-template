package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/leijux/go-clean-template/internal/entity"
	"github.com/leijux/go-clean-template/internal/usecase"
	"github.com/leijux/go-clean-template/internal/usecase/translation"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errInternalServErr = errors.New("internal server error")

func newTranslationUseCase(
	t *testing.T,
) (usecase.Translation, *MockTranslationRepo, *MockTranslationWebAPI, *MockTranslationCache) {
	t.Helper()

	ctrl := gomock.NewController(t)

	repo := NewMockTranslationRepo(ctrl)
	webAPI := NewMockTranslationWebAPI(ctrl)
	cache := NewMockTranslationCache(ctrl)

	useCase := translation.New(repo, webAPI, cache)

	return useCase, repo, webAPI, cache
}

func TestHistory(t *testing.T) {
	t.Parallel()

	t.Run("empty result", func(t *testing.T) {
		t.Parallel()

		uc, repo, _, _ := newTranslationUseCase(t)
		repo.EXPECT().GetHistory(gomock.Any(), "").Return(nil, nil)

		res, err := uc.History(context.Background(), "")

		require.Equal(t, entity.TranslationHistory{}, res)
		require.NoError(t, err)
	})

	t.Run("result with error", func(t *testing.T) {
		t.Parallel()

		uc, repo, _, _ := newTranslationUseCase(t)
		repo.EXPECT().GetHistory(gomock.Any(), "").Return(nil, errInternalServErr)

		res, err := uc.History(context.Background(), "")

		require.Equal(t, entity.TranslationHistory{}, res)
		require.ErrorIs(t, err, errInternalServErr)
	})
}

func TestTranslateCacheHit(t *testing.T) {
	t.Parallel()

	uc, _, _, cache := newTranslationUseCase(t)
	cache.EXPECT().Get(gomock.Any(), "", entity.Translation{}).
		Return(entity.Translation{Translation: "cached"}, true, nil)

	res, err := uc.Translate(context.Background(), "", entity.Translation{})

	require.EqualValues(t, entity.Translation{Translation: "cached"}, res)
	require.NoError(t, err)
}

func TestTranslateCacheMiss(t *testing.T) {
	t.Parallel()

	uc, repo, webAPI, cache := newTranslationUseCase(t)
	cache.EXPECT().Get(gomock.Any(), "", entity.Translation{}).Return(entity.Translation{}, false, nil)
	webAPI.EXPECT().Translate(gomock.Any(), entity.Translation{}).Return(entity.Translation{}, nil)
	repo.EXPECT().Store(gomock.Any(), "", entity.Translation{}).Return(nil)
	cache.EXPECT().Set(gomock.Any(), "", entity.Translation{}).Return(nil)

	res, err := uc.Translate(context.Background(), "", entity.Translation{})

	require.EqualValues(t, entity.Translation{}, res)
	require.NoError(t, err)
}

func TestTranslateCacheGetError(t *testing.T) {
	t.Parallel()

	uc, _, _, cache := newTranslationUseCase(t)
	cache.EXPECT().Get(gomock.Any(), "", entity.Translation{}).
		Return(entity.Translation{}, false, errInternalServErr)

	res, err := uc.Translate(context.Background(), "", entity.Translation{})

	require.EqualValues(t, entity.Translation{}, res)
	require.ErrorIs(t, err, errInternalServErr)
}

func TestTranslateWebAPIError(t *testing.T) {
	t.Parallel()

	uc, _, webAPI, cache := newTranslationUseCase(t)
	cache.EXPECT().Get(gomock.Any(), "", entity.Translation{}).Return(entity.Translation{}, false, nil)
	webAPI.EXPECT().Translate(gomock.Any(), entity.Translation{}).Return(entity.Translation{}, errInternalServErr)

	res, err := uc.Translate(context.Background(), "", entity.Translation{})

	require.EqualValues(t, entity.Translation{}, res)
	require.ErrorIs(t, err, errInternalServErr)
}

func TestTranslateRepoError(t *testing.T) {
	t.Parallel()

	uc, repo, webAPI, cache := newTranslationUseCase(t)
	cache.EXPECT().Get(gomock.Any(), "", entity.Translation{}).Return(entity.Translation{}, false, nil)
	webAPI.EXPECT().Translate(gomock.Any(), entity.Translation{}).Return(entity.Translation{}, nil)
	repo.EXPECT().Store(gomock.Any(), "", entity.Translation{}).Return(errInternalServErr)

	res, err := uc.Translate(context.Background(), "", entity.Translation{})

	require.EqualValues(t, entity.Translation{}, res)
	require.ErrorIs(t, err, errInternalServErr)
}

func TestTranslateCacheSetError(t *testing.T) {
	t.Parallel()

	uc, repo, webAPI, cache := newTranslationUseCase(t)
	cache.EXPECT().Get(gomock.Any(), "", entity.Translation{}).Return(entity.Translation{}, false, nil)
	webAPI.EXPECT().Translate(gomock.Any(), entity.Translation{}).Return(entity.Translation{}, nil)
	repo.EXPECT().Store(gomock.Any(), "", entity.Translation{}).Return(nil)
	cache.EXPECT().Set(gomock.Any(), "", entity.Translation{}).Return(errInternalServErr)

	res, err := uc.Translate(context.Background(), "", entity.Translation{})

	require.EqualValues(t, entity.Translation{}, res)
	require.ErrorIs(t, err, errInternalServErr)
}
