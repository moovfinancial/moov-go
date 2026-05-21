package mv2607

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type TransferClient struct {
	*moov.Client
}

func NewTransferClient(client *moov.Client) TransferClient {
	return TransferClient{Client: client}
}

func (t TransferClient) CreateTransfer(ctx context.Context, partnerAccountID string, transfer CreateTransfer, options ...moov.CreateTransferArgs) moov.CreateTransferGenericBuilder[Transfer] {
	return moov.CreateTransferGeneric[CreateTransfer, Transfer](ctx, t.Client, moov.Version2026_07, partnerAccountID, transfer, options...)
}

func (t TransferClient) ListTransfers(ctx context.Context, accountID string, filters ...moov.ListTransferFilter) ([]Transfer, error) {
	return moov.ListTransfersGeneric[Transfer](ctx, t.Client, moov.Version2026_07, accountID, filters...)
}

func (t TransferClient) GetTransfer(ctx context.Context, accountID, transferID string) (*Transfer, error) {
	return moov.GetTransferGeneric[Transfer](ctx, t.Client, moov.Version2026_07, accountID, transferID)
}

func (t TransferClient) RefundTransfer(ctx context.Context, partnerAccountID, transferID string, refund CreateRefund, options ...moov.CreateRefundArgs) (*Refund, *moov.RefundStarted, error) {
	return moov.RefundTransferGeneric[CreateRefund, Refund](ctx, t.Client, moov.Version2026_07, partnerAccountID, transferID, refund, options...)
}

func (t TransferClient) ListRefunds(ctx context.Context, accountID, transferID string) ([]Refund, error) {
	return moov.ListRefundsGeneric[Refund](ctx, t.Client, moov.Version2026_07, accountID, transferID)
}

func (t TransferClient) GetRefund(ctx context.Context, accountID, transferID, refundID string) (*Refund, error) {
	return moov.GetRefundGeneric[Refund](ctx, t.Client, moov.Version2026_07, accountID, transferID, refundID)
}

func (t TransferClient) ReverseTransfer(ctx context.Context, partnerAccountID, transferID string, reversal CreateReversal, options ...moov.CreateReversalArgs) (*CreatedReversal, error) {
	return moov.ReverseTransferGeneric[CreateReversal, CreatedReversal](ctx, t.Client, moov.Version2026_07, partnerAccountID, transferID, reversal, options...)
}
