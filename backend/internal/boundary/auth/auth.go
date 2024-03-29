package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/halimath/fate-core-remote-table/backend/internal/entity/id"
	"github.com/halimath/fate-core-remote-table/backend/internal/infra/config"
	"github.com/halimath/jose/jws"
	"github.com/halimath/jose/jwt"
	"github.com/halimath/kvlog"
	"github.com/labstack/echo/v4"
)

const (
	authTokenContextKey = "authToken"
	authTokenIssuer     = "fate-core-remote-table"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type Provider interface {
	CreateToken() (string, error)
	RenewToken(token string) (string, error)
	Authorize(token string) (id.ID, error)
}

func IsAuthorized(ctx echo.Context) bool {
	_, ok := UserID(ctx)
	return ok
}

func UserID(ctx echo.Context) (id.ID, bool) {
	s := ctx.Get(authTokenContextKey)
	if s == nil {
		return "", false
	}

	n, ok := s.(id.ID)

	return n, ok
}

func ExtractBearerToken(ctx echo.Context) (string, bool) {
	authHeader := ctx.Request().Header.Get("Authorization")
	if len(authHeader) == 0 {
		return "", false
	}

	if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return "", false
	}

	return authHeader[len("bearer "):], true

}

func Middleware(p Provider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			tokenString, ok := ExtractBearerToken(ctx)
			if !ok {
				return next(ctx)
			}

			sub, err := p.Authorize(tokenString)
			if err != nil {
				kvlog.L.Logs("invalidAuthToken", kvlog.WithErr(err))
				return next(ctx)
			}

			ctx.Set(authTokenContextKey, sub)

			return next(ctx)
		}
	}
}

type provider struct {
	signature jws.SignerVerifier
	tokenTTL  time.Duration
}

func (p *provider) CreateToken() (string, error) {
	return p.createToken(id.New())
}

func (p *provider) createToken(userID id.ID) (string, error) {
	token, err := jwt.Sign(p.signature, jwt.StandardClaims{
		ID:             uuid.New().String(),
		Subject:        userID.String(),
		Issuer:         authTokenIssuer,
		Audience:       []string{authTokenIssuer},
		ExpirationTime: time.Now().Add(p.tokenTTL).Unix(),
	})

	if err != nil {
		return "", err
	}

	return token.Compact(), nil
}

func (p *provider) RenewToken(tokenString string) (string, error) {
	n, err := p.authorize(tokenString, jwt.Signature(p.signature), jwt.Audience(authTokenIssuer), jwt.Issuer(authTokenIssuer))
	if err != nil {
		return "", err
	}

	return p.createToken(n)
}

func (p *provider) Authorize(tokenString string) (id.ID, error) {
	return p.authorize(tokenString, jwt.Signature(p.signature), jwt.Audience(authTokenIssuer), jwt.Issuer(authTokenIssuer), jwt.ExpirationTime(0))
}

func (p *provider) authorize(tokenString string, verifier ...jwt.Verifier) (id.ID, error) {
	token, err := jwt.Decode(tokenString)
	if err != nil {
		return "", ErrUnauthorized
	}

	if err := token.Verify(verifier...); err != nil {
		return "", ErrUnauthorized
	}

	return id.FromString(token.StandardClaims().Subject), nil
}

func Provide(cfg config.Config) Provider {
	return &provider{
		signature: jws.HS256([]byte(cfg.AuthTokenSecret)),
		tokenTTL:  cfg.AuthTokenTTL,
	}
}
