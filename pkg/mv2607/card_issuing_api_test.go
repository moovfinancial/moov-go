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

// newCardIssuingTestClient points a client at srv so tests can assert on the raw request.
func newCardIssuingTestClient(t *testing.T, srv *httptest.Server) mv2607.CardIssuingClient {
	t.Helper()

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	return mv2607.NewCardIssuingClient(client)
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

		categories, err := newCardIssuingTestClient(t, srv).ListMerchantCategories(context.Background())
		require.NoError(t, err)

		require.Equal(t, http.MethodGet, method)
		require.Equal(t, "/issuing/merchant-categories", path)
		require.Equal(t, moov.Version2026_07.String(), version)

		require.Equal(t, &mv2607.MerchantCategories{
			MerchantCategories: []mv2607.MerchantCategory{
				{Category: moov.IssuingMerchantCategory_RestaurantsAndDining, MCCs: []string{"5812", "5813"}},
				{Category: moov.IssuingMerchantCategory_Airlines, MCCs: []string{"4511", "3000-3299"}},
			},
		}, categories)
	})
}
