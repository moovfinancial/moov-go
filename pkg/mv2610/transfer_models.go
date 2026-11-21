package mv2610

import (
	"time"

	"github.com/moovfinancial/moov-go/pkg/moov"
)

// CreateTransfer is the v2026.10 request to create a transfer. Amounts are decimal strings, and
// rail-specific settings move from source.cardDetails and source.achDetails into options.
type CreateTransfer struct {
	Source      CreateTransferSource       `json:"source"`
	Destination CreateTransferDestination  `json:"destination"`
	Options     *CreateTransferRailOptions `json:"options,omitempty"`

	Amount         moov.AmountDecimal            `json:"amount"`
	FacilitatorFee *CreateTransferFacilitatorFee `json:"facilitatorFee,omitempty"`
	Description    string                        `json:"description,omitempty"`
	Metadata       map[string]string             `json:"metadata,omitempty"`
	ForeignID      *string                       `json:"foreignID,omitempty"`
	LineItems      *moov.CreateTransferLineItems `json:"lineItems,omitempty"`
	AmountDetails  *CreateTransferAmountDetails  `json:"amountDetails,omitempty"`
	FeePaidBy      *moov.TransferFeePaidBy       `json:"feePaidBy,omitempty"`
}

// CreateTransferSource is where the funds come from. Set one of paymentMethodID, transferID, or paymentToken.
type CreateTransferSource struct {
	TransferID      string `json:"transferID,omitempty"`
	PaymentMethodID string `json:"paymentMethodID,omitempty"`
	PaymentToken    string `json:"paymentToken,omitempty"`
}

type CreateTransferDestination struct {
	PaymentMethodID string `json:"paymentMethodID"`
}

// CreateTransferRailOptions holds rail-specific settings. Set only the key that matches the transfer's rail.
type CreateTransferRailOptions struct {
	CardPayment  *CreateCardPaymentOptions                   `json:"cardPayment,omitempty"`
	PushToCard   *moov.CreateTransfer_CardDetailsDestination `json:"pushToCard,omitempty"`
	PullFromCard *moov.CreateTransfer_CardDetailsSource      `json:"pullFromCard,omitempty"`
	AchDebit     *moov.CreateTransfer_AchDetailsSource       `json:"achDebit,omitempty"`
	AchCredit    *moov.CreateTransfer_AchDetailsBase         `json:"achCredit,omitempty"`
}

// CreateCardPaymentOptions has the same fields as CardPaymentOptions. CaptureMethod defaults to automatic.
type CreateCardPaymentOptions = CardPaymentOptions

// CreateTransferFacilitatorFee sets either Total or Markup.
type CreateTransferFacilitatorFee = FacilitatorFee

type CreateTransferAmountDetails = TransferAmountDetails

// Transfer is the v2026.10 transfer. Rail details move from source and destination into options and
// processingDetails, and refunds, disputes, and cancellations move to their own endpoints.
type Transfer struct {
	TransferID        string                    `json:"transferID"`
	TransferType      TransferType              `json:"transferType"`
	CreatedOn         time.Time                 `json:"createdOn"`
	CompletedOn       *time.Time                `json:"completedOn,omitempty"`
	Status            moov.TransferStatus       `json:"status"`
	FailureReason     *moov.FailureReason       `json:"failureReason,omitempty"`
	Amount            moov.AmountDecimal        `json:"amount"`
	Source            TransferSource            `json:"source"`
	Destination       TransferDestination       `json:"destination"`
	Options           TransferRailOptions       `json:"options"`
	ProcessingDetails TransferProcessingDetails `json:"processingDetails"`
	// Authorization is set only for an auth-capture card payment.
	Authorization *TransferAuthorization `json:"authorization,omitempty"`

	Description     string                  `json:"description,omitempty"`
	Metadata        map[string]string       `json:"metadata,omitempty"`
	FacilitatorFee  *FacilitatorFee         `json:"facilitatorFee,omitempty"`
	MoovFee         *moov.AmountDecimal     `json:"moovFee,omitempty"`
	MoovFeeDetails  *MoovFeeDetails         `json:"moovFeeDetails,omitempty"`
	MoovFees        []moov.MoovFee          `json:"moovFees,omitempty"`
	RefundedAmount  *moov.AmountDecimal     `json:"refundedAmount,omitempty"`
	DisputedAmount  *moov.AmountDecimal     `json:"disputedAmount,omitempty"`
	AmountDetails   *TransferAmountDetails  `json:"amountDetails,omitempty"`
	LineItems       *moov.TransferLineItems `json:"lineItems,omitempty"`
	GroupID         *string                 `json:"groupID,omitempty"`
	SweepID         *string                 `json:"sweepID,omitempty"`
	ScheduleID      *string                 `json:"scheduleID,omitempty"`
	OccurrenceID    *string                 `json:"occurrenceID,omitempty"`
	PaymentLinkCode *string                 `json:"paymentLinkCode,omitempty"`
	ForeignID       *string                 `json:"foreignID,omitempty"`
	InvoiceID       *string                 `json:"invoiceID,omitempty"`
}

