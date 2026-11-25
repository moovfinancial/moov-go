package mv2610

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

// TransferClient calls the transfers API pinned to moov.Version2026_10.
type TransferClient struct{ *moov.Client }

// NewTransferClient wraps client in a TransferClient.
func NewTransferClient(client *moov.Client) TransferClient {
	return TransferClient{Client: client}
}

func (t TransferClient) CreateTransfer(ctx context.Context, partnerAccountID string, transfer CreateTransfer, options ...moov.CreateTransferArgs) moov.CreateTransferGenericBuilder[Transfer] {
	return moov.CreateTransferGeneric[CreateTransfer, Transfer](ctx, t.Client, moov.Version2026_10, partnerAccountID, transfer, options...)
}

// ListTransfers lists transfers. Use moov.WithTransferTypes to filter by TransferType.
func (t TransferClient) ListTransfers(ctx context.Context, accountID string, filters ...moov.ListTransferFilter) ([]Transfer, error) {
	return moov.ListTransfersGeneric[Transfer](ctx, t.Client, moov.Version2026_10, accountID, filters...)
}

func (t TransferClient) GetTransfer(ctx context.Context, accountID, transferID string) (*Transfer, error) {
	return moov.GetTransferGeneric[Transfer](ctx, t.Client, moov.Version2026_10, accountID, transferID)
}

// FetchTransfers retrieves multiple transfers by their IDs, keyed by transfer ID.
func (t TransferClient) FetchTransfers(ctx context.Context, accountID string, request moov.FetchTransfersRequest) (map[string]Transfer, error) {
	return moov.FetchTransfersGeneric[Transfer](ctx, t.Client, moov.Version2026_10, accountID, request)
}

// ListTransferEvents lists the transfer's event timeline, oldest to newest.
// A transfer with no stored events returns an empty list.
func (t TransferClient) ListTransferEvents(ctx context.Context, accountID, transferID string) ([]TransferEvent, error) {
	return moov.ListTransferEventsGeneric[TransferEvent](ctx, t.Client, moov.Version2026_10, accountID, transferID)
}

// RefundTransfer returns the refund for a 200 and for a 202 response. Without
// moov.WithRefundWaitForRailResponse, the refund has no status or processing details yet.
func (t TransferClient) RefundTransfer(ctx context.Context, partnerAccountID, transferID string, refund CreateRefund, options ...moov.CreateRefundArgs) (*Refund, error) {
	return moov.CreateRefundGeneric[CreateRefund, Refund](ctx, t.Client, moov.Version2026_10, partnerAccountID, transferID, refund, options...)
}

func (t TransferClient) ListRefunds(ctx context.Context, accountID, transferID string) ([]Refund, error) {
	return moov.ListRefundsGeneric[Refund](ctx, t.Client, moov.Version2026_10, accountID, transferID)
}

func (t TransferClient) GetRefund(ctx context.Context, accountID, transferID, refundID string) (*Refund, error) {
	return moov.GetRefundGeneric[Refund](ctx, t.Client, moov.Version2026_10, accountID, transferID, refundID)
}

func (t TransferClient) ReverseTransfer(ctx context.Context, partnerAccountID, transferID string, reversal CreateReversal, options ...moov.CreateReversalArgs) (*CreatedReversal, error) {
	return moov.ReverseTransferGeneric[CreateReversal, CreatedReversal](ctx, t.Client, moov.Version2026_10, partnerAccountID, transferID, reversal, options...)
}

// CancelTransfer cancels a transfer. The client sends a random idempotency key unless
// moov.WithCancellationIdempotencyKey sets one.
func (t TransferClient) CancelTransfer(ctx context.Context, accountID, transferID string, cancellation CreateCancellation, options ...moov.CreateCancellationArgs) (*moov.Cancellation, error) {
	return moov.CancelTransferGeneric[CreateCancellation, moov.Cancellation](ctx, t.Client, moov.Version2026_10, accountID, transferID, cancellation, options...)
}
