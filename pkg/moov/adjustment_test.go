package moov_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

const adjustmentJSON = `{
	"adjustmentID": "adj-123",
	"walletID": "wallet-456",
	"createdOn": "2026-01-02T03:04:05Z",
	"amount": { "currency": "USD", "valueDecimal": "-12.50" }
}`

func TestListAdjustments(t *testing.T) {
	want := []moov.Adjustment{{
		AdjustmentID: "adj-123",
		WalletID:     "wallet-456",
		CreatedOn:    time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		Amount:       moov.AmountDecimal{Currency: "USD", ValueDecimal: "-12.50"},
	}}

	tests := []struct {
		name       string
		opts       []moov.AdjustmentListFilter
		response   string
		wantWallet string
		want       []moov.Adjustment
	}{
		{
			name:     "lists adjustments without a wallet filter",
			response: "[" + adjustmentJSON + "]",
			want:     want,
		},
		{
			name:       "filters by walletID",
			opts:       []moov.AdjustmentListFilter{moov.WithAdjustmentWalletID("wallet-456")},
			response:   "[" + adjustmentJSON + "]",
			wantWallet: "wallet-456",
			want:       want,
		},
		{
			name:     "returns an empty slice for an empty response",
			response: "[]",
			want:     []moov.Adjustment{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				method  string
				path    string
				wallet  string
				version string
				accept  string
			)

			// The handler runs on its own goroutine, so it records what it saw rather than
			// asserting: require.* calls t.FailNow, which is only valid on the test goroutine.
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				method = r.Method
				path = r.URL.Path
				wallet = r.URL.Query().Get("walletID")
				version = r.Header.Get(moov.VersionHeader)
				accept = r.Header.Get("Accept")

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.response))
			}))
			t.Cleanup(srv.Close)

			adjustments, err := newFilesTestClient(t, srv).ListAdjustments(BgCtx(), "account-123", tt.opts...)
			require.NoError(t, err)

			require.Equal(t, http.MethodGet, method)
			require.Equal(t, "/accounts/account-123/adjustments", path)
			require.Equal(t, tt.wantWallet, wallet)
			require.Equal(t, "application/json", accept)
			require.Empty(t, version)

			require.Equal(t, tt.want, adjustments)
		})
	}
}

func TestGetAdjustment(t *testing.T) {
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

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(adjustmentJSON))
	}))
	t.Cleanup(srv.Close)

	adjustment, err := newFilesTestClient(t, srv).GetAdjustment(BgCtx(), "account-123", "adj-123")
	require.NoError(t, err)

	require.Equal(t, http.MethodGet, method)
	require.Equal(t, "/accounts/account-123/adjustments/adj-123", path)
	require.Equal(t, "application/json", accept)
	require.Empty(t, version)

	require.Equal(t, "adj-123", adjustment.AdjustmentID)
	require.Equal(t, "wallet-456", adjustment.WalletID)
	require.Equal(t, moov.AmountDecimal{Currency: "USD", ValueDecimal: "-12.50"}, adjustment.Amount)
	require.Equal(t, time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC), adjustment.CreatedOn)
}
