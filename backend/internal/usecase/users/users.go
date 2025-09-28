package users

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (uc *UseCase) GetUserByJwtSub(ctx context.Context, sub uuid.UUID) (entity.User, error) {
	return uc.repo.GetUserByJwtSub(ctx, sub)
}

func (uc *UseCase) CreateUserFromJwt(ctx context.Context, sub uuid.UUID, claims jwt.Claims) (entity.User, error) {
	mapClaims := claims.(jwt.MapClaims)
	name := parseString(mapClaims, "name", "Unknown")
	return uc.repo.CreateUser(ctx, entity.User{
		Name:   name,
		JwtSub: sub,
	})
}

func (uc *UseCase) HasReadPermission(ctx context.Context, userID, tripID int32) (bool, error) {
	return uc.repo.HasReadPermission(ctx, userID, tripID)
}

func (uc *UseCase) HasWritePermission(ctx context.Context, userID, tripID int32) (bool, error) {
	return uc.repo.HasWritePermission(ctx, userID, tripID)
}

func (uc *UseCase) IsTripOwner(ctx context.Context, userID int32, tripID int32) (bool, error) {
	return uc.repo.IsTripOwner(ctx, userID, tripID)
}

func parseString(mapClaims map[string]any, key string, defaultValue string) string {
	raw, ok := mapClaims[key]
	if !ok {
		return defaultValue
	}

	str, ok := raw.(string)
	if !ok {
		return defaultValue
	}

	return str
}
