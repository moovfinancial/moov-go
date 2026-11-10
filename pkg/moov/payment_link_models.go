package moov

import "time"

type Email string

type PaymentLinkStatus string

const (
	PaymentLinkStatus_Active   PaymentLinkStatus = "active"
	PaymentLinkStatus_Used     PaymentLinkStatus = "used"
	PaymentLinkStatus_Disabled PaymentLinkStatus = "disabled"
	PaymentLinkStatus_Expired  PaymentLinkStatus = "expired"
)

type PaymentLinkType string

const (
	PaymentLinkType_Payment             PaymentLinkType = "payment"
	PaymentLinkType_Payout              PaymentLinkType = "payout"
	PaymentLinkType_InvoicePayment      PaymentLinkType = "invoice-payment"
	PaymentLinkType_CustomAmountPayment PaymentLinkType = "custom-amount-payment"
)

// ACHCompanyEntryDescription An optional override of the default NACHA company entry description for a transfer.
type ACHCompanyEntryDescription string

// ACHOriginatingCompanyName An optional override of the default NACHA company name for a transfer.
type ACHOriginatingCompanyName string

// ACHPaymentDetails Options for payment links used to collect an ACH payment.
type ACHPaymentDetails struct {
	CompanyEntryDescription *ACHCompanyEntryDescription `json:"companyEntryDescription,omitempty"`
	OriginatingCompanyName  *ACHOriginatingCompanyName  `json:"originatingCompanyName,omitempty"`
}

// AmountDecimalRange If either `minimum` or `maximum` is omitted, the range is "open" on that end: `minimum` specified: `amt >= minimum` `maximum` specified: `amt <= maximum` both specified: `minimum <= amt <= maximum`
type AmountDecimalRange struct {
	Minimum *AmountDecimal `json:"minimum,omitempty"`
	Maximum *AmountDecimal `json:"maximum,omitempty"`
}

// AmountDecimalRangeUpdate If either `minimum` or `maximum` is omitted, the range is "open" on that end: `minimum` specified: `amt >= minimum` `maximum` specified: `amt <= maximum` both specified: `minimum <= amt <= maximum`
type AmountDecimalRangeUpdate struct {
	Minimum *AmountDecimalUpdate `json:"minimum,omitempty"`
	Maximum *AmountDecimalUpdate `json:"maximum,omitempty"`
}

type AmountDecimalUpdate struct {
	Currency     *Currency `json:"currency,omitempty"`
	ValueDecimal string    `json:"valueDecimal,omitempty"`
}

// CallToAction The text to be displayed on web form's submit button. If set to "auto" the UI will automatically select between "pay" and "confirm" for payments and payouts respectively.
type CallToAction string

const (
	CallToAction_Pay       CallToAction = "pay"
	CallToAction_Book      CallToAction = "book"
	CallToAction_Subscribe CallToAction = "subscribe"
	CallToAction_Donate    CallToAction = "donate"
	CallToAction_Confirm   CallToAction = "confirm"
	CallToAction_Auto      CallToAction = "auto"
)

// CardDynamicDescriptor An optional override of the default card statement descriptor for a transfer. Accounts must be enabled by Moov to set this field.
type CardDynamicDescriptor string

// CardPaymentDetails Options for payment links used to collect a card payment.
type CardPaymentDetails struct {
	DynamicDescriptor *CardDynamicDescriptor `json:"dynamicDescriptor,omitempty"`
}

// CollectionPaymentMethodType Payment methods allowed for collecting a payment.
type CollectionPaymentMethodType string

const (
	CollectionPaymentMethodType_ApplePay        CollectionPaymentMethodType = "apple-pay"
	CollectionPaymentMethodType_CardPayment     CollectionPaymentMethodType = "card-payment"
	CollectionPaymentMethodType_ACHDebitCollect CollectionPaymentMethodType = "ach-debit-collect"
	CollectionPaymentMethodType_GooglePay       CollectionPaymentMethodType = "google-pay"
)

