package moov

import (
	"context"
	"net/http"
)

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

type UnderwritingClient[T any, V any] struct {
	Version string
}

func (uc UnderwritingClient[T, V]) Get(ctx context.Context, client Client, accountID string) (*V, error) {
	httpResp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathUnderwriting, accountID),
		Version(uc.Version),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[V](httpResp)
}

func (uc UnderwritingClient[T, V]) Upsert(ctx context.Context, client Client, accountID string, update T) (*V, error) {
	httpResp, err := client.CallHttp(ctx,
		Endpoint(http.MethodPut, pathUnderwriting, accountID),
		Version(uc.Version),
		AcceptJson(),
		JsonBody(update))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[V](httpResp)
}
