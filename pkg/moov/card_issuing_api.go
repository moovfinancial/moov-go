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
func (c Client) ListIssuedCards(ctx context.Context, accountID string, filters ...ListIssuedCardsFilter) ([]IssuedCard, error) {
	httpResp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathIssuedCards, accountID),
		prependArgs(filters, AcceptJson())...)
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

// ListIssuedCardAuthorizations retrieves all issued card authorizations for the given account
// https://docs.moov.io/api/money-movement/issuing/list-authorizations/
func (c Client) ListIssuedCardAuthorizations(ctx context.Context, accountID string, filters ...ListIssuedCardAuthorizationsFilter) ([]IssuedCardAuthorization, error) {
	httpResp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathIssuingAuthorizations, accountID),
		prependArgs(filters, AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[IssuedCardAuthorization](httpResp)
}

// GetIssuedCardAuthorization retrieves the details of an issued card authorization for the given account
// https://docs.moov.io/api/money-movement/issuing/get-authorization/
func (c Client) GetIssuedCardAuthorization(ctx context.Context, accountID string, authorizationID string) (*IssuedCardAuthorization, error) {
	httpResp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathIssuingAuthorization, accountID, authorizationID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[IssuedCardAuthorization](httpResp)
}

// ListIssuedCardAuthorizationEvents retrieves all events that affect an issued card authorization for the given account
// https://docs.moov.io/api/money-movement/issuing/list-authorization-events/
func (c Client) ListIssuedCardAuthorizationEvents(ctx context.Context, accountID string, authorizationID string, filters ...ListIssuedCardAuthorizationEventsFilter) ([]IssuedCardAuthorizationEvent, error) {
	httpResp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathIssuingAuthorizationEvents, accountID, authorizationID),
		prependArgs(filters, AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[IssuedCardAuthorizationEvent](httpResp)
}

// ListIssuedCardAuthorizations retrieves all issued card transactions for the given account
// https://docs.moov.io/api/money-movement/issuing/list-card-transactions/
func (c Client) ListIssuedCardTransactions(ctx context.Context, accountID string, filters ...ListIssuedCardTransactionsFilter) ([]IssuedCardTransaction, error) {
	httpResp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathIssuingTransactions, accountID),
		prependArgs(filters, AcceptJson())...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[IssuedCardTransaction](httpResp)
}

// GetIssuedCardAuthorization retrieves the details of an issued card transaction for the given account
// https://docs.moov.io/api/money-movement/issuing/get-card-transaction/
func (c Client) GetIssuedCardTransaction(ctx context.Context, accountID string, cardTransactionID string) (*IssuedCardTransaction, error) {
	httpResp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathIssuingTransaction, accountID, cardTransactionID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[IssuedCardTransaction](httpResp)
}
