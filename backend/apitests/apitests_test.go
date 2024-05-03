//go:generate oapi-codegen -package apitests -generate types,client -o client_gen_test.go ../../docs/api.yaml
package apitests

import (
	"context"
	"net/http"
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
	"github.com/halimath/fate-core-remote-table/backend/apitests/httpresponsewith"
	"github.com/halimath/fate-core-remote-table/backend/internal"
	"github.com/halimath/fixture"
)

func TestAPISceneario(t *testing.T) {
	fixture.With(t, new(fix)).
		Run("example_scenario", func(t *testing.T, f *fix) {
			gmClient := f.AuthorizedAPIClient(t)

			// Create session
			var sessionID string
			r, err := gmClient.CreateSession(f.ctx, CreateSession{
				Title: "Test Session",
			})
			expect.WithMessage(t, "gm: create session").That(
				is.NoError(err),
				expect.FailNow(httpresponsewith.SuccessfulStatusCode(r)),
				httpresponsewith.TextBody(r, &sessionID),
			)

			// Load session
			var session Session
			r, err = gmClient.GetSession(f.ctx, sessionID)
			expect.WithMessage(t, "gm: get session").That(
				is.NoError(err),
				expect.FailNow(httpresponsewith.SuccessfulStatusCode(r)),
				httpresponsewith.JSOnBody(r, &session),
			)

			// Create api client for player
			playerClient := f.AuthorizedAPIClient(t)

			// Create a PC for that player
			var pcID string
			r, err = playerClient.JoinSession(f.ctx, sessionID, JoinSession{
				Name: "Player One",
			})
			expect.WithMessage(t, "p1: join session").That(
				is.NoError(err),
				expect.FailNow(httpresponsewith.SuccessfulStatusCode(r)),
				httpresponsewith.TextBody(r, &pcID),
			)

			// Load session for GM and expect that player one is part of the session
			r, err = gmClient.GetSession(f.ctx, sessionID)
			expect.WithMessage(t, "gm: get session").
				That(
					is.NoError(err),
					expect.FailNow(httpresponsewith.SuccessfulStatusCode(r)),
					httpresponsewith.JSOnBody(r, &session),
				).
				That(
					is.DeepEqualTo(session, Session{
						Id:    sessionID,
						Title: "Test Session",
						Characters: []Character{
							{
								Id:         pcID,
								Name:       "Player One",
								Type:       CharacterTypePC,
								FatePoints: 0,
							},
						},
					}, is.ExcludeFields{"OwnerId"}),
				)

			// Increase fate points for player one
			r, err = gmClient.UpdateFatePoints(f.ctx, sessionID, session.Characters[0].Id, UpdateFatePoints{
				FatePointsDelta: 2,
			})
			expect.WithMessage(t, "gm: increment fate points").That(
				is.NoError(err),
				expect.FailNow(httpresponsewith.SuccessfulStatusCode(r)),
			)

			// Load session for Player and expect fate points to be 2
			r, err = playerClient.GetSession(f.ctx, sessionID)
			expect.WithMessage(t, "p1: get session").
				That(
					is.NoError(err),
					expect.FailNow(httpresponsewith.SuccessfulStatusCode(r)),
					httpresponsewith.JSOnBody(r, &session),
				).
				That(
					is.DeepEqualTo(session, Session{
						Id:    sessionID,
						Title: "Test Session",
						Characters: []Character{
							{
								Id:         pcID,
								Name:       "Player One",
								Type:       CharacterTypePC,
								FatePoints: 2,
							},
						},
					}, is.ExcludeFields{"OwnerId"}),
				)
		})
}

func withAuthorizationBearerToken(token string) RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}

type fix struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	apiClient *Client
}

const server = "http://localhost:8080/api"

func (f *fix) BeforeAll(t *testing.T) error {
	f.ctx, f.cancel = context.WithCancel(context.Background())
	f.done = make(chan struct{})

	go func() {
		defer close(f.done)
		t.Setenv("DEV_MODE", "1")
		internal.RunService(f.ctx)
	}()

	var err error
	f.apiClient, err = NewClient(server)

	return err
}

func (f *fix) AuthorizedAPIClient(t *testing.T) *Client {
	var authToken string
	r, err := f.apiClient.CreateAuthToken(f.ctx)
	expect.WithMessage(t, "auth token").That(expect.FailNow(
		is.NoError(err),
		httpresponsewith.SuccessfulStatusCode(r),
		httpresponsewith.TextBody(r, &authToken),
	))

	// t.Log(authToken)

	c, err := NewClient(server, WithRequestEditorFn(withAuthorizationBearerToken(authToken)))
	expect.That(t, expect.FailNow(is.NoError(err)))
	return c
}

func (f *fix) AfterAll(t *testing.T) error {
	// Stop the backend...
	f.cancel()
	// ... and wait for shutdown
	<-f.done

	return nil
}
