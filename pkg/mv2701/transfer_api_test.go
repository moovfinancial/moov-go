package mv2701_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2701"
)

func TestListTransferEvents(t *testing.T) {
	var path, version string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{
				"eventID": "evt-wire",
				"eventType": "wire-credit",
				"amount": {"currency": "USD", "valueDecimal": "50.00"},
				"details": {"wireCredit": {"status": "failed", "failureCode": "account-closed"}},
				"occurredOn": "2027-01-05T17:00:03Z"
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

	events, err := mv2701.NewTransferClient(client).ListTransferEvents(context.Background(), "account-123", "transfer-123")
	require.NoError(t, err)

	require.Equal(t, "/accounts/account-123/transfers/transfer-123/events", path)
	require.Equal(t, moov.Version2027_01.String(), version)

	failureCode := mv2701.WireFailureCode_AccountClosed
	require.Equal(t, []mv2701.TransferEvent{
		{
			EventID:   "evt-wire",
			EventType: mv2701.TransferEventType_WireCredit,
			Amount:    moov.AmountDecimal{Currency: "USD", ValueDecimal: "50.00"},
			Details: mv2701.TransferEventDetails{
				WireCredit: &mv2701.TransferEventWireCreditDetails{
					Status:      mv2701.WireTransactionStatus_Failed,
					FailureCode: &failureCode,
				},
			},
			OccurredOn: time.Date(2027, 1, 5, 17, 0, 3, 0, time.UTC),
		},
	}, events)
}
