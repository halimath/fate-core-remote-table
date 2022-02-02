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

func Middleware(p Provider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			authHeader := ctx.Request().Header.Get("Authorization")
			if len(authHeader) == 0 {
				kvlog.Debug(kvlog.Msg("no auth header given"))
				return next(ctx)
			}

			if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
				kvlog.Debug(kvlog.Msg("no bearer scheme"))
				return next(ctx)
			}

			tokenString := authHeader[len("bearer "):]

			sub, err := p.Authorize(tokenString)
			if err != nil {
				kvlog.Warn(kvlog.Evt("invalidAuthToken"), kvlog.Err(err))
				return echo.ErrForbidden
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
