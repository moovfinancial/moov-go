package moov

import "time"

// CreateTransfer struct for CreateTransfer
type CreateTransfer struct {
	Source         CreateTransfer_Source         `json:"source"`
	Destination    CreateTransfer_Destination    `json:"destination"`
	Amount         Amount                        `json:"amount"`
	SalesTaxAmount *Amount                       `json:"salesTaxAmount,omitempty"`
	FacilitatorFee CreateTransfer_FacilitatorFee `json:"facilitatorFee,omitempty"`
	// An optional description of the transfer for your own internal use.
	Description string `json:"description,omitempty"`
	// Free-form key-value pair list. Useful for storing information that is not captured elsewhere.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// CreateTransfer_Source Where funds for a transfer originate. For the source,
// you must include either a `paymentMethodID`,`transferID` or `paymentToken`.
// A `transferID` is used to create a [transfer
// group](https://docs.moov.io/guides/money-movement/transfer-groups/),
// associating the new transfer with a parent transfer.
type CreateTransfer_Source struct {
	TransferID      string                            `json:"transferID,omitempty"`
	PaymentMethodID string                            `json:"paymentMethodID,omitempty"`
	PaymentToken    string                            `json:"paymentToken,omitempty"`
	CardDetails     *CreateTransfer_CardDetailsSource `json:"cardDetails,omitempty"`
	AchDetails      *CreateTransfer_AchDetailsSource  `json:"achDetails,omitempty"`
}

// CreateTransfer_CardDetailsSource struct for CreateTransfer_CardDetailsSource
type CreateTransfer_CardDetailsSource struct {
	// An optional override of the default card statement descriptor for a transfer.
	DynamicDescriptor string             `json:"dynamicDescriptor,omitempty"`
	TransactionSource *TransactionSource `json:"transactionSource,omitempty"`
}

// CreateTransfer_AchDetailsSource struct for CreateTransfer_AchDetailsSource
type CreateTransfer_AchDetailsSource struct {
	// An optional override of the default NACHA company entry description for a transfer.
	CompanyEntryDescription string `json:"companyEntryDescription,omitempty"`
	// An optional override of the default NACHA company name for a transfer.
	OriginatingCompanyName string           `json:"originatingCompanyName,omitempty"`
	DebitHoldPeriod        *DebitHoldPeriod `json:"debitHoldPeriod,omitempty"`
	SecCode                *SecCode         `json:"secCode,omitempty"`
}

// CreateTransfer_Destination The final stage of a transfer and the ultimate recipient of the funds.
type CreateTransfer_Destination struct {
	PaymentMethodID string                                 `json:"paymentMethodID"`
	CardDetails     *CreateTransfer_CardDetailsDestination `json:"cardDetails,omitempty"`
	AchDetails      *CreateTransfer_AchDetailsBase         `json:"achDetails,omitempty"`
}

// CreateTransfer_CardDetailsDestination struct for CreateTransfer_CardDetailsDestination
type CreateTransfer_CardDetailsDestination struct {
	// An optional override of the default card statement descriptor for a transfer.
	DynamicDescriptor string `json:"dynamicDescriptor,omitempty"`
}

// CreateTransfer_AchDetailsBase If transfer involves ACH, override default card acceptance properties.
type CreateTransfer_AchDetailsBase struct {
	// An optional override of the default NACHA company entry description for a transfer.
	CompanyEntryDescription string `json:"companyEntryDescription,omitempty"`
	// An optional override of the default NACHA company name for a transfer.
	OriginatingCompanyName string `json:"originatingCompanyName,omitempty"`
}

// CreateTransfer_FacilitatorFee Total or markup fee.
type CreateTransfer_FacilitatorFee struct {
	// Total facilitator fee in cents. Only either `total` or `totalDecimal` can be set.
	Total *int64 `json:"total,omitempty"`
	// Same as `total`, but a decimal-formatted numerical string that represents up to 9 decimal place precision. Only either `total` or `totalDecimal` can be set. Set this field if you expect the fee to be in fractions of a cent.
	TotalDecimal *string `json:"totalDecimal,omitempty"`
	// Markup facilitator fee in cents. Only either `markup` or `markupDecimal` can be set.
	Markup *int64 `json:"markup,omitempty"`
	// Same as `markup`, but a decimal-formatted numerical string that represents up to 9 decimal place precision. Only either `markup` or `markupDecimal` can be set. Set this field if you expect the fee to be in fractions of a cent.
	MarkupDecimal *string `json:"markupDecimal,omitempty"`
}

// TransferStarted is where the request to create a transfer was recorded and kicked off but hasn't completed yet
type TransferStarted struct {
	// Identifier for the transfer.
	TransferID string    `json:"transferID,omitempty"`
	CreatedOn  time.Time `json:"createdOn,omitempty"`
}

// Transfer struct for Transfer
type Transfer struct {
	TransferID    string         `json:"transferID,omitempty"`
	CreatedOn     time.Time      `json:"createdOn,omitempty"`
	CompletedOn   *time.Time     `json:"completedOn,omitempty"`
	Status        TransferStatus `json:"status,omitempty"`
	FailureReason *FailureReason `json:"failureReason,omitempty"`
	Amount        Amount         `json:"amount,omitempty"`

	Source      TransferSource      `json:"source,omitempty"`
	Destination TransferDestination `json:"destination,omitempty"`

	// A description of the transfer.
	Description string `json:"description,omitempty"`
	// Free-form key-value pair list. Useful for storing information that is not captured elsewhere.
	Metadata       map[string]string  `json:"metadata,omitempty"`
	FacilitatorFee *GetFacilitatorFee `json:"facilitatorFee,omitempty"`
	// Fees charged to your platform account for transfers.
	MoovFee *int64 `json:"moovFee,omitempty"`
	// Same as `moovFee`, but a decimal-formatted numerical string that represents up to 9 decimal place precision.
	MoovFeeDecimal string          `json:"moovFeeDecimal,omitempty"`
	MoovFeeDetails *MoovFeeDetails `json:"moovFeeDetails,omitempty"`
	// Fees charged to your platform account for transfers.
	MoovFees []MoovFee `json:"moovFees,omitempty"`

	// ID for all transfers associated with a [transfer group](https://docs.moov.io/guides/money-movement/transfer-groups/).
	GroupID *string `json:"groupID,omitempty"`
	// ID of the associated sweep sweep.
	SweepID *string `json:"sweepID,omitempty"`
	// ID of the associated schedule.
	ScheduleID *string `json:"scheduleID,omitempty"`
	// ID of the associated occurrence.
	OccurrenceID *string `json:"occurrenceID,omitempty"`
	// The total refunded amount for a card transfer, representing one refunded amount, or multiple partial refunded amounts. Contains an integer value and its currency. See the `refunds` array for additional details.
	RefundedAmount *Amount `json:"refundedAmount,omitempty"`
	// A list of refunds for a card transfer.
	Refunds []Refund `json:"refunds,omitempty"`
	// The total disputed amount for a card transfer.
	DisputedAmount *Amount `json:"disputedAmount,omitempty"`
	// A list of disputes for a card transfer.
	Disputes []GetDispute `json:"disputes,omitempty"`
	// A list of cancellations for a transfer.
	Cancellations []Cancellation `json:"cancellations,omitempty"`

	// Optional sales tax amount. Transfer.Amount.Value should be inclusive of any sales tax and represents the total amount charged.
	SalesTaxAmount *Amount `json:"salesTaxAmount,omitempty"`

	PaymentLinkCode *string `json:"paymentLinkCode,omitempty"`
}

// Fees charged to accounts involved in the transfer.
type MoovFee struct {
	// ID of the account that fees were charged to.
	AccountID string `json:"accountID"`
	// List of fee IDs that sum to the totalAmount.
	FeeIDs []string `json:"feeIDs"`
	// The total amount of fees charged to the account.
	TotalAmount AmountDecimal `json:"totalAmount"`
	// Indicates whether the account charged was the partner, source, or destination of the transfer.
	TransferParty TransferParty `json:"transferParty"`
}

// AmountDecimal A decimal value representing money in a specific currency.
type AmountDecimal struct {
	// A 3-letter ISO 4217 currency code.
	Currency string `json:"currency"`
	// A decimal-formatted numerical string that represents up to 9 decimal place precision. For example, $12.987654321 is '12.987654321'.
	ValueDecimal string `json:"valueDecimal"`
}

// Amount A representation of money containing an integer value and its currency.
type Amount struct {
	// A 3-letter ISO 4217 currency code.
	Currency string `json:"currency,omitempty" otel:"currency"`
	// Quantity in the smallest unit of the specified currency. In USD this is cents, so $12.04 is 1204 and $0.99 would be 99.
	Value int64 `json:"value,omitempty" otel:"value"`
}

// GetFacilitatorFee Fee you charged your customer for the transfer.
type GetFacilitatorFee struct {
	// Total facilitator fee in cents.
	Total int64 `json:"total"`
	// Same as `total`, but a decimal-formatted numerical string that represents up to 9 decimal place precision.
	TotalDecimal string `json:"totalDecimal,omitempty"`
	// Markup facilitator fee in cents.
	Markup int64 `json:"markup,omitempty"`
	// Same as `markup`, but a decimal-formatted numerical string that represents up to 9 decimal place precision.
	MarkupDecimal string `json:"markupDecimal,omitempty"`
}

// MoovFeeDetails Processing and pass-through costs that add up to the moovFee.
type MoovFeeDetails struct {
	// Card scheme fees accrued during authorization and settlement. String type represents dollars with up to 9 decimal place precision.
	CardScheme string `json:"cardScheme,omitempty"`
	// Network interchange fee for Visa, Mastercard, or Discover. String type represents dollars with up to 9 decimal place precision.
	Interchange string `json:"interchange,omitempty"`
	// Network discount fee for American Express. String type represents dollars with up to 9 decimal place precision.
	Discount string `json:"discount,omitempty"`
	// Moov processing fee. String type represents dollars with up to 9 decimal place precision.
	MoovProcessing string `json:"moovProcessing"`
}

// Refund Details of a card refund.
type Refund struct {
	RefundID  string       `json:"refundID,omitempty"`
	CreatedOn time.Time    `json:"createdOn,omitempty"`
	UpdatedOn time.Time    `json:"updatedOn,omitempty"`
	Status    RefundStatus `json:"status,omitempty"`
	// This field is deprecated and will be removed in December 2023.
	FailureCode *CardFailureCode   `json:"failureCode,omitempty"`
	Amount      Amount             `json:"amount,omitempty"`
	CardDetails *RefundCardDetails `json:"cardDetails,omitempty"`
}

// RefundCardDetails struct for RefundCardDetails
type RefundCardDetails struct {
	Status      RefundCardStatus `json:"status,omitempty"`
	FailureCode *CardFailureCode `json:"failureCode,omitempty"`
	InitiatedOn *time.Time       `json:"initiatedOn,omitempty"`
	ConfirmedOn *time.Time       `json:"confirmedOn,omitempty"`
	SettledOn   *time.Time       `json:"settledOn,omitempty"`
	FailedOn    *time.Time       `json:"failedOn,omitempty"`
	CompletedOn *time.Time       `json:"completedOn,omitempty"`
}

// GetDispute Details of a card dispute.
type GetDispute struct {
	DisputeID string    `json:"disputeID,omitempty"`
	CreatedOn time.Time `json:"createdOn,omitempty"`
	Amount    Amount    `json:"amount,omitempty"`
}

type Cancellation struct {
	CancellationID string             `json:"cancellationID"`
	Status         CancellationStatus `json:"status"`
	CreatedOn      time.Time          `json:"createdOn"`
}

// TransferSource struct for TransferSource
type TransferSource struct {
	PaymentMethodID   string                    `json:"paymentMethodID,omitempty"`
	PaymentMethodType PaymentMethodType         `json:"paymentMethodType,omitempty"`
	Account           TransferAccount           `json:"account,omitempty"`
	BankAccount       *BankAccountPaymentMethod `json:"bankAccount,omitempty"`
	Wallet            *WalletPaymentMethod      `json:"wallet,omitempty"`
	Card              *CardPaymentMethod        `json:"card,omitempty"`
	ApplePay          *ApplePayPaymentMethod    `json:"applePay,omitempty"`
	AchDetails        *AchDetailsSource         `json:"achDetails,omitempty"`
	CardDetails       *CardDetails              `json:"cardDetails,omitempty"`
	TransferID        string                    `json:"transferID,omitempty"`
}

// TransferAccount struct for TransferAccount
type TransferAccount struct {
	// ID of account.
	AccountID string `json:"accountID,omitempty"`
	// Email address.
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// AchDetailsSource struct for AchDetailsSource
type AchDetailsSource struct {
	Status      AchStatus     `json:"status"`
	TraceNumber string        `json:"traceNumber"`
	Return      *AchException `json:"return,omitempty"`
	Correction  *AchException `json:"correction,omitempty"`
	// An optional override of the default NACHA company entry description for a transfer.
	CompanyEntryDescription string `json:"companyEntryDescription,omitempty"`
	// An optional override of the default NACHA company name for a transfer.
	OriginatingCompanyName string          `json:"originatingCompanyName,omitempty"`
	SecCode                SecCode         `json:"secCode,omitempty"`
	InitiatedOn            *time.Time      `json:"initiatedOn,omitempty"`
	OriginatedOn           *time.Time      `json:"originatedOn,omitempty"`
	CorrectedOn            *time.Time      `json:"correctedOn,omitempty"`
	ReturnedOn             *time.Time      `json:"returnedOn,omitempty"`
	CompletedOn            *time.Time      `json:"completedOn,omitempty"`
	CanceledOn             *time.Time      `json:"canceledOn,omitempty"`
	DebitHoldPeriod        DebitHoldPeriod `json:"debitHoldPeriod,omitempty"`
}

// AchException struct for AchException
type AchException struct {
	Code        string `json:"code,omitempty"`
	Reason      string `json:"reason,omitempty"`
	Description string `json:"description,omitempty"`
}

// TransferDestination struct for TransferDestination
type TransferDestination struct {
	PaymentMethodID   string                    `json:"paymentMethodID,omitempty"`
	PaymentMethodType PaymentMethodType         `json:"paymentMethodType,omitempty"`
	Account           TransferAccount           `json:"account,omitempty"`
	BankAccount       *BankAccountPaymentMethod `json:"bankAccount,omitempty"`
	Wallet            *WalletPaymentMethod      `json:"wallet,omitempty"`
	Card              *CardPaymentMethod        `json:"card,omitempty"`
	ApplePay          *ApplePayPaymentMethod    `json:"applePay,omitempty"`
	AchDetails        *AchDetails               `json:"achDetails,omitempty"`
	CardDetails       *CardDetails              `json:"cardDetails,omitempty"`
	RtpDetails        *RtpDetails               `json:"rtpDetails,omitempty"`
}

// AchDetails ACH specific details about the transaction.
type AchDetails struct {
	Status      AchStatus     `json:"status"`
	TraceNumber string        `json:"traceNumber"`
	Return      *AchException `json:"return,omitempty"`
	Correction  *AchException `json:"correction,omitempty"`
	// An optional override of the default NACHA company entry description for a transfer.
	CompanyEntryDescription string `json:"companyEntryDescription,omitempty"`
	// An optional override of the default NACHA company name for a transfer.
	OriginatingCompanyName string     `json:"originatingCompanyName,omitempty"`
	SecCode                SecCode    `json:"secCode,omitempty"`
	InitiatedOn            *time.Time `json:"initiatedOn,omitempty"`
	OriginatedOn           *time.Time `json:"originatedOn,omitempty"`
	CorrectedOn            *time.Time `json:"correctedOn,omitempty"`
	ReturnedOn             *time.Time `json:"returnedOn,omitempty"`
	CompletedOn            *time.Time `json:"completedOn,omitempty"`
	CanceledOn             *time.Time `json:"canceledOn,omitempty"`
}

// RtpDetails RTP specific details about the transaction.
type RtpDetails struct {
	Status RtpStatus `json:"status"`
	// Code returned by rail network on failure.
	NetworkResponseCode      *string         `json:"networkResponseCode,omitempty"`
	FailureCode              *RtpFailureCode `json:"failureCode,omitempty"`
	InitiatedOn              *time.Time      `json:"initiatedOn,omitempty"`
	CompletedOn              *time.Time      `json:"completedOn,omitempty"`
	FailedOn                 *time.Time      `json:"failedOn,omitempty"`
	AcceptedWithoutPostingOn *time.Time      `json:"acceptedWithoutPostingOn,omitempty"`
}

type patchTransfer struct {
	Metadata map[string]string `json:"metadata,omitempty"`
}

// CreateRefund Specifies a partial amount to refund. This request body is optional, an empty body will issue a refund for the full amount of the original transfer.
type CreateRefund struct {
	// Amount to refund in cents. If null, the original transfer's full amount will be refunded.
	Amount int64 `json:"amount,omitempty"`
}

type RefundStarted struct {
	Transfer
	RefundedTransferID string `otel:"refunded_transfer_id,omitempty"` // original transfer id
}

// CreateReversal struct for CreateReversal
type CreateReversal struct {
	// Amount to reverse in cents. If null, the original transfer's full amount will be reversed. Partial amounts will automatically trigger a refund instead of a cancellation.
	Amount int64 `json:"amount,omitempty"`
}

// CreatedReversal struct for CreatedReversal
type CreatedReversal struct {
	Cancellation *CreatedCancellation `json:"cancellation,omitempty"`
	Refund       *Refund              `json:"refund,omitempty"`
}

// CreatedCancellation struct for CreatedCancellation
type CreatedCancellation struct {
	Status    CancellationStatus `json:"status,omitempty"`
	CreatedOn time.Time          `json:"createdOn,omitempty"`
}

// CreateTransferOptions struct for CreateTransferOptions
type CreateTransferOptions struct {
	Source      CreateTransferOptionsTarget `json:"source"`
	Destination CreateTransferOptionsTarget `json:"destination"`
	Amount      Amount                      `json:"amount"`
}

// CreateTransferOptionsTarget struct for CreateTransferOptionsTarget
type CreateTransferOptionsTarget struct {
	AccountID       string `json:"accountID,omitempty"`
	PaymentMethodID string `json:"paymentMethodID,omitempty"`
}

// CreatedTransferOptions struct for CreatedTransferOptions
type TransferOptions struct {
	SourceOptions      []PaymentMethod `json:"sourceOptions,omitempty"`
	DestinationOptions []PaymentMethod `json:"destinationOptions,omitempty"`
}

/* ======== enumerations ======== */

// TransactionSource Specifies the nature and initiator of a transaction. Crucial for recurring and merchant-initiated transactions as per card scheme rules. Omit for customer-initiated e-commerce transactions.  - `first-recurring`: Initial transaction in a recurring series or saving a card for future merchant-initiated charges - `recurring`: Regular, merchant-initiated scheduled transactions - `unscheduled`: Non-regular, merchant-initiated transactions like account top-ups
type TransactionSource string

// // List of TransactionSource
const (
	TransactionSource_FirstRecurring TransactionSource = "first-recurring"
	TransactionSource_Recurring      TransactionSource = "recurring"
	TransactionSource_Unscheduled    TransactionSource = "unscheduled"
)

// DebitHoldPeriod An optional override of your default ACH hold period in banking days. The hold period must be longer than or equal to your default setting.
type DebitHoldPeriod string

// List of DebitHoldPeriod
const (
	DebitHoldPeriod_NoHold DebitHoldPeriod = "no-hold"
	DebitHoldPeriod_1Day   DebitHoldPeriod = "1-day"
	DebitHoldPeriod_2Days  DebitHoldPeriod = "2-days"
)

// SecCode Code used to identify the ACH authorization method.
type SecCode string

// List of SECCode
const (
	SecCode_WEB SecCode = "WEB"
	SecCode_PPD SecCode = "PPD"
	SecCode_CCD SecCode = "CCD"
	SecCode_TEL SecCode = "TEL"
)

// TransferStatus Current status of a transfer.
type TransferStatus string

// List of TransferStatus
const (
	TransferStatus_Created   TransferStatus = "created"
	TransferStatus_Pending   TransferStatus = "pending"
	TransferStatus_Completed TransferStatus = "completed"
	TransferStatus_Failed    TransferStatus = "failed"
	TransferStatus_Reversed  TransferStatus = "reversed"
	TransferStatus_Queued    TransferStatus = "queued"
	TransferStatus_Canceled  TransferStatus = "canceled"
)

// FailureReason Reason for a transfer's failure.
type FailureReason string

// List of FailureReason
const (
	FailureReason_Source_Payment_Error      FailureReason = "source-payment-error"
	FailureReason_Destination_Payment_Error FailureReason = "destination-payment-error"
	FailureReason_Wallet_Insufficient_Funds FailureReason = "wallet-insufficient-funds"
	FailureReason_Rejected_HighRisk         FailureReason = "rejected-high-risk"
	FailureReason_Processing_Error          FailureReason = "processing-error"
)

// RefundStatus the model 'RefundStatus'
type RefundStatus string

// List of RefundStatus
const (
	RefundStatus_Created   RefundStatus = "created"
	RefundStatus_Pending   RefundStatus = "pending"
	RefundStatus_Completed RefundStatus = "completed"
	RefundStatus_Failed    RefundStatus = "failed"
)

// CardFailureCode the model 'CardFailureCode'
type CardFailureCode string

// List of CardFailureCode
const (
	CardFailureCode_CallIssuer                 CardFailureCode = "call-issuer"
	CardFailureCode_DoNotHonor                 CardFailureCode = "do-not-honor"
	CardFailureCode_ProcessingError            CardFailureCode = "processing-error"
	CardFailureCode_InvalidTransaction         CardFailureCode = "invalid-transaction"
	CardFailureCode_InvalidAmount              CardFailureCode = "invalid-amount"
	CardFailureCode_NoSuchIssuer               CardFailureCode = "no-such-issuer"
	CardFailureCode_ReenterTransaction         CardFailureCode = "reenter-transaction"
	CardFailureCode_CVV_Mismatch               CardFailureCode = "cvv-mismatch"
	CardFailureCode_LostOrStolen               CardFailureCode = "lost-or-stolen"
	CardFailureCode_Insufficient_Funds         CardFailureCode = "insufficient-funds"
	CardFailureCode_InvalidCardNumber          CardFailureCode = "invalid-card-number"
	CardFailureCode_InvalidMerchant            CardFailureCode = "invalid-merchant"
	CardFailureCode_ExpiredCard                CardFailureCode = "expired-card"
	CardFailureCode_IncorrectPin               CardFailureCode = "incorrect-pin"
	CardFailureCode_TransactionNotAllowed      CardFailureCode = "transaction-not-allowed"
	CardFailureCode_SuspectedFraud             CardFailureCode = "suspected-fraud"
	CardFailureCode_AmountLimitedExceeded      CardFailureCode = "amount-limit-exceeded"
	CardFailureCode_VelocityLimitExceeded      CardFailureCode = "velocity-limit-exceeded"
	CardFailureCode_RevocationOfAauthorization CardFailureCode = "revocation-of-authorization"
	CardFailureCode_CardNotActivated           CardFailureCode = "card-not-activated"
	CardFailureCode_IssuerNotAvailable         CardFailureCode = "issuer-not-available"
	CardFailureCode_CouldNotRoute              CardFailureCode = "could-not-route"
	CardFailureCode_CardholderAccounterClosed  CardFailureCode = "cardholder-account-closed"
	CardFailureCode_DuplicateTransaction       CardFailureCode = "duplicate-transaction"
	CardFailureCode_UnknownIssue               CardFailureCode = "unknown-issue"
)

// RefundCardStatus Status of the refund.
type RefundCardStatus string

// List of RefundCardStatus
const (
	RefundCardStatus_Initiated RefundCardStatus = "initiated"
	RefundCardStatus_Confirmed RefundCardStatus = "confirmed"
	RefundCardStatus_Settled   RefundCardStatus = "settled"
	RefundCardStatus_Failed    RefundCardStatus = "failed"
	RefundCardStatus_Completed RefundCardStatus = "completed"
)

// AchStatus Status of the ACH lifecycle.
type AchStatus string

// List of ACHStatus
const (
	AchStatus_Initied    AchStatus = "initiated"
	AchStatus_Originated AchStatus = "originated"
	AchStatus_Corrected  AchStatus = "corrected"
	AchStatus_Returned   AchStatus = "returned"
	AchStatus_Completed  AchStatus = "completed"
)

// RtpStatus Status of the RTP lifecycle.
type RtpStatus string

// List of RTPStatus
const (
	RtpStatus_Initiated              RtpStatus = "initiated"
	RtpStatus_Completed              RtpStatus = "completed"
	RtpStatus_Failed                 RtpStatus = "failed"
	RtpStatus_AcceptedWithoutPosting RtpStatus = "accepted-without-posting"
)

// RtpFailureCode Status codes for RTP failures.
type RtpFailureCode string

// List of RTPFailureCode
const (
	RtpFailureCode_ProcessingError         RtpFailureCode = "processing-error"
	RtpFailureCode_InvalidAccount          RtpFailureCode = "invalid-account"
	RtpFailureCode_AccountClosed           RtpFailureCode = "account-closed"
	RtpFailureCode_AccountBlocked          RtpFailureCode = "account-blocked"
	RtpFailureCode_InvalidField            RtpFailureCode = "invalid-field"
	RtpFailureCode_TransactionNotSupported RtpFailureCode = "transaction-not-supported"
	RtpFailureCode_LimitExceeded           RtpFailureCode = "limit-exceeded"
	RtpFailureCode_InvalidAmount           RtpFailureCode = "invalid-amount"
	RtpFailureCode_CustomerDeceased        RtpFailureCode = "customer-deceased"
	RtpFailureCode_Other                   RtpFailureCode = "other"
)

// CancellationStatus Cancellation status.
type CancellationStatus string

// List of CancellationStatus
const (
	CancellationStatus_Pending   CancellationStatus = "pending"
	CancellationStatus_Completed CancellationStatus = "completed"
	CancellationStatus_Failed    CancellationStatus = "failed"
)

// TransferParty The party that was charged for the transfer.
type TransferParty string

// List of TransferParty
const (
	TransferParty_Source      TransferParty = "source"
	TransferParty_Destination TransferParty = "destination"
	TransferParty_Partner     TransferParty = "partner"
)
