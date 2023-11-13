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

type Wallet struct {
	WalletID string `json:"walletID,omitempty"`
}

type Card struct {
	CardID             string             `json:"cardID,omitempty"`
	Fingerprint        string             `json:"fingerprint,omitempty"`
	Brand              string             `json:"brand,omitempty"`
	CardType           string             `json:"cardType,omitempty"`
	LastFourCardNumber string             `json:"lastFourCardNumber,omitempty"`
	Bin                string             `json:"bin,omitempty"`
	Expiration         Expiration         `json:"expiration,omitempty"`
	HolderName         string             `json:"holderName,omitempty"`
	BillingAddress     Address            `json:"billingAddress,omitempty"`
	CardVerification   CardVerification   `json:"cardVerification,omitempty"`
	Issuer             string             `json:"issuer,omitempty"`
	IssuerCountry      string             `json:"issuerCountry,omitempty"`
	CardOnFile         bool               `json:"cardOnFile,omitempty"`
	MerchantAccountID  string             `json:"merchantAccountID,omitempty"`
	CardAccountUpdater CardAccountUpdater `json:"cardAccountUpdater,omitempty"`
	DomesticPushToCard string             `json:"domesticPushToCard,omitempty"`
}

type Expiration struct {
	Month string `json:"month,omitempty"`
	Year  string `json:"year,omitempty"`
}

type CardVerification struct {
	Cvv          string `json:"cvv,omitempty"`
	AddressLine1 string `json:"addressLine1,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
}

type CardAccountUpdater struct {
	UpdatedOn  time.Time `json:"updatedOn,omitempty"`
	UpdateType string    `json:"updateType,omitempty"`
}

type ApplePay struct {
	Brand           string     `json:"brand,omitempty"`
	CardType        string     `json:"cardType,omitempty"`
	CardDisplayName string     `json:"cardDisplayName,omitempty"`
	Fingerprint     string     `json:"fingerprint,omitempty"`
	Expiration      Expiration `json:"expiration,omitempty"`
	DynamicLastFour string     `json:"dynamicLastFour,omitempty"`
}

type AchDetails struct {
	Status      string `json:"status,omitempty"`
	TraceNumber string `json:"traceNumber,omitempty"`
	Return      struct {
		Code        string `json:"code,omitempty"`
		Reason      string `json:"reason,omitempty"`
		Description string `json:"description,omitempty"`
	}
	Correction struct {
		Code        string `json:"code,omitempty"`
		Reason      string `json:"reason,omitempty"`
		Description string `json:"description,omitempty"`
	}
	CompanyEntryDescription string           `json:"companyEntryDescription,omitempty"`
	OriginatingCompanyName  string           `json:"originatingCompanyName,omitempty"`
	StatusUpdates           ACHStatusUpdates `json:"statusUpdates,omitempty"`
	DebitHoldPeriod         string           `json:"debitHoldPeriod,omitempty"`
}

type CardDetails struct {
	Status                   string            `json:"status,omitempty"`
	FailureCode              string            `json:"failureCode,omitempty"`
	DynamicDescriptor        string            `json:"dynamicDescriptor,omitempty"`
	TransactionSource        string            `json:"transactionSource,omitempty"`
	InterchangeQualification string            `json:"interchangeQualification,omitempty"`
	StatusUpdates            CardStatusUpdates `json:"statusUpdates,omitempty"`
}

type CardStatusUpdates struct {
	Initiated time.Time `json:"initiated,omitempty"`
	Confirmed time.Time `json:"confirmed,omitempty"`
	Settled   time.Time `json:"settled,omitempty"`
	Failed    time.Time `json:"failed,omitempty"`
	Canceled  time.Time `json:"canceled,omitempty"`
	Completed time.Time `json:"completed,omitempty"`
}

type ACHStatusUpdates struct {
	Initiated  time.Time `json:"initiated,omitempty"`
	Originated time.Time `json:"originated,omitempty"`
	Corrected  time.Time `json:"corrected,omitempty"`
	Returned   time.Time `json:"returned,omitempty"`
	Completed  time.Time `json:"completed,omitempty"`
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
