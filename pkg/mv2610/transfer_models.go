package mv2610

import (
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

// TransferEvent is an event in a Transfer's ordered timeline.
// Amounts are positive display context and are not additive across events.
type TransferEvent struct {
	EventID    string               `json:"eventID"`
	EventType  TransferEventType    `json:"eventType"`
	Amount     moov.AmountDecimal   `json:"amount"`
	Details    TransferEventDetails `json:"details"`
	OccurredOn time.Time            `json:"occurredOn"`
}

// TransferEventType is the event family. It names the TransferEventDetails field that is set.
type TransferEventType string

// List of TransferEventType
const (
	TransferEventType_Transfer          TransferEventType = "transfer"
	TransferEventType_Authorization     TransferEventType = "authorization"
	TransferEventType_Capture           TransferEventType = "capture"
	TransferEventType_Refund            TransferEventType = "refund"
	TransferEventType_Cancellation      TransferEventType = "cancellation"
	TransferEventType_Dispute           TransferEventType = "dispute"
	TransferEventType_AchCredit         TransferEventType = "ach-credit"
	TransferEventType_AchDebit          TransferEventType = "ach-debit"
	TransferEventType_InstantBankCredit TransferEventType = "instant-bank-credit"
	TransferEventType_CardPayment       TransferEventType = "card-payment"
	TransferEventType_PushToCard        TransferEventType = "push-to-card"
	TransferEventType_PullFromCard      TransferEventType = "pull-from-card"
	TransferEventType_WalletCredit      TransferEventType = "wallet-credit"
	TransferEventType_WalletDebit       TransferEventType = "wallet-debit"
)

// TransferEventDetails holds the details for a Transfer timeline event.
// A valid event contains exactly one detail object that matches its event type.
type TransferEventDetails struct {
	Transfer          *TransferEventTransferDetails          `json:"transfer,omitempty"`
	Authorization     *TransferEventAuthorizationDetails     `json:"authorization,omitempty"`
	Capture           *TransferEventCaptureDetails           `json:"capture,omitempty"`
	Refund            *TransferEventRefundDetails            `json:"refund,omitempty"`
	Cancellation      *TransferEventCancellationDetails      `json:"cancellation,omitempty"`
	Dispute           *TransferEventDisputeDetails           `json:"dispute,omitempty"`
	AchCredit         *TransferEventAchDetails               `json:"achCredit,omitempty"`
	AchDebit          *TransferEventAchDetails               `json:"achDebit,omitempty"`
	InstantBankCredit *TransferEventInstantBankCreditDetails `json:"instantBankCredit,omitempty"`
	CardPayment       *TransferEventCardPaymentDetails       `json:"cardPayment,omitempty"`
	PushToCard        *TransferEventPushToCardDetails        `json:"pushToCard,omitempty"`
	PullFromCard      *TransferEventPullFromCardDetails      `json:"pullFromCard,omitempty"`
	WalletCredit      *TransferEventWalletTransactionDetails `json:"walletCredit,omitempty"`
	WalletDebit       *TransferEventWalletTransactionDetails `json:"walletDebit,omitempty"`
}

type TransferEventTransferDetails struct {
	Status moov.TransferStatus `json:"status"`
}

type TransferEventAuthorizationDetails struct {
	AuthorizationID string                           `json:"authorizationID"`
	Status          TransferEventAuthorizationStatus `json:"status"`
	FailureCode     *moov.CardFailureCode            `json:"failureCode,omitempty"`
}

type TransferEventAuthorizationStatus string

// List of TransferEventAuthorizationStatus
const (
	TransferEventAuthorizationStatus_Approved TransferEventAuthorizationStatus = "approved"
	TransferEventAuthorizationStatus_Declined TransferEventAuthorizationStatus = "declined"
	TransferEventAuthorizationStatus_Reversed TransferEventAuthorizationStatus = "reversed"
)

type TransferEventCaptureDetails struct {
	CaptureID   string                `json:"captureID"`
	Status      CaptureStatus         `json:"status"`
	FailureCode *moov.CardFailureCode `json:"failureCode,omitempty"`
}

type CaptureStatus string

// List of CaptureStatus
const (
	CaptureStatus_Created   CaptureStatus = "created"
	CaptureStatus_Pending   CaptureStatus = "pending"
	CaptureStatus_Completed CaptureStatus = "completed"
	CaptureStatus_Failed    CaptureStatus = "failed"
	CaptureStatus_Canceled  CaptureStatus = "canceled"
)

type TransferEventRefundDetails struct {
	RefundID    string                `json:"refundID"`
	CaptureID   *string               `json:"captureID,omitempty"`
	Status      moov.RefundStatus     `json:"status"`
	FailureCode *moov.CardFailureCode `json:"failureCode,omitempty"`
}

type TransferEventCancellationDetails struct {
	CancellationID string                  `json:"cancellationID"`
	Status         moov.CancellationStatus `json:"status"`
}

type TransferEventDisputeDetails struct {
	DisputeID string             `json:"disputeID"`
	CaptureID *string            `json:"captureID,omitempty"`
	Phase     moov.DisputePhase  `json:"phase"`
	Status    moov.DisputeStatus `json:"status"`
}

type TransferEventAchDetails struct {
	Status     moov.AchStatus     `json:"status"`
	Return     *moov.AchException `json:"return,omitempty"`
	Correction *moov.AchException `json:"correction,omitempty"`
}

type TransferEventInstantBankCreditDetails struct {
	Status      moov.InstantBankStatus       `json:"status"`
	FailureCode *moov.InstantBankFailureCode `json:"failureCode,omitempty"`
}

type TransferEventCardPaymentDetails struct {
	Status      moov.CardTransactionStatus `json:"status"`
	FailureCode *moov.CardFailureCode      `json:"failureCode,omitempty"`
}

type TransferEventPushToCardDetails struct {
	Status      PushToCardTransactionStatus `json:"status"`
	FailureCode *moov.CardFailureCode       `json:"failureCode,omitempty"`
}

type PushToCardTransactionStatus string

// List of PushToCardTransactionStatus
const (
	PushToCardTransactionStatus_Initiated PushToCardTransactionStatus = "initiated"
	PushToCardTransactionStatus_Deferred  PushToCardTransactionStatus = "deferred"
	PushToCardTransactionStatus_Canceled  PushToCardTransactionStatus = "canceled"
	PushToCardTransactionStatus_Failed    PushToCardTransactionStatus = "failed"
	PushToCardTransactionStatus_Completed PushToCardTransactionStatus = "completed"
)

type TransferEventPullFromCardDetails struct {
	Status      PullFromCardTransactionStatus `json:"status"`
	FailureCode *moov.CardFailureCode         `json:"failureCode,omitempty"`
}

type PullFromCardTransactionStatus string

// List of PullFromCardTransactionStatus
const (
	PullFromCardTransactionStatus_Initiated PullFromCardTransactionStatus = "initiated"
	PullFromCardTransactionStatus_Failed    PullFromCardTransactionStatus = "failed"
	PullFromCardTransactionStatus_Completed PullFromCardTransactionStatus = "completed"
)

type TransferEventWalletTransactionDetails struct {
	AccountID       string                     `json:"accountID"`
	WalletID        string                     `json:"walletID"`
	TransactionID   string                     `json:"transactionID"`
	TransactionType moov.WalletTransactionType `json:"transactionType"`
	CaptureID       *string                    `json:"captureID,omitempty"`
}
