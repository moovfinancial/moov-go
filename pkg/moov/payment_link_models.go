package moov

import "time"

// This file contains the data models for the Payment Links resource.
// See https://docs.moov.io/api/money-movement/payment-links/ for the API reference.

type CreateAccountType string

const (
	CreateAccountType_Individual CreateAccountType = "individual"
	CreateAccountType_Business   CreateAccountType = "business"
)

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

type ACHPaymentSettings struct {
	// The description that shows up on ACH transactions. This will default to the account's display name on account creation.
	CompanyName string `json:"companyName"`
}

// AmountDecimalRange If either `minimum` or `maximum` is omitted, the range is "open" on that end: `minimum` specified: `amt >= minimum` `maximum` specified: `amt <= maximum` both specified: `minimum <= amt <= maximum`
type AmountDecimalRange struct {
	// Minimum amount allowed in the range
	Minimum *AmountDecimal `json:"minimum,omitempty"`
	// Maximum amount allowed in the range
	Maximum *AmountDecimal `json:"maximum,omitempty"`
}

// AmountDecimalRangeUpdate If either `minimum` or `maximum` is omitted, the range is "open" on that end: `minimum` specified: `amt >= minimum` `maximum` specified: `amt <= maximum` both specified: `minimum <= amt <= maximum`
type AmountDecimalRangeUpdate struct {
	// Minimum amount allowed in the range
	Minimum *AmountDecimalUpdate `json:"minimum,omitempty"`
	// Maximum amount allowed in the range
	Maximum *AmountDecimalUpdate `json:"maximum,omitempty"`
}

type AmountDecimalUpdate struct {
	Currency *Currency `json:"currency,omitempty"`
	// A decimal-formatted numerical string that represents up to 9 decimal place precision. For example, $12.987654321 is '12.987654321'.
	ValueDecimal string `json:"valueDecimal,omitempty"`
}

type BirthDate struct {
	Day   int64 `json:"day"`
	Month int64 `json:"month"`
	Year  int64 `json:"year"`
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
	// The partner's Moov account ID.
	PartnerAccountID string `json:"partnerAccountID"`
	// The merchant's preferred payment method ID. Must be a wallet payment method.
	MerchantPaymentMethodID string `json:"merchantPaymentMethodID"`
	// The fixed amount of the payment link. In API versions before `2026.07.00`, this was a required field. In API version `2026.07.00` and beyond, this field is required for `payment` and `payout` links and must be omitted for `customAmountPayment` links, where the payor chooses the amount.
	Amount *Amount `json:"amount,omitempty"`
	// An optional limit on the number of times this payment link can be used. **For payouts, `maxUses` is always 1.**
	MaxUses int64 `json:"maxUses,omitempty"`
	// An optional expiration date for this payment link.
	ExpiresOn *time.Time                  `json:"expiresOn,omitempty"`
	Display   PaymentLinkDisplayOptions   `json:"display"`
	Customer  *PaymentLinkCustomerOptions `json:"customer,omitempty"`
	Payment   *PaymentLinkPaymentDetails  `json:"payment,omitempty"`
	Payout    *PaymentLinkPayoutDetails   `json:"payout,omitempty"`
	// Options for a custom amount payment link. Mutually exclusive with `payment` and `payout`.
	CustomAmountPayment *PaymentLinkCustomAmountPaymentDetails `json:"customAmountPayment,omitempty"`
	LineItems           *CreatePaymentLinkLineItems            `json:"lineItems,omitempty"`
	AmountDetails       *CreatePaymentLinkAmountDetails        `json:"amountDetails,omitempty"`
}

type CreatePaymentLinkAmountDetails struct {
	// The amount of tax applied to the payment link.
	Tax *AmountDecimal `json:"tax,omitempty"`
	// The amount of surcharge applied to the payment link.
	Surcharge *AmountDecimal `json:"surcharge,omitempty"`
}

