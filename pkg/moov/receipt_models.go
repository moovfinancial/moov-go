package moov

import (
	"time"
)

type CreateReceipt struct {
	Kind string `json:"kind"`

	// lookup
	ForTransferID   *string `json:"forTransferID"`
	ForScheduleID   *string `json:"forScheduleID"`
	ForOccurrenceID *string `json:"forOccurrenceID"`

	// recipient
	Email          *string `json:"email"`
	EmailAccountID *string `json:"emailAccountID"`
}

type Receipt struct {
	// ID is the unique identifier for the receipt
	ID string `json:"receiptID"`

	// AccountID is the ID of the account that created the receipt
	CreatedBy string `json:"createdBy"`

	// kind of receipt, simple sale, signature, business copy, etc...
	Kind string `json:"kind"`

	// ForID is the ID of the schedule or transfer or whatever that this receipt is watching and will trigger for
	ForTransferID   *string `json:"forTransferID"`
	ForScheduleID   *string `json:"forScheduleID"`
	ForOccurrenceID *string `json:"forOccurrenceID"`

	// recipient ...
	Email          *string `json:"email"`
	EmailAccountID *string `json:"emailAccountID"`

	// DisabledOn is the date and time the receipt was disabled
	DisabledOn *time.Time `json:"disabledOn"`

	// SentFor is the list of receipts that have been sent
	SentFor []SentReceipt `json:"sentFor"`
}

type SentReceipt struct {
	ReceiptID      string    `json:"receiptID"`
	IdempotencyKey string    `json:"idempotencyKey"`
	SentOn         time.Time `json:"sentOn"`
}
