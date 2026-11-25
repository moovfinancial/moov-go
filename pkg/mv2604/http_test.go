package mv2604_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type expectedRequest struct {
	method   string
	path     string
	jsonBody string
}

// newTestClient checks one HTTP request, or zero if want is nil.
func newTestClient(t *testing.T, want *expectedRequest, status int, response string) *moov.Client {
	t.Helper()

	type request struct {
		method, path, version, contentType, body string
		readErr                                  error
	}
	var mu sync.Mutex
	var requests []request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		mu.Lock()
		requests = append(requests, request{
			method:      r.Method,
			path:        r.URL.RequestURI(),
			version:     r.Header.Get("X-Moov-Version"),
			contentType: r.Header.Get("Content-Type"),
			body:        string(body),
			readErr:     err,
		})
		mu.Unlock()

		if response != "" {
			w.Header().Set("Content-Type", "application/json")
		}
		w.WriteHeader(status)
		if response != "" {
			if _, err := io.WriteString(w, response); err != nil {
				t.Errorf("write test response: %v", err)
			}
		}
	}))
	t.Cleanup(func() {
		server.Close()

		mu.Lock()
		received := requests
		mu.Unlock()

		if want == nil {
			require.Empty(t, received, "validation should stop before sending an HTTP request")
			return
		}
		require.Len(t, received, 1)
		got := received[0]
		require.NoError(t, got.readErr)
		require.Equal(t, want.method, got.method)
		require.Equal(t, want.path, got.path)
		require.Equal(t, "v2026.04.00", got.version)
		if want.jsonBody == "" {
			require.Empty(t, got.body)
			return
		}
		require.Equal(t, "application/json", got.contentType)
		require.JSONEq(t, want.jsonBody, got.body)
	})

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(server.URL, "http://")
	return client
}
