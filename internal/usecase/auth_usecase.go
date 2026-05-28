package usecase

import (
	"context"

	"github.com/CristianSsousa/saga-api/internal/domain"
	"github.com/CristianSsousa/saga-api/pkg/hash"
	jwtpkg "github.com/CristianSsousa/saga-api/pkg/jwt"
)

type AuthUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

func NewAuthUsecase(userRepo domain.UserRepository, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{userRepo: userRepo, jwtSecret: jwtSecret}
}

type RegisterInput struct {
	Email    string
	Username string
	Password string
}

type AuthOutput struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

func (uc *AuthUsecase) Register(ctx context.Context, input RegisterInput) (*AuthOutput, error) {
	existing, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err == nil && existing != nil {
		return nil, domain.ErrEmailAlreadyInUse
	}
	hashed, err := hash.Hash(input.Password)
	if err != nil {
		return nil, err
	}
	user := &domain.User{Email: input.Email, Username: input.Username, PasswordHash: hashed}
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	created, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	token, err := jwtpkg.Generate(created.ID, uc.jwtSecret)
	if err != nil {
		return nil, err
	}
	return &AuthOutput{Token: token, User: *created}, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, email, password string) (*AuthOutput, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	if !hash.Compare(user.PasswordHash, password) {
		return nil, domain.ErrInvalidCredentials
	}
	token, err := jwtpkg.Generate(user.ID, uc.jwtSecret)
	if err != nil {
		return nil, err
	}
	return &AuthOutput{Token: token, User: *user}, nil
}
