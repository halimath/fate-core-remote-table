package web

import (
	"errors"
	"net/http"

	"github.com/halimath/fate-core-remote-table/backend/internal/auth"
	"github.com/halimath/httputils/errmux"
	"github.com/halimath/httputils/response"
	"github.com/halimath/kvlog"
)

func newAuthMux(m *auth.TokenHandler) http.Handler {
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

	return mux
}
