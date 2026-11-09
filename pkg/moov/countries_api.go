package moov

import (
	"context"
	"net/http"
)

// GetAccountCountries retrieves the countries where an account can operate
// https://docs.moov.io/api/moov-accounts/countries/get/
func (c Client) GetAccountCountries(ctx context.Context, accountID string) (*AccountCountries, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathAccountCountries, accountID),
		AcceptJson(),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[AccountCountries](resp)
}

// AssignAccountCountries replaces the countries where an account can operate
// https://docs.moov.io/api/moov-accounts/countries/assign/
func (c Client) AssignAccountCountries(ctx context.Context, accountID string, countries AccountCountries) (*AccountCountries, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathAccountCountries, accountID),
		AcceptJson(),
		JsonBody(countries),
	)
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[AccountCountries](resp)
}
