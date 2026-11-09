package moov_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/moovfinancial/moov-go/internal/testtools"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func TestWalletMarshal(t *testing.T) {
	input := []byte(`{
			"walletID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
			"availableBalance": {
				"currency": "USD",
				"value": 1204,
				"valueDecimal": "12.987654321"
			}
		}`)

	wallet := new(moov.Wallet)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&wallet)
	require.NoError(t, err)

	require.Equal(t, "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", wallet.WalletID)
}

func TestListWallet(t *testing.T) {
	mc := NewTestClient(t)

	wallets, err := mc.ListWallets(context.Background(), testtools.MERCHANT_ID)
	NoResponseError(t, err)
	require.NotNil(t, wallets)

	for i := range wallets {
		t.Logf("wallet[%d]: %#v", i, wallets[i])
	}
}

func TestGetWallet(t *testing.T) {
	mc := NewTestClient(t)

	accountID := testtools.MERCHANT_ID
	walletID := testtools.MERCHANT_WALLET_ID

	wallet, err := mc.GetWallet(context.Background(), accountID, walletID)
	NoResponseError(t, err)
	require.NotNil(t, wallet)

	require.Equal(t, walletID, wallet.WalletID)
}

func TestListWalletTransactions(t *testing.T) {
	mc := NewTestClient(t)

	accountID := testtools.MERCHANT_ID
	walletID := testtools.MERCHANT_WALLET_ID

	transactions, err := mc.ListWalletTransactions(context.Background(), accountID, walletID)
	NoResponseError(t, err)
	require.NotEmpty(t, transactions)

	for i := range transactions {
		// Only check the first few transactions
		if i > 5 {
			break
		}
		require.Equal(t, walletID, transactions[i].WalletID)

		txn, err := mc.GetWalletTransaction(context.Background(), accountID, walletID, transactions[i].TransactionID)
		NoResponseError(t, err)
		require.Equal(t, transactions[i].TransactionID, txn.TransactionID)
	}
}