// CreatePaymentLink Request to create a new payment link. A payment link must include exactly one of `payment`, `payout`, or `customAmountPayment` details. For payout payment links, `maxUses` will automatically be set to 1, as these are intended for a one-time disbursement to a specific recipient.
type CreatePaymentLink struct {
	PartnerAccountID        string                                 `json:"partnerAccountID"`
	MerchantPaymentMethodID string                                 `json:"merchantPaymentMethodID"`
	Amount                  *Amount                                `json:"amount,omitempty"`
	MaxUses                 int64                                  `json:"maxUses,omitempty"`
	ExpiresOn               *time.Time                             `json:"expiresOn,omitempty"`
	Display                 PaymentLinkDisplayOptions              `json:"display"`
	Customer                *PaymentLinkCustomerOptions            `json:"customer,omitempty"`
	Payment                 *PaymentLinkPaymentDetails             `json:"payment,omitempty"`
	Payout                  *PaymentLinkPayoutDetails              `json:"payout,omitempty"`
	CustomAmountPayment     *PaymentLinkCustomAmountPaymentDetails `json:"customAmountPayment,omitempty"`
	LineItems               *CreatePaymentLinkLineItems            `json:"lineItems,omitempty"`
	AmountDetails           *CreatePaymentLinkAmountDetails        `json:"amountDetails,omitempty"`
}

type CreatePaymentLinkAmountDetails struct {
	Tax       *AmountDecimal `json:"tax,omitempty"`
	Surcharge *AmountDecimal `json:"surcharge,omitempty"`
}

// CreatePaymentLinkLineItem Represents a single line item in a payment link, including optional modifiers and quantity.
type CreatePaymentLinkLineItem struct {
	Name      string                            `json:"name"`
	BasePrice AmountDecimal                     `json:"basePrice"`
	Quantity  int64                             `json:"quantity"`
	Options   []CreatePaymentLinkLineItemOption `json:"options,omitempty"`
	ImageIDs  []ImageID                         `json:"imageIDs,omitempty"`
	ProductID string                            `json:"productID,omitempty"`
}

// CreatePaymentLinkLineItemOption Represents a modifier or option applied to a line item.
type CreatePaymentLinkLineItemOption struct {
	Name          string         `json:"name"`
	Quantity      int64          `json:"quantity"`
	PriceModifier *AmountDecimal `json:"priceModifier,omitempty"`
	ImageIDs      []ImageID      `json:"imageIDs,omitempty"`
	Group         string         `json:"group,omitempty"`
}

// CreatePaymentLinkLineItems An optional collection of line items for a payment link. When line items are provided, their total plus tax must equal the payment link amount.
type CreatePaymentLinkLineItems struct {
	Items []CreatePaymentLinkLineItem `json:"items"`
}

// CreatePaymentLinkLineItemsUpdate An optional collection of line items for a payment link. When line items are provided, their total plus tax must equal the payment link amount.
type CreatePaymentLinkLineItemsUpdate struct {
	Items []CreatePaymentLinkLineItem `json:"items,omitempty"`
}

// Currency A 3-letter ISO 4217 currency code.
type Currency string

// DisbursementPaymentMethodType Payment methods allowed for disbursing funds.
type DisbursementPaymentMethodType string

const (
	DisbursementPaymentMethodType_PushToCard        DisbursementPaymentMethodType = "push-to-card"
	DisbursementPaymentMethodType_RtpCredit         DisbursementPaymentMethodType = "rtp-credit"
	DisbursementPaymentMethodType_ACHCreditSameDay  DisbursementPaymentMethodType = "ach-credit-same-day"
	DisbursementPaymentMethodType_ACHCreditStandard DisbursementPaymentMethodType = "ach-credit-standard"
	DisbursementPaymentMethodType_PushToApplePay    DisbursementPaymentMethodType = "push-to-apple-pay"
	DisbursementPaymentMethodType_PushToGooglePay   DisbursementPaymentMethodType = "push-to-google-pay"
)

// ImageID Unique identifier for a image resource.
type ImageID string

// ImagePublicID A unique identifier for an image, used in public image links.
type ImagePublicID string

type PaymentLink struct {
	Code                    PaymentLinkCode                        `json:"code"`
	PaymentLinkType         PaymentLinkType                        `json:"paymentLinkType"`
	Mode                    Mode                                   `json:"mode"`
	Status                  PaymentLinkStatus                      `json:"status"`
	PartnerAccountID        string                                 `json:"partnerAccountID"`
	MerchantAccountID       string                                 `json:"merchantAccountID"`
	OwnerAccountID          string                                 `json:"ownerAccountID"`
	MerchantPaymentMethodID string                                 `json:"merchantPaymentMethodID"`
	Link                    string                                 `json:"link"`
	Amount                  *Amount                                `json:"amount,omitempty"`
	Uses                    int64                                  `json:"uses"`
	MaxUses                 int64                                  `json:"maxUses,omitempty"`
	LastUsedOn              *time.Time                             `json:"lastUsedOn,omitempty"`
	ExpiresOn               *time.Time                             `json:"expiresOn,omitempty"`
	Display                 PaymentLinkDisplayOptions              `json:"display"`
	Customer                PaymentLinkCustomerOptions             `json:"customer"`
	Payment                 *PaymentLinkPaymentDetails             `json:"payment,omitempty"`
	Payout                  *PaymentLinkPayoutDetails              `json:"payout,omitempty"`
	CustomAmountPayment     *PaymentLinkCustomAmountPaymentDetails `json:"customAmountPayment,omitempty"`
	LineItems               *PaymentLinkLineItems                  `json:"lineItems,omitempty"`
	CreatedOn               time.Time                              `json:"createdOn"`
	UpdatedOn               time.Time                              `json:"updatedOn"`
	DisabledOn              *time.Time                             `json:"disabledOn,omitempty"`
	AmountDetails           *PaymentLinkAmountDetails              `json:"amountDetails,omitempty"`
}

