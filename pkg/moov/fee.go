package moov

import (
	"context"
	"net/http"
)

// ListFees lists Fees for a Moov account
// https://docs.moov.io/api/moov-accounts/billing/get-fees/
func (c Client) ListFees(ctx context.Context, accountID string, filters ...FeeListFilter) ([]IncurredFee, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathFees, accountID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[IncurredFee](resp)
}

// FetchFees fetches Fees for a Moov account by feeIDs
// https://docs.moov.io/api/moov-accounts/billing/list-fees/
func (c Client) FetchFees(ctx context.Context, accountID string, request FeeFetchRequest) ([]IncurredFee, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathFeesFetch, accountID),
		AcceptJson(),
		JsonBody(request))
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[IncurredFee](resp)
}
