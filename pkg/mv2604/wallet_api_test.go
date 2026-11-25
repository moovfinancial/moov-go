package mv2604_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2604"
)

func TestUpdateWallet(t *testing.T) {
	client := mv2604.NewWalletClient(newTestClient(t, &expectedRequest{
		method:   http.MethodPatch,
		path:     "/accounts/a/wallets/w",
		jsonBody: `{"name":"reserve","description":null}`,
	}, http.StatusOK, `{"walletID":"w","name":"reserve"}`))

	name := "reserve"
	got, err := client.UpdateWallet(t.Context(), "a", "w", mv2604.UpdateWallet{
		Name:        &name,
		Description: moov.SetNull[string](),
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "w", got.WalletID)
	require.Equal(t, "reserve", got.Name)
}
