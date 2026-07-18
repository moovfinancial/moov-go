package mv2610

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
	return moov.CreateTransferGeneric[CreateTransfer, Transfer](ctx, t.Client, moov.Version2026_10, partnerAccountID, transfer, options...)
}

func (t TransferClient) ListTransfers(ctx context.Context, accountID string, filters ...moov.ListTransferFilter) ([]Transfer, error) {
	return moov.ListTransfersGeneric[Transfer](ctx, t.Client, moov.Version2026_10, accountID, filters...)
}

func (t TransferClient) GetTransfer(ctx context.Context, accountID, transferID string) (*Transfer, error) {
	return moov.GetTransferGeneric[Transfer](ctx, t.Client, moov.Version2026_10, accountID, transferID)
}

func (t TransferClient) ListCaptures(ctx context.Context, accountID, transferID string) ([]Capture, error) {
	return moov.ListCapturesGeneric[Capture](ctx, t.Client, moov.Version2026_10, accountID, transferID)
}

func (t TransferClient) GetCapture(ctx context.Context, accountID, transferID, captureID string) (*Capture, error) {
	return moov.GetCaptureGeneric[Capture](ctx, t.Client, moov.Version2026_10, accountID, transferID, captureID)
}

func (t TransferClient) RefundTransfer(ctx context.Context, partnerAccountID, transferID string, refund CreateRefund, options ...moov.CreateRefundArgs) (*CreateRefundResponse, *moov.RefundStarted, error) {
	return moov.RefundTransferGeneric[CreateRefund, CreateRefundResponse](ctx, t.Client, moov.Version2026_10, partnerAccountID, transferID, refund, options...)
}

func (t TransferClient) ListRefunds(ctx context.Context, accountID, transferID string) ([]CardAcquiringRefund, error) {
	return moov.ListRefundsGeneric[CardAcquiringRefund](ctx, t.Client, moov.Version2026_10, accountID, transferID)
}

func (t TransferClient) GetRefund(ctx context.Context, accountID, transferID, refundID string) (*CardAcquiringRefund, error) {
	return moov.GetRefundGeneric[CardAcquiringRefund](ctx, t.Client, moov.Version2026_10, accountID, transferID, refundID)
}

func (t TransferClient) ReverseTransfer(ctx context.Context, partnerAccountID, transferID string, reversal CreateReversal, options ...moov.CreateReversalArgs) (*Reversal, error) {
	return moov.ReverseTransferGeneric[CreateReversal, Reversal](ctx, t.Client, moov.Version2026_10, partnerAccountID, transferID, reversal, options...)
}
