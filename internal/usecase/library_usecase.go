package usecase

import (
	"context"

	"github.com/CristianSsousa/saga-api/internal/domain"
	"github.com/CristianSsousa/saga-api/internal/repository"
)

type LibraryUsecase struct {
	libraryRepo    domain.LibraryRepository
	mediaCacheRepo *repository.MediaCacheRepo
}

func NewLibraryUsecase(libraryRepo domain.LibraryRepository, mediaCacheRepo *repository.MediaCacheRepo) *LibraryUsecase {
	return &LibraryUsecase{
		libraryRepo:    libraryRepo,
		mediaCacheRepo: mediaCacheRepo,
	}
}

type AddToLibraryInput struct {
	Media  domain.Media
	Status domain.Status
}

// Add upserts the media into media_cache, then inserts into user_library.
func (uc *LibraryUsecase) Add(ctx context.Context, userID string, input AddToLibraryInput) (*domain.LibraryItem, error) {
	mediaID, err := uc.mediaCacheRepo.Upsert(ctx, input.Media)
	if err != nil {
		return nil, err
	}
	return uc.libraryRepo.Add(ctx, userID, mediaID, input.Status)
}

type UpdateLibraryInput struct {
	Status *domain.Status
	Rating *int
	Notes  *string
}

func (uc *LibraryUsecase) Update(ctx context.Context, id, userID string, input UpdateLibraryInput) (*domain.LibraryItem, error) {
	return uc.libraryRepo.Update(ctx, id, userID, input.Status, input.Rating, input.Notes)
}

func (uc *LibraryUsecase) Delete(ctx context.Context, id, userID string) error {
	return uc.libraryRepo.SoftDelete(ctx, id, userID)
}

func (uc *LibraryUsecase) List(ctx context.Context, userID string, filter domain.ListFilter) ([]domain.LibraryItem, error) {
	return uc.libraryRepo.List(ctx, userID, filter)
}

func (uc *LibraryUsecase) GetByID(ctx context.Context, id, userID string) (*domain.LibraryItem, error) {
	return uc.libraryRepo.FindByID(ctx, id, userID)
}