// CreatePaymentLinkLineItem Represents a single line item in a payment link, including optional modifiers and quantity.
type CreatePaymentLinkLineItem struct {
	// The name of the item.
	Name string `json:"name"`
	// The base price of the item before applying option modifiers.
	BasePrice AmountDecimal `json:"basePrice"`
	// The quantity of this item.
	Quantity int64 `json:"quantity"`
	// Optional list of modifiers applied to this item (e.g., toppings, upgrades, customizations).
	Options []CreatePaymentLinkLineItemOption `json:"options,omitempty"`
	// Optional list of images associated with this line item. This field is being deprecated in favor using the images associated with a productID and will soon be unsupported.
	ImageIDs []ImageID `json:"imageIDs,omitempty"`
	// Optional unique identifier associating the line item with a product. When provided, images associated with the product will be included on the line item. This does not populate other details of the line item.
	ProductID string `json:"productID,omitempty"`
}

// CreatePaymentLinkLineItemOption Represents a modifier or option applied to a line item.
type CreatePaymentLinkLineItemOption struct {
	// The name of the option or modifier.
	Name string `json:"name"`
	// The quantity of this option.
	Quantity int64 `json:"quantity"`
	// Optional price modification applied by this option. Can be positive, negative, or zero.
	PriceModifier *AmountDecimal `json:"priceModifier,omitempty"`
	// Optional list of images associated with this line item option. This field is being deprecated in favor using the images associated with a productID and will soon be unsupported.
	ImageIDs []ImageID `json:"imageIDs,omitempty"`
	// Optional group identifier to categorize related options (e.g., 'toppings').
	Group string `json:"group,omitempty"`
}

// CreatePaymentLinkLineItems An optional collection of line items for a payment link. When line items are provided, their total plus tax must equal the payment link amount.
type CreatePaymentLinkLineItems struct {
	// The list of line items.
	Items []CreatePaymentLinkLineItem `json:"items"`
}

