package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/httputils/errmux"
	"github.com/halimath/httputils/response"
	"github.com/halimath/kvlog"
)

func newAuthMux(m auth.TokenHandler) http.Handler {
	mux := errmux.NewServeMux()

	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) error {
		logger := kvlog.FromContext(r.Context())

		var token string
		var err error

		existingToken, ok := extractBearerToken(r)
		if ok {
			logger.Logs("renewToken")
			token, err = m.RenewToken(existingToken)
			if errors.Is(err, auth.ErrUnauthorized) {
				logger.Logs("token renewal failed. creating a fresh one")
				token, err = m.CreateToken()
			}
		} else {
			logger.Logs("createToken")
			token, err = m.CreateToken()
		}

		if err != nil {
			logger.Logs("error creating auth token", kvlog.WithErr(err))
			return err
		}

		return response.PlainText(w, r, token, response.StatusCode(http.StatusCreated))
	})

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) error {
		logger := kvlog.FromContext(r.Context())

		authToken, ok := extractBearerToken(r)
		if !ok {
			return sendUnauthorized(w, r)
		}

		authInfo, err := m.Authorize(authToken)
		if err != nil {
			logger.Logs("invalid authentication token", kvlog.WithErr(err))
			return sendForbidden(w, r)
		}

		return response.JSON(w, r, AuthenticationInfo{
			UserId:  authInfo.UserID,
			Expires: authInfo.Expires,
		})
	})

	return mux
}

func extractBearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return "", false
	}

	if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return "", false
	}

	return authHeader[len("bearer "):], true

}

func authMiddleware(m auth.TokenHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, ok := extractBearerToken(r)
			if !ok {
				sendUnauthorized(w, r)
				return
			}

			authInfo, err := m.Authorize(tokenString)
			if err != nil {
				kvlog.L.Logs("invalidAuthToken", kvlog.WithErr(err))
				sendForbidden(w, r)
				return
			}

			r = r.WithContext(auth.WithUserID(r.Context(), authInfo.UserID))

			next.ServeHTTP(w, r)
		})
	}
}

func sendUnauthorized(w http.ResponseWriter, r *http.Request) error {
	return response.Problem(w, r, response.ProblemDetails{
		Type:   "https://github.com/halimath/fate-table/problem/unauthorized",
		Title:  "Unauthorized",
		Status: http.StatusUnauthorized,
		Detail: "The request is missing an Authorization header or the authorization scheme is invalid.",
	})
}

func sendForbidden(w http.ResponseWriter, r *http.Request) error {
	return response.Problem(w, r, response.ProblemDetails{
		Type:   "https://github.com/halimath/fate-table/problem/forbidden",
		Title:  "Forbidden",
		Status: http.StatusForbidden,
		Detail: "The request's authorization token is either invalid or the user is not permitted to execute the operation.",
	})
}
