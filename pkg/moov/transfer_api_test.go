package moov_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func TestOfflineCreateTransferWaitForRailResponse(t *testing.T) {
	const accountID = "account-1"
	key := uuid.MustParse("73c25004-7a4f-4c6b-8186-f4cbf91bc432")
	request := moov.CreateTransfer{
		Source:      moov.CreateTransfer_Source{PaymentMethodID: "source-1"},
		Destination: moov.CreateTransfer_Destination{PaymentMethodID: "destination-1"},
		Amount:      moov.Amount{Currency: "USD", Value: 1234},
	}

	tests := []struct {
		name         string
		status       int
		response     string
		wantTransfer *moov.Transfer
		wantStarted  *moov.TransferStarted
		wantErr      error
	}{
		{
			name:     "completed",
			status:   http.StatusOK,
			response: `{"transferID":"transfer-1","status":"completed","amount":{"currency":"USD","value":1234}}`,
			wantTransfer: &moov.Transfer{
				TransferID: "transfer-1",
				Status:     moov.TransferStatus_Completed,
				Amount:     moov.Amount{Currency: "USD", Value: 1234},
			},
		},
		{
			name:     "started",
			status:   http.StatusAccepted,
			response: `{"transferID":"transfer-1","createdOn":"2026-11-25T12:34:56Z"}`,
			wantStarted: &moov.TransferStarted{
				TransferID: "transfer-1",
				CreatedOn:  time.Date(2026, time.November, 25, 12, 34, 56, 0, time.UTC),
			},
		},
		{
			name:     "idempotency conflict",
			status:   http.StatusConflict,
			response: `{"error":"idempotency key conflict"}`,
			wantErr:  moov.ErrXIdempotencyKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newTransferTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/accounts/account-1/transfers", r.URL.Path)
				assert.Equal(t, "rail-response", r.Header.Get("X-Wait-For"))
				assert.Equal(t, []string{key.String()}, r.Header.Values("X-Idempotency-Key"))

				var body moov.CreateTransfer
				if assert.NoError(t, json.NewDecoder(r.Body).Decode(&body)) {
					assert.Equal(t, request, body)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.response))
			}))

			transfer, started, err := client.CreateTransfer(t.Context(), accountID, request,
				moov.WithTransferIdempotencyKey(key)).WaitForRailResponse()

			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.wantTransfer, transfer)
			require.Equal(t, tt.wantStarted, started)
		})
	}
}

func TestOfflineListTransfersFilters(t *testing.T) {
	client := newTransferTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/accounts/account-1/transfers", r.URL.Path)
		assert.Equal(t, url.Values{
			"skip":          {"5"},
			"count":         {"2"},
			"transferTypes": {"card-payment,ach-debit"},
		}, r.URL.Query())

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"transferID":"transfer-1","status":"completed"},{"transferID":"transfer-2","status":"pending"}]`))
	}))

	transfers, err := client.ListTransfers(t.Context(), "account-1",
		moov.WithTransferSkip(5),
		moov.WithTransferCount(2),
		moov.WithTransferTypes([]string{"card-payment", "ach-debit"}),
	)
	require.NoError(t, err)
	require.Equal(t, []moov.Transfer{
		{TransferID: "transfer-1", Status: moov.TransferStatus_Completed},
		{TransferID: "transfer-2", Status: moov.TransferStatus_Pending},
	}, transfers)
}

func newTransferTestClient(t *testing.T, handler http.Handler) *moov.Client {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{Host: srv.Listener.Addr().String(), PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	return client
}
