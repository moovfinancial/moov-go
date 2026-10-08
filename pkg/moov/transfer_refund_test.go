package moov_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func TestRefundTransfer_Started(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"transferID":"refund-1"}`))
	}))
	t.Cleanup(srv.Close)

	client, err := moov.NewClient(
		moov.WithCredentials(moov.Credentials{PublicKey: "pk", SecretKey: "sk"}),
		moov.WithMoovURLScheme("http"),
	)
	require.NoError(t, err)
	client.Credentials.Host = strings.TrimPrefix(srv.URL, "http://")

	refund := moov.CreateRefund{Amount: 100}

	t.Run("RefundTransfer", func(t *testing.T) {
		completed, started, err := client.RefundTransfer(context.Background(), "account-1", "transfer-1", refund)
		require.NoError(t, err)
		require.Nil(t, completed)
		require.Equal(t, "refund-1", started.TransferID)
	})

	t.Run("RefundTransferGeneric", func(t *testing.T) {
		completed, started, err := moov.RefundTransferGeneric[moov.CreateRefund, moov.Refund](context.Background(), client, moov.Version2026_10, "account-1", "transfer-1", refund)
		require.NoError(t, err)
		require.Nil(t, completed)
		require.Equal(t, "refund-1", started.TransferID)
	})
}
