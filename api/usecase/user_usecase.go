package usecase

import (
	"context"

	"github.com/rssh-jp/test-api/api/domain"
)

// UserUsecase handles business logic for user operations
type UserUsecase interface {
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
	CreateUser(ctx context.Context, name, email string, age *int32) (*domain.User, error)
	UpdateUser(ctx context.Context, id int64, name, email *string, age *int32) (*domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type userUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(userRepo domain.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return u.userRepo.FindAll(ctx)
}

func (u *userUsecase) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *userUsecase) CreateUser(ctx context.Context, name, email string, age *int32) (*domain.User, error) {
	user := &domain.User{
		Name:  name,
		Email: email,
		Age:   age,
	}
	err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) UpdateUser(ctx context.Context, id int64, name, email *string, age *int32) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		user.Name = *name
	}
	if email != nil {
		user.Email = *email
	}
	if age != nil {
		user.Age = age
	}

	err = u.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, id int64) error {
	return u.userRepo.Delete(ctx, id)
}
