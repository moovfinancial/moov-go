package moov

import (
	"context"
	"net/http"
)

type UnderwritingClient[T any, V any] struct {
	Version Version
}

// GetUnderwriting returns the underwriting information for the given account.
func (uc UnderwritingClient[T, V]) Get(ctx context.Context, client Client, accountID string) (*V, error) {
	httpResp, err := client.CallHttp(ctx,
		Endpoint(http.MethodGet, pathUnderwriting, accountID),
		MoovVersion(uc.Version),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[V](httpResp)
}

// Upsert adds or updates underwriting information for the given account.
// Returns the underwriting information for the account.
func (uc UnderwritingClient[T, V]) Upsert(ctx context.Context, client Client, accountID string, update T) (*V, error) {
	httpResp, err := client.CallHttp(ctx,
		Endpoint(http.MethodPost, pathUnderwriting, accountID),
		MoovVersion(uc.Version),
		AcceptJson(),
		JsonBody(update))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[V](httpResp)
}

// Legacy

// Only use for Preversioned API calls. Use mvxxxx.Underwriting.Upsert(...) instead.
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

// Only use for Preversioned API calls. Use mvxxxx.Underwriting.Get(...) instead.
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
