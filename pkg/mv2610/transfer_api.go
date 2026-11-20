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

// ListTransferEvents lists the transfer's event timeline, oldest to newest.
// A transfer with no stored events returns an empty list.
func (t TransferClient) ListTransferEvents(ctx context.Context, accountID, transferID string) ([]TransferEvent, error) {
	return moov.ListTransferEventsGeneric[TransferEvent](ctx, t.Client, moov.Version2026_10, accountID, transferID)
}
