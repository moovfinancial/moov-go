package moov

import "time"

// WebhookEventType The type of event that occurred.
type WebhookEventType string

// List of WebhookEventType
const (
	WEBHOOKEVENTTYPE_STAR                         WebhookEventType = "*"
	WEBHOOKEVENTTYPE_ACCOUNT_CREATED              WebhookEventType = "account.created"
	WEBHOOKEVENTTYPE_ACCOUNT_UPDATED              WebhookEventType = "account.updated"
	WEBHOOKEVENTTYPE_ACCOUNT_DISCONNECTED         WebhookEventType = "account.disconnected"
	WEBHOOKEVENTTYPE_BALANCE_UPDATED              WebhookEventType = "balance.updated"
	WEBHOOKEVENTTYPE_BANK_ACCOUNT_CREATED         WebhookEventType = "bankAccount.created"
	WEBHOOKEVENTTYPE_BANK_ACCOUNT_UPDATED         WebhookEventType = "bankAccount.updated"
	WEBHOOKEVENTTYPE_BANK_ACCOUNT_DELETED         WebhookEventType = "bankAccount.deleted"
	WEBHOOKEVENTTYPE_CANCELLATION_CREATED         WebhookEventType = "cancellation.created"
	WEBHOOKEVENTTYPE_CANCELLATION_UPDATED         WebhookEventType = "cancellation.updated"
	WEBHOOKEVENTTYPE_CARD_AUTO_UPDATED            WebhookEventType = "card.autoUpdated"
	WEBHOOKEVENTTYPE_CAPABILITY_REQUESTED         WebhookEventType = "capability.requested"
	WEBHOOKEVENTTYPE_CAPABILITY_UPDATED           WebhookEventType = "capability.updated"
	WEBHOOKEVENTTYPE_DISPUTE_CREATED              WebhookEventType = "dispute.created"
	WEBHOOKEVENTTYPE_DISPUTE_UPDATED              WebhookEventType = "dispute.updated"
	WEBHOOKEVENTTYPE_INVOICE_CREATED              WebhookEventType = "invoice.created"
	WEBHOOKEVENTTYPE_INVOICE_UPDATED              WebhookEventType = "invoice.updated"
	WEBHOOKEVENTTYPE_NETWORK_ID_UPDATED           WebhookEventType = "networkID.updated"
	WEBHOOKEVENTTYPE_PAYMENT_METHOD_ENABLED       WebhookEventType = "paymentMethod.enabled"
	WEBHOOKEVENTTYPE_PAYMENT_METHOD_DISABLED      WebhookEventType = "paymentMethod.disabled"
	WEBHOOKEVENTTYPE_REFUND_CREATED               WebhookEventType = "refund.created"
	WEBHOOKEVENTTYPE_REFUND_UPDATED               WebhookEventType = "refund.updated"
	WEBHOOKEVENTTYPE_REPRESENTATIVE_CREATED       WebhookEventType = "representative.created"
	WEBHOOKEVENTTYPE_REPRESENTATIVE_UPDATED       WebhookEventType = "representative.updated"
	WEBHOOKEVENTTYPE_REPRESENTATIVE_DELETED       WebhookEventType = "representative.deleted"
	WEBHOOKEVENTTYPE_SWEEP_CREATED                WebhookEventType = "sweep.created"
	WEBHOOKEVENTTYPE_SWEEP_UPDATED                WebhookEventType = "sweep.updated"
	WEBHOOKEVENTTYPE_TERMINAL_APPLICATION_CREATED WebhookEventType = "terminalApplication.created"
	WEBHOOKEVENTTYPE_TERMINAL_APPLICATION_UPDATED WebhookEventType = "terminalApplication.updated"
	WEBHOOKEVENTTYPE_TICKET_CREATED               WebhookEventType = "ticket.created"
	WEBHOOKEVENTTYPE_TICKET_UPDATED               WebhookEventType = "ticket.updated"
	WEBHOOKEVENTTYPE_TICKET_MESSAGE_ADDED         WebhookEventType = "ticket.messageAdded"
	WEBHOOKEVENTTYPE_TRANSFER_CREATED             WebhookEventType = "transfer.created"
	WEBHOOKEVENTTYPE_TRANSFER_UPDATED             WebhookEventType = "transfer.updated"
	WEBHOOKEVENTTYPE_WALLET_CREATED               WebhookEventType = "wallet.created"
	WEBHOOKEVENTTYPE_WALLET_UPDATED               WebhookEventType = "wallet.updated"
	WEBHOOKEVENTTYPE_WALLET_TRANSACTION_UPDATED   WebhookEventType = "walletTransaction.updated"
	WEBHOOKEVENTTYPE_BILLING_STATEMENT_CREATED    WebhookEventType = "billingStatement.created"
)

const (
	WebhookStatusEnabled  WebhookStatus = "enabled"
	WebhookStatusDisabled WebhookStatus = "disabled"
)

type WebhookStatus string

type CreateWebhook struct {
	URL                  string             `json:"url"`
	Description          string             `json:"description"`
	Status               WebhookStatus      `json:"status"`
	EventTypeIdentifiers []WebhookEventType `json:"eventTypes"`
}

type UpdateWebhook struct {
	URL                  string             `json:"url"`
	Description          string             `json:"description"`
	Status               WebhookStatus      `json:"status"`
	EventTypeIdentifiers []WebhookEventType `json:"eventTypes"`
}

type Webhook struct {
	WebhookID   string             `json:"webhookID" otel:"webhook_id"`
	URL         string             `json:"url" otel:"url"`
	Description string             `json:"description,omitempty" otel:"description"`
	Status      WebhookStatus      `json:"status" otel:"status"`
	LastUsedOn  *time.Time         `json:"lastUsedOn,omitempty" otel:"last_used_on"`
	EventTypes  []WebhookEventType `json:"eventTypes" otel:"event_types"`
	CreatedOn   time.Time          `json:"createdOn" otel:"created_on"`
	UpdatedOn   time.Time          `json:"updatedOn" otel:"updated_on"`
}
