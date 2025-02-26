package moov

import (
	"time"
)

type CreateReceipt struct {
	Kind           string  `json:"kind"`
	ForID          string  `json:"forId"`
	Email          *string `json:"email"`
	EmailAccountID *string `json:"emailAccountId"`
}

type Receipt struct {
	// ID is the unique identifier for the receipt
	ID string `json:"id"`

	// AccountID is the ID of the account that created the receipt
	CreatedBy string `json:"createdBy"`

	// kind of receipt, simple sale, signature, business copy, etc...
	Kind string `json:"kind"`

	// ForID is the ID of the schedule or transfer that this receipt is for
	ForID string `json:"forId"`

	// recipient ...
	Email          *string `json:"email"`
	EmailAccountID *string `json:"emailAccountId"`

	// DisabledOn is the date and time the receipt was disabled
	DisabledOn *time.Time `json:"disabledOn"`

	// SentFor is the list of receipts that have been sent
	SentFor []SentReceipt `json:"sentFor"`
}

type SentReceipt struct {
	ReceiptID      string    `json:"receiptId"`
	IdempotencyKey string    `json:"idempotencyKey"`
	SentOn         time.Time `json:"sentOn"`
}
