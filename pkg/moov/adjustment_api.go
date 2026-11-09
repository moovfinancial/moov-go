package moov

import (
	"context"
	"net/http"
)

// ListAdjustments lists adjustments for a Moov account
// https://docs.moov.io/api/moov-accounts/adjustments/list/
func (c Client) ListAdjustments(ctx context.Context, accountID string, filters ...AdjustmentListFilter) ([]Adjustment, error) {
	args := prependArgs(filters, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathAdjustments, accountID), args...)
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Adjustment](resp)
}

// GetAdjustment retrieves an adjustment for a Moov account
// https://docs.moov.io/api/moov-accounts/adjustments/get/
func (c Client) GetAdjustment(ctx context.Context, accountID, adjustmentID string) (*Adjustment, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathAdjustment, accountID, adjustmentID),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Adjustment](resp)
}
