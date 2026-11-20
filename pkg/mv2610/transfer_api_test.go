package mv2610_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2610"
)

// The response covers optional nested details: an ACH return and a wallet captureID.
func TestListTransferEvents(t *testing.T) {
	var path, version string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"eventID": "evt-capture",
				"eventType": "capture",
				"amount": {"currency": "USD", "valueDecimal": "100.00"},
				"details": {"capture": {"captureID": "capt-1", "status": "completed"}},
				"occurredOn": "2026-08-25T16:01:01Z"
			},
			{
				"eventID": "evt-ach",
				"eventType": "ach-debit",
				"amount": {"currency": "USD", "valueDecimal": "100.00"},
				"details": {"achDebit": {"status": "returned", "return": {"code": "R01", "reason": "Insufficient funds"}}},
				"occurredOn": "2026-08-25T14:05:00Z"
			},
			{
				"eventID": "evt-wallet",
				"eventType": "wallet-credit",
				"amount": {"currency": "USD", "valueDecimal": "100.00"},
				"details": {"walletCredit": {"accountID": "acct-1", "walletID": "wallet-1", "transactionID": "txn-1", "transactionType": "payment", "captureID": "capt-1"}},
				"occurredOn": "2026-08-25T16:01:02Z"
			}
		]`))
	}))
	t.Cleanup(srv.Close)

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	transfers := mv2610.NewTransferClient(client)

	_, err = transfers.ListTransferEvents(context.Background(), "account-123", "")
	require.EqualError(t, err, "accountID and transferID are required")
	require.Empty(t, path, "an empty transferID must not send a request")

	events, err := transfers.ListTransferEvents(context.Background(), "account-123", "transfer-123")
	require.NoError(t, err)

	require.Equal(t, "/accounts/account-123/transfers/transfer-123/events", path)
	require.Equal(t, moov.Version2026_10.String(), version)

	amount := moov.AmountDecimal{Currency: "USD", ValueDecimal: "100.00"}
	captureID := "capt-1"
	require.Equal(t, []mv2610.TransferEvent{
		{
			EventID:   "evt-capture",
			EventType: mv2610.TransferEventType_Capture,
			Amount:    amount,
			Details: mv2610.TransferEventDetails{
				Capture: &mv2610.TransferEventCaptureDetails{CaptureID: "capt-1", Status: mv2610.CaptureStatus_Completed},
			},
			OccurredOn: time.Date(2026, 8, 25, 16, 1, 1, 0, time.UTC),
		},
		{
			EventID:   "evt-ach",
			EventType: mv2610.TransferEventType_AchDebit,
			Amount:    amount,
			Details: mv2610.TransferEventDetails{
				AchDebit: &mv2610.TransferEventAchDetails{
					Status: moov.AchStatus_Returned,
					Return: &moov.AchException{Code: "R01", Reason: "Insufficient funds"},
				},
			},
			OccurredOn: time.Date(2026, 8, 25, 14, 5, 0, 0, time.UTC),
		},
		{
			EventID:   "evt-wallet",
			EventType: mv2610.TransferEventType_WalletCredit,
			Amount:    amount,
			Details: mv2610.TransferEventDetails{
				WalletCredit: &mv2610.TransferEventWalletTransactionDetails{
					AccountID:       "acct-1",
					WalletID:        "wallet-1",
					TransactionID:   "txn-1",
					TransactionType: moov.WalletTransactionTypePayment,
					CaptureID:       &captureID,
				},
			},
			OccurredOn: time.Date(2026, 8, 25, 16, 1, 2, 0, time.UTC),
		},
	}, events)
}
