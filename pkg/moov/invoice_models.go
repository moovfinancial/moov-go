package moov

import "time"

type Invoice struct {
	InvoiceID         string           `json:"invoiceID"`
	InvoiceNumber     string           `json:"invoiceNumber"`
	Description       string           `json:"description"`
	CustomerAccountID string           `json:"customerAccountID"`
	PartnerAccountID  string           `json:"partnerAccountID,omitempty"`
	Status            InvoiceStatus    `json:"status"`
	LineItems         InvoiceLineItems `json:"lineItems"`
	SubtotalAmount    AmountDecimal    `json:"subtotalAmount"`
	TaxAmount         AmountDecimal    `json:"taxAmount"`
	TotalAmount       AmountDecimal    `json:"totalAmount"`
	PendingAmount     AmountDecimal    `json:"pendingAmount"`
	PaidAmount        AmountDecimal    `json:"paidAmount"`
	RefundedAmount    AmountDecimal    `json:"refundedAmount"`
	DisputedAmount    AmountDecimal    `json:"disputedAmount"`
	PaymentLinkCode   string           `json:"paymentLinkCode"`
	Payments          []InvoicePayment `json:"payments"`
	CreatedOn         time.Time        `json:"createdOn"`
	InvoiceDate       *time.Time       `json:"invoiceDate"`
	DueDate           *time.Time       `json:"dueDate"`
	SentOn            *time.Time       `json:"sentOn"`
	PaidOn            *time.Time       `json:"paidOn"`
	CanceledOn        *time.Time       `json:"canceledOn"`
}

type InvoiceStatus string

const (
	InvoiceStatusDraft          InvoiceStatus = "draft"
	InvoiceStatusUnpaid         InvoiceStatus = "unpaid"
	InvoiceStatusPaymentPending InvoiceStatus = "payment-pending"
	InvoiceStatusPaid           InvoiceStatus = "paid"
	InvoiceStatusOverdue        InvoiceStatus = "overdue"
	InvoiceStatusCanceled       InvoiceStatus = "canceled"
)

type InvoicePayment struct {
	Amount             AmountDecimal           `json:"amount"`
	InvoicePaymentID   string                  `json:"invoicePaymentID"`
	InvoicePaymentType InvoicePaymentType      `json:"invoicePaymentType"`
	Transfer           *InvoiceTransferPayment `json:"transfer"`
	External           *InvoiceExternalPayment `json:"external"`
}

// InvoicePaymentType represents the type of payment
type InvoicePaymentType string

const (
	InvoicePaymentTypeTransfer InvoicePaymentType = "transfer"
	InvoicePaymentTypeExternal InvoicePaymentType = "external"
)

// InvoiceTransferPayment represents a payment made via Moov transfer
type InvoiceTransferPayment struct {
	TransferID string `json:"transferID"`
}

// InvoiceExternalPayment represents a payment made outside Moov
type InvoiceExternalPayment struct {
	Description string     `json:"description"`
	ForeignID   string     `json:"foreignID,omitempty"`
	PaymentDate *time.Time `json:"paymentDate,omitempty"`
}

// CreateInvoice represents the request to create an invoice
type CreateInvoice struct {
	CustomerAccountID string                 `json:"customerAccountID"`
	Description       string                 `json:"description"`
	LineItems         CreateInvoiceLineItems `json:"lineItems"`
	InvoiceDate       *time.Time             `json:"invoiceDate,omitempty"`
	DueDate           *time.Time             `json:"dueDate,omitempty"`
	TaxAmount         *AmountDecimal         `json:"taxAmount,omitempty"`
}

// UpdateInvoice represents the request to update an invoice
type UpdateInvoice struct {
	Description *string                 `json:"description,omitempty"`
	LineItems   *CreateInvoiceLineItems `json:"lineItems,omitempty"`
	InvoiceDate *time.Time              `json:"invoiceDate,omitempty"`
	DueDate     *time.Time              `json:"dueDate,omitempty"`
	Status      *InvoiceStatus          `json:"status,omitempty"`
	TaxAmount   *AmountDecimal          `json:"taxAmount,omitempty"`
}

// CreateInvoicePayment represents the request to create an external payment for an invoice
type CreateInvoicePayment struct {
	Amount      AmountDecimal `json:"amount"`
	ForeignID   *string       `json:"foreignID,omitempty"`
	Description *string       `json:"description,omitempty"`
	PaymentDate *time.Time    `json:"paymentDate,omitempty"`
}

// InvoiceLineItems represents a collection of line items for an invoice
type InvoiceLineItems struct {
	Items []InvoiceLineItem `json:"items"`
}

// InvoiceLineItem represents a single item in an invoice
type InvoiceLineItem struct {
	Name      string                         `json:"name"`
	BasePrice AmountDecimal                  `json:"basePrice"`
	Quantity  int32                          `json:"quantity"`
	ProductID *string                        `json:"productID,omitempty"`
	Options   []InvoiceLineItemOption        `json:"options,omitempty"`
	Images    []InvoiceLineItemImageMetadata `json:"images,omitempty"`
}

// InvoiceLineItemOption represents a modifier or option applied to a line item
type InvoiceLineItemOption struct {
	Name          string                         `json:"name"`
	Quantity      int32                          `json:"quantity"`
	PriceModifier *AmountDecimal                 `json:"priceModifier,omitempty"`
	Group         *string                        `json:"group,omitempty"`
	Images        []InvoiceLineItemImageMetadata `json:"images,omitempty"`
}

// InvoiceLineItemImageMetadata represents metadata for an image associated with a line item
type InvoiceLineItemImageMetadata struct {
	ImageID  string  `json:"imageID"`
	AltText  *string `json:"altText,omitempty"`
	Link     string  `json:"link"`
	PublicID string  `json:"publicID"`
}

// CreateInvoiceLineItems represents a collection of line items for creating/updating an invoice
type CreateInvoiceLineItems struct {
	Items []CreateInvoiceLineItem `json:"items"`
}

// CreateInvoiceLineItem represents a single item in an invoice creation/update request
type CreateInvoiceLineItem struct {
	Name      string                        `json:"name"`
	BasePrice AmountDecimal                 `json:"basePrice"`
	Quantity  int32                         `json:"quantity"`
	ProductID *string                       `json:"productID,omitempty"`
	Options   []CreateInvoiceLineItemOption `json:"options,omitempty"`
	Images    []CreateInvoiceLineItemImage  `json:"images,omitempty"`
}

// CreateInvoiceLineItemOption represents a modifier or option in an invoice creation/update request
type CreateInvoiceLineItemOption struct {
	Name          string                       `json:"name"`
	Quantity      int32                        `json:"quantity"`
	PriceModifier *AmountDecimal               `json:"priceModifier,omitempty"`
	Group         *string                      `json:"group,omitempty"`
	Images        []CreateInvoiceLineItemImage `json:"images,omitempty"`
}

// CreateInvoiceLineItemImage represents an image to associate with a line item
type CreateInvoiceLineItemImage struct {
	ImageID string `json:"imageID,omitempty"`
}
