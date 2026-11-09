package moov_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2607"
)

func TestListIssuingMerchantCategoriesGeneric(t *testing.T) {
	t.Run("decodes the merchant categories envelope", func(t *testing.T) {
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

		categories, err := moov.ListIssuingMerchantCategoriesGeneric[mv2607.MerchantCategories](context.Background(), newCardIssuingTestClient(t, srv), moov.Version2026_07)
		require.NoError(t, err)

		require.Equal(t, http.MethodGet, method)
		require.Equal(t, "/issuing/merchant-categories", path)
		require.Equal(t, moov.Version2026_07.String(), version)

		require.Len(t, categories.MerchantCategories, 2)
		require.Equal(t, moov.IssuingMerchantCategory_Airlines, categories.MerchantCategories[1].Category)
		require.Equal(t, []string{"4511", "3000-3299"}, categories.MerchantCategories[1].MCCs)
	})

	t.Run("returns an error for a forbidden response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		t.Cleanup(srv.Close)

		_, err := moov.ListIssuingMerchantCategoriesGeneric[mv2607.MerchantCategories](context.Background(), newCardIssuingTestClient(t, srv), moov.Version2026_07)
		require.Error(t, err)
	})
}
