package mv2604_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2604"
)

func TestCreateInvoiceNotification(t *testing.T) {
	var (
		method  string
		path    string
		version string
		body    []byte
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		method = r.Method
		path = r.URL.Path
		version = r.Header.Get(moov.VersionHeader)
		body = b

		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(srv.Close)

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	invoice := mv2604.NewInvoiceClient(client)

	err = invoice.CreateInvoiceNotification(context.Background(), "account-123", "invoice-456", mv2604.CreateInvoiceNotification{
		Channels: []mv2604.NotificationChannel{mv2604.NotificationChannelEmail},
	})
	require.NoError(t, err)

	require.Equal(t, http.MethodPost, method)
	require.Equal(t, "/accounts/account-123/invoices/invoice-456/notifications", path)
	require.Equal(t, moov.Version2026_04.String(), version)
	require.JSONEq(t, `{"channels":["email"]}`, string(body))
}
