package mv2610

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

// CardIssuingClient calls the card issuing API pinned to moov.Version2026_10.
type CardIssuingClient struct{ *moov.Client }

// NewCardIssuingClient wraps client in a CardIssuingClient.
func NewCardIssuingClient(client *moov.Client) CardIssuingClient {
	return CardIssuingClient{Client: client}
}

// ListIssuedCardActivity lists issued card activity for the given account.
func (c CardIssuingClient) ListIssuedCardActivity(ctx context.Context, accountID string, filters ...moov.ListIssuedCardActivityFilter) ([]IssuedCardActivity, error) {
	return moov.ListIssuedCardActivityGeneric[IssuedCardActivity](ctx, c.Client, moov.Version2026_10, accountID, filters...)
}
