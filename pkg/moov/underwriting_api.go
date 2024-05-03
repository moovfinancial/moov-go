package moov

import (
	"context"
	"net/http"
)

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

func (c Client) GetUnderwriting(ctx context.Context, accountID string) (*Underwriting, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathUnderwriting, accountID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Underwriting](resp)
}
