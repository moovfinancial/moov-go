package moov_test

import (
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_PaymentMethods(t *testing.T) {
	mc := NewTestClient(t)

	account := getLincolnBank(t, mc)

	paymentMethods, err := mc.ListPaymentMethods(BgCtx(), account.AccountID)

	t.Run("list", func(t *testing.T) {
		NoResponseError(t, err)
		require.NotEmpty(t, paymentMethods)
	})

	t.Run("list with method type filter", func(t *testing.T) {
		filtered, err := mc.ListPaymentMethods(BgCtx(),
			account.AccountID,
			moov.WithPaymentMethodType(string(paymentMethods[0].PaymentMethodType)),
		)
		NoResponseError(t, err)
		require.NotEmpty(t, filtered)
	})

	t.Run("list with sourceId filter", func(t *testing.T) {

		// find any sourceId we can look up an account by
		sourceId := ""
		for _, pm := range paymentMethods {
			if sourceId != "" {
				continue
			}

			if pm.Wallet != nil {
				sourceId = pm.Wallet.WalletID
			} else if pm.ApplePay != nil {
				sourceId = pm.BankAccount.BankAccountID
			} else if pm.Card != nil {
				sourceId = pm.Card.CardID
			}
		}
		require.NotEmpty(t, sourceId)

		filtered, err := mc.ListPaymentMethods(BgCtx(),
			account.AccountID,
			moov.WithPaymentMethodSourceID(sourceId),
		)
		NoResponseError(t, err)
		require.NotEmpty(t, filtered)
	})

	t.Run("get", func(t *testing.T) {
		cap, err := mc.GetPaymentMethod(BgCtx(), account.AccountID, paymentMethods[0].PaymentMethodID)
		NoResponseError(t, err)
		require.NotNil(t, cap)
	})
}
