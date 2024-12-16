package moov

import (
	"context"
	"net/http"
)

// These API's currently require contacting support with a use-case to enable them.

// List your applications so you can get the ID's to create keys on them.
func (c Client) ListApplications(ctx context.Context) ([]Application, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathApplications),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Application](resp)
}

// This allows for creating an API key from another API key.
func (c Client) CreateApplicationKeys(ctx context.Context, applicationID string, key CreateApplicationKey) (*ApplicationKeyWithSecret, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathApplicationKeys, applicationID),
		AcceptJson(),
		JsonBody(key))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[ApplicationKeyWithSecret](resp)
}
