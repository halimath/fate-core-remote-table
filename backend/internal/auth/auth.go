// Package auth provides functions and type to implement the cross-cutting concern of user authentication.
// It provides functions to store and retrieve a user's ID from a context.Context as well as a TokenHandler
// to create signed tokens describing a user's id and verify those tokens.
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

// IsAuthorized returns true, if ctx contains a valid user's ID.
func IsAuthorized(ctx context.Context) bool {
	_, ok := UserID(ctx)
	return ok
}

// UserID retrieves the user's ID from ctx and returns it including an ok flag
// which is false, if ctx contains no valid user ID.
func UserID(ctx context.Context) (string, bool) {
	s := ctx.Value(userIDContextKey)
	if s == nil {
		return "", false
	}

	n, ok := s.(string)

	return n, ok
}

// WithUserID creates a new context with ctx as it's parent holding userID as
// an authenticated user's ID.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

var (
	// ErrUnauthorized is a sentinel error value returned when a token operation
	// is aborted due to an authorization error.
	ErrUnauthorized = errors.New("unauthorized")
)

// TokenHandler defines a type for creating, verifying and renewing authentication
// tokens.
type TokenHandler struct {
	signature jws.SignerVerifier
	tokenTTL  time.Duration
}

// CreateToken creates a new token for a new, randomly generated user id. It
// returns the encoded token.
func (h *TokenHandler) CreateToken() (string, error) {
	return h.createToken(id.New())
}

func (h *TokenHandler) createToken(userID string) (string, error) {
	token, err := jwt.Sign(h.signature, jwt.StandardClaims{
		ID:             id.New(),
		Subject:        userID,
		Issuer:         authTokenIssuer,
		Audience:       []string{authTokenIssuer},
		ExpirationTime: time.Now().Add(h.tokenTTL).Unix(),
	})

	if err != nil {
		return "", err
	}

	return token.Compact(), nil
}

// RenewToken creates a new token holding the same user's ID as tokenString. It only does so, if tokenString
// is a valid token. For this validity check, the token's TTL will be handled with an increased leeway. It
// returns a new token or an error. If tokenString cannot be verified, the returned error is ErrUnauthorized.
func (h *TokenHandler) RenewToken(tokenString string) (string, error) {
	n, err := authorize(tokenString, jwt.Signature(h.signature), jwt.Audience(authTokenIssuer), jwt.Issuer(authTokenIssuer), jwt.ExpirationTime(h.tokenTTL))
	if err != nil {
		return "", err
	}

	return h.createToken(n)
}

// Authorize authorizes tokenString and returns the encoded user's ID or an error.
func (h *TokenHandler) Authorize(tokenString string) (string, error) {
	return authorize(tokenString, jwt.Signature(h.signature), jwt.Audience(authTokenIssuer), jwt.Issuer(authTokenIssuer), jwt.ExpirationTime(0))
}

func authorize(tokenString string, verifier ...jwt.Verifier) (string, error) {
	token, err := jwt.Decode(tokenString)
	if err != nil {
		return "", ErrUnauthorized
	}

	if err := token.Verify(verifier...); err != nil {
		return "", ErrUnauthorized
	}

	return token.StandardClaims().Subject, nil
}

func Provide(cfg config.Config) *TokenHandler {
	return &TokenHandler{
		signature: jws.HS256([]byte(cfg.AuthTokenSecret)),
		tokenTTL:  cfg.AuthTokenTTL,
	}
}
