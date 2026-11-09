package moov_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func TestGetMerchantProcessingAgreement(t *testing.T) {
	pdf := []byte("%PDF-1.4 merchant processing agreement")

	t.Run("downloads the agreement as PDF bytes", func(t *testing.T) {
		var (
			method  string
			path    string
			version string
			accept  string
		)

		// The handler runs on its own goroutine, so it records what it saw rather than
		// asserting: require.* calls t.FailNow, which is only valid on the test goroutine.
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			method = r.Method
			path = r.URL.Path
			version = r.Header.Get(moov.VersionHeader)
			accept = r.Header.Get("Accept")

			w.Header().Set("Content-Type", "application/pdf")
			_, _ = w.Write(pdf)
		}))
		t.Cleanup(srv.Close)

		agreement, err := newFilesTestClient(t, srv).GetMerchantProcessingAgreement(BgCtx(), "account-123")
		require.NoError(t, err)

		require.Equal(t, http.MethodGet, method)
		require.Equal(t, "/accounts/account-123/merchant-agreement", path)
		require.Equal(t, "application/pdf", accept)
		require.Empty(t, version)

		require.Equal(t, pdf, []byte(agreement))
	})

	t.Run("returns a typed not-found error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		t.Cleanup(srv.Close)

		agreement, err := newFilesTestClient(t, srv).GetMerchantProcessingAgreement(BgCtx(), "account-123")
		require.Nil(t, agreement)

		var httpErr moov.HttpCallResponse
		require.ErrorAs(t, err, &httpErr)
		require.Equal(t, moov.StatusNotFound, httpErr.Status())
	})
}
