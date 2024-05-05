package httpresponsewith

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
)

func SuccessfulStatusCode(r *http.Response) expect.ExpectFunc {
	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

		if r.StatusCode < http.StatusOK || r.StatusCode >= http.StatusMultipleChoices {
			b, _ := io.ReadAll(r.Body)
			defer r.Body.Close()
			t.Errorf("unsuccessful status code: %d; response was: %s", r.StatusCode, string(b))
		}
	})
}

func Header(r *http.Response, h string, values ...string) expect.ExpectFunc {
	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

		got, ok := r.Header[h]
		if !ok {
			t.Errorf("expected response to contain header %s but none found", h)
			return
		}

		for _, val := range values {
			if !slices.Contains(got, val) {
				t.Errorf("expected response to contain header %s with value %q but that value was not found", h, val)
			}
		}
	})
}

func TextBody(r *http.Response, body *string) expect.ExpectFunc {
	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		expect.WithMessage(t, "response body", expect.FailNow(is.NoError(err)))

		*body = string(data)
	})
}

func JSOnBody(r *http.Response, body any) expect.ExpectFunc {
	return expect.ExpectFunc(func(t expect.TB) {
		t.Helper()

		e := expect.WithMessage(t, "response body")
		e.That(Header(r, "Content-Type", "application/json"))

		defer r.Body.Close()
		data, err := io.ReadAll(r.Body)
		e.That(expect.FailNow(is.NoError(err)))

		err = json.Unmarshal(data, body)
		e.That(expect.FailNow(is.NoError(err)))
	})
}
