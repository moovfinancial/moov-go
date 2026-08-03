package transfers

import (
	"context"
	"testing"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/moovfinancial/moov-go/pkg/mv2604"
	"github.com/stretchr/testify/require"
)

const (
	stagingPartnerAccountID = "db04bf9d-91f6-4206-ba38-6844636532ad"
	stagingWalletPMID       = "76dbae77-8839-4a5e-93fb-877a1041b09e"
)

func TestCreateAndPatchTransfer(t *testing.T) {
	mc, err := moov.NewClient()
	require.NoError(t, err)
	transferClientV2604 := mv2604.NewTransferClient(mc)

	var (
		ctx                   = context.Background()
		partnerAccountID      = stagingPartnerAccountID
		sourcePaymentMethodID = "155221bd-abf7-4a35-afe4-744646094280"
		destPaymentMethodID   = stagingWalletPMID
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

func TestCreatePublicTransferTypes(t *testing.T) {
	testCases := []struct {
		name                       string
		sourcePaymentMethodID      string
		destinationPaymentMethodID string
		cancel                     bool
	}{
		{
			name:                       "wallet-to-bank",
			sourcePaymentMethodID:      stagingWalletPMID,
			destinationPaymentMethodID: "3d8c580e-424b-4bf4-8255-544f94231de2",
		},
		{
			name:                       "bank-to-bank",
			sourcePaymentMethodID:      "155221bd-abf7-4a35-afe4-744646094280",
			destinationPaymentMethodID: "a65511b8-9a12-43d9-8b68-d3ee8df3faa9",
			cancel:                     true,
		},
		{
			name:                       "wallet-to-wallet",
			sourcePaymentMethodID:      "a5ab9116-26b2-42a0-a18e-5320c4d67b0e",
			destinationPaymentMethodID: "13af02c1-6115-48a9-a19b-f0d48e17b130",
		},
		{
			name:                       "wallet-to-rtp",
			sourcePaymentMethodID:      "a5ab9116-26b2-42a0-a18e-5320c4d67b0e",
			destinationPaymentMethodID: "f9823aca-f3a2-4122-a967-c28d177931ad",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc, err := moov.NewClient()
			require.NoError(t, err)

			ctx := context.Background()
			started, err := mc.CreateTransfer(ctx, stagingPartnerAccountID, moov.CreateTransfer{
				Source:      moov.CreateTransfer_Source{PaymentMethodID: tc.sourcePaymentMethodID},
				Destination: moov.CreateTransfer_Destination{PaymentMethodID: tc.destinationPaymentMethodID},
				Amount:      moov.Amount{Currency: "USD", Value: 100},
				Description: "moov-go auth-capture staging test: " + tc.name,
			}).Started()
			require.NoError(t, err)
			require.NotEmpty(t, started.TransferID)
			t.Logf("Created %s transfer: %s", tc.name, started.TransferID)

			if tc.cancel {
				cancellation, err := mc.CancelTransfer(ctx, stagingPartnerAccountID, started.TransferID)
				require.NoError(t, err)
				require.NotEmpty(t, cancellation.CancellationID)
				require.NotEqual(t, moov.CancellationStatus_Failed, cancellation.Status)
				t.Logf("Canceled transfer %s with cancellation %s (%s)", started.TransferID, cancellation.CancellationID, cancellation.Status)
			}
		})
	}
}

func TestCreateCardToWalletAndRefund(t *testing.T) {
	mc, err := moov.NewClient()
	require.NoError(t, err)

	ctx := context.Background()
	transfer, started, err := mc.CreateTransfer(ctx, stagingPartnerAccountID, moov.CreateTransfer{
		Source: moov.CreateTransfer_Source{
			PaymentMethodID: "f55ceaad-2f5b-4d76-8973-14a68798fe33",
			CardDetails: &moov.CreateTransfer_CardDetailsSource{
				DynamicDescriptor: "AUTHCAPTURE TEST",
			},
		},
		Destination: moov.CreateTransfer_Destination{PaymentMethodID: stagingWalletPMID},
		Amount:      moov.Amount{Currency: "USD", Value: 100},
		Description: "moov-go auth-capture staging test: card-to-wallet",
	}).WaitForRailResponse()
	require.NoError(t, err)
	require.Nil(t, started, "timed out waiting for card transfer rail response: %+v", started)
	require.NotNil(t, transfer)
	require.NotEmpty(t, transfer.TransferID)
	t.Logf("Created card-to-wallet transfer: %s (%s)", transfer.TransferID, transfer.Status)

	refund, refundStarted, err := mc.RefundTransfer(
		ctx,
		stagingPartnerAccountID,
		transfer.TransferID,
		moov.CreateRefund{Amount: 100},
		moov.WithRefundWaitForRailResponse(),
	)
	require.NoError(t, err)
	require.Nil(t, refundStarted, "timed out waiting for refund rail response: %+v", refundStarted)
	require.NotNil(t, refund)
	require.NotEmpty(t, refund.RefundID)
	t.Logf("Refunded transfer %s with refund %s (%s)", transfer.TransferID, refund.RefundID, refund.Status)
}
