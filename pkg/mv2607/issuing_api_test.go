package mv2607_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2607"
)

// newIssuingTestClient points a client at srv so tests can assert on the raw request.
func newIssuingTestClient(t *testing.T, srv *httptest.Server) mv2607.IssuingClient {
	t.Helper()

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	return mv2607.NewIssuingClient(client)
}

func TestListMerchantCategories(t *testing.T) {
	t.Run("returns categories with their MCCs", func(t *testing.T) {
		var (
			method  string
			path    string
			version string
		)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			method = r.Method
			path = r.URL.Path
			version = r.Header.Get(moov.VersionHeader)

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"merchantCategories": [
					{"category": "restaurants-and-dining", "mccs": ["5812", "5813"]},
					{"category": "airlines", "mccs": ["4511", "3000-3299"]}
				]
			}`))
		}))
		t.Cleanup(srv.Close)

		categories, err := newIssuingTestClient(t, srv).ListMerchantCategories(context.Background())
		require.NoError(t, err)

		require.Equal(t, http.MethodGet, method)
		require.Equal(t, "/issuing/merchant-categories", path)
		require.Equal(t, moov.Version2026_07.String(), version)

		require.Len(t, categories, 2)
		require.Equal(t, moov.IssuingMerchantCategory_RestaurantsAndDining, categories[0].Category)
		require.Equal(t, []string{"5812", "5813"}, categories[0].MCCs)
		require.Equal(t, moov.IssuingMerchantCategory_Airlines, categories[1].Category)
		require.Equal(t, []string{"4511", "3000-3299"}, categories[1].MCCs)
	})

	t.Run("returns an empty slice for an empty merchantCategories array", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"merchantCategories": []}`))
		}))
		t.Cleanup(srv.Close)

		categories, err := newIssuingTestClient(t, srv).ListMerchantCategories(context.Background())
		require.NoError(t, err)

		require.NotNil(t, categories)
		require.Empty(t, categories)
	})
}
