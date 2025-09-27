package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"kompass/internal/entity"
	"kompass/pkg/postgres"
	"kompass/pkg/sqlc"
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
		return nil, fmt.Errorf("get users: %w", err)
	}

	return mapUsers(users), nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int32) (entity.User, error) {
	user, err := r.Queries.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, fmt.Errorf("get user [id=%d]: %w", id, err)
	}

	return mapUser(user), nil
}

func (r *UserRepo) GetUserByJwtSub(ctx context.Context, sub uuid.UUID) (entity.User, error) {
	user, err := r.Queries.GetUserByJwtSub(ctx, sub)
	if err != nil {
		return entity.User{}, fmt.Errorf("get user [sub=%s]: %w", sub.String(), err)
	}

	return mapUser(user), nil
}

func (r *UserRepo) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	userID, err := r.Queries.InsertUser(ctx, sqlc.InsertUserParams{
		Name:   user.Name,
		JwtSub: user.JwtSub,
	})
	if err != nil {
		return entity.User{}, fmt.Errorf("insert user: %w", err)
	}

	return r.GetUserByID(ctx, userID)
}

func mapUsers(users []sqlc.User) []entity.User {
	result := []entity.User{}
	for _, user := range users {
		result = append(result, mapUser(user))
	}
	return result
}

func mapUser(user sqlc.User) entity.User {
	return entity.User{
		ID:     user.ID,
		Name:   user.Name,
		JwtSub: user.JwtSub,
	}
}
