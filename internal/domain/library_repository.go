package domain

import "context"

type ListFilter struct {
	MediaType *MediaType
	Status    *Status
	Cursor    string
	Limit     int
}

type LibraryRepository interface {
	Add(ctx context.Context, userID string, mediaID string, status Status) (*LibraryItem, error)
	Update(ctx context.Context, id string, userID string, status *Status, rating *int, notes *string) (*LibraryItem, error)
	SoftDelete(ctx context.Context, id string, userID string) error
	List(ctx context.Context, userID string, filter ListFilter) ([]LibraryItem, error)
	FindByID(ctx context.Context, id string, userID string) (*LibraryItem, error)
}
