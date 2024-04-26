package auth

import (
	"context"
	"errors"
	"time"

	"github.com/halimath/fate-core-remote-table/backend/internal/id"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/jose/jws"
	"github.com/halimath/jose/jwt"
)

type userIDContextKeyType string

const userIDContextKey userIDContextKeyType = "userID"
const authTokenIssuer = "fate-table"

var (
	ErrUnauthorized = errors.New("unauthorized")
)

func IsAuthorized(ctx context.Context) bool {
	_, ok := UserID(ctx)
	return ok
}

func UserID(ctx context.Context) (string, bool) {
	s := ctx.Value(userIDContextKey)
	if s == nil {
		return "", false
	}

	n, ok := s.(string)

	return n, ok
}

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

type Manager struct {
	signature jws.SignerVerifier
	tokenTTL  time.Duration
}

func (m *Manager) CreateToken() (string, error) {
	return m.createToken(id.New())
}

func (m *Manager) createToken(userID string) (string, error) {
	token, err := jwt.Sign(m.signature, jwt.StandardClaims{
		ID:             id.New(),
		Subject:        userID,
		Issuer:         authTokenIssuer,
		Audience:       []string{authTokenIssuer},
		ExpirationTime: time.Now().Add(m.tokenTTL).Unix(),
	})

	if err != nil {
		return "", err
	}

	return token.Compact(), nil
}

func (m *Manager) RenewToken(tokenString string) (string, error) {
	n, err := m.authorize(tokenString, jwt.Signature(m.signature), jwt.Audience(authTokenIssuer), jwt.Issuer(authTokenIssuer))
	if err != nil {
		return "", err
	}

	return m.createToken(n)
}

func (m *Manager) Authorize(tokenString string) (string, error) {
	return m.authorize(tokenString, jwt.Signature(m.signature), jwt.Audience(authTokenIssuer), jwt.Issuer(authTokenIssuer), jwt.ExpirationTime(0))
}

func (m *Manager) authorize(tokenString string, verifier ...jwt.Verifier) (string, error) {
	token, err := jwt.Decode(tokenString)
	if err != nil {
		return "", ErrUnauthorized
	}

	if err := token.Verify(verifier...); err != nil {
		return "", ErrUnauthorized
	}

	return token.StandardClaims().Subject, nil
}

func Provide(cfg config.Config) *Manager {
	return &Manager{
		signature: jws.HS256([]byte(cfg.AuthTokenSecret)),
		tokenTTL:  cfg.AuthTokenTTL,
	}
}
