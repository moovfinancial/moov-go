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
	Status                   string            `json:"status,omitempty"`
	FailureCode              string            `json:"failureCode,omitempty"`
	DynamicDescriptor        string            `json:"dynamicDescriptor,omitempty"`
	TransactionSource        string            `json:"transactionSource,omitempty"`
	InterchangeQualification string            `json:"interchangeQualification,omitempty"`
	StatusUpdates            CardStatusUpdates `json:"statusUpdates,omitempty"`
	InitiatedOn              *time.Time        `json:"initiatedOn,omitempty"`
	CompletedOn              *time.Time        `json:"completedOn,omitempty"`
	ConfirmedOn              *time.Time        `json:"confirmedOn,omitempty"`
}

type CardStatusUpdates struct {
	Initiated time.Time `json:"initiated,omitempty"`
	Confirmed time.Time `json:"confirmed,omitempty"`
	Settled   time.Time `json:"settled,omitempty"`
	Failed    time.Time `json:"failed,omitempty"`
	Canceled  time.Time `json:"canceled,omitempty"`
	Completed time.Time `json:"completed,omitempty"`
}

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
