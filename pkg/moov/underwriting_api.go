package moov

import (
	"context"
	"net/http"
)

// Not able to use method generics, this is the closest I can get
type UnderwritingClient[T any, V any] struct {
	client  *Client // embedded client for making HTTP requests
	version string
}

func (c Client) Underwriting() UnderwritingClient[UpdateUnderwriting, Underwriting] {
	return UnderwritingClient[UpdateUnderwriting, Underwriting]{
		client:  &c,
		version: "",
	}
}

func (c Client) UnderwritingV2507() UnderwritingClient[UpsertUnderwriting, UnderwritingV2507] {
	return UnderwritingClient[UpsertUnderwriting, UnderwritingV2507]{
		client:  &c,
		version: Q3_2025,
	}
}

func (uc UnderwritingClient[T, V]) Get(ctx context.Context, accountID string) (*V, error) {
	var zero *V
	httpResp, err := uc.client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathUnderwriting, accountID),
		Version(uc.version),
		AcceptJson())
	if err != nil {
		return zero, err
	}

	resp, err := CompletedObjectOrError[V](httpResp)
	if err != nil {
		return zero, err
	}

	return resp, nil
}

// generics with function
type UnderwritingParams struct {
	version string
}

type UnderwritingOption func(p *UnderwritingParams)

func WithVersion(version string) UnderwritingOption {
	return func(p *UnderwritingParams) {
		p.version = version
	}
}

func UnderwritingGet[T any](ctx context.Context, client Client, accountID string, options ...UnderwritingOption) (*T, error) {
	var zero *T

	params := prepareOptions(options...)

	httpResp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathUnderwriting, accountID),
		Version(params.version),
		AcceptJson())
	if err != nil {
		return zero, err
	}

	resp, err := CompletedObjectOrError[T](httpResp)
	if err != nil {
		return zero, err
	}

	return resp, nil
}

func UnderwritingUpsert[T any, X any](ctx context.Context, client Client, accountID string, underwriting T, options ...UnderwritingOption) (*X, error) {
	var zero *X
	params := prepareOptions(options...)
	httpResp, err := client.CallHttp(ctx,
		Endpoint(http.MethodPut, pathUnderwriting, accountID),
		Version(params.version),
		AcceptJson(),
		JsonBody(underwriting))
	if err != nil {
		return nil, err
	}

	resp, err := CompletedObjectOrError[X](httpResp)
	if err != nil {
		return zero, err
	}

	return resp, nil
}

func prepareOptions(options ...UnderwritingOption) *UnderwritingParams {
	params := &UnderwritingParams{}
	for _, option := range options {
		option(params)
	}

	if params.version == "" {
		params.version = PreVersioning
	}

	return params
}

// UpsertUnderwriting adds or updates underwriting information for the given account.
// The account must have a description and an MCC set to create underwriting.
// Returns the underwriting information for the account.
func (c Client) UpsertUnderwriting(ctx context.Context, accountID string, underwriting UpdateUnderwriting) (*Underwriting, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathUnderwriting, accountID),
		AcceptJson(),
		JsonBody(underwriting))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Underwriting](resp)
}

// UpsertUnderwriting adds or updates underwriting information for the given account.
// Returns the underwriting information for the account.
func (c Client) UpsertUnderwritingV2507(ctx context.Context, accountID string, underwriting UpsertUnderwriting) (*UnderwritingV2507, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathUnderwriting, accountID),
		Version(Q3_2025),
		AcceptJson(),
		JsonBody(underwriting))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[UnderwritingV2507](resp)
}

// GetUnderwriting returns the underwriting information for the given account.
func (c Client) GetUnderwriting(ctx context.Context, accountID string) (*Underwriting, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathUnderwriting, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Underwriting](resp)
}

func (c Client) GetUnderwritingV2507(ctx context.Context, accountID string) (*UnderwritingV2507, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathUnderwriting, accountID),
		Version(Q3_2025),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[UnderwritingV2507](resp)
}