// TransferType is the rail and direction used to move funds for a transfer.
type TransferType string

// List of TransferType
const (
	TransferType_CardPayment         TransferType = "card-payment"
	TransferType_PushToCard          TransferType = "push-to-card"
	TransferType_PullFromCard        TransferType = "pull-from-card"
	TransferType_AchDebit            TransferType = "ach-debit"
	TransferType_AchCredit           TransferType = "ach-credit"
	TransferType_AchDebitToAchCredit TransferType = "ach-debit-to-ach-credit"
	TransferType_InstantBankCredit   TransferType = "instant-bank-credit"
	TransferType_Wallet              TransferType = "wallet"
)

type TransferSource struct {
	TransferID        string                          `json:"transferID,omitempty"`
	PaymentMethodID   string                          `json:"paymentMethodID,omitempty"`
	PaymentMethodType moov.PaymentMethodType          `json:"paymentMethodType,omitempty"`
	Account           *moov.TransferAccount           `json:"account,omitempty"`
	BankAccount       *moov.BankAccountPaymentMethod  `json:"bankAccount,omitempty"`
	Wallet            *moov.WalletPaymentMethod       `json:"wallet,omitempty"`
	Card              *moov.CardPaymentMethod         `json:"card,omitempty"`
	ApplePay          *moov.ApplePayPaymentMethod     `json:"applePay,omitempty"`
	GooglePay         *moov.GooglePayPaymentMethod    `json:"googlePay,omitempty"`
	TerminalCard      *moov.TerminalCardPaymentMethod `json:"terminalCard,omitempty"`
}

type TransferDestination struct {
	PaymentMethodID   string                         `json:"paymentMethodID,omitempty"`
	PaymentMethodType moov.PaymentMethodType         `json:"paymentMethodType,omitempty"`
	Account           *moov.TransferAccount          `json:"account,omitempty"`
	BankAccount       *moov.BankAccountPaymentMethod `json:"bankAccount,omitempty"`
	Wallet            *moov.WalletPaymentMethod      `json:"wallet,omitempty"`
	Card              *moov.CardPaymentMethod        `json:"card,omitempty"`
	ApplePay          *moov.ApplePayPaymentMethod    `json:"applePay,omitempty"`
	GooglePay         *moov.GooglePayPaymentMethod   `json:"googlePay,omitempty"`
}

// TransferRailOptions are the rail-specific settings the transfer was created with. Only the key for the transfer's rail is set.
type TransferRailOptions struct {
	CardPayment  *CardPaymentOptions  `json:"cardPayment,omitempty"`
	PushToCard   *PushToCardOptions   `json:"pushToCard,omitempty"`
	PullFromCard *PullFromCardOptions `json:"pullFromCard,omitempty"`
	AchDebit     *AchDebitOptions     `json:"achDebit,omitempty"`
	AchCredit    *AchCreditOptions    `json:"achCredit,omitempty"`
}

type CardPaymentOptions struct {
	DynamicDescriptor *string                 `json:"dynamicDescriptor,omitempty"`
	TransactionSource *moov.TransactionSource `json:"transactionSource,omitempty"`
	CaptureMethod     *CaptureMethod          `json:"captureMethod,omitempty"`
	// IsExtendedAuthorizationRequested asks the card network for an extended authorization hold.
	IsExtendedAuthorizationRequested *bool `json:"isExtendedAuthorizationRequested,omitempty"`
	// IsEstimatedAuthorization flags the authorization amount as an estimate.
	IsEstimatedAuthorization *bool `json:"isEstimatedAuthorization,omitempty"`
}

// CaptureMethod selects whether a card payment captures automatically or waits for a capture request.
type CaptureMethod string

// List of CaptureMethod
const (
	CaptureMethod_Automatic CaptureMethod = "automatic"
	CaptureMethod_Manual    CaptureMethod = "manual"
)

type PushToCardOptions struct {
	DynamicDescriptor   *string    `json:"dynamicDescriptor,omitempty"`
	ScheduledDeliveryOn *time.Time `json:"scheduledDeliveryOn,omitempty"`
}

