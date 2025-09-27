package util

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/wiremock/go-wiremock"
	"kompass/integration-test/client/api"
	"math/big"
	"net/http"
	"testing"
)

type UserToken string
type UserName string

func (u UserToken) Bearerauth(context.Context, api.OperationName) (api.Bearerauth, error) {
	return api.Bearerauth{
		Token: string(u),
	}, nil
}

func GeneratePrivateKeyAndJwkStub(t testing.TB) (*rsa.PrivateKey, *wiremock.StubRule) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	stubRule := wiremock.Get(wiremock.URLPathEqualTo("/auth/jwks.json")).
		WillReturnResponse(
			wiremock.NewResponse().
				WithJSONBody(createJwkSet(privateKey)).
				WithHeader("Content-Type", "application/json").
				WithStatus(http.StatusOK),
		)

	return privateKey, stubRule
}

func GenerateJwtForUser(t testing.TB, user UserName, privateKey *rsa.PrivateKey) UserToken {
	t.Helper()

	subUuid := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(user))

	token := jwt.Token{
		Header: map[string]any{
			"typ": "JWT",
			"alg": jwt.SigningMethodRS256.Alg(),
			"kid": "test-kid",
		},
		Claims: jwt.MapClaims{
			"sub":  subUuid.String(),
			"name": string(user),
		},
		Method: jwt.SigningMethodRS256,
	}
	signedJwt, err := token.SignedString(privateKey)
	assert.NoError(t, err)

	return UserToken(signedJwt)
}

func createJwkSet(privateKey *rsa.PrivateKey) map[string]interface{} {
	n := privateKey.PublicKey.N
	e := privateKey.PublicKey.E

	nB64 := base64.RawURLEncoding.EncodeToString(n.Bytes())
	eB64 := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(e)).Bytes())

	jwk := map[string]interface{}{
		"kty": "RSA",
		"kid": "test-kid",
		"n":   nB64,
		"e":   eB64,
	}

	return map[string]interface{}{
		"keys": []interface{}{jwk},
	}
}
