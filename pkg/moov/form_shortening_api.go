package moov

import (
	"context"
	"net/http"
)

// ListIndustries returns a list of industries that can be used for form shortening / auto-fill.
// https://docs.moov.io/api/v2026.01.00/enrichment/form-shortening/list-industries/
func (c Client) ListIndustries(ctx context.Context) (*Industries, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathIndustries),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Industries](resp)
}

// EnrichBusinessProfile returns publicly available business information for the given email address,
// which can be used to pre-fill form fields and reduce user input errors.
// https://docs.moov.io/api/v2026.01.00/enrichment/form-shortening/enrich-business-profile/
func (c Client) EnrichBusinessProfile(ctx context.Context, email string) (*EnrichedBusinessProfile, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathEnrichmentProfile),
		AcceptJson(),
		callBuilderFn(func(call *callBuilder) error {
			call.params["email"] = email
			return nil
		}))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[EnrichedBusinessProfile](resp)
}
