package domain

import "context"

type SearchService interface {
	Search(ctx context.Context, query string, mediaType MediaType) ([]Media, error)
}
