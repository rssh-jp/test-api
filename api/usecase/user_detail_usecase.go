package usecase

import (
	"context"

	"github.com/rssh-jp/test-api/api/domain"
)

type UserDetailUsecase interface {
	GetUserDetailByID(ctx context.Context, id int64) (*domain.UserDetail, error)
	GetUserDetailByUsername(ctx context.Context, username string) (*domain.UserDetail, error)
}

type userDetailUsecase struct {
	repo domain.UserDetailRepository
}

// NewUserDetailUsecase creates a new user detail usecase
func NewUserDetailUsecase(repo domain.UserDetailRepository) UserDetailUsecase {
	return &userDetailUsecase{repo: repo}
}

func (u *userDetailUsecase) GetUserDetailByID(ctx context.Context, id int64) (*domain.UserDetail, error) {
	return u.repo.FindDetailByID(ctx, id)
}

func (u *userDetailUsecase) GetUserDetailByUsername(ctx context.Context, username string) (*domain.UserDetail, error) {
	return u.repo.FindDetailByUsername(ctx, username)
}