type PullFromCardOptions struct {
	DynamicDescriptor *string                 `json:"dynamicDescriptor,omitempty"`
	TransactionSource *moov.TransactionSource `json:"transactionSource,omitempty"`
}

type AchDebitOptions struct {
	CompanyEntryDescription *string                 `json:"companyEntryDescription,omitempty"`
	OriginatingCompanyName  *string                 `json:"originatingCompanyName,omitempty"`
	SecCode                 *moov.SecCode           `json:"secCode,omitempty"`
	DebitHoldPeriod         *moov.DebitHoldPeriod   `json:"debitHoldPeriod,omitempty"`
	Addenda                 []moov.AchAddendaRecord `json:"addenda,omitempty"`
}

type AchCreditOptions struct {
	CompanyEntryDescription *string                 `json:"companyEntryDescription,omitempty"`
	OriginatingCompanyName  *string                 `json:"originatingCompanyName,omitempty"`
	Addenda                 []moov.AchAddendaRecord `json:"addenda,omitempty"`
}

// TransferProcessingDetails reports how the rail processed the transfer. Only the key for the transfer's rail is set.
type TransferProcessingDetails struct {
	CardPayment       *CardPaymentTransferProcessingDetails       `json:"cardPayment,omitempty"`
	PushToCard        *PushToCardTransferProcessingDetails        `json:"pushToCard,omitempty"`
	PullFromCard      *PullFromCardTransferProcessingDetails      `json:"pullFromCard,omitempty"`
	AchDebit          *AchDebitTransferProcessingDetails          `json:"achDebit,omitempty"`
	AchCredit         *AchCreditTransferProcessingDetails         `json:"achCredit,omitempty"`
	InstantBankCredit *InstantBankCreditTransferProcessingDetails `json:"instantBankCredit,omitempty"`
}

type CardPaymentTransferProcessingDetails struct {
	Status                   *moov.CardTransactionStatus `json:"status,omitempty"`
	AuthorizationCode        *string                     `json:"authorizationCode,omitempty"`
	NetworkTransactionID     *string                     `json:"networkTransactionID,omitempty"`
	FailureCode              *moov.CardFailureCode       `json:"failureCode,omitempty"`
	RetrievalReferenceNumber *string                     `json:"retrievalReferenceNumber,omitempty"`
}

type PushToCardTransferProcessingDetails struct {
	Status               moov.CardTransactionStatus `json:"status"`
	AuthorizationCode    *string                    `json:"authorizationCode,omitempty"`
	NetworkTransactionID *string                    `json:"networkTransactionID,omitempty"`
	NetworkResponseCode  *string                    `json:"networkResponseCode,omitempty"`
	FailureCode          *moov.CardFailureCode      `json:"failureCode,omitempty"`
}

type PullFromCardTransferProcessingDetails struct {
	Status               moov.CardTransactionStatus `json:"status"`
	AuthorizationCode    *string                    `json:"authorizationCode,omitempty"`
	NetworkTransactionID *string                    `json:"networkTransactionID,omitempty"`
	NetworkResponseCode  *string                    `json:"networkResponseCode,omitempty"`
	FailureCode          *moov.CardFailureCode      `json:"failureCode,omitempty"`
}

type AchDebitTransferProcessingDetails struct {
	Status      moov.AchStatus     `json:"status"`
	TraceNumber string             `json:"traceNumber"`
	Return      *moov.AchException `json:"return,omitempty"`
	Correction  *moov.AchException `json:"correction,omitempty"`
}

type AchCreditTransferProcessingDetails struct {
	Status      moov.AchStatus     `json:"status"`
	TraceNumber string             `json:"traceNumber"`
	Return      *moov.AchException `json:"return,omitempty"`
	Correction  *moov.AchException `json:"correction,omitempty"`
}

type InstantBankCreditTransferProcessingDetails struct {
	Status              moov.InstantBankStatus       `json:"status"`
	Network             moov.InstantBankNetwork      `json:"network"`
	NetworkResponseCode *string                      `json:"networkResponseCode,omitempty"`
	FailureCode         *moov.InstantBankFailureCode `json:"failureCode,omitempty"`
	EndToEndID          *string                      `json:"endToEndID,omitempty"`
}

