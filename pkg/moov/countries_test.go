package moov_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func TestGetAccountCountries(t *testing.T) {
	var (
		method  string
		path    string
		version string
		accept  string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		accept = r.Header.Get("Accept")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"countries":["United States","Canada"]}`))
	}))
	t.Cleanup(srv.Close)

	countries, err := newFilesTestClient(t, srv).GetAccountCountries(BgCtx(), "account-123")
	require.NoError(t, err)

	require.Equal(t, http.MethodGet, method)
	require.Equal(t, "/accounts/account-123/countries", path)
	require.Equal(t, "application/json", accept)
	require.Equal(t, moov.Version2026_07.String(), version)

	require.Equal(t, []string{"United States", "Canada"}, countries.Countries)
}

func TestAssignAccountCountries(t *testing.T) {
	var (
		method    string
		path      string
		version   string
		accept    string
		body      moov.AccountCountries
		decodeErr error
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		accept = r.Header.Get("Accept")
		decodeErr = json.NewDecoder(r.Body).Decode(&body)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"countries":["United States","Canada"]}`))
	}))
	t.Cleanup(srv.Close)

	assigned, err := newFilesTestClient(t, srv).AssignAccountCountries(
		BgCtx(), "account-123", moov.AccountCountries{Countries: []string{"United States", "Canada"}})
	require.NoError(t, err)
	require.NoError(t, decodeErr)

	require.Equal(t, http.MethodPut, method)
	require.Equal(t, "/accounts/account-123/countries", path)
	require.Equal(t, "application/json", accept)
	require.Equal(t, moov.Version2026_07.String(), version)

	require.Equal(t, []string{"United States", "Canada"}, body.Countries)
	require.Equal(t, []string{"United States", "Canada"}, assigned.Countries)
}
