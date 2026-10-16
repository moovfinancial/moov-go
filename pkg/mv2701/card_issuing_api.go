package mv2701

import (
	"context"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

// CardIssuingClient calls the card issuing API pinned to moov.Version2027_01.
type CardIssuingClient struct{ *moov.Client }

// NewCardIssuingClient wraps client in a CardIssuingClient.
func NewCardIssuingClient(client *moov.Client) CardIssuingClient {
	return CardIssuingClient{Client: client}
}

// CreateIssuedCard requests a new card be issued for the given account.
// https://docs.moov.io/api/money-movement/issuing/create/
func (c CardIssuingClient) CreateIssuedCard(ctx context.Context, accountID string, card CreateIssuedCard) (*moov.IssuedCard, error) {
	return moov.CreateIssuedCardGeneric[CreateIssuedCard, moov.IssuedCard](ctx, c.Client, moov.Version2027_01, accountID, card)
}
