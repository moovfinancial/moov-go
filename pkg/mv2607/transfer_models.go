package mv2607

import "github.com/moovfinancial/moov-go/pkg/moov"

type CreateTransfer struct {
	Source      moov.CreateTransfer_Source      `json:"source"`
	Destination moov.CreateTransfer_Destination `json:"destination"`

	Amount         moov.Amount                        `json:"amount"`
	SalesTaxAmount *moov.Amount                       `json:"salesTaxAmount,omitempty"`
	AmountDetails  *CreateTransferAmountDetails       `json:"amountDetails,omitempty"`
	FacilitatorFee moov.CreateTransfer_FacilitatorFee `json:"facilitatorFee,omitempty"`
	Description    string                             `json:"description,omitempty"`
	Metadata       map[string]string                  `json:"metadata,omitempty"`
	ForeignID      *string                            `json:"foreignID,omitempty"`
	LineItems      *moov.CreateTransferLineItems      `json:"lineItems,omitempty"`
}

type CreateTransferAmountDetails struct {
	TipAmount *moov.AmountDecimal `json:"tip,omitempty"`
	// Surcharge is an optional fee added to the transfer amount, typically used to cover additional costs such as card scheme fees. Transfer.Amount should be inclusive of any surcharge and represents the total amount charged to the customer. Note that surcharges must be enabled for the merchant account and are only valid on credit card payments.
	Surcharge *moov.AmountDecimal `json:"surcharge,omitempty"`
}

type Transfer struct {
	moov.Transfer
	AmountDetails *TransferAmountDetails `json:"amountDetails,omitempty"`
	Refunds       []Refund               `json:"refunds,omitempty"`
}

type TransferAmountDetails struct {
	TipAmount *moov.AmountDecimal `json:"tip,omitempty"`
	Surcharge *moov.AmountDecimal `json:"surcharge,omitempty"`
}

type CreateRefund struct {
	Amount        int64                `json:"amount,omitempty"`
	AmountDetails *RefundAmountDetails `json:"amountDetails,omitempty"`
}

type CreateReversal struct {
	Amount        int64                `json:"amount,omitempty"`
	AmountDetails *RefundAmountDetails `json:"amountDetails,omitempty"`
}

type Refund struct {
	moov.Refund
	AmountDetails *RefundAmountDetails `json:"amountDetails,omitempty"`
}

type RefundAmountDetails struct {
	// Surcharge is the amount of the original transfer's surcharge that was refunded.
	Surcharge *moov.AmountDecimal `json:"surcharge,omitempty"`
}

type CreatedReversal struct {
	Cancellation *moov.CreatedCancellation `json:"cancellation,omitempty"`
	Refund       *Refund                   `json:"refund,omitempty"`
}

type RefundStarted struct {
	Transfer
	RefundedTransferID string `otel:"refunded_transfer_id,omitempty"`
}

type TransferStarted = moov.TransferStarted
type CreateTransferOptions = moov.CreateTransferOptions
type CreateTransferOptionsTarget = moov.CreateTransferOptionsTarget
