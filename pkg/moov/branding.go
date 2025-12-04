package moov

import (
	"context"
	"net/http"
)

type Brand struct {
	Colors struct {
		Dark struct {
			Accent string `json:"accent"`
		} `json:"dark"`
		Light struct {
			Accent string `json:"accent"`
		} `json:"light"`
	} `json:"colors"`
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

// PatchAccountBranding updates branding (colors) for the specified account.
// https://docs.moov.io/api/enrichment/branding/patch/
func (c Client) PatchAccountBranding(ctx context.Context, accountID string, colorBrand Brand) (*Brand, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPatch, pathBrandings, accountID),
		AcceptJson(),
		JsonBody(colorBrand))
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
