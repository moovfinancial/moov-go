package moov

import (
	"context"
	"errors"
	"net/http"
	"time"
)

var (
	ErrNoCardUpdateFilters = errors.New("no card update filters provided")
	ErrUpdateCardConflict  = errors.New("attempting to update an existing disabled card")
)

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

type CreateCard struct {
	CardNumber        string     `json:"cardNumber,omitempty"`
	CardCvv           string     `json:"cardCvv,omitempty"`
	Expiration        Expiration `json:"expiration,omitempty"`
	HolderName        string     `json:"holderName,omitempty"`
	BillingAddress    Address    `json:"billingAddress,omitempty"`
	CardOnFile        bool       `json:"cardOnFile,omitempty"`
	MerchantAccountID string     `json:"merchantAccountID,omitempty"`
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
		return nil, ErrNoAccount
	case StatusStateConflict:
		return nil, ErrDuplicateLinkCard
	default:
		return nil, resp.Error()
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
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodGet, "/accounts/%s/cards/%s", accountID, cardID), AcceptJson())
	if err != nil {
		return nil, err
	}

	return CompletedObjectOrError[Card](resp)
}

type CardUpdateFilter func(*cardPatch) error

type cardPatch struct {
	CardCvv        string     `json:"cardCvv,omitempty"`
	Expiration     Expiration `json:"expiration,omitempty"`
	BillingAddress Address    `json:"billingAddress,omitempty"`
	CardOnFile     bool       `json:"cardOnFile,omitempty"`
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
func WithCardBillingAddress(address Address) CardUpdateFilter {
	return func(card *cardPatch) error {
		card.BillingAddress = address
		return nil
	}
}

// WithCardExpiration sets the expiration date for the card
func WithCardExpiration(expiration Expiration) CardUpdateFilter {
	return func(card *cardPatch) error {
		card.Expiration = expiration
		return nil
	}
}

// WithCardCvv sets the CVV for the card
func WithCardCVV(cvv string) CardUpdateFilter {
	return func(card *cardPatch) error {
		card.CardCvv = cvv
		return nil
	}
}

// WithCardOnFile sets the card on file for the card boolean
func WithCardOnFile(cardOnFile bool) CardUpdateFilter {
	return func(card *cardPatch) error {
		card.CardOnFile = cardOnFile
		return nil
	}
}

// UpdateCard Update a linked card and/or resubmit it for verification.
// If a value is provided for CVV, a new verification ($0 authorization) will be submitted for the card. Updating the expiration date or address will update the information stored on file for the card but will not be verified
// https://docs.moov.io/api/#tag/Cards/operation/updateCard
func (c Client) UpdateCard(ctx context.Context, accountID string, cardID string, opt1 CardUpdateFilter, opts ...CardUpdateFilter) (*Card, error) {
	// Create a new CardPost payload and apply any filters
	opts = append([]CardUpdateFilter{opt1}, opts...)
	payload, err := applyCardUpdateFilters(opts...)
	if err != nil {
		return nil, err
	}

	resp, err := c.CallHttp(ctx, Endpoint("/accounts/%s/cards/%s", accountID, cardID), AcceptJson(), JsonBody(payload))
	if err != nil {
		return nil, err
	}

	switch resp.Status() {
	case StatusCompleted:
		return UnmarshalObjectResponse[Card](resp)
	case StatusStateConflict:
		return nil, ErrUpdateCardConflict
	case StatusFailedValidation:
		return nil, ErrCardDataInvalid
	default:
		return nil, resp.Error()
	}
}

// DisableCard disables a card associated with a Moov account
// https://docs.moov.io/api/#tag/Cards/operation/deleteCard
func (c Client) DisableCard(ctx context.Context, accountID string, cardID string) error {
	resp, err := c.CallHttp(ctx, Endpoint(http.MethodDelete, "/accounts/%s/cards/%s", accountID, cardID))
	if err != nil {
		return err
	}

	return CompletedNilOrError(resp)
}
