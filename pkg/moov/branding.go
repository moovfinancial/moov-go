package moov

import (
	"context"
	"net/http"
)

type BrandColor struct {
	Accent string `json:"accent"`
}
type BrandColors struct {
	Dark  BrandColor `json:"dark"`
	Light BrandColor `json:"light"`
}
type Brand struct {
	Colors BrandColors `json:"colors"`
}

// CreateAccountBranding creates branding (colors) for the specified account.
// https://docs.moov.io/api/enrichment/branding/create/
func (c Client) CreateAccountBranding(ctx context.Context, accountID string, colorBrand Brand) (*Brand, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathBrandings, accountID),
		AcceptJson(),
		JsonBody(colorBrand))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Brand](resp)
}

// GetAccountBranding retrieves branding (colors) for the specified account.
// https://docs.moov.io/api/enrichment/branding/get/
func (c Client) GetAccountBranding(ctx context.Context, accountID string) (*Brand, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathBrandings, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Brand](resp)
}

// UpsertAccountBranding updates branding (colors) for the specified account.
// https://docs.moov.io/api/enrichment/branding/put/
func (c Client) UpsertAccountBranding(ctx context.Context, accountID string, colorBrand Brand) (*Brand, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathBrandings, accountID),
		AcceptJson(),
		JsonBody(colorBrand))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Brand](resp)
}
