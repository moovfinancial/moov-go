package moov

import "time"

// PaymentMethod A method of moving money
type PaymentMethod struct {
	// UUID v4
	PaymentMethodID   string                    `json:"paymentMethodID,omitempty"`
	PaymentMethodType PaymentMethodType         `json:"paymentMethodType,omitempty"`
	Wallet            *WalletPaymentMethod      `json:"wallet,omitempty"`
	BankAccount       *BankAccountPaymentMethod `json:"bankAccount,omitempty"`
	Card              *CardPaymentMethod        `json:"card,omitempty"`
	ApplePay          *ApplePayPaymentMethod    `json:"applePay,omitempty"`
}

// PaymentMethodType The payment method type that represents a payment rail and directionality
type PaymentMethodType string

// List of PaymentMethodType
const (
	PaymentMethodType_MoovWallet        PaymentMethodType = "moov-wallet"
	PaymentMethodType_AchDebitFund      PaymentMethodType = "ach-debit-fund"
	PaymentMethodType_AchDebitCollect   PaymentMethodType = "ach-debit-collect"
	PaymentMethodType_AchCreditStandard PaymentMethodType = "ach-credit-standard"
	PaymentMethodType_AchCreditSameDay  PaymentMethodType = "ach-credit-same-day"
	PaymentMethodType_RtpCredit         PaymentMethodType = "rtp-credit"
	PaymentMethodType_CardPayment       PaymentMethodType = "card-payment"
	PaymentMethodType_ApplePay          PaymentMethodType = "apple-pay"
	PaymentMethodType_PushToCard        PaymentMethodType = "push-to-card"
	PaymentMethodType_PullFromCard      PaymentMethodType = "pull-from-card"
)

// WalletPaymentMethod A Moov wallet to store funds for transfers.
type WalletPaymentMethod struct {
	// UUID v4
	WalletID string `json:"walletID,omitempty"`
}

// BankAccountPaymentMethod Describes a bank account on a Moov account.
type BankAccountPaymentMethod struct {
	// UUID v4
	BankAccountID string `json:"bankAccountID,omitempty"`
	// Once the bank account is linked, we don't reveal the full bank account number. The fingerprint acts as a way to identify whether two linked bank accounts are the same.
	Fingerprint           string                  `json:"fingerprint,omitempty"`
	Status                BankAccountStatus       `json:"status,omitempty"`
	HolderName            string                  `json:"holderName,omitempty"`
	HolderType            HolderType              `json:"holderType,omitempty"`
	BankName              string                  `json:"bankName,omitempty"`
	BankAccountType       BankAccountType         `json:"bankAccountType,omitempty"`
	RoutingNumber         string                  `json:"routingNumber,omitempty"`
	LastFourAccountNumber string                  `json:"lastFourAccountNumber,omitempty"`
	UpdatedOn             time.Time               `json:"updatedOn,omitempty"`
	StatusReason          BankAccountStatusReason `json:"statusReason,omitempty"`
	ExceptionDetails      *ExceptionDetails       `json:"exceptionDetails,omitempty"`
}

// Card Describes a card on a Moov account.
type CardPaymentMethod struct {
	// UUID v4
	CardID string `json:"cardID,omitempty"`
	// Uniquely identifies a linked payment card or token. For Apple Pay, the fingerprint is based on the tokenized card number and may vary based on the user's device. This field can be used to identify specific payment methods across multiple accounts on your platform.
	Fingerprint        string            `json:"fingerprint,omitempty"`
	Brand              CardBrand         `json:"brand,omitempty"`
	CardType           CardType          `json:"cardType,omitempty"`
	LastFourCardNumber string            `json:"lastFourCardNumber,omitempty"`
	Bin                string            `json:"bin,omitempty"`
	Expiration         CardExpiration    `json:"expiration,omitempty"`
	HolderName         string            `json:"holderName,omitempty"`
	BillingAddress     Address           `json:"billingAddress,omitempty"`
	CardVerification   CardVerifications `json:"cardVerification,omitempty"`
	Issuer             string            `json:"issuer,omitempty"`
	IssuerCountry      string            `json:"issuerCountry,omitempty"`
	// Indicates cardholder has authorized card to be stored for future payments.
	CardOnFile bool `json:"cardOnFile,omitempty"`
	// ID of the Moov account acting as a merchant or other entity authorized to store the card.  Defaults to your platform account ID if cardOnFile is set to true and no other account is provided.
	MerchantAccountID    string               `json:"merchantAccountID,omitempty"`
	CardAccountUpdater   CardAccountUpdater   `json:"cardAccountUpdater,omitempty"`
	DomesticPushToCard   DomesticPushToCard   `json:"domesticPushToCard,omitempty"`
	DomesticPullFromCard DomesticPullFromCard `json:"domesticPullFromCard,omitempty"`
	// Includes any payment methods generated for a newly linked card, removing the need to  call the List Payment Methods endpoint following a successful Link Card request.  **NOTE: This field is only populated for Link Card requests made with the `X-Wait-For` header.**
	PaymentMethods []BasicPaymentMethod `json:"paymentMethods,omitempty"`
}

