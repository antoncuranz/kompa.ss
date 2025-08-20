package users

import (
	"context"
	"kompass/internal/entity"
	"kompass/internal/repo"
)

type UseCase struct {
	repo repo.UserRepo
}

func New(r repo.UserRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) GetUserByID(ctx context.Context, id int32) (entity.User, error) {
	return uc.repo.GetUserByID(ctx, id)
}

func (uc *UseCase) GetUsers(ctx context.Context) ([]entity.User, error) {
	return uc.repo.GetUsers(ctx)
}
