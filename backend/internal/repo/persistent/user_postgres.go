package persistent

import (
	"backplate/internal/entity"
	"backplate/pkg/postgres"
	"backplate/pkg/sqlc"
	"context"
)

type UserRepo struct {
	*sqlc.Queries
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{sqlc.New(pg.Pool)}
}

func (r *UserRepo) GetUsers(ctx context.Context) ([]entity.User, error) {
	users, err := r.Queries.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	return mapUsers(users), nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int32) (entity.User, error) {
	user, err := r.Queries.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return mapUser(user), nil
}

func mapUsers(users []sqlc.User) []entity.User {
	var result []entity.User
	for _, user := range users {
		result = append(result, mapUser(user))
	}
	return result
}

func mapUser(user sqlc.User) entity.User {
	return entity.User{
		ID:   user.ID,
		Name: user.Name,
	}
}