// ApplePayResponse Describes an Apple Pay token on a Moov account.
type ApplePayPaymentMethod struct {
	Brand    CardBrand `json:"brand,omitempty"`
	CardType CardType  `json:"cardType,omitempty"`
	// User-friendly name of the tokenized card returned by Apple.  It usually contains the brand and the last four digits of the underlying card for example, \"Visa 1256\".  There is no standard format.
	CardDisplayName string `json:"cardDisplayName,omitempty"`
	// Uniquely identifies a linked payment card or token. For Apple Pay, the fingerprint is based on the tokenized card number and may vary based on the user's device. This field can be used to identify specific payment methods across multiple accounts on your platform.
	Fingerprint string         `json:"fingerprint,omitempty"`
	Expiration  CardExpiration `json:"expiration,omitempty"`
	// The last four digits of the Apple Pay token, which may differ from the tokenized card's last four digits
	DynamicLastFour string `json:"dynamicLastFour,omitempty"`
}

// buuuuunch of card stuff below....

// CardBrand The card brand.
type CardBrand string

// List of CardBrand
const (
	CardBrand_AmericanExpress CardBrand = "American Express"
	CardBrand_Discover        CardBrand = "Discover"
	CardBrand_Mastercard      CardBrand = "Mastercard"
	CardBrand_Visa            CardBrand = "Visa"
)

// CardType The type of the card.
type CardType string

// List of CardType
const (
	CardType_Dedit   CardType = "debit"
	CardType_Credit  CardType = "credit"
	CardType_Prepaid CardType = "prepaid"
	CardType_Unknown CardType = "unknown"
)

// CardExpiration The expiration date of the linked card or token.
type CardExpiration struct {
	Month string `json:"month,omitempty"`
	Year  string `json:"year,omitempty"`
}

// CardVerifications The results of submitting cardholder data to a card network for verification.
type CardVerifications struct {
	Cvv          CardVerificationResult `json:"cvv,omitempty"`
	AddressLine1 CardVerificationResult `json:"addressLine1,omitempty"`
	PostalCode   CardVerificationResult `json:"postalCode,omitempty"`
}

// CardVerificationResult the model 'CardVerificationResult'
type CardVerificationResult string

// List of CardVerificationResult
const (
	CardVerificationResult_NoMatch     CardVerificationResult = "noMatch"
	CardVerificationResult_Match       CardVerificationResult = "match"
	CardVerificationResult_NotChecked  CardVerificationResult = "notChecked"
	CardVerificationResult_Unavailable CardVerificationResult = "unavailable"
)

// DomesticPushToCard Indicates which level of domestic push-to-card transfer is supported by the card, if any.
type DomesticPushToCard string

// List of DomesticPushToCard
const (
	DomesticPushToCard_NotSupported DomesticPushToCard = "not-supported"
	DomesticPushToCard_Standard     DomesticPushToCard = "standard"
	DomesticPushToCard_FastFunds    DomesticPushToCard = "fast-funds"
	DomesticPushToCard_Unknown      DomesticPushToCard = "unknown"
)

// DomesticPullFromCard Indicates if the card supports domestic push-to-card transfer.
type DomesticPullFromCard string

// List of DomesticPullFromCard
const (
	DomesticPullFromCard_NotSupported DomesticPullFromCard = "not-supported"
	DomesticPullFromCard_Supported    DomesticPullFromCard = "supported"
	DomesticPullFromCard_Unknown      DomesticPullFromCard = "unknown"
)

// BasicPaymentMethod struct for BasicPaymentMethod
type BasicPaymentMethod struct {
	// UUID v4
	PaymentMethodID   string            `json:"paymentMethodID,omitempty"`
	PaymentMethodType PaymentMethodType `json:"paymentMethodType,omitempty"`
}
