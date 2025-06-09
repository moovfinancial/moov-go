package moov

import (
	"context"
	"net/http"
)

// WithInstitutionName filters institutions by their name
func WithInstitutionName(name string) ListInstitutionsFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["name"] = name
		return nil
	})
}

// WithInstitutionRoutingNumber filters institutions by their name
func WithInstitutionRoutingNumber(routingNumber string) ListInstitutionsFilter {
	return callBuilderFn(func(call *callBuilder) error {
		call.params["routingNumber"] = routingNumber
		return nil
	})
}

// WithInstitutionLimit filters institutions by their name
func WithInstitutionLimit(limit int) ListInstitutionsFilter {
	return Limit(limit)
}

type ListInstitutionsFilter callArg

// SearchInstitutions will return financial institution information for ACH, RTP, and Wire payment rails based on name and routing number searches.
func (c Client) SearchInstitutions(ctx context.Context, opts ...ListInstitutionsFilter) (*InstitutionsSearchResponse, error) {
	args := prependArgs(opts, AcceptJson())
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathInstitutions), args...)
	if err != nil {
		return nil, err
	}
	return CompletedObjectOrError[InstitutionsSearchResponse](resp)
}