type PaymentLinkAmountDetails struct {
	Tax       *AmountDecimal `json:"tax,omitempty"`
	Surcharge *AmountDecimal `json:"surcharge,omitempty"`
}

// PaymentLinkCode Unique code identifying this payment link.
type PaymentLinkCode string

// PaymentLinkCustomAmountPaymentDetails A custom amount payment link shares all the options of a `payment` link, but the payor chooses how much to pay rather than the merchant fixing the amount. The amount may optionally be constrained to a range.
type PaymentLinkCustomAmountPaymentDetails struct {
	AllowedMethods   []CollectionPaymentMethodType `json:"allowedMethods"`
	CardDetails      *CardPaymentDetails           `json:"cardDetails,omitempty"`
	ACHDetails       *ACHPaymentDetails            `json:"achDetails,omitempty"`
	Metadata         map[string]string             `json:"metadata,omitempty"`
	AmountRange      *AmountDecimalRange           `json:"amountRange,omitempty"`
	SuggestedAmounts []AmountDecimal               `json:"suggestedAmounts,omitempty"`
}

// PaymentLinkCustomAmountPaymentDetailsUpdate A custom amount payment link shares all the options of a `payment` link, but the payor chooses how much to pay rather than the merchant fixing the amount. The amount may optionally be constrained to a range.
type PaymentLinkCustomAmountPaymentDetailsUpdate struct {
	AllowedMethods   []CollectionPaymentMethodType `json:"allowedMethods,omitempty"`
	CardDetails      *CardPaymentDetails           `json:"cardDetails,omitempty"`
	ACHDetails       *ACHPaymentDetails            `json:"achDetails,omitempty"`
	Metadata         map[string]string             `json:"metadata,omitempty"`
	AmountRange      *AmountDecimalRangeUpdate     `json:"amountRange,omitempty"`
	SuggestedAmounts []AmountDecimal               `json:"suggestedAmounts,omitempty"`
}

