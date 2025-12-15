package moov

import (
	"context"
	"net/http"
)

// CreateIssuedCard requests a new card be issued for the given account
// https://docs.moov.io/api/money-movement/issuing/create/
func (c Client) CreateIssuedCard(ctx context.Context, accountID string, card CreateIssuedCard) (*IssuedCard, error) {
	httpResp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathIssuedCards, accountID),
		AcceptJson(),
		JsonBody(card))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[IssuedCard](httpResp)
}

// ListIssuedCards retrieves all issued cards for the given account
// https://docs.moov.io/api/money-movement/issuing/list/
func (c Client) ListIssuedCards(ctx context.Context, accountID string) ([]IssuedCard, error) {
	httpResp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathIssuedCards, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[IssuedCard](httpResp)
}

// GetIssuedCard retrieves the specified issued card for the given account
// https://docs.moov.io/api/money-movement/issuing/get/
func (c Client) GetIssuedCard(ctx context.Context, accountID string, cardID string) (*IssuedCard, error) {
	httpResp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathIssuedCard, accountID, cardID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[IssuedCard](httpResp)
}

// UpdateIssuedCard updates a specified issued card for the given account
// https://docs.moov.io/api/money-movement/issuing/update/
func (c Client) UpdateIssuedCard(ctx context.Context, accountID string, cardID string, update UpdateIssuedCard) error {
	httpResp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPatch, pathIssuedCard, accountID, cardID),
		AcceptJson(),
		JsonBody(update))
	if err != nil {
		return err
	}

	return CompletedNilOrError(httpResp)
}
