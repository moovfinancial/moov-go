package moov_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

func Test_Transfers(t *testing.T) {
	mc := NewTestClient(t)

	account := getLincolnBank(t, mc)

	options, err := mc.TransferOptions(BgCtx(), moov.CreateTransferOptions{
		Source: moov.CreateTransferOptionsTarget{
			AccountID: account.AccountID,
		},
		Destination: moov.CreateTransferOptionsTarget{
			AccountID: FACILITATOR_ID,
		},
		Amount: moov.Amount{
			Currency: "USD",
			Value:    1,
		},
	})

	NoResponseError(t, err)
	source, dest := paymentMethodsFromOptions(t, options, moov.PaymentMethodType_AchDebitFund, moov.PaymentMethodType_MoovWallet)

	t.Run("make async transfer", func(t *testing.T) {
		started, err := mc.CreateTransfer(BgCtx(), moov.CreateTransfer{
			Source: moov.CreateTransfer_Source{
				PaymentMethodID: source,
			},
			Destination: moov.CreateTransfer_Destination{
				PaymentMethodID: dest,
			},
			Amount: moov.Amount{
				Currency: "usd",
				Value:    1,
			},
		}).Started()
		NoResponseError(t, err)

		// We made an async transfer, so completed should be nil, while started not nil
		require.NotNil(t, started)
	})

	t.Run("make sync transfer", func(t *testing.T) {
		completed, started, err := mc.CreateTransfer(BgCtx(), moov.CreateTransfer{
			Source: moov.CreateTransfer_Source{
				PaymentMethodID: source,
			},
			Destination: moov.CreateTransfer_Destination{
				PaymentMethodID: dest,
			},
			Amount: moov.Amount{
				Currency: "usd",
				Value:    1,
			},
		}).WaitForRailResponse()
		NoResponseError(t, err)

		// We made a sync transfer, so completed should be not nil, while started is nil
		require.NotNil(t, completed)
		require.Nil(t, started)
	})

	t.Run("list transfers", func(t *testing.T) {
		transfers, err := mc.ListTransfers(BgCtx(),
			moov.WithTransferAccountIDs([]string{
				account.AccountID,
			}),
			// moov.WithTransferCount(10),
			// moov.WithTransferStatus(string(moov.TransferStatus_Completed)),
		)
		require.NoError(t, err)
		require.Greater(t, len(transfers), 0)
	})
}
