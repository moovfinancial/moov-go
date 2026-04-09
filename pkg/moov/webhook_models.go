package moov

import (
	"time"
)

type EventType string

const (
	EventTypeAccountCreated           EventType = "account.created"
	EventTypeAccountDisconnected      EventType = "account.disconnected"
	EventTypeAccountUpdated           EventType = "account.updated"
	EventTypeBalanceUpdated           EventType = "balance.updated"
	EventTypeBankAccountCreated       EventType = "bankAccount.created"
	EventTypeBankAccountDeleted       EventType = "bankAccount.deleted"
	EventTypeBankAccountUpdated       EventType = "bankAccount.updated"
	EventTypeBillingStatementCreated  EventType = "billingStatement.created"
	EventTypeCancellationCreated      EventType = "cancellation.created"
	EventTypeCancellationUpdated      EventType = "cancellation.updated"
	EventTypeCardAutoUpdated          EventType = "card.autoUpdated"
	EventTypeCapabilityRequested      EventType = "capability.requested"
	EventTypeCapabilityUpdated        EventType = "capability.updated"
	EventTypeDisputeCreated           EventType = "dispute.created"
	EventTypeDisputeUpdated           EventType = "dispute.updated"
	EventTypeInvoiceCreated           EventType = "invoice.created"
	EventTypeInvoiceUpdated           EventType = "invoice.updated"
	EventTypeNetworkIDUpdated         EventType = "networkID.updated"
	EventTypePaymentMethodDisabled    EventType = "paymentMethod.disabled"
	EventTypePaymentMethodEnabled     EventType = "paymentMethod.enabled"
	EventTypeRefundCreated            EventType = "refund.created"
	EventTypeRefundUpdated            EventType = "refund.updated"
	EventTypeRepresentativeCreated    EventType = "representative.created"
	EventTypeRepresentativeDeleted    EventType = "representative.deleted"
	EventTypeRepresentativeUpdated    EventType = "representative.updated"
	EventTypeSweepCreated             EventType = "sweep.created"
	EventTypeSweepUpdated             EventType = "sweep.updated"
	EventTypeTestPing                 EventType = "event.test"
	EventTypeTicketCreated            EventType = "ticket.created"
	EventTypeTicketUpdated            EventType = "ticket.updated"
	EventTypeTicketMessageAdded       EventType = "ticket.messageAdded"
	EventTypeTransferCreated          EventType = "transfer.created"
	EventTypeTransferUpdated          EventType = "transfer.updated"
	EventTypeWalletCreated            EventType = "wallet.created"
	EventTypeWalletUpdated            EventType = "wallet.updated"
	EventTypeWalletTransactionUpdated EventType = "walletTransaction.updated"
)

const (
	WebhookStatusEnabled  WebhookStatus = "enabled"
	WebhookStatusDisabled WebhookStatus = "disabled"
)

type WebhookStatus string

// CreateWebhook is the payload for creating a new webhook.
type CreateWebhook struct {
	// URL is the destination for webhook events. Must be a valid URL.
	URL string `json:"url"`
	// Description is an optional description for the webhook.
	Description string `json:"description"`
	// Status is the webhook status (enabled or disabled).
	Status WebhookStatus `json:"status"`
	// EventTypes is the list of event types to subscribe to.
	EventTypes []EventType `json:"eventTypes"`
}

// UpdateWebhook is the payload for updating an existing webhook.
type UpdateWebhook struct {
	// URL is the destination for webhook events. Must be a valid URL.
	URL string `json:"url"`
	// Description is an optional description for the webhook.
	Description string `json:"description"`
	// Status is the webhook status (enabled or disabled).
	Status WebhookStatus `json:"status"`
	// EventTypes is the list of event types to subscribe to.
	EventTypes []EventType `json:"eventTypes"`
}

type Webhook struct {
	WebhookID   string        `json:"webhookID"`
	URL         string        `json:"url"`
	Description string        `json:"description,omitempty"`
	Status      WebhookStatus `json:"status"`
	LastUsedOn  *time.Time    `json:"lastUsedOn,omitempty"`
	EventTypes  []EventType   `json:"eventTypes"`
	CreatedOn   time.Time     `json:"createdOn"`
	UpdatedOn   time.Time     `json:"updatedOn"`
}

// WebhookSecret contains the signing secret for verifying webhook payloads.
type WebhookSecret struct {
	Secret string `json:"secret"`
}

// WebhookPing is returned when pinging a webhook endpoint.
type WebhookPing struct {
	// The webhook that was pinged.
	Webhook Webhook `json:"webhook"`
	// The request body sent to the target URL. It will contain an event type of `event.test` and an empty (null) data payload.
	RequestBodySent map[string]any `json:"requestBodySent"`
	// The response status code after sending a ping event to the URL.
	ResponseStatusCode int32 `json:"responseStatusCode"`
}

// WebhookEventType describes a webhook event type that can be subscribed to.
type WebhookEventType struct {
	EventTypeID string    `json:"eventTypeID"`
	Type        EventType `json:"type"`
	Description string    `json:"description"`
}
