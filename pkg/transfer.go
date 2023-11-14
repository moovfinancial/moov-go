package moov

import "time"

type SynchronousTransfer struct {
	TransferID     string            `json:"transferID,omitempty"`
	CreatedOn      time.Time         `json:"createdOn,omitempty"`
	CompletedOn    time.Time         `json:"completedOn,omitempty"`
	Status         string            `json:"status,omitempty"`
	FailureReason  string            `json:"failureReason,omitempty"`
	Amoount        Amount            `json:"amount,omitempty"`
	Description    string            `json:"description,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	FacilitatorFee FacilitatorFee    `json:"facilitatorFee,omitempty"`
	MoovFee        int               `json:"moovFee,omitempty"`
	MoovFeeDecimal string            `json:"moovFeeDecimal,omitempty"`
	MoovFeeDetails MoovFeeDetails    `json:"moovFeeDetails,omitempty"`
	GroupID        string            `json:"groupID,omitempty"`
	RefundedAmount Amount            `json:"refundedAmount,omitempty"`
	Refunds        []Refund          `json:"refunds,omitempty"`
	DisputedAmount Amount            `json:"disputedAmount,omitempty"`
	Disputes       []Dispute         `json:"disputes,omitempty"`
	Source         Source            `json:"source,omitempty"`
	Destination    Destination       `json:"destination,omitempty"`
}

type AsynchronousTransfer struct {
	TransferID string    `json:"transferID,omitempty"`
	CreatedOn  time.Time `json:"createdOn,omitempty"`
}

type Amount struct {
	Currency string `json:"currency,omitempty"`
	Value    int    `json:"value,omitempty"`
}

type FacilitatorFee struct {
	Total         int    `json:"total,omitempty"`
	TotalDecimal  string `json:"totalDecimal,omitempty"`
	Markup        int    `json:"markup,omitempty"`
	MarkupDecimal string `json:"markupDecimal,omitempty"`
}

type MoovFeeDetails struct {
	CardScheme     string `json:"cardScheme,omitempty"`
	Interchange    string `json:"interchange,omitempty"`
	MoovProcessing string `json:"moovProcessing,omitempty"`
}

type Refund struct {
	RefundID    string      `json:"refundID,omitempty"`
	CreatedOn   time.Time   `json:"createdOn,omitempty"`
	UpdatedOn   time.Time   `json:"updatedOn,omitempty"`
	Status      string      `json:"status,omitempty"`
	FailureCode string      `json:"failureCode,omitempty"`
	Amount      Amount      `json:"amount,omitempty"`
	CardDetails CardDetails `json:"cardDetails,omitempty"`
}

type Dispute struct {
	DisputeID string    `json:"disputeID,omitempty"`
	CreatedOn time.Time `json:"createdOn,omitempty"`
	Amount    Amount    `json:"amount,omitempty"`
}

type Source struct {
	PaymentMethodID   string          `json:"paymentMethodID,omitempty"`
	PaymentMethodType string          `json:"paymentMethodType,omitempty"`
	Account           TransferAccount `json:"account,omitempty"`
	BankAccount       BankAccount     `json:"bankAccount,omitempty"`
	Wallet            Wallet          `json:"wallet,omitempty"`
	Card              Card            `json:"card,omitempty"`
	ApplePay          ApplePay        `json:"applePay,omitempty"`
	AchDetails        AchDetails      `json:"achDetails,omitempty"`
	CardDetails       CardDetails     `json:"cardDetails,omitempty"`
	TransferID        string          `json:"transferID,omitempty"`
}

type TransferAccount struct {
	AccountID   string `json:"accountID,omitempty"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type Destination struct {
	PaymentMethodID   string          `json:"paymentMethodID,omitempty"`
	PaymentMethodType string          `json:"paymentMethodType,omitempty"`
	Account           TransferAccount `json:"account,omitempty"`
	BankAccount       BankAccount     `json:"bankAccount,omitempty"`
	Wallet            Wallet          `json:"wallet,omitempty"`
	Card              Card            `json:"card,omitempty"`
	ApplePay          ApplePay        `json:"applePay,omitempty"`
	AchDetails        AchDetails      `json:"achDetails,omitempty"`
	CardDetails       CardDetails     `json:"cardDetails,omitempty"`
}

// CreateTransfer creates a new transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/createTransfer

// ListTransfers lists all transfers
// https://docs.moov.io/api/index.html#tag/Transfers/operation/listTransfers

// GetTransfer retrieves a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/getTransfer

// UpdateTransferMetaData updates the metadata for a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/patchTransfer

// TransferOptions lists all transfer options between a source and destination
// https://docs.moov.io/api/#tag/Transfers/operation/createTransferOptions

// RefundTransfer refunds a transfer
// https://docs.moov.io/api/#tag/Transfers/operation/refundTransfer

// ListRefunds lists all refunds for a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/getRefunds

// GetRefund retrieves a refund for a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/getRefund

// ReverseTransfer reverses a transfer
// https://docs.moov.io/api/index.html#tag/Transfers/operation/reverseTransfer
