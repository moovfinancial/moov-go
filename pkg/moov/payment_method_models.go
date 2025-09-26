package moov

// PaymentMethod A method of moving money
type PaymentMethod struct {
	PaymentMethodID   string                    `json:"paymentMethodID,omitempty"`
	PaymentMethodType PaymentMethodType         `json:"paymentMethodType,omitempty"`
	Wallet            *WalletPaymentMethod      `json:"wallet,omitempty"`
	BankAccount       *BankAccountPaymentMethod `json:"bankAccount,omitempty"`
	Card              *CardPaymentMethod        `json:"card,omitempty"`
	ApplePay          *ApplePayPaymentMethod    `json:"applePay,omitempty"`
}

// BasicPaymentMethod struct for BasicPaymentMethod
type BasicPaymentMethod struct {
	PaymentMethodID   string            `json:"paymentMethodID,omitempty"`
	PaymentMethodType PaymentMethodType `json:"paymentMethodType,omitempty"`
}

// PaymentMethodType The payment method type that represents a payment rail and directionality
type PaymentMethodType string

// List of PaymentMethodType
const (
	PaymentMethodType_MoovWallet         PaymentMethodType = "moov-wallet"
	PaymentMethodType_AchDebitFund       PaymentMethodType = "ach-debit-fund"
	PaymentMethodType_AchDebitCollect    PaymentMethodType = "ach-debit-collect"
	PaymentMethodType_AchCreditStandard  PaymentMethodType = "ach-credit-standard"
	PaymentMethodType_AchCreditSameDay   PaymentMethodType = "ach-credit-same-day" // #nosec G101
	PaymentMethodType_RtpCredit          PaymentMethodType = "rtp-credit"
	PaymentMethodType_CardPayment        PaymentMethodType = "card-payment"
	PaymentMethodType_ApplePay           PaymentMethodType = "apple-pay"
	PaymentMethodType_PushToCard         PaymentMethodType = "push-to-card"
	PaymentMethodType_PullFromCard       PaymentMethodType = "pull-from-card"
	PaymentMethodType_CardPresentPayment PaymentMethodType = "card-present-payment"
)

// WalletPaymentMethod A Moov wallet to store funds for transfers.
type WalletPaymentMethod struct {
	WalletID         string     `json:"walletID,omitempty"`
	PartnerAccountID string     `json:"partnerAccountID,omitempty"`
	WalletType       WalletType `json:"walletType,omitempty"`
}

type BankAccountPaymentMethod BankAccount

// Card Describes a card on a Moov account.
type CardPaymentMethod struct {
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

type CardUpdateType string

// List of CardUpdateType
const (
	CardUpdateType_AccountClosed     CardUpdateType = "account-closed"
	CardUpdateType_ContactCardHolder CardUpdateType = "contact-cardholder"
	CardUpdateType_ExpirationUpdate  CardUpdateType = "expiration-update"
	CardUpdateType_NumberUpdate      CardUpdateType = "number-update"
)

// CardExpiration The expiration date of the linked card or token.
type CardExpiration struct {
	Month string `json:"month,omitempty"`
	Year  string `json:"year,omitempty"`
}

// CardVerifications The results of submitting cardholder data to a card network for verification.
type CardVerifications struct {
	Cvv          CardVerificationResult  `json:"cvv,omitempty"`
	AddressLine1 CardVerificationResult  `json:"addressLine1,omitempty"`
	PostalCode   CardVerificationResult  `json:"postalCode,omitempty"`
	AccountName  AccountNameVerification `json:"accountName,omitempty"`
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

type AccountNameVerification struct {
	FirstName  CardVerificationResult `json:"firstName,omitempty"`
	LastName   CardVerificationResult `json:"lastName,omitempty"`
	MiddleName CardVerificationResult `json:"middleName,omitempty"`
	FullName   CardVerificationResult `json:"fullName,omitempty"`
}

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

type EntryMode string

const (
	EntryModeContactless EntryMode = "contactless"
)

// TerminalCard - Describes payment card details captured with tap or in-person payment.
type TerminalCardPaymentMethod struct {
	// How the card information was entered into the point of sale terminal.
	EntryMode EntryMode `json:"entryMode,omitempty"`
	// The card brand.
	Brand CardBrand `json:"brand,omitempty"`
	Bin   string    `json:"bin,omitempty"`
	// The type of the card.
	CardType CardType `json:"cardType,omitempty"`
	// The expiration date of the card or token.
	Expiration CardExpiration `json:"expiration,omitempty"`
	// Uniquely identifies a terminal card.
	Fingerprint string `json:"fingerprint,omitempty"`
	// The name of the cardholder as it appears on the card.
	HolderName string `json:"holderName,omitempty"`
	// Financial institution that issued the card.
	Issuer string `json:"issuer,omitempty"`
	// Country where the card was issued.
	IssuerCountry string `json:"issuerCountry,omitempty"`
	// Last four digits of the card number
	LastFourCardNumber string `json:"lastFourCardNumber,omitempty"`
	// Identifier for the point of sale terminal application.
	ApplicationID string `json:"applicationID,omitempty"`
	// Name label for the point of sale terminal application.
	ApplicationName string `json:"applicationName,omitempty"`
}
