package moov_test

import (
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/internal/testtools"
	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func Test_Receipts(t *testing.T) {
	mc := NewTestClient(t)

	customer := CreateTemporaryTestAccount(t, mc, createTestIndividualAccount())
	customerCard := createTemporaryCard(t, mc, customer.AccountID)

	var transfer *moov.TransferStarted
	var receipt *moov.Receipt

	t.Run("make async transfer", func(t *testing.T) {
		transfer = createReceiptTestTransfer(t, mc, customerCard.PaymentMethods[0].PaymentMethodID)
	})

	t.Run("create receipt", func(t *testing.T) {
		receipt = createTestReceipt(t, mc, transfer.TransferID)
	})

	t.Run("list receipts", func(t *testing.T) {
		requireTestReceiptListed(t, mc, transfer.TransferID, receipt)
	})
}

func createReceiptTestTransfer(t *testing.T, mc *moov.Client, paymentMethodID string) *moov.TransferStarted {
	t.Helper()

	transfer, err := mc.CreateTransfer(BgCtx(), testtools.PARTNER_ID, moov.CreateTransfer{
		Source: moov.CreateTransfer_Source{
			PaymentMethodID: paymentMethodID,
		},
		Destination: moov.CreateTransfer_Destination{
			PaymentMethodID: testtools.MERCHANT_WALLET_PM_ID,
		},
		Amount: moov.Amount{
			Currency: "usd",
			Value:    1,
		},
	}).Started()
	NoResponseError(t, err)
	require.NotNil(t, transfer)

	require.Eventually(t, func() bool {
		_, err := mc.GetTransfer(BgCtx(), testtools.PARTNER_ID, transfer.TransferID)
		return err == nil
	}, 3*time.Second, 250*time.Millisecond)

	return transfer
}

func createTestReceipt(t *testing.T, mc *moov.Client, transferID string) *moov.Receipt {
	t.Helper()

	receipts, err := mc.CreateReceipt(BgCtx(), moov.CreateReceipt{
		Kind:          "sale.customer.v1",
		ForTransferID: &transferID,
		Email:         moov.PtrOf("noreply@moov.io"),
	})
	require.NoError(t, err)
	require.Len(t, receipts, 1)
	require.NotEmpty(t, receipts[0].ID)

	t.Logf("receipt: %+v\n", &receipts[0])
	return &receipts[0]
}

func requireTestReceiptListed(t *testing.T, mc *moov.Client, transferID string, receipt *moov.Receipt) {
	t.Helper()

	receipts, err := mc.ListReceipts(BgCtx(), moov.ReceiptByTransferID(transferID))
	require.NoError(t, err)
	require.Len(t, receipts, 1)

	for i := range receipts {
		receipts[i].SentFor = nil
	}

	require.Contains(t, receipts, *receipt)
}
