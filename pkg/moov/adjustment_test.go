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

func TestAdjustments(t *testing.T) {
	want := moov.Adjustment{
		AdjustmentID: "adj-123",
		WalletID:     "wallet-456",
		CreatedOn:    time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC),
		Amount:       moov.AmountDecimal{Currency: "USD", ValueDecimal: "-12.50"},
	}

	tests := []struct {
		name       string
		get        bool
		path       string
		opts       []moov.AdjustmentListFilter
		response   string
		wantWallet string
	}{
		{
			name:     "lists adjustments without a wallet filter",
			path:     "/accounts/account-123/adjustments",
			response: "[" + adjustmentJSON + "]",
		},
		{
			name:       "filters by walletID",
			path:       "/accounts/account-123/adjustments",
			opts:       []moov.AdjustmentListFilter{moov.WithAdjustmentWalletID("wallet-456")},
			response:   "[" + adjustmentJSON + "]",
			wantWallet: "wallet-456",
		},
		{
			name:     "gets an adjustment",
			get:      true,
			path:     "/accounts/account-123/adjustments/adj-123",
			response: adjustmentJSON,
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

			client := newFilesTestClient(t, srv)

			if tt.get {
				adjustment, err := client.GetAdjustment(BgCtx(), "account-123", "adj-123")
				require.NoError(t, err)
				require.Equal(t, want, *adjustment)
			} else {
				adjustments, err := client.ListAdjustments(BgCtx(), "account-123", tt.opts...)
				require.NoError(t, err)
				require.Equal(t, []moov.Adjustment{want}, adjustments)
			}

			require.Equal(t, http.MethodGet, method)
			require.Equal(t, tt.path, path)
			require.Equal(t, tt.wantWallet, wallet)
			require.Equal(t, "application/json", accept)
			require.Empty(t, version)
		})
	}
}