// TransferAuthorization holds the authorization and capture amounts for an auth-capture card payment.
type TransferAuthorization struct {
	AuthorizationID string             `json:"authorizationID"`
	RequestedAmount moov.AmountDecimal `json:"requestedAmount"`
	// AuthorizedAmount is the hold the issuer approved.
	AuthorizedAmount moov.AmountDecimal `json:"authorizedAmount"`
	// CapturedAmount is the sum of the captures that did not fail and were not canceled.
	CapturedAmount moov.AmountDecimal `json:"capturedAmount"`
	// CapturableAmount is what remains of the authorization after captures and authorization cancellations.
	CapturableAmount moov.AmountDecimal `json:"capturableAmount"`
	// AmountDetails is the tip, tax, and surcharge the card network authorized. It does not change after authorization.
	AmountDetails *TransferAmountDetails `json:"amountDetails,omitempty"`
	ExpiresOn     *time.Time             `json:"expiresOn,omitempty"`
	// IsExtendedAuthorizationGranted is set only when options.cardPayment.isExtendedAuthorizationRequested is true.
	IsExtendedAuthorizationGranted *bool `json:"isExtendedAuthorizationGranted,omitempty"`
}

type TransferAmountDetails struct {
	Tip       *moov.AmountDecimal `json:"tip,omitempty"`
	Tax       *moov.AmountDecimal `json:"tax,omitempty"`
	Surcharge *moov.AmountDecimal `json:"surcharge,omitempty"`
}

type FacilitatorFee struct {
	Total  *moov.AmountDecimal `json:"total,omitempty"`
	Markup *moov.AmountDecimal `json:"markup,omitempty"`
}

// MoovFeeDetails lists the processing and pass-through costs that add up to the transfer's moovFee.
type MoovFeeDetails struct {
	CardScheme     *moov.AmountDecimal `json:"cardScheme,omitempty"`
	Interchange    *moov.AmountDecimal `json:"interchange,omitempty"`
	Discount       *moov.AmountDecimal `json:"discount,omitempty"`
	MoovProcessing moov.AmountDecimal  `json:"moovProcessing"`
}

// AuthorizationStatus is the status of a card authorization on a transfer.
type AuthorizationStatus string

// List of AuthorizationStatus
const (
	AuthorizationStatus_Approved AuthorizationStatus = "approved"
	AuthorizationStatus_Declined AuthorizationStatus = "declined"
	AuthorizationStatus_Reversed AuthorizationStatus = "reversed"
)

type CaptureStatus string

// List of CaptureStatus
const (
	CaptureStatus_Created   CaptureStatus = "created"
	CaptureStatus_Pending   CaptureStatus = "pending"
	CaptureStatus_Completed CaptureStatus = "completed"
	CaptureStatus_Failed    CaptureStatus = "failed"
	CaptureStatus_Canceled  CaptureStatus = "canceled"
)

// CreateRefund refunds a card payment. Without Amount, the request refunds the original transfer's full amount.
type CreateRefund struct {
	Amount *moov.AmountDecimal `json:"amount,omitempty"`
	// CaptureID refunds one capture of an auth-capture card payment.
	CaptureID     *string              `json:"captureID,omitempty"`
	AmountDetails *RefundAmountDetails `json:"amountDetails,omitempty"`
}

type Refund struct {
	RefundID          string                  `json:"refundID"`
	CreatedOn         time.Time               `json:"createdOn"`
	UpdatedOn         time.Time               `json:"updatedOn"`
	Status            moov.RefundStatus       `json:"status"`
	Amount            moov.AmountDecimal      `json:"amount"`
	CaptureID         *string                 `json:"captureID,omitempty"`
	AmountDetails     *RefundAmountDetails    `json:"amountDetails,omitempty"`
	ProcessingDetails RefundProcessingDetails `json:"processingDetails"`
}

type RefundAmountDetails struct {
	Surcharge *moov.AmountDecimal `json:"surcharge,omitempty"`
}

type RefundProcessingDetails struct {
	CardPayment *CardPaymentRefundProcessingDetails `json:"cardPayment,omitempty"`
}

type CardPaymentRefundProcessingDetails struct {
	NetworkTransactionID *string               `json:"networkTransactionID,omitempty"`
	FailureCode          *moov.CardFailureCode `json:"failureCode,omitempty"`
}

type CreateReversal struct {
	// Amount is optional. Without Amount, the request reverses the full amount of the transfer.
	Amount        *moov.AmountDecimal  `json:"amount,omitempty"`
	AmountDetails *RefundAmountDetails `json:"amountDetails,omitempty"`
}

// CreatedReversal has either a Cancellation or a Refund, depending on how the transfer was reversed.
type CreatedReversal struct {
	Cancellation *moov.Cancellation `json:"cancellation,omitempty"`
	Refund       *Refund            `json:"refund,omitempty"`
}

type CreateCancellation struct {
	// Amount voids part of an auth-capture card payment and cannot exceed capturableAmount.
	// Without Amount, the request voids the full capturable amount.
	Amount *moov.AmountDecimal `json:"amount,omitempty"`
}
