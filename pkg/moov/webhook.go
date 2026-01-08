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

// PingWebhook sends a test event to the webhook URL to verify it is working.
// The ping will send an event with type `event.test` and a null data payload.
func (c *Client) PingWebhook(ctx context.Context, webhookID string) (*WebhookPing, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodPost, pathWebhookPing, webhookID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[WebhookPing](resp)
}

// GetWebhookSecret retrieves the signing secret for a webhook.
// Use this secret to verify the signature of incoming webhook payloads.
func (c *Client) GetWebhookSecret(ctx context.Context, webhookID string) (*WebhookSecret, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathWebhookSecret, webhookID),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[WebhookSecret](resp)
}

// ListWebhookEventTypes returns all available webhook event types that can be subscribed to.
func (c *Client) ListWebhookEventTypes(ctx context.Context) ([]WebhookEventType, error) {
	resp, err := c.CallHttp(ctx,
		Endpoint(http.MethodGet, pathEventTypes),
		AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[WebhookEventType](resp)
}
