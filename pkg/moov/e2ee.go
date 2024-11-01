package moov

import (
	"context"
	"net/http"

	jose "github.com/go-jose/go-jose/v4"
)

type EndToEndToken struct {
	Token string `json:"token"`
}

func (c Client) GenerateEndToEndPublicKey(ctx context.Context) (*jose.JSONWebKey, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathEndToEndPublicKey))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return UnmarshalObjectResponse[jose.JSONWebKey](resp)
	default:
		return nil, resp
	}
}

func (c Client) TestEndToEndToken(ctx context.Context, token string) error {
	body := &struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathEndToEndTokenTest), JsonBody(body))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
