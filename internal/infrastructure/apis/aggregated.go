package apis

import (
	"context"
	"fmt"

	"github.com/CristianSsousa/saga-api/internal/domain"
)

type Config struct {
	TMDBApiKey string
	RAWGApiKey string
}

type AggregatedSearchService struct {
	tmdb        *TMDBClient
	rawg        *RAWGClient
	openLibrary *OpenLibraryClient
	aniList     *AniListClient
}

func NewAggregatedSearchService(cfg Config) *AggregatedSearchService {
	return &AggregatedSearchService{
		tmdb:        NewTMDBClient(cfg.TMDBApiKey),
		rawg:        NewRAWGClient(cfg.RAWGApiKey),
		openLibrary: NewOpenLibraryClient(),
		aniList:     NewAniListClient(),
	}
}

func (s *AggregatedSearchService) Search(ctx context.Context, query string, mediaType domain.MediaType) ([]domain.Media, error) {
	switch mediaType {
	case domain.MediaTypeMovie, domain.MediaTypeTV:
		return s.tmdb.Search(ctx, query, mediaType)
	case domain.MediaTypeGame:
		return s.rawg.Search(ctx, query, mediaType)
	case domain.MediaTypeBook:
		return s.openLibrary.Search(ctx, query, mediaType)
	case domain.MediaTypeManga:
		return s.aniList.Search(ctx, query, mediaType)
	default:
		return nil, fmt.Errorf("unsupported media type: %s", mediaType)
	}
}
