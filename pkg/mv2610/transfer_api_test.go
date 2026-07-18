package mv2610_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2610"
)

func TestTransferClient(t *testing.T) {
	var requests []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, moov.Version2026_10.String(), r.Header.Get(moov.VersionHeader))
		requests = append(requests, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/accounts/account-123/transfers":
			if r.Method == http.MethodPost {
				_, _ = w.Write([]byte(transferJSON))
			} else {
				_, _ = w.Write([]byte("[" + transferJSON + "]"))
			}
		case "/accounts/account-123/transfers/transfer-123":
			_, _ = w.Write([]byte(transferJSON))
		case "/accounts/account-123/transfers/transfer-123/captures":
			_, _ = w.Write([]byte("[" + captureJSON + "]"))
		case "/accounts/account-123/transfers/transfer-123/captures/capture-123":
			_, _ = w.Write([]byte(captureJSON))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")
	transfers := mv2610.NewTransferClient(client)

	created, _, err := transfers.CreateTransfer(context.Background(), "account-123", mv2610.CreateTransfer{
		Source:      mv2610.CreateTransferSource{},
		Destination: mv2610.CreateTransferDestination{},
		Amount:      mv2610.AmountDecimal{Currency: "USD", ValueDecimal: "12.34"},
	}).WaitForRailResponse()
	require.NoError(t, err)
	require.Equal(t, "12.34", created.Amount.ValueDecimal)

	listed, err := transfers.ListTransfers(context.Background(), "account-123")
	require.NoError(t, err)
	require.Equal(t, "transfer-123", listed[0].TransferID)

	got, err := transfers.GetTransfer(context.Background(), "account-123", "transfer-123")
	require.NoError(t, err)
	require.Equal(t, "USD", got.Amount.Currency)

	captures, err := transfers.ListCaptures(context.Background(), "account-123", "transfer-123")
	require.NoError(t, err)
	require.Equal(t, "capture-123", captures[0].CaptureID)
	require.Equal(t, "5.67", captures[0].Amount.ValueDecimal)

	capture, err := transfers.GetCapture(context.Background(), "account-123", "transfer-123", "capture-123")
	require.NoError(t, err)
	require.Equal(t, mv2610.CaptureStatus("completed"), capture.Status)

	require.Equal(t, []string{
		"POST /accounts/account-123/transfers",
		"GET /accounts/account-123/transfers",
		"GET /accounts/account-123/transfers/transfer-123",
		"GET /accounts/account-123/transfers/transfer-123/captures",
		"GET /accounts/account-123/transfers/transfer-123/captures/capture-123",
	}, requests)
}

const transferJSON = `{
	"transferID":"transfer-123",
	"amount":{"currency":"USD","valueDecimal":"12.34"},
	"createdOn":"2026-07-18T00:00:00Z",
	"options":{},
	"processingDetails":{}
}`

const captureJSON = `{
	"captureID":"capture-123",
	"amount":{"currency":"USD","valueDecimal":"5.67"},
	"isFinal":true,
	"status":"completed",
	"createdOn":"2026-07-18T00:00:00Z"
}`
