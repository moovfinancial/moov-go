package mhooks

import (
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

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
	EventTypeTestPing                 EventType = "event.test"
	EventTypeTransferCreated          EventType = "transfer.created"
	EventTypeTransferUpdated          EventType = "transfer.updated"
	EventTypeWalletTransactionUpdated EventType = "walletTransaction.updated"
)

type AccountCreated struct {
	// ID of the account
	AccountID string `json:"accountID"`
}

type AccountDeleted struct {
	// ID of the account
	AccountID string `json:"accountID"`
}

type AccountUpdated struct {
	// ID of the account
	AccountID string `json:"accountID"`
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
}

type CapabilityUpdated struct {
	Capability moov.CapabilityName `json:"capabilityID"`
	// ID of the account requesting the capability
	AccountID string `json:"accountID"`
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

type TestPing struct {
	Ping bool `json:"ping"`
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
