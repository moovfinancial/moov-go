package moov

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Card struct {
	CardID               string             `json:"cardID,omitempty"`
	Fingerprint          string             `json:"fingerprint,omitempty"`
	Brand                string             `json:"brand,omitempty"`
	CardCategory         string             `json:"cardCategory,omitempty"`
	CardType             string             `json:"cardType,omitempty"`
	LastFourCardNumber   string             `json:"lastFourCardNumber,omitempty"`
	Bin                  string             `json:"bin,omitempty"`
	Expiration           Expiration         `json:"expiration,omitempty"`
	HolderName           string             `json:"holderName,omitempty"`
	BillingAddress       Address            `json:"billingAddress,omitempty"`
	CardVerification     CardVerification   `json:"cardVerification,omitempty"`
	Issuer               string             `json:"issuer,omitempty"`
	IssuerCountry        string             `json:"issuerCountry,omitempty"`
	IssuerURL            string             `json:"issuerURL,omitempty"`
	IssuerPhone          string             `json:"issuerPhone,omitempty"`
	Commercial           *bool              `json:"commercial,omitempty"`
	Regulated            *bool              `json:"regulated,omitempty"`
	CardOnFile           bool               `json:"cardOnFile,omitempty"`
	MerchantAccountID    string             `json:"merchantAccountID,omitempty"`
	CardAccountUpdater   CardAccountUpdater `json:"cardAccountUpdater,omitempty"`
	DomesticPushToCard   string             `json:"domesticPushToCard,omitempty"`
	DomesticPullFromCard string             `json:"domesticPullFromCard,omitempty"`
	PaymentMethods       []PaymentMethod    `json:"paymentMethods,omitempty"`
}

type Expiration struct {
	Month string `json:"month,omitempty"`
	Year  string `json:"year,omitempty"`
}

type CardVerification struct {
	Cvv          string                  `json:"cvv,omitempty"`
	AddressLine1 string                  `json:"addressLine1,omitempty"`
	PostalCode   string                  `json:"postalCode,omitempty"`
	AccountName  AccountNameVerification `json:"accountName,omitempty"`
}

type CardAccountUpdater struct {
	UpdatedOn  time.Time `json:"updatedOn,omitempty"`
	UpdateType string    `json:"updateType,omitempty"`
}

type CardDetails struct {
	Status            CardTransactionStatus `json:"status,omitempty"`
	FailureCode       string                `json:"failureCode,omitempty"`
	DynamicDescriptor string                `json:"dynamicDescriptor,omitempty"`
	TransactionSource string                `json:"transactionSource,omitempty"`
	FeeProgram        *string               `json:"feeProgram,omitempty"`
	InitiatedOn       *time.Time            `json:"initiatedOn,omitempty"`
	ConfirmedOn       *time.Time            `json:"confirmedOn,omitempty"`
	SettledOn         *time.Time            `json:"settledOn,omitempty"`
	FailedOn          *time.Time            `json:"failedOn,omitempty"`
	CanceledOn        *time.Time            `json:"canceledOn,omitempty"`
	CompletedOn       *time.Time            `json:"completedOn,omitempty"`

	// Deprecated: This field has been deprecated, use `FeeProgram` instead.
	InterchangeQualification string `json:"interchangeQualification,omitempty"`
}

// CardTransactionStatus represents the status of a card transaction within a Transfer
type CardTransactionStatus string

const (
	// Transaction has been initiated
	CardTransactionStatus_Initiated CardTransactionStatus = "initiated"

	// Transaction has been authorized by the card network
	CardTransactionStatus_Confirmed CardTransactionStatus = "confirmed"

	// Transaction settled with issuer; Moov wallet will be credited by 1 PM ET on a banking day, or the next banking day if it is a weekend or a holiday
	CardTransactionStatus_Settled CardTransactionStatus = "settled"

	// Funds have been credited to the merchant Moov wallet
	CardTransactionStatus_Completed CardTransactionStatus = "completed"

	// Transaction was successfully cancelled and authorization has been reversed
	CardTransactionStatus_Canceled CardTransactionStatus = "canceled"

	// Transaction failed; specific failure reason will be in cardDetails.failureCode
	CardTransactionStatus_Failed CardTransactionStatus = "failed"
)

type CreateCard struct {
	CardNumber        string     `json:"cardNumber,omitempty"`
	CardCvv           string     `json:"cardCvv,omitempty"`
	Expiration        Expiration `json:"expiration,omitempty"`
	HolderName        string     `json:"holderName,omitempty"`
	VerifyName        bool       `json:"verifyName,omitempty"`
	BillingAddress    Address    `json:"billingAddress,omitempty"`
	CardOnFile        bool       `json:"cardOnFile,omitempty"`
	MerchantAccountID string     `json:"merchantAccountID,omitempty"`

	EndToEndToken *EndToEndToken `json:"e2ee,omitempty"`
}

