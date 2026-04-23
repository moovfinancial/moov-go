package moov

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHTTPCallResponse(t *testing.T) {
	t.Run("400", func(t *testing.T) {
		resp := &httpCallResponse{
			resp: &http.Response{
				StatusCode: http.StatusBadRequest,
			},
			body: []byte(`{"error":"X-Idempotency-Key HTTP header must be a valid UUID"}`),
		}
		expected := strings.TrimSpace(`
error from moov - status: bad_request http.request_id:  http.status_code: 400
  X-Idempotency-Key HTTP header must be a valid UUID
`)
		require.Equal(t, expected, resp.Error())
	})

	t.Run("409 and 422", func(t *testing.T) {
		resp := &httpCallResponse{
			resp: &http.Response{
				StatusCode: http.StatusUnprocessableEntity,
			},
			body: []byte(`{"profile":{"business":{"taxID":{"ein":{"number":"must be a valid employer identification number"}}}}}`),
		}
		expected := "error from moov - status: failed_validation http.request_id:  http.status_code: 422 - profile.business.taxID.ein.number: must be a valid employer identification number"
		require.Equal(t, expected, resp.Error())
	})
}

func TestCallHttp_AuthHeader(t *testing.T) {
	type capture struct {
		auth string
	}

	newClient := func(t *testing.T, cfg ...ClientConfigurable) (*Client, *capture) {
		t.Helper()
		cap := &capture{}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cap.auth = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
		}))
		t.Cleanup(srv.Close)

		host := strings.TrimPrefix(srv.URL, "http://")
		cfg = append(cfg, WithMoovURLScheme("http"))

		c, err := NewClient(cfg...)
		require.NoError(t, err)
		c.Credentials.Host = host
		return c, cap
	}

	t.Run("Credentials.Token sends Bearer and omits Basic", func(t *testing.T) {
		c, cap := newClient(t)
		c = c.WithBearerToken("abc")
		_, err := c.CallHttp(context.Background(), Endpoint(http.MethodGet, "/ping"))
		require.NoError(t, err)
		require.Equal(t, "Bearer abc", cap.auth)
	})

	t.Run("WithToken option sets the bearer", func(t *testing.T) {
		c, cap := newClient(t)
		c = c.WithBearerToken("xyz")
		_, err := c.CallHttp(context.Background(), Endpoint(http.MethodGet, "/ping"))
		require.NoError(t, err)
		require.Equal(t, "Bearer xyz", cap.auth)
	})

	t.Run("falls back to Basic auth when no token is set", func(t *testing.T) {
		c, cap := newClient(t, WithCredentials(Credentials{PublicKey: "pk", SecretKey: "sk"}))
		_, err := c.CallHttp(context.Background(), Endpoint(http.MethodGet, "/ping"))
		require.NoError(t, err)
		require.True(t, strings.HasPrefix(cap.auth, "Basic "), "expected Basic auth, got %q", cap.auth)
	})
}

func TestCredentials_Validate(t *testing.T) {
	t.Run("token-only is valid", func(t *testing.T) {
		c := Credentials{Token: "abc"}
		require.NoError(t, c.Validate())
	})

	t.Run("public+secret is valid", func(t *testing.T) {
		c := Credentials{PublicKey: "pk", SecretKey: "sk"}
		require.NoError(t, c.Validate())
	})

	t.Run("empty is invalid", func(t *testing.T) {
		c := Credentials{}
		require.ErrorIs(t, c.Validate(), ErrCredentialsNotSet)
	})

	t.Run("partial basic creds is invalid", func(t *testing.T) {
		require.ErrorIs(t, (&Credentials{PublicKey: "pk"}).Validate(), ErrCredentialsNotSet)
		require.ErrorIs(t, (&Credentials{SecretKey: "sk"}).Validate(), ErrCredentialsNotSet)
	})
}
