package moov

import (
	"encoding/json"
	"errors"
	"fmt"
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

type CardPost struct {
	CardNumber        string     `json:"cardNumber,omitempty"`
	CardCvv           string     `json:"cardCvv,omitempty"`
	Expiration        Expiration `json:"expiration,omitempty"`
	HolderName        string     `json:"holderName,omitempty"`
	BillingAddress    Address    `json:"billingAddress,omitempty"`
	CardOnFile        bool       `json:"cardOnFile,omitempty"`
	MerchantAccountID string     `json:"merchantAccountID,omitempty"`
}

type CardPatch struct {
	CardCvv        string     `json:"cardCvv,omitempty"`
	Expiration     Expiration `json:"expiration,omitempty"`
	BillingAddress Address    `json:"billingAddress,omitempty"`
	CardOnFile     bool       `json:"cardOnFile,omitempty"`
}

// CreateCard creates a new card for the given customer linked to their account
// https://docs.moov.io/api/#tag/Cards/operation/card
func (c Client) CreateCard(accountID string, card CardPost) (Card, error) {
	respCard := Card{}

	url := fmt.Sprintf("%s/%s", baseURL, fmt.Sprintf(pathCards, accountID))
	header := map[string]string{"X-Wait-For": "payment-method"}
	body, statusCode, err := c.GetHTTPResponse(http.MethodPost, url, card, header)

	if err != nil {
		return respCard, err
	}

	switch statusCode {
	case http.StatusOK:
		// card created
		err = json.Unmarshal(body, &respCard)
		if err != nil {
			return respCard, err
		}
		return respCard, nil
	case http.StatusUnauthorized:
		return respCard, ErrAuthCredentialsNotSet
	case http.StatusNotFound:
		return respCard, ErrNoAccount
	case http.StatusConflict:
		return respCard, ErrDuplicateLinkCard
	case http.StatusUnprocessableEntity:
		return respCard, ErrCardDataInvalid
	case http.StatusTooManyRequests:
		return respCard, ErrRateLimit
	}
	return respCard, ErrDefault(statusCode)
}

// ListCards lists all cards for the given customer Moov account
// https://docs.moov.io/api/#tag/Cards/operation/listCards
func (c Client) ListCards(accountID string) ([]Card, error) {
	var resCards []Card
	url := fmt.Sprintf("%s/%s", baseURL, fmt.Sprintf(pathCards, accountID))

	body, statusCode, err := c.GetHTTPResponse(http.MethodGet, url, nil, nil)
	if err != nil {
		return resCards, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resCards)
		if err != nil {
			return resCards, err
		}
		return resCards, nil
	case http.StatusUnauthorized:
		return resCards, ErrAuthCredentialsNotSet
	case http.StatusNotFound:
		return resCards, ErrNoAccount
	case http.StatusTooManyRequests:
		return resCards, ErrRateLimit
	}
	return resCards, ErrDefault(statusCode)
}

// GetCard retrieves a card for the given customer Moov account
// https://docs.moov.io/api/#tag/Cards/operation/getCard
func (c Client) GetCard(accountID string, cardID string) (Card, error) {
	resCard := Card{}
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(pathCards, accountID), cardID)

	body, statusCode, err := c.GetHTTPResponse(http.MethodGet, url, nil, nil)
	if err != nil {
		return resCard, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resCard)
		if err != nil {
			return resCard, err
		}
		return resCard, nil
	case http.StatusUnauthorized:
		return resCard, ErrAuthCredentialsNotSet
	case http.StatusNotFound:
		return resCard, ErrNoAccount
	case http.StatusTooManyRequests:
		return resCard, ErrRateLimit
	}
	return resCard, ErrDefault(statusCode)
}

type CardUpdateFilter func(*CardPatch) error

func applyCardUpdateFilters(opts ...CardUpdateFilter) (*CardPatch, error) {
	card := &CardPatch{}
	// Error if no opts provided
	if len(opts) == 0 {
		return card, ErrNoCardUpdateFilters
	}
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
	return func(card *CardPatch) error {
		card.BillingAddress = address
		return nil
	}
}

// WithCardExpiration sets the expiration date for the card
func WithCardExpiration(expiration Expiration) CardUpdateFilter {
	return func(card *CardPatch) error {
		card.Expiration = expiration
		return nil
	}
}

// WithCardCvv sets the CVV for the card
func WithCardCVV(cvv string) CardUpdateFilter {
	return func(card *CardPatch) error {
		card.CardCvv = cvv
		return nil
	}
}

// WithCardOnFile sets the card on file for the card boolean
func WithCardOnFile(cardOnFile bool) CardUpdateFilter {
	return func(card *CardPatch) error {
		card.CardOnFile = cardOnFile
		return nil
	}
}

// UpdateCard Update a linked card and/or resubmit it for verification.
// If a value is provided for CVV, a new verification ($0 authorization) will be submitted for the card. Updating the expiration date or address will update the information stored on file for the card but will not be verified
// https://docs.moov.io/api/#tag/Cards/operation/updateCard
func (c Client) UpdateCard(accountID string, cardID string, opt1 CardUpdateFilter, opts ...CardUpdateFilter) (Card, error) {
	resCard := Card{}
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(pathCards, accountID), cardID)
	// Create a new CardPost payload and apply any filters
	opts = append([]CardUpdateFilter{opt1}, opts...)
	payload, err := applyCardUpdateFilters(opts...)
	if err != nil {
		return resCard, err
	}
	body, statusCode, err := c.GetHTTPResponse(http.MethodPatch, url, payload, nil)

	if err != nil {
		return resCard, err
	}

	switch statusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &resCard)
		if err != nil {
			return resCard, err
		}
		return resCard, nil
	case http.StatusUnauthorized:
		return resCard, ErrAuthCredentialsNotSet
	case http.StatusNotFound:
		return resCard, ErrNoAccount
	case http.StatusConflict:
		return resCard, ErrUpdateCardConflict
	case http.StatusUnprocessableEntity:
		// TODO: parse error message from the body. See https://docs.moov.io/api/sources/cards/update/
		return resCard, ErrCardDataInvalid
	case http.StatusTooManyRequests:
		return resCard, ErrRateLimit
	}
	return resCard, ErrDefault(statusCode)
}

// DisableCard disables a card associated with a Moov account
// https://docs.moov.io/api/#tag/Cards/operation/deleteCard
func (c Client) DisableCard(accountID string, cardID string) error {
	url := fmt.Sprintf("%s/%s/%s", baseURL, fmt.Sprintf(pathCards, accountID), cardID)

	_, statusCode, err := c.GetHTTPResponse(http.MethodDelete, url, nil, nil)
	if err != nil {
		return err
	}

	switch statusCode {
	case http.StatusNoContent:
		// card deleted
		return nil
	case http.StatusUnauthorized:
		return ErrAuthCredentialsNotSet
	case http.StatusNotFound:
		return ErrNoAccount
	case http.StatusTooManyRequests:
		return ErrRateLimit
	}
	return nil
}
