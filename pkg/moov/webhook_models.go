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

type UpsertWebhook struct {
	URL         string        `json:"url"`
	Description string        `json:"description"`
	Status      WebhookStatus `json:"status"`
	EventTypes  []EventType   `json:"eventTypes"`
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
