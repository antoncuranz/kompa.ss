package middleware

import (
	"fmt"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"kompass/config"
	"kompass/internal/usecase"
)

const AuthorizationHeader = "Authorization"
const BearerScheme = "Bearer"

func JwtMiddleware(cfg config.Auth) func(c *fiber.Ctx) error {
	if cfg.NoAuthUserOverride != "" {
		return func(ctx *fiber.Ctx) error {
			user := cfg.NoAuthUserOverride
			subUuid := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(user))

			tokenOverride := &jwt.Token{
				Claims: jwt.MapClaims{
					"sub":  subUuid.String(),
					"name": user,
				},
			}

			ctx.Locals("jwt", tokenOverride)
			return ctx.Next()
		}
	}

	return jwtware.New(jwtware.Config{
		JWKSetURLs:  []string{cfg.JwksUrl},
		AuthScheme:  BearerScheme,
		TokenLookup: "header:" + AuthorizationHeader,
		ContextKey:  "jwt",
	})
}

func RetrieveOrCreateUserMiddleware(uc usecase.Users) func(c *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		token := ctx.Locals("jwt").(*jwt.Token)
		claims := token.Claims
		sub, err := claims.GetSubject()
		if err != nil {
			return fmt.Errorf("get subject from claims: %w", err)
		}
		subUuid, err := uuid.Parse(sub)
		if err != nil {
			return fmt.Errorf("parse sub as uuid: %w", err)
		}

		user, err := uc.GetUserByJwtSub(ctx.UserContext(), subUuid)
		if err != nil {
			user, err = uc.CreateUserFromJwt(ctx.UserContext(), subUuid, claims)
			if err != nil {
				return fmt.Errorf("create user from claims: %w", err)
			}
		}

		ctx.Locals("user", user)
		return ctx.Next()
	}
}
