package mhooks

import "github.com/moovfinancial/moov-go/pkg/moov"

type EventType string

const (
	EventTypeAccountCreated           EventType = "account.created"
	EventTypeAccountDeleted           EventType = "account.deleted"
	EventTypeAccountUpdated           EventType = "account.updated"
	EventTypeBalanceUpdated           EventType = "balance.updated"
	EventTypeBankAccountCreated       EventType = "bankAccount.created"
	EventTypeBankAccountDeleted       EventType = "bankAccount.deleted"
	EventTypeBankAccountUpdated       EventType = "bankAccount.updated"
	EventTypeCardAutoUpdated          EventType = "card.autoUpdated"
	EventTypeCapabilityRequested      EventType = "capability.requested"
	EventTypeCapabilityUpdated        EventType = "capability.updated"
	EventTypeDisputeCreated           EventType = "dispute.created"
	EventTypeDisputeUpdated           EventType = "dispute.updated"
	EventTypeNetworkIDUpdated         EventType = "networkID.updated"
	EventTypePaymentMethodDisabled    EventType = "paymentMethod.disabled"
	EventTypePaymentMethodEnabled     EventType = "paymentMethod.enabled"
	EventTypeRefundCreated            EventType = "refund.created"
	EventTypeRefundUpdated            EventType = "refund.updated"
	EventTypeRepresentativeCreated    EventType = "representative.created"
	EventTypeRepresentativeDeleted    EventType = "representative.deleted"
	EventTypeRepresentativeUpdated    EventType = "representative.updated"
	EventTypeTransferCreated          EventType = "transfer.created"
	EventTypeTransferUpdated          EventType = "transfer.updated"
	EventTypeWalletTransactionUpdated EventType = "walletTransaction.updated"
)

type AccountCreated struct {
	// ID of the account
	AccountID string `json:"accountID"`
}

type TransferCreated struct {
	// ID of the facilitator account
	AccountID string `json:"accountID"`
	// ID of the transfer
	TransferID string `json:"transferID"`
	// Status of the transfer
	Status moov.TransferStatus `json:"status"`
}

type TransferUpdated struct {
	// ID of the facilitator account
	AccountID string `json:"accountID,omitempty"`
	// ID of the transfer
	TransferID string `json:"transferID,omitempty"`
	// Status of the transfer
	Status      string               `json:"status,omitempty"`
	Source      PaymentMethodPartial `json:"source,omitempty"`
	Destination PaymentMethodPartial `json:"destination,omitempty"`
}

type PaymentMethodPartial struct {
	// ID of the account
	AccountID string `json:"accountID,omitempty"`
	// ID of the payment method
	PaymentMethodID string `json:"paymentMethodID,omitempty"`
}

// TODO(vince,4/25/2024): I'll add the rest of the models in upcoming PRs.
