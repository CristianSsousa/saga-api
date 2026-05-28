package usecase

import (
	"context"

	"github.com/CristianSsousa/saga-api/internal/domain"
	"github.com/CristianSsousa/saga-api/internal/repository"
)

type SearchUsecase struct {
	searchSvc      domain.SearchService
	mediaCacheRepo *repository.MediaCacheRepo
}

func NewSearchUsecase(searchSvc domain.SearchService, mediaCacheRepo *repository.MediaCacheRepo) *SearchUsecase {
	return &SearchUsecase{searchSvc: searchSvc, mediaCacheRepo: mediaCacheRepo}
}

func (uc *SearchUsecase) Search(ctx context.Context, query string, mediaType domain.MediaType) ([]domain.Media, error) {
	return uc.searchSvc.Search(ctx, query, mediaType)
}
