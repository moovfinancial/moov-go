package moov_test

import (
	"context"
	"testing"

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

	t.Cleanup(func() {
		if createdWebhook != nil {
			mc.DeleteWebhook(context.Background(), createdWebhook.WebhookID)
			createdWebhook = nil
		}
	})

	t.Run("list webhooks", func(t *testing.T) {
		webhooks, err := mc.ListWebhooks(context.Background())
		require.NoError(t, err)

		var gotWebhookID string
		for _, webhook := range webhooks {
			if webhook.WebhookID == createdWebhook.WebhookID {
				gotWebhookID = createdWebhook.WebhookID
			}
		}

		if gotWebhookID == "" {
			require.Failf(t, "webhook %s not found in list", createdWebhook.WebhookID)
		}
	})

	t.Run("get webhook", func(t *testing.T) {
		webhook, err := mc.GetWebhook(context.Background(), createdWebhook.WebhookID)
		require.NoError(t, err)
		require.NotNil(t, webhook)
		require.Equal(t, createdWebhook.WebhookID, webhook.WebhookID)
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
