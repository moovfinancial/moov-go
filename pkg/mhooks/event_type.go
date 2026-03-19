package mhooks

import (
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

type EventType = moov.EventType

const (
	EventTypeAccountCreated           = moov.EventTypeAccountCreated
	EventTypeAccountDisconnected      = moov.EventTypeAccountDisconnected
	EventTypeAccountUpdated           = moov.EventTypeAccountUpdated
	EventTypeBalanceUpdated           = moov.EventTypeBalanceUpdated
	EventTypeBankAccountCreated       = moov.EventTypeBankAccountCreated
	EventTypeBankAccountDeleted       = moov.EventTypeBankAccountDeleted
	EventTypeBankAccountUpdated       = moov.EventTypeBankAccountUpdated
	EventTypeBillingStatementCreated  = moov.EventTypeBillingStatementCreated
	EventTypeCancellationCreated      = moov.EventTypeCancellationCreated
	EventTypeCancellationUpdated      = moov.EventTypeCancellationUpdated
	EventTypeCardAutoUpdated          = moov.EventTypeCardAutoUpdated
	EventTypeCapabilityRequested      = moov.EventTypeCapabilityRequested
	EventTypeCapabilityUpdated        = moov.EventTypeCapabilityUpdated
	EventTypeDisputeCreated           = moov.EventTypeDisputeCreated
	EventTypeDisputeUpdated           = moov.EventTypeDisputeUpdated
	EventTypeInvoiceCreated           = moov.EventTypeInvoiceCreated
	EventTypeInvoiceUpdated           = moov.EventTypeInvoiceUpdated
	EventTypeNetworkIDUpdated         = moov.EventTypeNetworkIDUpdated
	EventTypePaymentMethodDisabled    = moov.EventTypePaymentMethodDisabled
	EventTypePaymentMethodEnabled     = moov.EventTypePaymentMethodEnabled
	EventTypeRefundCreated            = moov.EventTypeRefundCreated
	EventTypeRefundUpdated            = moov.EventTypeRefundUpdated
	EventTypeRepresentativeCreated    = moov.EventTypeRepresentativeCreated
	EventTypeRepresentativeDeleted    = moov.EventTypeRepresentativeDeleted
	EventTypeRepresentativeUpdated    = moov.EventTypeRepresentativeUpdated
	EventTypeSweepCreated             = moov.EventTypeSweepCreated
	EventTypeSweepUpdated             = moov.EventTypeSweepUpdated
	EventTypeTestPing                 = moov.EventTypeTestPing
	EventTypeTicketCreated            = moov.EventTypeTicketCreated
	EventTypeTicketUpdated            = moov.EventTypeTicketUpdated
	EventTypeTicketMessageAdded       = moov.EventTypeTicketMessageAdded
	EventTypeTransferCreated          = moov.EventTypeTransferCreated
	EventTypeTransferUpdated          = moov.EventTypeTransferUpdated
	EventTypeWalletCreated            = moov.EventTypeWalletCreated
	EventTypeWalletUpdated            = moov.EventTypeWalletUpdated
	EventTypeWalletTransactionUpdated = moov.EventTypeWalletTransactionUpdated
)

type AccountCreated struct {
	// ID of the account
	AccountID string `json:"accountID"`
	ForeignID string `json:"foreignID,omitempty"`
}

type AccountDisconnected struct {
	// ID of the account
	AccountID string `json:"accountID"`
	ForeignID string `json:"foreignID,omitempty"`
}

type AccountUpdated struct {
	// ID of the account
	AccountID string `json:"accountID"`
	ForeignID string `json:"foreignID,omitempty"`
}

type BalanceUpdated struct {
	// ID of the Account associated with the wallet
	AccountID string `json:"accountID"`
	// ID of the Wallet
	WalletID string `json:"walletID"`
}

type BankAccountCreated struct {
	// ID of the bank account
	BankAccountID string `json:"bankAccountID"`
	// ID of the account where the bank account was created
	AccountID string `json:"accountID"`
}

type BankAccountDeleted struct {
	// ID of the bank account
	BankAccountID string `json:"bankAccountID"`
	// ID of the account where the bank account was deleted
	AccountID string `json:"accountID"`
}

type BankAccountUpdated struct {
	// ID of the bank account
	BankAccountID string `json:"bankAccountID"`
	// ID of the account where the bank account was updated
	AccountID string `json:"accountID"`
	// Status of the bank account
	Status           moov.BankAccountStatus       `json:"status"`
	StatusReason     moov.BankAccountStatusReason `json:"statusReason"`
	ExceptionDetails *moov.ExceptionDetails       `json:"exceptionDetails,omitempty"`
}

type BillingStatementCreated struct {
	StatementID string `json:"statementID"`
}

type CancellationCreated struct {
	CancellationID string                  `json:"cancellationID"`
	TransferID     string                  `json:"transferID"`
	Status         moov.CancellationStatus `json:"status"`
}

type CancellationUpdated struct {
	CancellationID string                  `json:"cancellationID"`
	TransferID     string                  `json:"transferID"`
	Status         moov.CancellationStatus `json:"status"`
}

type CardAutoUpdated struct {
	// ID of the merchant's Account associated with the updated Card
	AccountID string `json:"accountID"`
	// ID of the Card that was updated
	CardID     string              `json:"cardID"`
	UpdateType moov.CardUpdateType `json:"updateType"`
}

type CapabilityRequested struct {
	Capability moov.CapabilityName `json:"capabilityID"`
	// ID of the account requesting the capability
	AccountID string `json:"accountID"`
	ForeignID string `json:"foreignID,omitempty"`
}

type CapabilityUpdated struct {
	Capability moov.CapabilityName `json:"capabilityID"`
	// ID of the account requesting the capability
	AccountID string `json:"accountID"`
	ForeignID string `json:"foreignID,omitempty"`
	// Status of the capability
	Status moov.CapabilityStatus `json:"status"`
}

type DisputeCreated struct {
	// ID of the merchant's Account associated with the disputed transaction.
	AccountID string `json:"accountID"`
	// ID of the disputed transfer.
	TransferID string `json:"transferID"`
	// ID of the disputed transaction.
	TransactionID string `json:"transactionID"`
	// ID of the dispute.
	DisputeID string       `json:"disputeID"`
	Status    moov.Dispute `json:"status"`
	Phase     moov.Dispute `json:"phase"`
}

type DisputeUpdated struct {
	// ID of the merchant's Account associated with the disputed transaction.
	AccountID string `json:"accountID"`
	// ID of the disputed transfer.
	TransferID string `json:"transferID"`
	// ID of the disputed transaction.
	TransactionID string `json:"transactionID"`
	// ID of the dispute.
	DisputeID string             `json:"disputeID"`
	Status    moov.DisputeStatus `json:"status"`
	Phase     moov.DisputePhase  `json:"phase"`
}
type InvoiceCreated struct {
	AccountID string `json:"accountID"`
	InvoiceID string `json:"invoiceID"`
}

type InvoiceUpdated struct {
	AccountID string             `json:"accountID"`
	InvoiceID string             `json:"invoiceID"`
	Status    moov.InvoiceStatus `json:"status"`
}

type NetworkIDUpdated struct {
	// ID of account.
	AccountID     string     `json:"accountID"`
	VisaMid       *string    `json:"visaMid,omitempty"`
	MastercardMid *string    `json:"mastercardMid,omitempty"`
	DiscoverMid   *string    `json:"discoverMid,omitempty"`
	AmexMid       *string    `json:"amexMid,omitempty"`
	UpdatedOn     *time.Time `json:"updatedOn,omitempty"`
}

type PaymentMethodDisabled struct {
	// ID of the payment method
	PaymentMethodID string `json:"paymentMethodID"`
	// ID of the account
	AccountID string `json:"accountID"`
	// ID of the bank account, card, or wallet
	SourceID string `json:"sourceID"`
}

type PaymentMethodEnabled struct {
	// ID of the payment method
	PaymentMethodID string `json:"paymentMethodID"`
	// ID of the account
	AccountID string `json:"accountID"`
	// ID of the bank account, card, or wallet
	SourceID string `json:"sourceID"`
}

type RefundCreated struct {
	// ID of the merchant's Account associated with the refund transfer
	AccountID string `json:"accountID"`
	// ID of the original transfer
	TransferID string `json:"transferID"`
	// ID of the refund transfer
	RefundID string `json:"refundID"`
}

type RefundUpdated struct {
	// ID of the merchant's Account associated with the refund transfer
	AccountID string `json:"accountID"`
	// ID of the original transfer
	TransferID string `json:"transferID"`
	// ID of the refund transfer
	RefundID string `json:"refundID"`
	// Status of the refund transfer
	Status moov.RefundStatus `json:"status"`
}

type RepresentativeCreated struct {
	// ID of the representative
	RepresentativeID string `json:"representativeID"`
	// ID of the account
	AccountID string `json:"accountID"`
}

type RepresentativeDeleted struct {
	// ID of the representative
	RepresentativeID string `json:"representativeID"`
	// ID of the account
	AccountID string `json:"accountID"`
}

type RepresentativeUpdated struct {
	// ID of the representative
	RepresentativeID string `json:"representativeID"`
	// ID of the account
	AccountID string `json:"accountID"`
}

type SweepCreated struct {
	// ID of the sweep
	SweepID string `json:"sweepID"`
	// ID of the Wallet
	WalletID string `json:"walletID"`
}

// SweepUpdated struct for SweepUpdated
type SweepUpdated struct {
	// ID of the sweep
	SweepID string `json:"sweepID"`
	// ID of the Wallet
	WalletID string           `json:"walletID"`
	Status   moov.SweepStatus `json:"status"`
	// ID of the transfer
	TransferID *string `json:"transferID,omitempty"`
}

type TestPing struct {
	Ping bool `json:"ping"`
}

type TicketCreated struct {
	// ID of the account
	AccountID string `json:"accountID"`
	// ID of the ticket
	TicketID  string `json:"ticketID"`
	ForeignID string `json:"foreignID,omitempty"`
}

type TicketUpdated struct {
	// ID of the account
	AccountID string `json:"accountID"`
	// ID of the ticket
	TicketID  string `json:"ticketID"`
	ForeignID string `json:"foreignID,omitempty"`
	// Status of the ticket
	Status moov.TicketStatus `json:"status"`
}

type TicketMessageAdded struct {
	// ID of the account
	AccountID string `json:"accountID"`
	// ID of the ticket
	TicketID  string `json:"ticketID"`
	ForeignID string `json:"foreignID,omitempty"`
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
	AccountID string `json:"accountID"`
	// ID of the transfer
	TransferID string `json:"transferID"`
	// Status of the transfer
	Status      TransferUpdatedStatus `json:"status"`
	Source      PaymentMethodPartial  `json:"source"`
	Destination PaymentMethodPartial  `json:"destination"`
}

type TransferUpdatedStatus string

const (
	TransferUpdatedStatus_Pending   TransferUpdatedStatus = "pending"
	TransferUpdatedStatus_Completed TransferUpdatedStatus = "completed"
	TransferUpdatedStatus_Failed    TransferUpdatedStatus = "failed"
	TransferUpdatedStatus_Reversed  TransferUpdatedStatus = "reversed"
	TransferUpdatedStatus_Queued    TransferUpdatedStatus = "queued"
	TransferUpdatedStatus_Canceled  TransferUpdatedStatus = "canceled"

	TransferUpdatedStatus_SourceCompleted  TransferUpdatedStatus = "source.completed"
	TransferUpdatedStatus_SourceCorrected  TransferUpdatedStatus = "source.corrected"
	TransferUpdatedStatus_SourceInitiated  TransferUpdatedStatus = "source.initiated"
	TransferUpdatedStatus_SourceReturned   TransferUpdatedStatus = "source.returned"
	TransferUpdatedStatus_SourceOriginated TransferUpdatedStatus = "source.originated"
	TransferUpdatedStatus_SourceConfirmed  TransferUpdatedStatus = "source.confirmed"
	TransferUpdatedStatus_SourceSettled    TransferUpdatedStatus = "source.settled"
	TransferUpdatedStatus_SourceFailed     TransferUpdatedStatus = "source.failed"
	TransferUpdatedStatus_SourceCanceled   TransferUpdatedStatus = "source.canceled"

	TransferUpdatedStatus_DestinationCompleted              TransferUpdatedStatus = "destination.completed"
	TransferUpdatedStatus_DestinationCorrected              TransferUpdatedStatus = "destination.corrected"
	TransferUpdatedStatus_DestinationInitiated              TransferUpdatedStatus = "destination.initiated"
	TransferUpdatedStatus_DestinationReturned               TransferUpdatedStatus = "destination.returned"
	TransferUpdatedStatus_DestinationOriginated             TransferUpdatedStatus = "destination.originated"
	TransferUpdatedStatus_DestinationConfirmed              TransferUpdatedStatus = "destination.confirmed"
	TransferUpdatedStatus_DestinationFailed                 TransferUpdatedStatus = "destination.failed"
	TransferUpdatedStatus_DestinationAcceptedWithoutPosting TransferUpdatedStatus = "destination.accepted-without-posting"
)

type PaymentMethodPartial struct {
	// ID of the account
	AccountID string `json:"accountID"`
	// ID of the payment method
	PaymentMethodID string `json:"paymentMethodID"`
}

type WalletCreated struct {
	// ID of the account
	AccountID string `json:"accountID"`
	// ID of the wallet
	WalletID string `json:"walletID"`
}

type WalletUpdated struct {
	// ID of the account
	AccountID string `json:"accountID"`
	// ID of the wallet
	WalletID string `json:"walletID"`
	// Status of the wallet
	Status moov.WalletStatus `json:"status"`
}

type WalletTransactionUpdated struct {
	// ID of the account associated with the wallet
	AccountID string `json:"accountID"`
	// ID of the Wallet
	WalletID string `json:"walletID"`
	// ID of the wallet transaction
	TransactionID string `json:"transactionID"`
	// Status of the wallet transaction
	Status moov.WalletTransactionStatus `json:"status"`
	// Available balance of the wallet.
	AvailableBalance *moov.AvailableBalance `json:"availableBalance,omitempty"`
}
