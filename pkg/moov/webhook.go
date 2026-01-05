package moov

import (
	"context"
	"net/http"
)

func (c *Client) CreateWebhook(ctx context.Context, webhook CreateWebhook) (*Webhook, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathWebhooks),
		AcceptJson(),
		JsonBody(webhook))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Webhook](resp)
}

func (c *Client) ListWebhooks(ctx context.Context) ([]Webhook, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathWebhooks),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Webhook](resp)
}

func (c *Client) GetWebhook(ctx context.Context, webhookID string) (*Webhook, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathWebhook, webhookID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Webhook](resp)
}

func (c *Client) UpdateWebhook(ctx context.Context, webhookID string, webhook UpdateWebhook) (*Webhook, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPut, pathWebhook, webhookID),
		AcceptJson(),
		JsonBody(webhook))
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Webhook](resp)
}

func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) error {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodDelete, pathWebhook, webhookID),
		AcceptJson())
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
