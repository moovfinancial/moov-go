package mv2610

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type CardIssuingClient struct{ *moov.Client }

func NewCardIssuingClient(client *moov.Client) CardIssuingClient {
	return CardIssuingClient{Client: client}
}

func (c CardIssuingClient) ListIssuedCardActivity(ctx context.Context, accountID string, filters ...moov.ListIssuedCardActivityFilter) ([]IssuedCardActivity, error) {
	return moov.ListIssuedCardActivityGeneric[IssuedCardActivity](ctx, c.Client, moov.Version2026_10, accountID, filters...)
}