// CreatePaymentLinkLineItemsUpdate An optional collection of line items for a payment link. When line items are provided, their total plus tax must equal the payment link amount.
type CreatePaymentLinkLineItemsUpdate struct {
	// The list of line items.
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

type IndividualName struct {
	// The individual's first given name.
	FirstName string `json:"firstName"`
	// The individual's second given name, if any.
	MiddleName string `json:"middleName,omitempty"`
	// The individual's family name.
	LastName string `json:"lastName"`
	// Suffix of a given name.
	Suffix string `json:"suffix,omitempty"`
}

// ManualTermsOfService Describes the acceptance of the Terms of Service. All data is required, and must be from the user.
type ManualTermsOfService struct {
	// The date and time the terms of service were accepted.
	AcceptedDate time.Time `json:"acceptedDate"`
	// The IP address (IPv4 or IPv6) from which the terms of service were accepted.
	AcceptedIP string `json:"acceptedIP"`
	// The user-agent of the user making the request.
	AcceptedUserAgent string `json:"acceptedUserAgent"`
	AcceptedDomain    string `json:"acceptedDomain"`
}

type PaymentLink struct {
	Code            PaymentLinkCode   `json:"code"`
	PaymentLinkType PaymentLinkType   `json:"paymentLinkType"`
	Mode            Mode              `json:"mode"`
	Status          PaymentLinkStatus `json:"status"`
	// The partner's Moov account ID.
	PartnerAccountID string `json:"partnerAccountID"`
	// The merchant's Moov account ID.
	MerchantAccountID string `json:"merchantAccountID"`
	// The payment link's owner's Moov account ID.
	OwnerAccountID string `json:"ownerAccountID"`
	// The merchant's preferred payment method ID. Must be a wallet payment method.
	MerchantPaymentMethodID string `json:"merchantPaymentMethodID"`
	// Link to the payment landing page for this payment link.
	Link string `json:"link"`
	// The fixed amount of the payment link. In API versions before `2026.07.00`, this was a required field. In API version `2026.07.00` and beyond, this field is present for `payment` and `payout` links and omitted for `customAmountPayment` links, where the payor chooses the amount.
	Amount *Amount `json:"amount,omitempty"`
	// The number of times this payment link has been used.
	Uses int64 `json:"uses"`
	// An optional limit on the number of times this payment link can be used. **For payouts, `maxUses` is always 1.**
	MaxUses int64 `json:"maxUses,omitempty"`
	// The timestamp when this payment link was last used.
	LastUsedOn *time.Time `json:"lastUsedOn,omitempty"`
	// An optional expiration date for this payment link.
	ExpiresOn           *time.Time                             `json:"expiresOn,omitempty"`
	Display             PaymentLinkDisplayOptions              `json:"display"`
	Customer            PaymentLinkCustomerOptions             `json:"customer"`
	Payment             *PaymentLinkPaymentDetails             `json:"payment,omitempty"`
	Payout              *PaymentLinkPayoutDetails              `json:"payout,omitempty"`
	CustomAmountPayment *PaymentLinkCustomAmountPaymentDetails `json:"customAmountPayment,omitempty"`
	LineItems           *PaymentLinkLineItems                  `json:"lineItems,omitempty"`
	CreatedOn           time.Time                              `json:"createdOn"`
	UpdatedOn           time.Time                              `json:"updatedOn"`
	DisabledOn          *time.Time                             `json:"disabledOn,omitempty"`
	AmountDetails       *PaymentLinkAmountDetails              `json:"amountDetails,omitempty"`
}

type PaymentLinkAmountDetails struct {
	// The amount of tax applied to the payment link.
	Tax *AmountDecimal `json:"tax,omitempty"`
	// The amount of surcharge applied to the payment link.
	Surcharge *AmountDecimal `json:"surcharge,omitempty"`
}

// PaymentLinkCode Unique code identifying this payment link.
type PaymentLinkCode string

// PaymentLinkCustomAmountPaymentDetails A custom amount payment link shares all the options of a `payment` link, but the payor chooses how much to pay rather than the merchant fixing the amount. The amount may optionally be constrained to a range.
type PaymentLinkCustomAmountPaymentDetails struct {
	// A list of payment methods that should be supported for this payment link.
	AllowedMethods []CollectionPaymentMethodType `json:"allowedMethods"`
	CardDetails    *CardPaymentDetails           `json:"cardDetails,omitempty"`
	ACHDetails     *ACHPaymentDetails            `json:"achDetails,omitempty"`
	// Optional free-form metadata for the transfer.
	Metadata map[string]string `json:"metadata,omitempty"`
	// The minimum and maximum amounts the payor can specify. When omitted, the payor may enter any amount. Both bounds must use USD.
	AmountRange *AmountDecimalRange `json:"amountRange,omitempty"`
	// Optional preset amounts displayed to the payor. Each suggested amount must use USD and, when an `amountRange` is set, fall within it.
	SuggestedAmounts []AmountDecimal `json:"suggestedAmounts,omitempty"`
}

// PaymentLinkCustomAmountPaymentDetailsUpdate A custom amount payment link shares all the options of a `payment` link, but the payor chooses how much to pay rather than the merchant fixing the amount. The amount may optionally be constrained to a range.
type PaymentLinkCustomAmountPaymentDetailsUpdate struct {
	// A list of payment methods that should be supported for this payment link.
	AllowedMethods []CollectionPaymentMethodType `json:"allowedMethods,omitempty"`
	CardDetails    *CardPaymentDetails           `json:"cardDetails,omitempty"`
	ACHDetails     *ACHPaymentDetails            `json:"achDetails,omitempty"`
	// Optional free-form metadata for the transfer.
	Metadata map[string]string `json:"metadata,omitempty"`
	// The minimum and maximum amounts the payor can specify. When omitted, the payor may enter any amount. Both bounds must use USD.
	AmountRange *AmountDecimalRangeUpdate `json:"amountRange,omitempty"`
	// Optional preset amounts displayed to the payor. Each suggested amount must use USD and, when an `amountRange` is set, fall within it.
	SuggestedAmounts []AmountDecimal `json:"suggestedAmounts,omitempty"`
}

type PaymentLinkCustomerOptions struct {
	// If true, a billing address is required when completing the payment form.
	RequireAddress bool `json:"requireAddress,omitempty"`
	// If true, a phone number is required when completing the payment form.
	RequirePhone bool `json:"requirePhone,omitempty"`
	// If true, tipping is enabled on the payment form. Defaults to false.
	TippingEnabled bool `json:"tippingEnabled,omitempty"`
	// Optional free-form metadata for the Moov account that will represent this customer.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// PaymentLinkDisplayOptions Customizable display options for a payment link.
type PaymentLinkDisplayOptions struct {
	// The payment page title displayed to the user.
	Title string `json:"title"`
	// A payment description displayed to the user.
	Description  string       `json:"description"`
	CallToAction CallToAction `json:"callToAction"`
}

// PaymentLinkDisplayOptionsUpdate Customizable display options for a payment link.
type PaymentLinkDisplayOptionsUpdate struct {
	// The payment page title displayed to the user.
	Title string `json:"title,omitempty"`
	// A payment description displayed to the user.
	Description  string        `json:"description,omitempty"`
	CallToAction *CallToAction `json:"callToAction,omitempty"`
}

// PaymentLinkLineItem Represents a single line item in a payment link, including optional modifiers and quantity.
type PaymentLinkLineItem struct {
	// The name of the item.
	Name string `json:"name"`
	// The base price of the item before applying option modifiers.
	BasePrice AmountDecimal `json:"basePrice"`
	// The quantity of this item.
	Quantity int64 `json:"quantity"`
	// Optional list of modifiers applied to this item (e.g., toppings, upgrades, customizations).
	Options []PaymentLinkLineItemOption `json:"options,omitempty"`
	// Optional list of images associated with this line item.
	Images []PaymentLinkLineItemImageMetadata `json:"images,omitempty"`
	// Optional unique identifier associating the line item with a product.
	ProductID string `json:"productID,omitempty"`
}

type PaymentLinkLineItemImageMetadata struct {
	ImageID ImageID `json:"imageID"`
	// Alternative text for the image.
	AltText string `json:"altText,omitempty"`
	// The image's public URL.
	Link     string        `json:"link"`
	PublicID ImagePublicID `json:"publicID"`
}

// PaymentLinkLineItemOption Represents a modifier or option applied to a line item.
type PaymentLinkLineItemOption struct {
	// The name of the option or modifier.
	Name string `json:"name"`
	// The quantity of this option.
	Quantity int64 `json:"quantity"`
	// Optional price modification applied by this option. Can be positive, negative, or zero.
	PriceModifier *AmountDecimal `json:"priceModifier,omitempty"`
	// Optional list of images associated with this line item option.
	Images []PaymentLinkLineItemImageMetadata `json:"images,omitempty"`
	// Optional group identifier to categorize related options (e.g., 'toppings').
	Group string `json:"group,omitempty"`
}

// PaymentLinkLineItems An optional collection of line items for a payment link. When line items are provided, their total plus tax must equal the payment link amount.
type PaymentLinkLineItems struct {
	// The list of line items.
	Items []PaymentLinkLineItem `json:"items"`
}

// PaymentLinkPaymentDetails Options for payment links used to collect payment.
type PaymentLinkPaymentDetails struct {
	// A list of payment methods that should be supported for this payment link.
	AllowedMethods []CollectionPaymentMethodType `json:"allowedMethods"`
	CardDetails    *CardPaymentDetails           `json:"cardDetails,omitempty"`
	ACHDetails     *ACHPaymentDetails            `json:"achDetails,omitempty"`
	// Optional free-form metadata for the transfer.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// PaymentLinkPaymentDetailsUpdate Options for payment links used to collect payment.
type PaymentLinkPaymentDetailsUpdate struct {
	// A list of payment methods that should be supported for this payment link.
	AllowedMethods []CollectionPaymentMethodType `json:"allowedMethods,omitempty"`
	CardDetails    *CardPaymentDetails           `json:"cardDetails,omitempty"`
	ACHDetails     *ACHPaymentDetails            `json:"achDetails,omitempty"`
	// Optional free-form metadata for the transfer.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// PaymentLinkPayoutDetails Options for payout links used to send a payout.
type PaymentLinkPayoutDetails struct {
	// A list of payment methods that should be supported for this payment link.
	AllowedMethods []DisbursementPaymentMethodType `json:"allowedMethods"`
	Recipient      PayoutRecipient                 `json:"recipient"`
	// Optional free-form metadata for the transfer.
	Metadata map[string]string `json:"metadata,omitempty"`
	// Delivery options for `push-to-card` and `push-to-apple-pay` payouts.
	PushOptions *PushOptions `json:"pushOptions,omitempty"`
	//             Indicates which party pays the fee, keyed by `PayoutFeePaidByKey`. If keys are not set, the default is `source`. Possible `PayoutFeePaidByKey` keys: `instant-push-to-card`, `deferred-push-to-card`, `instant-push-to-apple-pay`, `deferred-push-to-apple-pay`, `rtp-credit`, `ach-credit-same-day`, `ach-credit-standard`, `push-to-google-pay`
	FeePaidBy map[string]string `json:"feePaidBy,omitempty"`
}

// PaymentLinkPayoutDetailsUpdate Options for payout links used to send a payout.
type PaymentLinkPayoutDetailsUpdate struct {
	// A list of payment methods that should be supported for this payment link.
	AllowedMethods []DisbursementPaymentMethodType `json:"allowedMethods,omitempty"`
	Recipient      *PayoutRecipient                `json:"recipient,omitempty"`
	// Optional free-form metadata for the transfer.
	Metadata map[string]string `json:"metadata,omitempty"`
	// Delivery options for `push-to-card` and `push-to-apple-pay` payouts.
	PushOptions *PushOptionsUpdate `json:"pushOptions,omitempty"`
	//             Indicates which party pays the fee, keyed by `PayoutFeePaidByKey`. If keys are not set, the default is `source`. Possible `PayoutFeePaidByKey` keys: `instant-push-to-card`, `deferred-push-to-card`, `instant-push-to-apple-pay`, `deferred-push-to-apple-pay`, `rtp-credit`, `ach-credit-same-day`, `ach-credit-standard`, `push-to-google-pay`
	FeePaidBy map[string]string `json:"feePaidBy,omitempty"`
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
	// Delivery speeds the recipient may choose from for `push-to-card` and `push-to-apple-pay`. Include `instant` to allow immediate delivery, `deferred` to allow delayed delivery, or both to let the recipient choose at checkout.
	AllowedSpeeds []PushDeliverySpeed `json:"allowedSpeeds"`
	// Delay before delivering a deferred payout. Required when `allowedSpeeds` includes `deferred`. Accepted values are `24h` or `48h`.
	DeferredBy string `json:"deferredBy,omitempty"`
}

type PushOptionsUpdate struct {
	// Delivery speeds the recipient may choose from for `push-to-card` and `push-to-apple-pay`. Include `instant` to allow immediate delivery, `deferred` to allow delayed delivery, or both to let the recipient choose at checkout.
	AllowedSpeeds []PushDeliverySpeed `json:"allowedSpeeds,omitempty"`
	// Delay before delivering a deferred payout. Required when `allowedSpeeds` includes `deferred`. Accepted values are `24h` or `48h`.
	DeferredBy string `json:"deferredBy,omitempty"`
}

// Settings User provided settings to manage an account.
type Settings struct {
	CardPayment *CardPaymentSettings `json:"cardPayment,omitempty"`
	ACHPayment  *ACHPaymentSettings  `json:"achPayment,omitempty"`
}

type UpdatePaymentLink struct {
	Amount    *AmountUpdate                    `json:"amount,omitempty"`
	ExpiresOn *Nullable[time.Time]             `json:"expiresOn,omitempty"`
	Display   *PaymentLinkDisplayOptionsUpdate `json:"display,omitempty"`
	Customer  *PaymentLinkCustomerOptions      `json:"customer,omitempty"`
	Payment   *PaymentLinkPaymentDetailsUpdate `json:"payment,omitempty"`
	Payout    *PaymentLinkPayoutDetailsUpdate  `json:"payout,omitempty"`
	// Options for a custom amount payment link. A payment link's type cannot be changed after creation.
	CustomAmountPayment *PaymentLinkCustomAmountPaymentDetailsUpdate `json:"customAmountPayment,omitempty"`
	LineItems           *CreatePaymentLinkLineItemsUpdate            `json:"lineItems,omitempty"`
	AmountDetails       *UpdatePaymentLinkAmountDetails              `json:"amountDetails,omitempty"`
}

type UpdatePaymentLinkAmountDetails struct {
	// The amount of tax applied to the payment link. Nullable: send SetNull to clear.
	Tax *Nullable[AmountDecimalUpdate] `json:"tax,omitempty"`
	// The amount of surcharge applied to the payment link. Nullable: send SetNull to clear.
	Surcharge *Nullable[AmountDecimalUpdate] `json:"surcharge,omitempty"`
}

// AmountUpdate
type AmountUpdate struct {
	Currency *Currency `json:"currency,omitempty"`
	// Quantity in the smallest unit of the specified currency. In USD this is cents, for example, $12.04 is 1204 and $0.99 is 99.
	Value int64 `json:"value,omitempty"`
}
