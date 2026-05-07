package transfers

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2604"
	"github.com/stretchr/testify/require"
)

func TestCreateAndPatchTransfer(t *testing.T) {
	mc, err := moov.NewClient()
	require.NoError(t, err)
	transferClientV2604 := mv2604.NewTransferClient(mc)

	var (
		ctx                   = context.Background()
		partnerAccountID      = "ebbf46c6-122a-4367-bc45-7dd555e1d3b9"
		sourcePaymentMethodID = "b9a7153a-ad9b-4ada-b31b-f230fb8d1e8b" // ach-debit-fund
		destPaymentMethodID   = "67ebda6c-de48-474c-b49d-2cd3aa7d3f92" // moov-wallet
		initialForeignID      = "external-ref-123"
	)

	started, err := mc.CreateTransfer(ctx, partnerAccountID, moov.CreateTransfer{
		Source:      moov.CreateTransfer_Source{PaymentMethodID: sourcePaymentMethodID},
		Destination: moov.CreateTransfer_Destination{PaymentMethodID: destPaymentMethodID},
		Amount:      moov.Amount{Currency: "USD", Value: 100},
		Metadata:    map[string]string{"foo": "bar"},
		ForeignID:   &initialForeignID,
	}).Started()
	require.NoError(t, err)
	require.NotEmpty(t, started.TransferID)
	t.Logf("Created transfer: %+v", started)

	transfer, err := mc.GetTransfer(ctx, partnerAccountID, started.TransferID)
	require.NoError(t, err)
	require.NotNil(t, transfer.ForeignID)
	require.NotEmpty(t, transfer.Metadata)

	t.Run("v2604.PatchTransfer unsets foreignID and metadata", func(t *testing.T) {
		patched, err := transferClientV2604.PatchTransfer(ctx, partnerAccountID, started.TransferID, mv2604.PatchTransfer{
			ForeignID: moov.SetNull[string](),
			Metadata:  moov.SetNull[map[string]string](),
		})
		require.NoError(t, err)
		require.Nil(t, patched.ForeignID)
		require.Empty(t, patched.Metadata)
		t.Logf("unset foreignID and metadata: %+v", patched)

		fetched, err := mc.GetTransfer(ctx, partnerAccountID, started.TransferID)
		require.NoError(t, err)
		require.Nil(t, fetched.ForeignID)
		require.Empty(t, fetched.Metadata)
		t.Logf("got transfer with unset foreignID and metadata: %+v", fetched)
	})
}