type PaymentLinkCustomerOptions struct {
	RequireAddress bool              `json:"requireAddress,omitempty"`
	RequirePhone   bool              `json:"requirePhone,omitempty"`
	TippingEnabled bool              `json:"tippingEnabled,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// PaymentLinkDisplayOptions Customizable display options for a payment link.
type PaymentLinkDisplayOptions struct {
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	CallToAction CallToAction `json:"callToAction"`
}

// PaymentLinkDisplayOptionsUpdate Customizable display options for a payment link.
type PaymentLinkDisplayOptionsUpdate struct {
	Title        string        `json:"title,omitempty"`
	Description  string        `json:"description,omitempty"`
	CallToAction *CallToAction `json:"callToAction,omitempty"`
}

// PaymentLinkLineItem Represents a single line item in a payment link, including optional modifiers and quantity.
type PaymentLinkLineItem struct {
	Name      string                             `json:"name"`
	BasePrice AmountDecimal                      `json:"basePrice"`
	Quantity  int64                              `json:"quantity"`
	Options   []PaymentLinkLineItemOption        `json:"options,omitempty"`
	Images    []PaymentLinkLineItemImageMetadata `json:"images,omitempty"`
	ProductID string                             `json:"productID,omitempty"`
}

type PaymentLinkLineItemImageMetadata struct {
	ImageID  ImageID       `json:"imageID"`
	AltText  string        `json:"altText,omitempty"`
	Link     string        `json:"link"`
	PublicID ImagePublicID `json:"publicID"`
}

// PaymentLinkLineItemOption Represents a modifier or option applied to a line item.
type PaymentLinkLineItemOption struct {
	Name          string                             `json:"name"`
	Quantity      int64                              `json:"quantity"`
	PriceModifier *AmountDecimal                     `json:"priceModifier,omitempty"`
	Images        []PaymentLinkLineItemImageMetadata `json:"images,omitempty"`
	Group         string                             `json:"group,omitempty"`
}

// PaymentLinkLineItems An optional collection of line items for a payment link. When line items are provided, their total plus tax must equal the payment link amount.
type PaymentLinkLineItems struct {
	Items []PaymentLinkLineItem `json:"items"`
}

// PaymentLinkPaymentDetails Options for payment links used to collect payment.
type PaymentLinkPaymentDetails struct {
	AllowedMethods []CollectionPaymentMethodType `json:"allowedMethods"`
	CardDetails    *CardPaymentDetails           `json:"cardDetails,omitempty"`
	ACHDetails     *ACHPaymentDetails            `json:"achDetails,omitempty"`
	Metadata       map[string]string             `json:"metadata,omitempty"`
}

// PaymentLinkPaymentDetailsUpdate Options for payment links used to collect payment.
type PaymentLinkPaymentDetailsUpdate struct {
	AllowedMethods []CollectionPaymentMethodType `json:"allowedMethods,omitempty"`
	CardDetails    *CardPaymentDetails           `json:"cardDetails,omitempty"`
	ACHDetails     *ACHPaymentDetails            `json:"achDetails,omitempty"`
	Metadata       map[string]string             `json:"metadata,omitempty"`
}

// PaymentLinkPayoutDetails Options for payout links used to send a payout.
type PaymentLinkPayoutDetails struct {
	AllowedMethods []DisbursementPaymentMethodType `json:"allowedMethods"`
	Recipient      PayoutRecipient                 `json:"recipient"`
	Metadata       map[string]string               `json:"metadata,omitempty"`
	PushOptions    *PushOptions                    `json:"pushOptions,omitempty"`
	FeePaidBy      map[string]string               `json:"feePaidBy,omitempty"`
}

// PaymentLinkPayoutDetailsUpdate Options for payout links used to send a payout.
type PaymentLinkPayoutDetailsUpdate struct {
	AllowedMethods []DisbursementPaymentMethodType `json:"allowedMethods,omitempty"`
	Recipient      *PayoutRecipient                `json:"recipient,omitempty"`
	Metadata       map[string]string               `json:"metadata,omitempty"`
	PushOptions    *PushOptionsUpdate              `json:"pushOptions,omitempty"`
	FeePaidBy      map[string]string               `json:"feePaidBy,omitempty"`
}

// PayoutRecipient Specify the intended recipient of the payout. Either `email` or `phone` must be specified, but not both. This information will be used to authenticate the end user when they follow the payment link.
type PayoutRecipient struct {
	Email *Email `json:"email,omitempty"`
	Phone *Phone `json:"phone,omitempty"`
}

// PushDeliverySpeed Delivery speed options for push-to-card payouts.
type PushDeliverySpeed string

const (
	PushDeliverySpeed_Instant  PushDeliverySpeed = "instant"
	PushDeliverySpeed_Deferred PushDeliverySpeed = "deferred"
)

type PushOptions struct {
	AllowedSpeeds []PushDeliverySpeed `json:"allowedSpeeds"`
	DeferredBy    string              `json:"deferredBy,omitempty"`
}

type PushOptionsUpdate struct {
	AllowedSpeeds []PushDeliverySpeed `json:"allowedSpeeds,omitempty"`
	DeferredBy    string              `json:"deferredBy,omitempty"`
}

type UpdatePaymentLink struct {
	Amount              *AmountUpdate                                `json:"amount,omitempty"`
	ExpiresOn           *Nullable[time.Time]                         `json:"expiresOn,omitempty"`
	Display             *PaymentLinkDisplayOptionsUpdate             `json:"display,omitempty"`
	Customer            *PaymentLinkCustomerOptions                  `json:"customer,omitempty"`
	Payment             *PaymentLinkPaymentDetailsUpdate             `json:"payment,omitempty"`
	Payout              *PaymentLinkPayoutDetailsUpdate              `json:"payout,omitempty"`
	CustomAmountPayment *PaymentLinkCustomAmountPaymentDetailsUpdate `json:"customAmountPayment,omitempty"`
	LineItems           *CreatePaymentLinkLineItemsUpdate            `json:"lineItems,omitempty"`
	AmountDetails       *UpdatePaymentLinkAmountDetails              `json:"amountDetails,omitempty"`
}

type UpdatePaymentLinkAmountDetails struct {
	Tax       *Nullable[AmountDecimalUpdate] `json:"tax,omitempty"`
	Surcharge *Nullable[AmountDecimalUpdate] `json:"surcharge,omitempty"`
}

type AmountUpdate struct {
	Currency *Currency `json:"currency,omitempty"`
	Value    int64     `json:"value,omitempty"`
}
