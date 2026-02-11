package moov

import (
	"context"
	"net/http"
)

// GetFees lists Fees for a Moov account
// https://docs.moov.io/api/moov-accounts/billing/get-fees/
func (c Client) GetFees(ctx context.Context, accountID string, filters ...FeeGetFilter) ([]IncurredFee, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathFees, accountID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[IncurredFee](resp)
}

// ListFees fetches Fees for a Moov account by feeIDs
// https://docs.moov.io/api/moov-accounts/billing/list-fees/
func (c Client) ListFees(ctx context.Context, accountID string, request FeeListRequest) ([]IncurredFee, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathFeesFetch, accountID),
		AcceptJson(),
		JsonBody(request))
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[IncurredFee](resp)
}

// ListFeeRevenue retrieves revenue generated from merchant fees (for partners)
// https://docs.moov.io/api/moov-accounts/partner-billing/list-fee-revenue/
func (c Client) ListFeeRevenue(ctx context.Context, accountID string, filters ...FeeRevenueFilter) ([]IncurredFee, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathFeeRevenue, accountID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[IncurredFee](resp)
}