// CreateCard creates a new card for the given customer linked to their account
// https://docs.moov.io/api/#tag/Cards/operation/card
func (c Client) CreateCard(ctx context.Context, accountID string, card CreateCard) (*Card, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPost, pathCards, accountID), AcceptJson(), JsonBody(card), WaitFor("payment-method"))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return UnmarshalObjectResponse[Card](resp)
	case StatusNotFound:
		return nil, errors.Join(ErrAccountNotFound, resp)
	case StatusStateConflict:
		return nil, errors.Join(ErrAlreadyExists, resp)
	default:
		return nil, resp
	}
}

// ListCards lists all cards for the given customer Moov account
// https://docs.moov.io/api/#tag/Cards/operation/listCards
func (c Client) ListCards(ctx context.Context, accountID string) ([]Card, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathCards, accountID))
	if err != nil {
		return nil, err
	}

	return CompletedListOrError[Card](resp)
}

// GetCard retrieves a card for the given customer Moov account
// https://docs.moov.io/api/#tag/Cards/operation/getCard
func (c Client) GetCard(ctx context.Context, accountID string, cardID string) (*Card, error) {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, pathCard, accountID, cardID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Card](resp)
}

type CardUpdateFilter func(*cardPatch) error

type cardPatch struct {
	EndToEndToken *EndToEndToken `json:"e2ee,omitempty"`

	CardCvv        *string       `json:"cardCvv,omitempty"`
	Expiration     *Expiration   `json:"expiration,omitempty"`
	BillingAddress *AddressPatch `json:"billingAddress,omitempty"`
	CardOnFile     *bool         `json:"cardOnFile,omitempty"`
	HolderName     *string       `json:"holderName,omitempty"`
	VerifyName     *bool         `json:"verifyName,omitempty"`
}

func applyCardUpdateFilters(opts ...CardUpdateFilter) (*cardPatch, error) {
	card := &cardPatch{}
	// apply each filter to the card
	for _, opt := range opts {
		if err := opt(card); err != nil {
			return card, err
		}
	}
	return card, nil
}

// WithCardBillingAddress sets the billing address for the card
func WithCardBillingAddress(address AddressPatch) CardUpdateFilter {
	return func(card *cardPatch) error {
		card.BillingAddress = &address
		return nil
	}
}

// WithCardExpiration sets the expiration date for the card
func WithCardExpiration(expiration Expiration) CardUpdateFilter {
	return func(card *cardPatch) error {
		card.Expiration = &expiration
		return nil
	}
}

// WithCardCvv sets the CVV for the card
func WithCardCVV(cvv string) CardUpdateFilter {
	return func(card *cardPatch) error {
		card.CardCvv = &cvv
		return nil
	}
}

// WithCardOnFile sets the card on file for the card boolean
func WithCardOnFile(cardOnFile bool) CardUpdateFilter {
	return func(card *cardPatch) error {
		card.CardOnFile = &cardOnFile
		return nil
	}
}

// WithHolderName sets the holderName for the card
func WithHolderName(holderName string) CardUpdateFilter {
	return func(card *cardPatch) error {
		card.HolderName = &holderName
		return nil
	}
}

// WithVerifyName sets the verifyName flag for the card update request
func WithVerifyName(verifyName bool) CardUpdateFilter {
	return func(card *cardPatch) error {
		card.VerifyName = &verifyName
		return nil
	}
}

// WithVerifyName sets the verifyName flag for the card update request
func WithEndToEndToken(token EndToEndToken) CardUpdateFilter {
	return func(card *cardPatch) error {
		card.EndToEndToken = &token
		return nil
	}
}

// UpdateCard performs a partial update on a card associated with a Moov account. Only fields that
// are populated (non-nil) in the request will be updated. If a value is provided for CVV, a new
// verification ($0 authorization) will be submitted for the card. Updating the expiration date or
// address will update the information stored on file for the card but will not be verified
// https://docs.moov.io/api/#tag/Cards/operation/updateCard
func (c Client) UpdateCard(ctx context.Context, accountID string, cardID string, opt1 CardUpdateFilter, opts ...CardUpdateFilter) (*Card, error) {
	// Create a new CardPatch payload and apply any filters
	opts = append([]CardUpdateFilter{opt1}, opts...)
	payload, err := applyCardUpdateFilters(opts...)
	if err != nil {
		return nil, err
	}

	resp, err := c.CallHttp(ctx, Endpoint(http.MethodPatch, pathCard, accountID, cardID), AcceptJson(), JsonBody(payload))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return UnmarshalObjectResponse[Card](resp)
	default:
		return nil, resp
	}
}

// DisableCard disables a card associated with a Moov account
// https://docs.moov.io/api/#tag/Cards/operation/deleteCard
func (c Client) DisableCard(ctx context.Context, accountID string, cardID string) error {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodDelete, pathCard, accountID, cardID))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
