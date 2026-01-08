package moov_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
	"github.com/stretchr/testify/require"
)

func TestCreateWebhook(t *testing.T) {
	mc := NewTestClient(t)

	createWebhook := moov.CreateWebhook{
		URL:         "https://example.com/webhook",
		Description: "Test webhook",
		Status:      moov.WebhookStatusEnabled,
		EventTypes: []moov.EventType{
			moov.EventTypeAccountCreated,
			moov.EventTypeAccountDisconnected,
			moov.EventTypeAccountUpdated,
		},
	}

	createdWebhook, err := mc.CreateWebhook(context.Background(), createWebhook)
	require.NoError(t, err)
	require.NotNil(t, createdWebhook)

	normalizeWebhook(createdWebhook)

	t.Cleanup(func() {
		if createdWebhook != nil {
			mc.DeleteWebhook(context.Background(), createdWebhook.WebhookID)
			createdWebhook = nil
		}
	})

	t.Run("list webhooks", func(t *testing.T) {
		webhooks, err := mc.ListWebhooks(context.Background())
		for i := range webhooks {
			normalizeWebhook(&webhooks[i])
		}

		require.NoError(t, err)
		require.NotNil(t, webhooks)
		require.Contains(t, webhooks, *createdWebhook)
	})

	t.Run("get webhook", func(t *testing.T) {
		webhook, err := mc.GetWebhook(context.Background(), createdWebhook.WebhookID)
		normalizeWebhook(webhook)
		require.NoError(t, err)
		require.NotNil(t, webhook)
		require.Equal(t, createdWebhook, webhook)
	})

	t.Run("get webhook secret", func(t *testing.T) {
		secret, err := mc.GetWebhookSecret(context.Background(), createdWebhook.WebhookID)
		require.NoError(t, err)
		require.NotEmpty(t, secret.Secret)
	})

	t.Run("update webhook", func(t *testing.T) {
		webhook, err := mc.UpdateWebhook(context.Background(), createdWebhook.WebhookID, moov.UpdateWebhook{
			URL:         "https://example.com/webhook-new",
			Description: "Test webhook new",
			Status:      moov.WebhookStatusEnabled,
			EventTypes: []moov.EventType{
				moov.EventTypeAccountCreated,
				moov.EventTypeAccountDisconnected,
				moov.EventTypeAccountUpdated,
			},
		})
		normalizeWebhook(webhook)
		require.NoError(t, err)
		require.NotNil(t, webhook)
		require.Equal(t, "https://example.com/webhook-new", webhook.URL)
		require.Equal(t, "Test webhook new", webhook.Description)
	})

	t.Run("delete webhook", func(t *testing.T) {
		err := mc.DeleteWebhook(context.Background(), createdWebhook.WebhookID)
		require.NoError(t, err)
	})
}

func normalizeWebhook(webhook *moov.Webhook) {
	if webhook == nil {
		return
	}

	webhook.CreatedOn = webhook.CreatedOn.Truncate(time.Microsecond)
	webhook.UpdatedOn = webhook.UpdatedOn.Truncate(time.Microsecond)
	slices.Sort(webhook.EventTypes)
}
